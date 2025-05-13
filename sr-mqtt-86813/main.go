package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	retry "github.com/avast/retry-go/v4"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

// config represents our application configuration.
// It will be populated via environment variables.
type config struct {
	// Comma seperated list of hostname:port pairs
	Brokers []string `default:"localhost:1883" required:"true" desc:"The list of brokers in the cluster" validate:"dive,hostname_port"`

	// The username and password to use to connect to NATS.
	// Again they are populated via environment variables.
	Username string `default:"" required:"false" desc:"username to authenticate with" validate:"required,required_with=Password"`
	Password string `default:"" required:"false" desc:"password to authenticate with" validate:"required,required_with=Username"`

	Interval time.Duration `default:"3s" required:"true" desc:"interval between published messages" validate:"required"`

	Timeouts struct {
		Operations time.Duration `default:"3s" required:"true" desc:"timeout for waitin on operations" split_words:"true" validate:"required"`
		Connection time.Duration `default:"30s" required:"true" desc:"timeout for waitin on operations" split_words:"true" validate:"required"`
	}
}

var (
	// When set to true via command line flags, usage information will be shown
	// and the application will exit.
	usage bool

	// Our actual configuration
	cfg config

	// connectHandler is a callback which will be executed when
	// the client established a connection.
	connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {

		// As soon as we  have a connection, we also
		// want to subscribe to all of our demo topic.
		st := client.Subscribe("demo/#", 1, subscriptionMessageHandler)
		if st.WaitTimeout(cfg.Timeouts.Operations) && st.Error() != nil {
			slog.Error("Subscribing", "err", st.Error())
		}

		c := client.OptionsReader()
		slog.Info("Connected", "servers", c.Servers(), "client", c.ClientID())
	}

	// connectLostHandler is a callback executed when the client loses the connection
	// to the MQTT cluster.
	connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		slog.Warn("Connection lost", "error", err)
	}

	// defaultMessagePubHandler is invoked in case a message is not handled by a dedicated
	// handler tied to the subscription.
	defaultMessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		defer msg.Ack()
		slog.Info("Received message", "topic", msg.Topic(), "message", msg.Payload())
	}

	// subscriptionMessageHandler is the one wee will be using for our subscription.
	subscriptionMessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		defer msg.Ack()
		slog.Info("Received message",
			"handler", "subscriptionMessageHandler",
			"topic", msg.Topic(),
			"id", msg.MessageID(),
			"content", msg.Payload())
	}

	clientopts *mqtt.ClientOptions
)

func init() {

	// Setup the application basics...

	// ...a flag that lets us print out the usage
	// information.
	flag.BoolVar(&usage, "usage", false, "print usage")

	// The static options for our MQTT client.
	clientopts = mqtt.NewClientOptions().
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetClientID("mqtt_demo_client").
		SetResumeSubs(true).
		SetCleanSession(false).
		SetProtocolVersion(4).
		SetOnConnectHandler(connectHandler).
		SetConnectionLostHandler(connectLostHandler).
		SetDefaultPublishHandler(defaultMessagePubHandler).
		SetOrderMatters(true)
}

func main() {

	// parse the flags...
	flag.Parse()

	if usage {
		// If the usage flag was added
		// we only want to see the config options.
		envconfig.Usage("mqtt", &cfg)
		return
	}

	//... and the environment variables.
	envconfig.MustProcess("mqtt", &cfg)

	// We validate out configuration.
	if err := validator.New().Struct(cfg); err != nil {
		slog.Error("Validation failed", "error", err)
		os.Exit(1)
	}

	// Now we have sanitized our input
	// we can set up the client.
	for _, broker := range cfg.Brokers {
		clientopts.AddBroker(fmt.Sprintf("tcp://%s", broker))
	}

	// Authentication is optional, but highly recommended.
	if cfg.Username != "" && cfg.Password != "" {
		clientopts.SetUsername(cfg.Username)
		clientopts.SetPassword(cfg.Password)
	}

	// Create our client...
	client := mqtt.NewClient(clientopts)
	// ...and wait for it to connect AND that we did not
	// have an error during connecetion (say authentication errors)
	retry.Do(
		func() error {
			if token := client.Connect(); token.WaitTimeout(cfg.Timeouts.Connection) && token.Error() != nil {
				slog.Error("connecting to the brokers", "brokers", cfg.Brokers, "error", token.Error())
				return fmt.Errorf("connecting to brokers: %s", token.Error())
			}
			return nil
		},
		retry.Attempts(3),
		retry.DelayType(retry.BackOffDelay),
		retry.OnRetry(
			func(attempt uint, err error) {
				slog.Info("retrying connection to brokers", "brokers", cfg.Brokers, "attempt", attempt, "error", err)
			}),
	)

	defer client.Disconnect(500)

	// We want to be able to listen to signals.
	// This includes CTL+C
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	for {
		select {

		// We received a stopping signal
		case <-ctx.Done():
			slog.Warn("Received shuttdown signal")
			return

		case t := <-time.NewTicker(5 * time.Second).C:
			slog.Info("Trying to send message")
			token := client.Publish("demo/foo", 1, true, fmt.Sprintf("Testmessage: %s", t.Format(time.RFC3339)))
			if token.Wait() && token.Error() != nil {
				slog.Error("Error while publishing message", "err", token.Error())
			}
		}
	}

}

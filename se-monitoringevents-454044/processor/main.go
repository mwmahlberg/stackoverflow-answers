package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// loglevel is  just a wrapper around slog.Level to implement the UnmarshalText method
// for parsing log levels from environment variables.
type loglevel slog.Level

// UnmarshalText implements the encoding.TextUnmarshaler interface for loglevel.
func (l *loglevel) UnmarshalText(text []byte) error {
	switch string(text) {
	case "debug":
		*l = loglevel(slog.LevelDebug)
	case "info":
		*l = loglevel(slog.LevelInfo)
	case "warn":
		*l = loglevel(slog.LevelWarn)
	case "error":
		*l = loglevel(slog.LevelError)
	default:
		return fmt.Errorf("invalid log level: %s", text)
	}
	return nil
}

// timestamp is a wrapper around time.Time to implement the UnmarshalJSON method
// for parsing timestamps from JSON data.
// It converts the timestamp from nanoseconds since epoch to time.Time.
type timestamp time.Time

// UnmarshalJSON implements the json.Unmarshaler interface for timestamp.
func (t *timestamp) UnmarshalJSON(data []byte) error {
	var ts int64
	if err := json.Unmarshal(data, &ts); err != nil {
		return err
	}
	*t = timestamp(time.Unix(0, ts))
	return nil
}

// config holds the configuration for the processor service.
// It uses the envconfig package to load configuration from environment variables.
type config struct {
	LogLevel loglevel `split_words:"true" default:"info" desc:"Log level (debug, info, warn, error)"`
	Nats     struct {
		Host   string `split_words:"true" required:"true" default:"localhost:4222" desc:"NATS server host and port"`
		Client struct {
			Name string `split_words:"true" default:"processor" desc:"NATS client name"`
		}
		Stream struct {
			ClientName string `split_words:"true" default:"processor" desc:"NATS JetStream client name"`
			Name       string `split_words:"true" default:"ping-results" desc:"NATS JetStream stream name"`
		}
	}
}

// pingresult represents the structure of the ping result message.
// It is the Go representation of the data sent by the Telegraf ping plugin via CloudEvents.
type pingresult struct {
	Fields struct {
		AverageResponseMS   float64 `json:"average_response_ms"`
		MaximumResponseMS   float64 `json:"maximum_response_ms"`
		MinimumResponseMS   float64 `json:"minimum_response_ms"`
		PacketsReceived     int     `json:"packets_received"`
		PacketsTransmitted  int     `json:"packets_transmitted"`
		PercentPacketLoss   float64 `json:"percent_packet_loss"`
		ResultCode          int     `json:"result_code"`
		StandardDeviationMS float64 `json:"standard_deviation_ms"`
		TTL                 int     `json:"ttl"`
	}
	Name      string            `json:"name"`
	Tags      map[string]string `json:"tags"`
	Timestamp timestamp         `json:"timestamp"`
}

// processMsg is a function that processes incoming messages from the NATS JetStream.
// It unmarshals the CloudEvent data into a pingresult struct and sends it to the downstream channel.
// If the message cannot be processed, it sends a negative acknowledgment (Nak) to the NATS server.
// This way, the message will be retried later and is not lost.
func processMsg(downstream chan *pingresult) func(msg jetstream.Msg) {
	return func(msg jetstream.Msg) {
		evt := cloudevents.NewEvent()
		if err := evt.UnmarshalJSON(msg.Data()); err != nil {
			msg.Nak()
			slog.Error("Failed to unmarshal CloudEvent", "error", err)
			return
		}
		r := new(pingresult)
		if err := evt.DataAs(r); err != nil {
			msg.Nak()
			slog.Error("Failed to unmarshal CloudEvent data", "error", err)
			return
		}
		slog.Debug("Received CloudEvent", "event", evt.Type(),
			"source", evt.Source(),
			"url", r.Tags["url"],
			"average", r.Fields.AverageResponseMS,
			"packageloss", r.Fields.PercentPacketLoss,
			"subject", msg.Subject(),
		)
		msg.Ack()
		downstream <- r
	}
}

func main() {
	var cfg config
	envconfig.Usage("processor", &cfg)

	// Load configuration from environment variables
	if err := envconfig.Process("processor", &cfg); err != nil {
		slog.Error("Failed to process environment variables", "error", err)
		os.Exit(1)
	}

	// Make sure we have a decent logger
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(cfg.LogLevel),
	})))

	slog.Info("Starting processor service", "cfg", cfg)

	// We need to handle OS signals to gracefully shut down the service
	// This is important for long-running services to avoid data loss
	// and to clean up resources properly.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Connect to the NATS server
	nc, err := nats.Connect(fmt.Sprintf("nats://%s", cfg.Nats.Host), nats.Name(cfg.Nats.Client.Name))
	if err != nil {
		slog.Error("Failed to connect to NATS server", "error", err)
		os.Exit(1)
	}
	defer nc.Close()

	// Create a JetStream context...
	js, err := jetstream.New(nc)
	if err != nil {
		slog.Error("Failed to create JetStream context", "error", err)
		os.Exit(1)
	}

	// ... and connect to the stream of incoming CloudEvents
	stream, err := js.Stream(ctx, cfg.Nats.Stream.Name)
	if err != nil {
		if err != jetstream.ErrStreamNotFound {
			slog.Error("Failed to get stream", "error", err)
			os.Exit(1)
		}
		slog.Error("Stream not found!", "error", err)
		os.Exit(1)
	}

	// Create a consumer for the stream...
	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable: "telegraf"})
	if err != nil {
		slog.Error("Failed to create consumer", "error", err)
		os.Exit(1)
	}

	// ... and consume messages from the stream.
	procChan := make(chan *pingresult, 64)
	// Note that we hand over a callback function to the consumer (the result of processMsg)
	// This function will be called for each message received from the stream.
	cons.Consume(processMsg(procChan))

	// All we have to do now is to process the messages received from the stream.
	// Of course, we could also do this in the callback function, but this way we can
	// separate the processing logic from the message receiving logic.
	go func() {
		for {
			select {
			case msg := <-procChan:
				// Here you can add your processing logic
				// For example, you can send the result to another stream or store it in a database
				slog.Info("Processing ping result", "result", msg)

			case <-ctx.Done():
				// THis will be called when the service received SIGINT or SIGTERM
				slog.Info("Shutting down processor service")
				return
			}
		}
	}()

}

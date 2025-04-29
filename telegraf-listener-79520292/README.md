telegraf-listener-79520292
==========================

This directory contains the code for [my answer][a] to the "question"
[Telegraf HTTP Listener Plugin Returning 400 "Bad Request" for JSON Data with
Measurement Key][q].

Usage
-----

```plaintext
git clone https://github.com/mwmahlberg/stackoverflow-answers.git mwmahlberg-so-answers
cd mwmahlberg-so-answers/telegraf-listener-79520292
```

### Use the agent based setup

```
docker-compose --profile agent up
```

### Use the plugin based setup

```
docker-compose --profile plugin up
```

[q]: https://stackoverflow.com/questions/79520292/telegraf-http-listener-plugin-returning-400-bad-request-for-json-data-with-mea
[a]: https://stackoverflow.com/a/79598177/1296707

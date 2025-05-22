sf-fluentd-logging-1181735
==========================

This folder contains the code to [my answer][myanswer] to the "question"
[docker-compose setup with containers logging to fluent-bit over fluentd driver][question] on [serverfault.com][sf].

- [Usage](#usage)
  - [Clone the git repo](#clone-the-git-repo)
  - [Use FluentBit for log collection](#use-fluentbit-for-log-collection)
  - [Use Grafana Alloy for log collection](#use-grafana-alloy-for-log-collection)


Usage
-----

### Clone the git repo

```plaintext
git clone https://github.com/mwmahlberg/stackoverflow-answers.git mwmahlberg-so-answers
cd mwmahlberg-so-answers/
```

### Use FluentBit for log collection

```plaintext
cd fluentbit
docker compose pull
docker compose up
``` 

### Use Grafana Alloy for log collection

```plaintext
cd fluentbit
docker compose pull
docker compose up
```

[myanswer]: https://serverfault.com/a/1184113/238425
[question]: https://serverfault.com/questions/1181735/docker-compose-setup-with-containers-logging-to-fluent-bit-over-fluentd-driver
[sf]: https://serverfault.com
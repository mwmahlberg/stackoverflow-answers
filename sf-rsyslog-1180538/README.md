sf-rsyslog-1180538
================

This folder contains the code to [my answer][myanswer] to the "question"
[rsyslog: how do I explicitly set a PRI value in a template][question].

Structure
---------

```plaintext
.
├── README.md
├── docker-compose.yaml
├── fluentbit
│   ├── fluentbit.yaml
│   └── priToLevel.lua
├── logger
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
└── rsyslog
    ├── Dockerfile
    └── rsyslog.conf
```

* `fluentbit` contains the configuration files for... well, [fluentbit](https://docs.fluentbit.io/manual)
* `logger` contains the source code for a simple log generator.
* `rsyslog` contains the source filles for a simple docker image with [rsyslog](https://www.rsyslog.com)

Usage
-----

```none
git clone https://github.com/mwmahlberg/stackoverflow-answers.git mwmahlberg-so-answers
cd mwmahlberg-so-answers/sf-rsyslog-1180538
docker-compose build
docker-compose pull
docker-compose up -d
docker-compose logs -f rsyslog
```

[myanswer]: https://serverfault.com/a/1180590/238425
[question]: https://serverfault.com/questions/1180538/rsyslog-how-do-i-explicitly-set-a-pri-value-in-a-template
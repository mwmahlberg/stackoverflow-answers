sf-rsyslog-1180538
================

This folder contains the code to [my answer][myanswer] to the "question"
[rsyslog: how do I explicitly set a PRI value in a template][question]

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

[myanswer]: https://stackoverflow.com/a/79587160/1296707
[question]: https://serverfault.com/questions/1180538/rsyslog-how-do-i-explicitly-set-a-pri-value-in-a-template
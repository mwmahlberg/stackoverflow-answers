se-monitoringevents-454044
==========================

Usage
-----

```none
git clone https://github.com/mwmahlberg/stackoverflow-answers.git mwmahlberg-so-answers
cd mwmahlberg-so-answers/se-monitoringevents-454044/
docker-compose up -d
docker-compose logs -f processor
```

Build
-----

### Prerequisites

* [ko](https://ko.build)
* Go 1.24

### Build the image

```none
cd processor
KO_REGISTRY=$YOUR_REGISTRY/$YOUR_DOCKER_USER_OR_ORG ko build
```

### Build the processor locally

```none
cd processor
go build [-o processor]
```

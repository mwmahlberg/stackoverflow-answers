dockerize-templates-39243389
============================

This directory contains the code for my [answer][a] to the "question"
[Templating config file with docker][q].

Usage
-----

The images are available on DockerHub. You should not need to build them.

### Embedded Template

```plaintext
docker run \
--env LANGS="de, tlh" \
--rm \
docker.io/mwmahlberg/so-dockerize-templates-39243389:embed \
cat /etc/sphinx.cfg
```

### Template NOT embedded

```plaintext
git clone https://github.com/mwmahlberg/stackoverflow-answers.git mwmahlberg-so-answers
cd mwmahlberg-so-answers/dockerize-templates-39243389
docker run \
--env LANGS="de, tlh" \
--rm \
-v $(pwd)/sphinx.tmpl:/usr/local/share/app/sphinx.cfg `# note the different name` \
docker.io/mwmahlberg/so-dockerize-templates-39243389:no-tmpl \
-template /usr/local/share/app/:/etc/ `# Note the source and target directories` \
cat /etc/sphinx.cfg
```

Build
------------

### Prerequisites

* A working container engine. The [Makefile](./Makefile) is likely to run only
   with docker or podman.
* make. Strictly speaking optional, but not using it makes builds tedious.

### Make targets

| Target     | Description                              | Default            |
| ---------- | ---------------------------------------- | :----------------- |
| `all`      | builds all images                        | :white_check_mark: |
| `no-embed` | builds  image without  embedded template | :x:                |
| `embed`    | builds image with embedded template      | :x:                |
| `clean`    | removes temporary build files            | :x:                |
| `squeaky`  | `clean` + removes images                 | :x:                |

### Environment Variables

| Name          | Description                                                                         | default value                                           |
| ------------- | ----------------------------------------------------------------------------------- | ------------------------------------------------------- |
| `DOCKER`      | Path to `docker` binary                                                             | Path to docker binary within `$PATH`, none if not found |
| `REGISTRY`    | docker registry to  use                                                             | `docker.io`                                             |
| `ORG`         | the docker organisation, most likely your username                                  | `$USER`                                                 |
| `REPO_PREFIX` | well... a prrefix for `$REPOSITORY`                                                 | ""                                                      |
| `REPOSITORY`  | the dockerhub repository                                                            | name of the directory containing the Makefile           |
| `IMAGE`       | the fully qualified image name to use. Overrides `REGISTRY`, `ORG` and `REPOSITORY` | `$(REGISTRY)/$(ORG)/$(REPO_PREFIX)${REPO}`              |

You can either set them per run

```shell
REPO_PREFIX="my-" make
```

or use an exported value

```shell
# set this for example in ~/.bashrc
export REPO_PREFIX="coolest-"
make
```

[a]: https://stackoverflow.com/a/79603794/1296707
[q]: https://stackoverflow.com/questions/39243389/templating-config-file-with-docker

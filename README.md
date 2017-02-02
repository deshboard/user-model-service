# User Model Service


## Setup

Before starting any Docker containers you might want to setup your own port forwarding rules based on `docker-compose.override.yml.example`.

``` bash
$ cp docker-compose.override.yml.example docker-compose.override.yml
```

Run the following commands to prepare and start the Docker environment:

``` bash
$ make setup
$ make start
```

Then the following ones to install dependencies locally:

``` bash
$ make install
```

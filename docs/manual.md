# User Manual

Hello is a simple hello world application that demonstrates the full dev workflow for 
testing, code reviewing, documentation and releasing software.

This user manual is 'meta' too, it mostly exists to demonstrate the structure of the documentaton.

We use [MkDocs](http://mkdocs.org) for our documentation workflow. 

You can serve and live reload documentation on ``localhost:8080`` while working on it

```bash
$ pip install mkdocs
$ make doc
```

## Glossary

Glossary usually consists of the terms that reader needs to know before starting to use the software, abbreviations or special software.

**Greeting**

Hello server stores greetings in a database, e.g. `Hello` or `Hola`. Greetings are identified by IDs, e.g. `hello.us` or `hello.sp`.

**Backend**

Backend is usually a database backend, such as Cassandra or Etcd

**Etcd**

Etcd is a discovery and presence database released by CoreOS and is used by hello app as a simple configuration backend.

**CLI**

CLI is a command line interface, hello app provides a special command line tool called `hctl`


## Installation

Installation section is very important for users to adopt the software. Ideally it should be as extensive as possible and should feature
various platforms.

The quickest way to install hello is to download the latest release:

```bash
$ export VERSION=v0.0.1-alpha.1

# download the latest release
$ cd /tmp
$ wget https://github.com/gravitational/hello/releases/download/$VERSION/hello-$VERSION-linux-amd64.tar.gz
$ tar -xzf hello-$VERSION-linux-amd64.tar.gz

# start a hello server listening on localhost:23456 
# and using etcd as a configuration backend
$ /tmp/hello -addr=localhost:23456\
            -backend=etcd\
            -backendConfig='{
                   "nodes": ["http://localhost:4001"], 
                   "key": "/hello"}'

```

## Overview

It is sometimes useful go give a full overivew of the operation of the software. 

Diagrams here may be appropriate.

![Overview](/img/Hello.png)

We are using [Draw IO](http://draw.io) to design diagrams and are following common style to ensure consistency in the docs.

You can open up the template with common collors and styles [here](/img/Hello.xml) using draw.io and start editing from there.

## Manage greetings

*Note*: This section is an application specific section and usually explains how to use the software in detail.

Users can manage greetings via API or CLI tool.

**Upsert a greeting**

**Important:** Make sure your samples actually work. The best way to ensure this is to execute them while writing.

```bash
# upsert a greeting via CLI
$ hctl -hello=http://localhost:23456 greeting upsert -id=hello.us -val=Hello

# upsert a greeting via API
$ curl -v -X POST -d prompt=hello.us -d value=Howdy http://localhost:23456/v1/greetings
```

**Get a greeting by ID**

```bash
# CLI
$ hctl -hello=http://localhost:23456 greeting get -id=hello.us

# API
$ curl http://localhost:23456/v1/greetings/hello.us
```

**Delete a greeting by ID**

```bash
# CLI
$ hctl -hello=http://localhost:23456 greeting delete -id=hello.us

# API
$ curl -X DELETE http://localhost:23456/v1/greetings/hello.us
```

## Saying Hello

```bash
# Say hello to Dog
$ hctl -hello=http://localhost:23456 hello -id=hello.us -name=Dog
OK: Hello, Dog!

# Say hello via API
curl -v -X POST -d prompt=hello.us -d name=Dog http://localhost:23456/v1/hello
{"val":"Hello, Dog!"}
```

## Operation

Operation section is important to understand what options are needed to run the service in production.

Here are the flags for running a hello server binary:

```bash
# addr sets a host and port for Hello server
-addr=localhost:23456 

# log output, the current supported are 'console' or 'syslog'
-log=console

# logging severity threshold, e.g. 'INFO', 'WARN' or 'ERROR'
-logSeverity=INFO

# backend type, currently only 'etcd'
-backend=etcd

# backendConfig is a backend-specific configuration string, e.g.
# etcd configuration server list and key
-backendConfig='{
   "nodes": ["http://localhost:4001"], 
   "key": "/teleport"}'
```

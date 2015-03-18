# Developers Guide

## Introduction

This guide will introduce you to the full development cycle, tooling used at Gravitational Inc.

## Development flow

We are following [Github Flow](https://guides.github.com/introduction/flow/). 
Please take some time to view it, it really takes 5 minutes and very simple.

## Managing dependencies

If you are introducing new dependency, you should vendor it if the project is the executable using [Godeps](https://github.com/tools/godep).
Note that we are using ``godep save -r ./..`` to save all dependencies in the Godeps folder. 
Here's the [article](http://blog.vulcanproxy.com/deps/) that explains vendoring with Godeps in more detail.

**Note** Libraries should not use vendoring, because end-user executables will re-vendor all deps anyways.


## Testing

### Test suite

We use [Gocheck](https://labix.org/gocheck) for writing rich test suites.

Check out a [complete example](https://github.com/gravitational/hello/blob/master/hello_test.go) of the test suite using the framework.

#### Mocking database

There's no way to easilly mock anything in Go, and some devs may be confused by this fact. However, you have Go interfaces to implement 
common functions for accessing database and implement in-memory backend.

Here's an example of how do we do it for Etcd.

** Backend interface**

We start at setting up the proper interface that describes the feature set for CRUD operations and some error functions:

```c
// GreetingBackend is an interface to the backend (usually a database)
// that provides some storage functionality.
type GreetingBackend interface {
    // UpsertGreeting updates or inserts the greeting into the database
    UpsertGreeting(id, val string) error
    // GetGreeting returns a greeting stored in a database by it's id
    GetGreeting(id string) (string, error)
    // DeleteGreeting deletes greeting by ID
    DeleteGreeting(id string) error
    // Close closes all resources associated with this backend
    Close() error
}

// NotFoundError returns whenever the greeting requested is not found
type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("greeting with id '%v' not found", e.ID)
}
```

See the complete example [here](https://github.com/gravitational/hello/blob/master/backend/backend.go)

** Test suite **

To guarantee that all backends provide the same behavior, we write an acceptance suite, that is implementation agnostic:

```c
package test

import (
   "testing"
   . "github.com/gravitational/hello/Godeps/_workspace/src/gopkg.in/check.v1"
   "github.com/gravitational/hello/backend"
)

func TestBackend(t *testing.T) { TestingT(t) }

type BackendSuite struct {
   B backend.GreetingBackend
}

// GreetingsCRUD tests simple CRUD cycle for greetings
func (s *BackendSuite) GreetingCRUD(c *C) {
   c.Assert(s.B.UpsertGreeting("hello.us", "Hello"), IsNil)
}
```

Note, that this is not a test suite, but rather a library for writing tests with common test-cases. 

See the complete example [here](https://github.com/gravitational/hello/blob/master/backend/test/suite.go)

** Memory backend and Etcd backend**

Once we get the acceptance suite, we can implement the Memory backend and Etcd backends:

* [Memory Backend](https://github.com/gravitational/hello/blob/master/backend/membk)
* [Etcd Backend](https://github.com/gravitational/hello/blob/master/backend/etcdbk)

Now, when writing tests for Memory Backend, you can reuse the test-suite we wrote before:

```c
type MemSuite struct {
   bk *MemBackend
   // this is the test suite for acceptance
   suite test.BackendSuite 
}

func (s *MemSuite) SetUpTest(c *C) {
  s.bk = New()
  // init the test suite with the memory backend
  s.suite.B = s.bk 
}

func (s *MemSuite) TearDownTest(c *C) {
  c.Assert(s.bk.Close(), IsNil)
}

func (s *MemSuite) TestGretingsCRUD(c *C) {
   s.suite.GreetingCRUD(c)
}
```

So all we do in Memory backend test suite is using the acceptance suite. 
We can add backend-specific initialization to the test suite too, see for example:

```c
func (s *EtcdSuite) SetUpSuite(c *C) {
   nodes_string := os.Getenv("TEST_ETCD_NODES")
   if nodes_string == "" {
     // Skips the entire suite
     c.Skip("This test requires etcd")
     return
   }
   s.nodes = strings.Split(nodes_string, ",")
}
```

In this example we execute the test only if the test caller has set up path to the running Etcd cluster.

[Full example](https://github.com/gravitational/hello/blob/master/backend/etcdbk/etcd_test.go#L29)

### Executing tests

Executing tests in Go is a matter of calling the go test command line tool:

```bash
go test -v ./... -cover
```

Each project provides a make file that has common targets used by CI/CD:

```bash
# executes tests that don't require any external deps
# usually they use Memory backends to mock data base
make test

# executes tests that imply that all requirements for 
# this project are met, e.g. Etcd and Cassandra are running
make test-with-deps 

# executes coverage tests for package name
make cover-package p=./package-name

# executes coverage tests for package name assuming external deps are met
make cover-package-with-deps p=./package-name

# executes benchmark for the package
make bench-package p=./package-name
```

Take a closer look at [Hello's makefile](https://github.com/gravitational/hello/blob/master/Makefile) for more details.

## CI pipeline

For CI we are using [Shippable](www.shippable.com). Here's how you set up your repo to work with shippable.

### Shippable.yaml

Put shippable.yaml file in the root of your source code repository:

```yaml
language: go

go:
  - 1.4  # build matrix, include desired version of go

build_image: shippableimages/ubuntu1404_go # docker build image

before_install:
  - source $HOME/.gvm/scripts/gvm  # go version manager
  - gvm install go$SHIPPABLE_GO_VERSION  # installs the right go version
  - gvm use go$SHIPPABLE_GO_VERSION  # sets it up
  - export GOPATH=/root/workspace/   # setting up proper GOPATH
  - go get golang.org/x/tools/cmd/cover  # installs coverage and
  - go get golang.org/x/tools/cmd/vet    # vet in case if not vendored

script:
  - make test   # See Makefile for details
```

### Shippable.com

Turn on the repository builds in shippable:

![Overview](/img/shippable.png)

## Code reviews

We are using [github flow](https://guides.github.com/introduction/flow/) for the code reviewing. 
Once you are ready to submit your changes for the code review, update the PR mentioning the code-reviewers 
and asking them for review explicitly. It is our common goal to provide timely and good code reviews to each other, so make sure
code reviewing is a high priority for you.

Here are some examples of good and bad code review requests and responses:

** Less-than-ideal PR **

![badpr](/img/badpr.png)

Let's take a closer look at what should be changed with this PR and the code review:

* PR is poorly documented, it's hard to tell why this change took place
* It does not refer to an issue too
* The two commits in this PR are not necessary to reflect the change logic, they should be squashed
* Code reviewer uses obscene words and attacks the developer directly, instead of providing polite, constructive and friendly critique on the solution
* PR does not pass the CI and the tests are broken

** Good PR **

![goodpr](/img/goodpr.png)

This PR and code review is better:

* It has the ``Purpose`` and ``Implementation`` sections completed that help to understand the change
* It refers to a github issue for more context
* The commits are split logically, commit with Godeps is a separate one to simplify code review
* Code reviewer is polite, focuses on the solution and provides constructive feedback
* The tests are passing and CI tool gives a green light


## Releases

Once code reviewers replied with their "looks good to me", you can merge PR into master branch.

### Branches and tags

Master branch is used for the dev team to see the latest changes that made it into the code. 
We are not using branches to deploy or release software, instead, we are using [github releases](https://github.com/blog/1547-release-your-software).
Releases version should follow [SemVer](http://semver.org/). Also, make sure the CHANGELOG.md is up-to-date too.

### Release script

We currently use Docker for releasing software, here's the step by step guide for relase:

**Step 1. Base image**

We will build a base Docker image, gravitational/release. This image contains your private credentials and should be kept locally on your machine:

```
FROM golang:1.4

RUN go get github.com/aktau/github-release
ENV GITHUB_TOKEN <github release token>
ENV USERNAME <github username>
ENV PASSWORD <github password>
```

You can generate new token for your deployment on github ``Applications`` section:

![token](/img/token.png)

Once it's all done, feel free to build a base build image:

```bash
sudo docker build -t gravitational/release .
```

If the build is successfull, you will see an image listed in your docker images:

```bash
$ sudo docker images

REPOSITORY              TAG                 IMAGE ID            CREATED             VIRTUAL SIZE
<none>                  <none>              99c612fee9d6        22 hours ago        523.3 MB
gravitational/release   latest              f2e60ece82c2        22 hours ago        523.3 MB

```

** Step 2. Trigger release **

The second Docker file, that triggers the release is located in ``release folder`` and in fact is very simple:

```bash
FROM gravitational/release

ADD . /tmp
RUN bash /tmp/release.sh v0.0.1-alpha.1 master CHANGELOG.md
```

It pulls the base image that we've just build and triggers the build script supplying the release version via command line arguments. It will cross compile the binaries and push them to github.

```bash
$ cd $GOPATH/hello/releases
$ sudo docker build .
```
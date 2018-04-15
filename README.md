[![Build Status](http://zanetworkercicd.eu.ngrok.io/api/badges/zanetworker/docktorino/status.svg?branch=master)](http://zanetworkercicd.eu.ngrok.io/api/badges/zanetworker/docktorino/status.svg?branch=master)


# Docktorino 

Docktorino is a real-time continuous testing tool for your docker builds, helping you containerize with confidence!

When building Docker images it is sometimes tricky to assert the image behavior, for example, whether the contents of the image you built is correct, or that commands can execute correctly inside your container ( maybe you forgot to set your binary's PATH).

Docktrino makes these types of assertions easy to define when building your Docker images. It then listens to your Docker builds and triggers these tests on the spot to notify you if you have done something wrong, or if your container image is misbehaving.

Check blog post [here](http://www.adelzaalouk.me/2018/docktorino/). 

<!-- TOC -->

- [Docktorino](#docktorino)
	- [Getting Started](#getting-started)
		- [Clone and build](#clone-and-build)
		- [Download from the releases](#download-from-the-releases)
	- [Usage](#usage)
	- [Contributing](#contributing)
	- [Authors](#authors)
	- [License](#license)

<!-- /TOC -->

## Getting Started

To get started with `Docktorino`, you can download the corresponding binary for your OS (Darwin, Linux, Windows) from the releases section. Or you can clone this repository and build the project locally.

### Clone and build

You need to install `make` and `Go` on your system before proceeding.

```bash
git clone https://github.com/zanetworker/docktorino.git
cd docktorino

# build docktorino binary if you have go installed
make OS=<darwin|linux|windows> install

# execute docktorino for command overview
docktorino

# build docktorino binary 
make OS=<darwin|linux|windows> dry

# execute docktorino command for overview
./docktorino
```

### Download from the releases

The binary is available on the releases section of this repo. Please note that it has been only tested on OSX.


## Usage

Docktorino uses project [Dockument](https://github.com/zanetworker/dockument) to fetch important information about the image in questions, in particularly, testing info. an Example Dockerfile would look as follows: 

```dockerfile
FROM golang:alpine as builder
 
###  Describe container tests based on container structure tests from google!
LABEL api.TEST.command=""\
      api.TEST.command.name="go version"\
      api.TEST.command.command="go"\
      api.TEST.command.args="version"\
      api.TEST.command.expectedOutput="go version"

# Add file and check it's contents!
LABEL api.TEST.fileExistence=""\
      api.TEST.fileExistence.name="Dockumentation Check"\
      api.TEST.fileExistence.path="/dockumentation.md"\
      api.TEST.fileExistence.shouldExist="true"\
      api.TEST.fileExistence.permissions=""
ADD dockumentation.md /dockumentation.md


# Assert that environment variables are correctly defined!
LABEL api.TEST.metadata=""\
      api.TEST.metadata.env="GOPATH=/go,PATH=/go/bin:/usr/local/go/bin:$PATH"\
      api.TEST.metadata.exposedPorts=""\
      api.TEST.metadata.volumes=""\
      api.TEST.metadata.cmd=""\
      api.TEST.metadata.workdir=""
ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

``` 

This Dockerfile defines three tests: 
- Command test to check the output of `go version`. 
- a FileExistence tests to make sure that a document we add to the image is where it should be.
- a metadata tests to make sure that Environment variables are correctly defined within the image.

These tests styles are based on the [GCP container-structure-tests](https://github.com/GoogleCloudPlatform/container-structure-test) project tests. 

Currently there is only one command available, this command starts `Docktorino` for the image of choice. For example `docktorino -i "repo/image:tag"` will start listening for new builds for "repo/image:tag" and trigger the tests defined in the Dockerfile.

```
Usage:
  docktorino start [flags]

Flags:
  -h, --help           help for start
  -i, --image string   the image you wish to trigger tests for
  -q, --quiet          quiet surpress output
  -v, --verbose        verbose testing output

Global Flags:
      --home string   location of your docktorino config. Overrides $DOCKTORINO_HOME  (default "/Users/adelias/.Docktorino"
``` 

For a complete demo please check this demo [video](https://youtu.be/lU7hpP2nfPw):

[![Docktorino Demo](./demo/2018-04-04-11-05-16.png)](https://youtu.be/lU7hpP2nfPw)

<!-- ![](./demo/docktrino_demo.gif) -->

## Contributing

<!-- [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) -->
Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Authors

See also the list of [contributors](https://github.com/zanetworker/dockument/graphs/contributors) who participated in this project.

## License

This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details

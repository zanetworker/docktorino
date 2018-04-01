FROM golang:alpine as builder
 
###  Describe container tests based on container structure tests from google
LABEL api.TEST.command=""\
      api.TEST.command.name="go version"\
      api.TEST.command.command="go"\
      api.TEST.command.args="version"\
      api.TEST.command.expectedOutput="go version"

# LABEL api.TEST.fileContent=""\
#       api.TEST.fileContent.name="Debian Sources"\
#       api.TEST.fileContent.path="/etc/apt/sources.list"\
#       api.TEST.fileContent.expectedContents="['.*httpredir\\.debian\\.org.*']"\
#       api.TEST.fileContent.excludedContents="['.*gce_debian_mirror.*']"


# Add file and check it's contents
LABEL api.TEST.fileExistence=""\
      api.TEST.fileExistence.name="Dockumentation Check"\
      api.TEST.fileExistence.path="/dockumentation.md"\
      api.TEST.fileExistence.shouldExist="true"\
      api.TEST.fileExistence.permissions=""
ADD dockumentation.md /dockumentation.md


LABEL api.TEST.metadata=""\
      api.TEST.metadata.env="GOPATH=/go,PATH=/go/bin:/usr/local/go/bin:$PATH"\
      api.TEST.metadata.exposedPorts=""\
      api.TEST.metadata.volumes=""\
      api.TEST.metadata.cmd=""\
      api.TEST.metadata.workdir=""
ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go



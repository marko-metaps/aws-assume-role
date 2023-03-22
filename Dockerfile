FROM golang

ENV GOPATH=

ENTRYPOINT ["go", "build", "-o", "bin/aws-assume-role"]

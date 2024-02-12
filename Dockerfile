#FROM ubuntu:latest
#LABEL authors="ivan"
#
#ENTRYPOINT ["top", "-b"]

FROM golang:1.21

WORKDIR /server
COPY server/ ./


## Download all the dependencies
#RUN go get -d -v ./...
#
## Install the package
#RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8081 8082

# Run the executable
CMD ["go run server.go"]
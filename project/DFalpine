FROM golang:1.9.3-alpine3.7 AS artifact
WORKDIR /go
ADD . /go
RUN apk add --no-cache git
ENV GOBIN "/go/bin"
RUN go env
RUN go get -v && rm bin/go
RUN go build -v -o bow
RUN chmod +x bow

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /evedel
COPY --from=artifact /go/bow /evedel/bow
COPY ./templates /evedel/templates
COPY ./resources /evedel/resources
CMD	 ["/evedel/bow"]

FROM golang:alpine
MAINTAINER Jermine.hu@qq.com
ENV APP_HOME /go/src/bow
WORKDIR $APP_HOME
COPY . $APP_HOME
RUN go env && go build -v -o bow && chmod +x bow && pwd && ls
FROM alpine:edge
RUN apk add --no-cache ca-certificates
WORKDIR /evedel
# Get a file from first floor image
COPY --from=0 /go/src/bow/bow .
COPY templates /evedel/templates
COPY resources /evedel/resources
CMD	 ["./bow"]

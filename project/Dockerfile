FROM golang:alpine
MAINTAINER Jermine.hu@qq.com
ENV APP_HOME /go/src/bow
WORKDIR $APP_HOME
COPY . $APP_HOME
RUN go env && go build -v -o bow && chmod +x bow && pwd && ls
FROM alpine:edge
RUN apk add --no-cache ca-certificates
# Get a file from first floor image
COPY --from=0 $APP_HOME/bow /evedel/bow
COPY templates /evedel/templates
COPY resources /evedel/resources
WORKDIR /evedel
CMD	 ./bow

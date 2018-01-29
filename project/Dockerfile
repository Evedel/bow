FROM golang:alpine
MAINTAINER Jermine.hu@qq.com
ENV APP_HOME /go/src/bow
WORKDIR $APP_HOME
COPY . $APP_HOME
RUN go env && go build -v -o bow && chmod +x bow
FROM     alpine:edge
RUN apk add --no-cache ca-certificates
ADD       bow /evedel/bow
ADD	      templates /evedel/templates
ADD	      resources /evedel/resources
WORKDIR	  /evedel
CMD	      ./bow

FROM     alpine:edge
RUN apk add --no-cache ca-certificates
ADD       main /evedel/bow
ADD	      templates /evedel/templates
ADD	      resources /evedel/resources
WORKDIR	  /evedel
CMD	      ./bow

FROM      debian:jessie
RUN       apt-get update && apt-get install -y ca-certificates
ADD       main /evedel/bow
ADD	      templates /evedel/templates
ADD	      resources /evedel/resources
WORKDIR	  /evedel
CMD	      ./bow

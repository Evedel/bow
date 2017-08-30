FROM	ubuntu:16.04

RUN		apt-get update &&\
			apt-get -y upgrade &&\
			apt-get install -y git nano wget &&\
			wget https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz &&\
			tar -xvf go1.9.linux-amd64.tar.gz &&\
			rm go1.9.linux-amd64.tar.gz &&\
			mv go /usr/local &&\
			apt-get clean autoremove

WORKDIR	/app

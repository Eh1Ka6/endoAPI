# EndoAPI

Entry level web API for Endocode. Probably nothing surprising but still fell creative

## Prerequisites

You must have a recent version of go installed. 
This API is meant to run on docker. If you have'nt installed it yet take a look at :   
 

https://docs.docker.com/install/linux/docker-ce/debian/

 
## Getting Started

To get started clone the repository in your home directory and enter the cloned repository

```
git clone https://github.com/Eh1Ka6/endoAPI.git && cd endoAPI
```

### Installing

Simply run make to run the tests and build the API

```
make 
```

### Running the api 

The App take's a listening port as argument (default 8080) and allow you  to  define some environment var

```
./endoApi -p 8080 GOPATH=/GO
```
### Deployment on docker

The Makefile is written to  support docker-build parameter. You'll probably need root privileges if your user is not in the docker group 

```
sudo su && make docker-build

```
## Versioning

Request the localhost:{yourport}/version to  get the actual build version 

## Authors

* **Ehua Kassi** -https://github.com/Eh1Ka6



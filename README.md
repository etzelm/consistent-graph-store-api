# Scalable, Fault Tolerant, &amp; Consistent Graph Store API

## Input Format Specifications

- graph names
  - charset: [a-zA-Z0-9_] i.e. Alphanumeric including underscore, and case-sensitive 
  - size:    1 to 250 characters

- vertex names
  - charset: [a-zA-Z0-9_] i.e. Alphanumeric including underscore, and case-sensitive 
  - size:    1 to 250 characters

- edge names
  - charset: [a-zA-Z0-9_] i.e. Alphanumeric including underscore, and case-sensitive 
  - size:    1 to 250 characters

## Environment Variables Used

- _"SERVERS"_ is used to keep track of all other active server hosts in our system
- _"IP"_ is used to store the docker network ip used for system inter-communication
- _"PORT"_ is used to store the local network port exposed by docker for the user
- _"R"_ is used to store the maximum number of server hosts a partition can be assigned

## Generate gservice with protoc

- protoc -I gservice/ gservice/gservice.proto --go_out=plugins=grpc:gservic

## Example Docker Commands

Starting a system with 4 active server hosts and a maximum partition size of 2:

- docker run -p 3001:3000 --ip=10.0.0.21:3000 --net=mynet -e IP="10.0.0.21" -e PORT="3001" -e R=2 -e SERVERS="10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000" mycontainer
- docker run -p 3002:3000 --ip=10.0.0.22:3000 --net=mynet -e IP="10.0.0.22" -e PORT="3002" -e R=2 -e SERVERS="10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000" mycontainer
- docker run -p 3003:3000 --ip=10.0.0.23:3000 --net=mynet -e IP="10.0.0.23" -e PORT="3003" -e R=2 -e SERVERS="10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000" mycontainer
- docker run -p 3004:3000 --ip=10.0.0.24:3000 --net=mynet -e IP="10.0.0.24" -e PORT="3004" -e R=2 -e SERVERS="10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000" mycontainer

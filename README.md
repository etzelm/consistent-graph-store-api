# Scalable, Fault Tolerant, &amp; Consistent Graph Store API

## Introduction

  The goal of this project is to provide a REST-accessible graph storage service that 
runs on port 3000 and is available as a resource named gs. For example, the service 
would listen at http://server-hostname:3000/gs. We want to develop distributed system 
software to support this service so that it can store an amount of data that would 
not normally fit onto a single machine system. To accomplish this, we will simulate 
our server code as if it is being run on multiple, separate hosts simultaneously, 
using Docker to provide this functionality. A single server host in our system stores 
only a certain subset of the graphs stored in the system as a whole. We also have 
them keep track of a list of all the other server hostnames in the known system so 
that they can forward requests they receive for graphs that arent stored locally for 
them. The plan is to distribute graphs among partitions that each have an active 
amount of server hosts assigned to them based on the total number of server hosts 
that exist in the system at the time of observation. This way each server host in a 
partition can store the same subset of graphs assigned to that partition, providing 
a measurable amount of fault-tolerance to the user if one of those hosts happens to 
crash or experience a network partition. 

  Scalability is achieved by allowing for the
user to change the system environment by adding or removing server hosts, based on 
their needs, using API calls which then have our distributed system software 
automatically reshuffle our partitioning and graph distribution across all active 
server hosts to attain maximum fault-tolerance and minimize access latency. To ensure 
strong consistency among server hosts in a partition that stores the same subset of 
graphs in our system, we will use an algorithm called Raft that uses a 2 phase commit 
sequence and timers to achieve consensus on a total causal order over any value given 
to us by the user. Due to the CAP theorem, we know that using partitions to attain 
fault tolerance means we cannot have a graph store that is both highly available and 
strongly consistent. In this project, we will favour strong consistency over having 
our system be highly available, meaning our service should only return responses to 
requests if it can guarantee that it is using the most recent data available to it.

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

- _"PARTITIONS"_ is used to keep track of all other active server hosts in our system
- _"IP"_ is used to store the docker network ip/port used for system inter-communication
- _"PORT"_ is used to store the local network port exposed by docker for the user
- _"R"_ is used to store the maximum number of server hosts a partition can be assigned

## Example Docker Commands

Starting a system with 4 active server hosts and a maximum partition size of 2:

- docker run -p 3001:3000 --ip=10.0.0.21:3000 --net=mynet -e IP="10.0.0.21:3000" -e PORT="3001" -e R=2 -e PARTITIONS="10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000" mycontainer
- docker run -p 3002:3000 --ip=10.0.0.22:3000 --net=mynet -e IP="10.0.0.22:3000" -e PORT="3002" -e R=2 -e PARTITIONS="10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000" mycontainer
- docker run -p 3003:3000 --ip=10.0.0.23:3000 --net=mynet -e IP="10.0.0.23:3000" -e PORT="3003" -e R=2 -e PARTITIONS="10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000" mycontainer
- docker run -p 3004:3000 --ip=10.0.0.24:3000 --net=mynet -e IP="10.0.0.24:3000" -e PORT="3004" -e R=2 -e PARTITIONS="10.0.0.21:3000,10.0.0.22:3000,10.0.0.23:3000,10.0.0.24:3000" mycontainer

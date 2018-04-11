# Scalable, Fault Tolerant, &amp; Consistent Graph Store API

## Introduction

The goal of this project is to provide a REST-accessible graph storage service that 
runs on port 3000 and is available as a resource named gs. For example, the service 
would listen at http://server-hostname:3000/gs. We want to develop distributed system 
software to support this service so that it can store the amount of data that would 
not normally fit onto a single-machine system. To accomplish this, we will simulate 
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
crash or experience a network partition. Scalability is achieved by allowing for the
user to change the system environment by adding or removing server hosts based on 
their needs using API calls which then have our distributed system software 
automatically reshuffle our partitioning and graph distribution across all active 
server hosts to attain maximum fault-tolerance and minimize access latency. To ensure 
strong consistency among server hosts in a partition that stores the same subset of 
graphs in our system, we will use an algorithm called Raft that uses a 2 phase commit 
sequence and timers to achieve consensus on a total order over any value given to us 
by the user.

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
- _"IP"_ is used to store the docker network ip and port used to inter-communicate
- _"PORT"_ is used to store the local network port exposed by server node for the user

## Partitioning Algorithms Implemented

temp

## Consistentency Algorithms Implemented

We plan to implement the RAFT algorithm to ensure consistency among the data stored 
in our distributed system.

## Technologies Used

gRPC/protocol buffers, gin server code
badges(pictures) to come later
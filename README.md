# Scalable, Fault Tolerant, &amp; Consistent Graph Store API

## Introduction

The goal of this project is to provide a REST-accessible graph storage service that 
runs on port 3000 and is available as a resource named gs. For example, the service 
would listen at http://server-hostname:3000/gs. We want to develop distributed system 
software to support this service so that it can store the amount of data that would 
not normally fit onto a single-machine system. To accomplish this, we will have a 
single server host in our system store only a certain subset of the graphs stored in 
the system as a whole. We will also have them keep track of a list of all the other 
server hostnames in the known system so that they can forward requests they receive 
for graphs that arent stored locally for them. The plan is to distribute graphs among 
partitions that each have an active amount of server hosts assigned to them based on 
the total number of server hosts that exist in the system at the time of observation. 
This way each server host in a partition can store the same subset of graphs assigned
to that partition, providing a measurable amount of fault-tolerance to the user if 
one of those hosts happens to crash.

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

## Partitioning Algorithms Implemented

temp

## Consistentency Algorithms Implemented

We plan to implement the RAFT algorithm to ensure consistency among the data stored 
in our distributed system.

## Technologies Used

gRPC/protocol buffers, gin server code
badges(pictures) to come later
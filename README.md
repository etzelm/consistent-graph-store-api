# Scalable &amp; Highly Consistent(CAP Theorem) Graph Store API

## Introduction

The goal of this project is to provide a REST-accessible graph storage service that 
runs on port 3000 and is available as a resource named gs, for example the service 
listens at http://server-hostname:3000/gs. We want to develop distributed system 
software to support this service so that it can store the amount of data that would 
not normally fit onto a single-machine system.

## Input Format Specifications
- vertices
  - charset: [a-zA-Z0-9_] i.e. Alphanumeric including underscore, and case-sensitive 
  - size:    1 to 250 characters

- edges
  - charset: [a-zA-Z0-9_] i.e. Alphanumeric including underscore, and case-sensitive 
  - size:    1 to 250 characters
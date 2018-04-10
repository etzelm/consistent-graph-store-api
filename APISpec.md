# Scalable &amp; Highly Consistent(CAP Theorem) Graph Store API

## Graph Based Functionality

1. PUT localhost:3000/gs -d "graph=g1&vertices=[v1,v2]&edge=e1&causal_payload=<payload>"
    - case: 'e1' does not exist
      - status code : 201
      - response type : application/json
      - response body:
<pre>
{
      "msg": "success"
}
</pre>

1. PUT localhost:3000/gs -d "graph=g1&vertices=[v1,v2]&edge=e1"
    - case: 'e1' exists
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
      "msg": "already existed"
}
</pre>

2. GET localhost:3000/gs?graph=g1
    - case: 'g1' does not exist
      - status code : 404
      - response type : application/json
      - response body:

<pre>
{
      "msg" : "error",
      "error" : "key does not exist"
}
</pre>

2. GET localhost:3000/gs?graph=g1
    - case: 'g1' exists
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
      "msg" : "success",
      "vertices" : [v1,v2],
      "edges" : [e1]
}
</pre>

3. DELETE localhost:3000/gs?graph=g1
    - case: 'g1' does not exist
      - status code : 404
      - response type : application/json
      - response body:
      
<pre>
{
      "msg" : "error",
      "error" : "key does not exist"
}
</pre>

3. DELETE localhost:3000/gs?graph=g1
    - case: 'g1' exists
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
      "msg" : "success"
}
</pre>

## Server Based Functionality

1. PUT, GET, DELETE 
    - case: main instance is down 
      - status code : 404
      - response type : application/json
      - response body:

<pre>
{
      "msg" : "error",
      "error" : "service is not available"
}
</pre>

2. PUT localhost:8081/gs/view_update -d "ip_port=10.0.0.22:8080&type=add"
    - case: adding a server node:
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
     "msg": "success"
}
</pre>

2. PUT localhost:8081/gs/view_update -d "ip_port=10.0.0.20:8080&type=remove"
    - case: removing a server node:
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
     "msg": "success"
}
</pre>
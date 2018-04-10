# Scalable &amp; Highly Consistent(CAP Theorem) Graph Store API

## Graph Based Functionality

1. PUT localhost:3000/gs -d "graph=g1&vertices=[v1,v2]&edge=e1&vector_clock=6.2.9.1"
    - case: 'e1' does not exist
      - status code : 201
      - response type : application/json
      - response body:
<pre>
{
      "msg": "success",
      "part": 2,
      "vector": "6.2.9.1",
      "time": "1248425146"
}
</pre>

1. PUT localhost:3000/gs -d "graph=g1&vertices=[v1,v2]&edge=e1&vector_clock=6.2.9.1"
    - case: 'e1' exists
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
      "msg": "already existed",
      "part": 2,
      "vector": "6.2.9.1",
      "time": "1248425146"
}
</pre>

2. GET localhost:3000/gs?graph=g1&vector_clock=6.2.9.1
    - case: 'g1' does not exist
      - status code : 404
      - response type : application/json
      - response body:

<pre>
{
      "msg": "error",
      "error": "key does not exist",
      "part": 2,
      "vector": "6.2.9.1",
      "time": "1248425146"
}
</pre>

2. GET localhost:3000/gs?graph=g1&vector_clock=6.2.9.1
    - case: 'g1' exists
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
      "msg": "success",
      "vertices": [v1,v2,v3],
      "edges": [[v1,v2],[v1,v3]],
      "part": 2,
      "vector": "6.2.9.1",
      "time": "1248425146"
}
</pre>

3. DELETE localhost:3000/gs?graph=g1&vector_clock=6.2.9.1
    - case: 'g1' does not exist
      - status code : 404
      - response type : application/json
      - response body:
      
<pre>
{
      "msg": "error",
      "error": "key does not exist",
      "part": 2,
      "vector": "6.2.9.1",
      "time": "1248425146"
}
</pre>

3. DELETE localhost:3000/gs?graph=g1&vector_clock=6.2.9.1
    - case: 'g1' exists
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
      "msg": "success",
      "part": 2,
      "vector": "6.2.9.1",
      "time": "1248425146"
}
</pre>

## Server Based Functionality

1. PUT, GET, DELETE 
    - case: a queried instance is down 
      - status code : 404
      - response type : application/json
      - response body:

<pre>
{
      "msg": "error",
      "error": "service is not available"
}
</pre>

2. PUT localhost:8081/gs/view_update -d "ip_port=10.0.0.22:8080&type=add"
    - case: adding a server node:
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
     "msg": "success",
     "part": 2,
     "part_count": 3
}
</pre>

2. PUT localhost:8081/gs/view_update -d "ip_port=10.0.0.20:8080&type=remove"
    - case: removing a server node:
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
     "msg": "success",
     "part_count": 2
}
</pre>

3. GET localhost:8081/gs/partition 
    - case: returning partition the node belongs to
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
    "msg": "success",
    "part": 3,
}
</pre>

4. GET localhost:8081/gs/all_partitions
    - case: returning all partitions in the system
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
    "msg": "success",
    "part_list": [0,1,2,3]
}
</pre>

5. GET localhost:8081/gs/partition_members?partition=2
    - case: returning all nodes in the given partition
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
    "msg": "success",
    "part_memb": ["10.0.0.21:8080", "10.0.0.22:8080", "10.0.0.23:8080"]
}
</pre>

6. GET localhost:8081/gs/graph_count
    - case: returning all nodes in the given partition
      - status code : 200
      - response type : application/json
      - response body:

<pre>
{
    "msg": "success",
    "count": 6
}
</pre>

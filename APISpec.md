# Scalable &amp; Highly Consistent(CAP Theorem) Graph Store API

1. PUT localhost:3000/gs -d "graph=g1&vertices=[v1,v2]&edge=e1"
    - case 'e1' does not exist
      - status code : 201
      - response type : application/json
      - response body:
<pre>
{
      "replaced": 0, // 1 if an existing key's val was replaced
      "msg": "success"
}
</pre>

1. PUT localhost:3000/gs -d "graph=g1&vertices=[v1,v2]&edge=e1"
    - case 'e1' exists
	    - status code : 200
	    - response type : application/json
	    - response body:

<pre>
{
      "replaced": 1, // 0 if key did not exist
      "msg": "success"
}
</pre>

2. GET localhost:3000/gs?graph=g1
    - case 'g1' does not exist
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
    - case 'g1' exists
      - status code : 200
      - response type : application/json
      - response body:

<pre>
			{
				"msg" : "success"
		 	}
</pre>

3. DELETE localhost:3000/gs?graph=g1
    - case 'g1' does not exist
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
    - case 'g1' exists
      - status code : 200
      - response type : application/json
      - response body:

<pre>
			{
				"msg" : "success"
		 	}
</pre>

4. PUT, GET, DELETE 
    - case the main instance is down 
      - status code : 404
      - response type : application/json
      - response body:

<pre>
			{
				"msg" : "error",
				"error" : "service is not available"
		 	}
</pre>
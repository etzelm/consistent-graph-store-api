# Scalable &amp; Highly Consistent(CAP Theorem) Graph Store API

1. PUT localhost:3000/graph -d "vertices=[v1,v2]&edge=e1"
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

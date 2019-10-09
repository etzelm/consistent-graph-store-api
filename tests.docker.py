from multiprocessing import Pool
import random
import string
from collections import Counter
import json
import os
import subprocess
import requests as req
import time
import traceback
import sys

NODE_COUNTER = 2
PRINT_HTTP_REQUESTS = False
PRINT_HTTP_RESPONSES = False
AVAILABILITY_THRESHOLD = 1 
TB = 5

HEADER = '\033[95m'
OKBLUE = '\033[94m'
OKGREEN = '\033[92m'
FAIL = '\033[91m'
ENDC = '\033[0m'

class Node:
    
    def __init__(self, access_port, ip, node_id):
        self.access_port = access_port
        self.ip = ip
        self.id = node_id

    def __repr__(self):
        return self.ip

def generate_ip_port():
    global NODE_COUNTER
    NODE_COUNTER += 1
    ip = '10.0.0.' + str(NODE_COUNTER)
    port = str(8080 + NODE_COUNTER)
    return ip, port            

def start_gs(num_nodes, container_name, R=2, net='net', sudo='sudo'):
    ip_ports = []
    for i in range(1, num_nodes+1):
        ip, port = generate_ip_port()
        ip_ports.append((ip, port))
    servers = ','.join([ip+":8080" for ip, _ in ip_ports])
    nodes = []
    print "Starting server nodes"
    for ip, port in ip_ports:
        cmd_str = sudo + ' docker run -d -p ' + port + ":8080 --net=" + net + " -e R=" + str(R) + " --ip=" + ip + " -e SERVERS=\"" + servers + "\" -e IP=\"" + ip + "\" -e PORT=\"8080\" " + container_name
        print cmd_str
        node_id = subprocess.check_output(cmd_str, shell=True).rstrip('\n')
        nodes.append(Node(port, ip, node_id))
    time.sleep(5)
    return nodes

def start_new_node(container_name, R=2, net='net', sudo='sudo'):
    ip, port = generate_ip_port()
    cmd_str = sudo + ' docker run -d -p ' + port + ":8080 --net=" + net + " -e R=" + str(R) + " --ip=" + ip + " -e IP=\"" + ip + "\" -e PORT=\"8080\" " + container_name
    print cmd_str
    node_id = subprocess.check_output(cmd_str, shell=True).rstrip('\n')
    time.sleep(5)
    return Node(port, ip, node_id)    

def stop_all_nodes(sudo):                                           
    # running_containers = subprocess.check_output([sudo, 'docker',  'ps', '-q'])
    # if len(running_containers):
    print "Stopping all nodes"
    os.system(sudo + " docker kill $(" + sudo + " docker ps -q)") 

def stop_node(node, sudo='sudo'):
    cmd_str = sudo + " docker kill %s" % node.id
    print cmd_str
    os.system(cmd_str)
    time.sleep(0.5)

def find_node(nodes, ip_port):
    ip = ip_port.split(":")[0]
    for n in nodes:
        if n.ip == ip:
            return n
    return None

def disconnect_node(node, network, sudo):
    cmd_str = sudo + " docker network disconnect " + network + " " + node.id
    print cmd_str
    time.sleep(0.5)
    os.system(cmd_str)
    time.sleep(0.5)

def connect_node(node, network, sudo):
    cmd_str = sudo + " docker network connect " + network + " --ip=" + node.ip + ' ' + node.id
    print cmd_str
   # r = subprocess.check_output(cmd_str.split())
   # print r
    time.sleep(0.5)
    os.system(cmd_str)
    time.sleep(0.5)

def add_node_to_gs(hostname, cur_node, new_node):
    d = None
    put_str = "http://" + hostname + ":" + str(cur_node.access_port) + "/gs/change_view"
    data = {'ip_port':new_node.ip + ":8080", 'type':'add'}
    try:
        if PRINT_HTTP_REQUESTS:
            print "PUT request:" + put_str + " data field " + str(data)
        r = req.put(put_str, data=data)
        if PRINT_HTTP_RESPONSES:
            print "Response:", r.text, r.status_code
        d = r.json()
        if r.status_code not in [200, 201, '200', '201']:
            raise Exception("Error, status code %s is not 200 or 201" % r.status_code)
        for field in ['msg', 'partitionID', 'number_of_partitions']:
            if not d.has_key(field):
                raise Exception("Field \"" + field + "\" is not present in response " + str(d))
    except Exception as e:
        print "ERROR IN ADDING A NODE TO THE KEY-VALUE STORE:",
        print e
    return d


def delete_node_from_gs(hostname, cur_node, node_to_delete):
    d = None
    put_str = "http://" + hostname + ":" + str(cur_node.access_port) + "/gs/change_view"
    data = {'ip_port':node_to_delete.ip + ":8080", 'type':'remove'}
    try:
        if PRINT_HTTP_REQUESTS:
            print "PUT request: " + put_str + " data field " + str(data)
        r = req.put(put_str, data=data)
        if PRINT_HTTP_RESPONSES:
            print "Response:", r.text, r.status_code
        d = r.json()
        if r.status_code not in [200, 201, '200', '201']:
            raise Exception("Error, status code %s is not 200 or 201" % r.status_code)
        for field in ['msg', 'number_of_partitions']:
            if not d.has_key(field):
                raise Exception("Field \"" + field + "\" is not present in response " + str(d))
    except Exception as e:
        print "ERROR IN DELETING A NODE TO THE KEY-VALUE STORE:",
        print e
    return d

def get_all_partitions_ids(node):
    get_str = "http://" + hostname + ":" + str(node.access_port) + "/gs/all_partitions"
    try:
        if PRINT_HTTP_REQUESTS:
            print "Get request: " + get_str
        r = req.get(get_str)
        if PRINT_HTTP_RESPONSES:
            print "Response:", r.text, r.status_code
        d = r.json()
        for field in ['msg', 'partitionID_list']:
            if not d.has_key(field):
                raise Exception("Field \"" + field + "\" is not present in response " + str(d))
    except Exception as e:
        print "THE FOLLOWING GET REQUEST RESULTED IN AN ERROR: ",
        print get_str 
        print e
    return d['partitionID_list'] # returns the current partition ID list of the gs

def get_partitionID_for_node(node):
    get_str = "http://" + hostname + ":" + str(node.access_port) + "/gs/partition"
    try:
        if PRINT_HTTP_REQUESTS:
            print "Get request: " + get_str
        r = req.get(get_str)
        if PRINT_HTTP_RESPONSES:
            print "Response:", r.text, r.status_code
        d = r.json()
        for field in ['msg', 'partitionID']:
            if not d.has_key(field):
                raise Exception("Field \"" + field + "\" is not present in response " + str(d))
    except Exception as e:
        print "THE FOLLOWING GET REQUEST RESULTED IN AN ERROR: ",
        print get_str
        print e
    return d['partitionID']    

def get_partition_members(node, partitionID):
    get_str = "http://" + hostname + ":" + str(node.access_port) + "/gs/partition_members?partitionID=" + str(partitionID)
    d = None
    try:
        if PRINT_HTTP_REQUESTS:
            print "Get request: " + get_str
        r = req.get(get_str)
        if PRINT_HTTP_RESPONSES:
            print "Response:", r.text, r.status_code
        d = r.json()
        for field in ['msg', 'partition_members']:
            if not d.has_key(field):
                raise Exception("Field \"" + field + "\" is not present in response " + str(d))
    except Exception as e:
        print "THE FOLLOWING GET REQUEST RESULTED IN AN ERROR: ",
        print get_str
        print e
    return d['partition_members']    

if __name__ == "__main__":
    container_name = 'graphstore'
    hostname = '172.17.0.1'
    network = 'mynet'
    sudo = ''
    tests_to_run = [1,2] #  

    if 1 in tests_to_run:
        try: # Test 1
            test_description = "Test 1: Basic functionality for obtaining information about partitions; tests the following GET requests get_all_partitions_ids, get_partition_members and get_partitionID."
            print HEADER + "" + test_description  + ENDC
            nodes = start_gs(4, container_name, R=2, net=network, sudo=sudo)
            partitionID_list =  get_all_partitions_ids(nodes[0])
            if len(partitionID_list) != 2:
                raise Exception("ERROR: the number of partitions should be 2")
            
            print OKBLUE + "Obtaining partition members for partition " + str(0)  + ENDC
            members = get_partition_members(nodes[0], 0)
            if len(members) != 2:
                raise Exception("ERROR: the size of a partition %d should be 2, but it is %d" % (0, len(members)))
            
            part_nodes = []
            for ip_port in members:
                n = find_node(nodes, ip_port)
                if n is None:
                    raise Exception("ERROR: mismatch in the node ids (likely bug in the test script)")
                part_nodes.append(n)
            print OKBLUE + "Asking nodes directly about their partition id. Information should be consistent" + ENDC
            for i in range(len(part_nodes)):
                part_id = get_partitionID_for_node(part_nodes[i])
                if part_id != 0:
                    raise Exception("ERRR: inconsistent information about partition ids!")
            print OKBLUE + "Ok, killing all the nodes in the partition " + str(0) + ENDC
            for node in part_nodes:
                stop_node(node, sudo=sudo)
            other_nodes = [n for n in nodes if n not in part_nodes]

            print OKGREEN + "OK, functionality for obtaining information about partitions looks good!" + ENDC
        except Exception as e:
            print FAIL + "Exception in test 1" + ENDC
            print FAIL + str(e) + ENDC
            traceback.print_exc(file=sys.stdout)
        stop_all_nodes(sudo)

    if 2 in tests_to_run:
        try: # Test 2
            test_description = "Test2: Node additions/deletions. A gs consists of 2 partitions with 2 replicas each. I add 3 new nodes. The number of partitions should become 4. Then I delete a node.The number of partitions should become 3. I then delete 2 more nodes. Now the number of partitions should be back to 2." 
            print HEADER + "" + test_description  + ENDC
            print 
            print OKBLUE + "Starting gs ..." + ENDC
            nodes = start_gs(4, container_name, R=2, net=network, sudo=sudo)

            print OKBLUE + "Adding 3 nodes" + ENDC
            n1 = start_new_node(container_name, R=2, net=network, sudo=sudo)
            n2 = start_new_node(container_name, R=2, net=network, sudo=sudo)
            n3 = start_new_node(container_name, R=2, net=network, sudo=sudo)

            resp_dict = add_node_to_gs(hostname, nodes[0], n1)
            number_of_partitions = resp_dict.get('number_of_partitions')
            if number_of_partitions != 3:
                print FAIL + "ERROR: the number of partitions should be 3, but it is " + str(number_of_partitions) + ENDC
            else:
                print OKGREEN + "OK, the number of partitions is 3" + ENDC
            resp_dict = add_node_to_gs(hostname, nodes[2], n2)
            number_of_partitions = resp_dict.get('number_of_partitions')
            if number_of_partitions != 3:
                print FAIL + "ERROR: the number of partitions should be 3, but it is " + str(number_of_partitions) + ENDC
            else:
                print OKGREEN + "OK, the number of partitions is 3" + ENDC
            resp_dict = add_node_to_gs(hostname, n1, n3)
            number_of_partitions = resp_dict.get('number_of_partitions')
            if number_of_partitions != 4:
                print FAIL + "ERROR: the number of partitions should be 4, but it is " + str(number_of_partitions) + ENDC
            else:
                print OKGREEN + "OK, the number of partitions is 4" + ENDC

            print OKBLUE + "Deleting nodes ..." + ENDC
            resp_dict = delete_node_from_gs(hostname, n3, nodes[0])
            number_of_partitions = resp_dict.get('number_of_partitions')
            if number_of_partitions != 3:
                print FAIL + "ERROR: the number of partitions should be 3, but it is " + str(number_of_partitions) + ENDC
            else:
                print OKGREEN + "OK, the number of partitions is 3" + ENDC
            resp_dict = delete_node_from_gs(hostname, n3, nodes[2])
            number_of_partitions = resp_dict.get('number_of_partitions')
            if number_of_partitions != 3:
                print FAIL + "ERROR: the number of partitions should be 3, but it is " + str(number_of_partitions) + ENDC
            else:
                print OKGREEN + "OK, the number of partitions is 3" + ENDC
            resp_dict = delete_node_from_gs(hostname, n3, n2)
            number_of_partitions = resp_dict.get('number_of_partitions')
            if number_of_partitions != 2:
                print FAIL + "ERROR: the number of partitions should be 2, but it is " + str(number_of_partitions) + ENDC
            else:
                print OKGREEN + "OK, the number of partitions is 2" + ENDC
            print OKBLUE + "Stopping the gs" + ENDC
        except Exception as e:
            print FAIL + "Exception in test 2" + ENDC
            print FAIL + str(e) + ENDC
            traceback.print_exc(file=sys.stdout)
        stop_all_nodes(sudo)            

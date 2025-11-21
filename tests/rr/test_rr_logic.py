# script used to test for round robin load balancing 

import requests
import socket
import time


def test_round_robin_logic(tries = 40):
    # obtain the list of backends by querying the internal dns of docker
    backends = socket.gethostbyname_ex("backends")
    print(f"the list of backends i have got from the dns server is {backends}")
    n = len(backends)
    i = 0
    while i < tries:
        time.sleep(5)
        print("sending request to load balancer")
        apiURL = "http://load_balancer:8080/lb"
        resp = requests.get(apiURL)
        print(f"the response from the backend is {resp.content}")
        i+=1

tries = 50
test_round_robin_logic(tries)
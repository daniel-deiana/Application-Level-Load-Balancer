import pytest
import requests
import time

def test_register_backend():
    while True:
        time.sleep(4)
        response = requests.get("http://load_balancer:8080/lb")
        print(response.content)
        print(response.status_code)


test_register_backend()
import requests
import time
import pytest
import os


def getPort():
    variables = {}
    fileName = os.path.dirname(os.path.realpath(__file__)) + \
        "/.env.functional_test"
    with open(fileName) as envFile:
        for line in envFile:
            name, var = line.partition("=")[::2]
            variables[name.strip()] = var.strip()
        return variables["USER_DB_PORT"]


class HTTPConnector():
    def __init__(self):
        self.URL = "http://127.0.0.1:" + getPort()
        connected = False
        timeout = 15
        while timeout > 0:
            try:
                r = self.GET("/", "")
                if r.status_code == 200:
                    connected = True
                    break
            except Exception:
                timeout -= 1
                time.sleep(1)

        if connected is False:
            raise Exception("Cannot connect to test server")

    def GET(self, address, params):
        url = self.URL + address
        return requests.get(url=url, params=params)

    def POST(self, address, json):
        url = self.URL + address
        return requests.post(url=url, json=json)


@pytest.fixture
def httpConnection():
    return HTTPConnector()

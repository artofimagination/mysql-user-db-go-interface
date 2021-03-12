import pytest
import json


def addUser(data, httpConnection):
    if "user" in data:
        try:
            r = httpConnection.POST("/add-user", data["user"])
        except Exception:
            pytest.fail("Failed to send POST request")
            return None

        if r.status_code != 201:
            pytest.fail(f"Failed to run test.\nDetails: {r.text}")
            return None

        response = getResponse(r.text)
        if response is None:
            return None
        return response["ID"]
    else:
        return data["user_id"]


def addProduct(data, userUUID, httpConnection):
    if "product" in data:
        dataToSend = dict()
        dataToSend["product"] = data["product"]
        if userUUID is None:
            pytest.fail("Missing user test data")
            return None
        dataToSend["user"] = userUUID
        try:
            r = httpConnection.POST("/add-product", dataToSend)
        except Exception:
            pytest.fail("Failed to send POST request")
            return None

        if r.status_code != 201:
            pytest.fail(f"Failed to run test.\nDetails: {r.text}")
            return None

        response = getResponse(r.text)
        if response is None:
            return None
        return response["ID"]
    else:
        return data["product_id"]


def addProject(data, userUUID, productUUID, httpConnection):
    if "project" in data:
        dataToSend = dict()
        dataToSend["project"] = data["project"]
        if userUUID is None:
            pytest.fail("Missing user test data")
            return None
        dataToSend["owner_id"] = userUUID
        if productUUID is None:
            pytest.fail("Missing project test data")
            return None
        dataToSend["product_id"] = productUUID
        try:
            r = httpConnection.POST("/add-project", dataToSend)
        except Exception:
            pytest.fail("Failed to send POST request")
            return None

        if r.status_code != 201:
            pytest.fail(f"Failed to run test.\nDetails: {r.text}")
            return None

        response = getResponse(r.text)
        if response is None:
            return None
        return response["ID"]
    else:
        return data["id"]


def addProjects(data, userUUID, productUUID, httpConnection):
    uuidList = list()
    if "project" in data:
        for element in data["project"]:
            if "name" in element:
                dataToSend = dict()
                dataToSend["project"] = element
                if userUUID is None:
                    pytest.fail("Missing user test data")
                    return None
                dataToSend["owner_id"] = userUUID
                if productUUID is None:
                    pytest.fail("Missing product test data")
                    return None
                dataToSend["product_id"] = productUUID
                try:
                    r = httpConnection.POST("/add-project", dataToSend)
                except Exception:
                    pytest.fail("Failed to send POST request")
                    return None

                if r.status_code != 201:
                    pytest.fail(f"Failed to run test.\nDetails: {r.text}")
                    return None

                response = getResponse(r.text)
                if response is None:
                    return None
                uuidList.append(response["ID"])
            else:
                uuidList.append(element["id"])
    return uuidList


def getResponse(responseText):
    response = json.loads(responseText)
    if "error" in response:
        error = response["error"]
        pytest.fail(f"Failed to run test.\nDetails: {error}")
        return None
    return response["data"]

import pytest
import json
import base64


def addUser(data, httpConnection):
    if "user" in data:
        try:
            r = httpConnection.POST("/add-user", data["user"])
        except Exception:
            pytest.fail("Failed to send POST request")
            return None

        response = getResponse(r.text)
        if response is None:
            return None
        return response["id"]
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

        response = getResponse(r.text)
        if response is None:
            return None
        return response["id"]
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

        response = getResponse(r.text)
        if response is None:
            return None
        if r.status_code != 201:
            pytest.fail(f"Failed to run test.\nDetails: {response}")
            return None

        return response["id"]
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

                response = getResponse(r.text)
                if response is None:
                    return None
                uuidList.append(response["id"])
            else:
                uuidList.append(element["id"])
    return uuidList


def addProjectViewer(data, httpConnection):
    if "viewer_id" in data:
        userUUID = addUser(data, httpConnection)
        if userUUID is None:
            return None

        productUUID = addProduct(data, userUUID, httpConnection)
        if productUUID is None:
            return None

        projectUUID = addProject(
          data,
          userUUID,
          productUUID,
          httpConnection)
        if projectUUID is None:
            return None

        dataToSend = dict()
        dataToSend["project_id"] = projectUUID
        dataToSend["viewer_id"] = data["viewer_id"]
        dataToSend["user_id"] = userUUID
        dataToSend["is_owner"] = data["is_owner"]

        try:
            r = httpConnection.POST("/add-project-viewer", dataToSend)
        except Exception:
            pytest.fail("Failed to send POST request")
            return
        return r


# getResponse unwraps the data/error from json response.
# @expected shall be set to None only if
# the response result is just to generate a component for a test
# but not actually returning a test result.
def getResponse(responseText, expected=None):
    response = json.loads(responseText)
    if "error" in response and response["error"] != "":
        error = response["error"]
        if "error" not in expected or \
                ("error" in expected and error != expected["error"]):
            pytest.fail(f"Failed to run test.\nReturned: {error}\n")
        return None
    return response["data"]


def convertPasswdToBase64(passwdString):
    message_bytes = passwdString.encode('ascii')
    return base64.b64encode(message_bytes)

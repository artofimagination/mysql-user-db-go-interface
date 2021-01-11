import pytest
import json

def addUser(data, httpConnection):
  if "user" in data:
    try:
      r = httpConnection.POST("/add-user", data["user"])
    except Exception as e:
      pytest.fail(f"Failed to send POST request")
      return None

    if r.status_code != 201:
      pytest.fail(f"Failed to run test.\nDetails: {r.text}")
      return None

    response = json.loads(r.text)
    return response["ID"]
  else:
    return data["user_id"]

def addProduct(data, userUUID, httpConnection):
  if "product" in data:
    dataToSend = dict()
    dataToSend["product"] = data["product"]
    if userUUID is None:
      pytest.fail(f"Missing user test data")
      return None
    dataToSend["user"] = userUUID
    try:
      r = httpConnection.POST("/add-product", dataToSend)
    except Exception as e:
      pytest.fail(f"Failed to send POST request")
      return None

    if r.status_code != 201:
      pytest.fail(f"Failed to run test.\nDetails: {r.text}")
      return None

    response = json.loads(r.text)
    return response["ID"]
  else:
    return data["product_id"]

def addProject(data, userUUID, productUUID, httpConnection):
  if "project" in data:
    dataToSend = dict()
    dataToSend["project"] = data["project"]
    if userUUID is None:
      pytest.fail(f"Missing user test data")
      return None
    dataToSend["owner_id"] = userUUID
    if productUUID is None:
      pytest.fail(f"Missing project test data")
      return None
    dataToSend["product_id"] = productUUID
    try:
      r = httpConnection.POST("/add-project", dataToSend)
    except Exception as e:
      pytest.fail(f"Failed to send POST request")
      return None

    if r.status_code != 201:
      pytest.fail(f"Failed to run test.\nDetails: {r.text}")
      return None

    response = json.loads(r.text)
    return response["ID"]
  else:
    return data["id"]
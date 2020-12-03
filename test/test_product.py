import pytest
import json
from functionalTest import httpConnection

dataColumns = ("data", "expected")
createTestData = [
    (# Input data
      {
      "product": {
        "name": "testProduct",
        "public": "true"
      },
      "user": {
        "name": "testUserOwner",
        "email": "testEmailOwner",
        "password": "testPassword"
      }
    },
    # Expected
    { 
      "Name":"testProduct",
      "Public":"true",
    }),

    # Input data
    ({
      "product": {
        "name": "testProduct",
        "public": "true"
      },
      "user": {
        "name": "testUserOwner",
        "email": "testEmailOwner",
        "password": "testPassword"
      }
    },
    # Expected
    "Product with this name already exists")
]

ids=['No existing product', 'Existing product']

@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_CreateProduct(httpConnection, data, expected):
  dataToSend = dict()
  print(data)
  if "user" in data:
    try:
      r = httpConnection.POST("/add-user", data["user"])
    except:
      pytest.fail(f"Failed to send POST request")
      return

    if r.status_code != 201:
      pytest.fail(f"Failed to add product.\nDetails: {r.text}")
      return

    response = json.loads(r.text)
    
    dataToSend["product"] = data["product"]
    dataToSend["user"] = response["ID"]
  else:
    dataToSend = data


  print(dataToSend)
  try:
    r = httpConnection.POST("/add-product", dataToSend)
  except Exception as e:
    pytest.fail(f"Failed to send POST request")
    return

  if r.status_code == 201:
    response = json.loads(r.text)
    if response["Name"] != expected["Name"] or \
      response["Public"] != expected["Public"]:
      pytest.fail(f"Test failed\nReturned: {response}\nExpected: {expected}")
      return
  elif r.status_code == 202:
    if r.text != expected:
      pytest.fail(f"Request failed\nStatus code: {r.status_code}\nReturned: {r.text}\nExpected: {expected}")
    return
  else:
    pytest.fail(f"Request failed\nStatus code: {r.status_code}\nDetails: {r.text}")
    return

createTestData = [
    ({
      'name': 'testProductGet',
      'public': 'true',
    },
    { 
      'name': 'testProductGet',
      'public': 'true',
    }),
    ({
      "id": "c34a7368-344a-11eb-adc1-0242ac120002"
    },
    "The selected product not found")
]

ids=['Existing product', 'No existing product']

@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetProduct(httpConnection, data, expected):
  uuid = ""
  if "name" in data:
    try:
      r = httpConnection.POST("/add-product", data)
    except Exception as e:
      pytest.fail(f"Failed to send POST request")
      return

    if r.status_code != 201:
      pytest.fail(f"Failed to add product.\nDetails: {r.text}")
      return

    response = json.loads(r.text)
    uuid = response["ID"]
  else:
    uuid = data["id"]
  
  try:
    r = httpConnection.GET("/get-product", {"id": uuid})
  except Exception as e:
    pytest.fail(f"Failed to send GET request")
    return

  if r.status_code == 200:
    response = json.loads(r.text)
    try:
      if response["Name"] != expected["Name"] or \
        response["Public"] != expected["Public"]:
        pytest.fail(f"Test failed\nReturned: {response}\nExpected: {expected}")
        return
    except Exception as e:
      pytest.fail(f"Failed to compare results.\nDetails: {e}")
      return
  elif r.status_code == 202:
    if r.text != expected:
      pytest.fail(f"Request failed\nStatus code: {r.status_code}\nReturned: {r.text}\nExpected: {expected}")
  else:
    pytest.fail(f"Request failed\nStatus code: {r.status_code}\nDetails: {r.text}")
    return
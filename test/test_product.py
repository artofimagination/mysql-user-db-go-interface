import pytest
import json
from functionalTest import httpConnection

dataColumns = ("data", "expected")
createTestData = [
    (# Input data
      {
      "product": {
        "name": "testProduct",
        "public": True
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
      "Public": True,
    }),

    # Input data
    ({
      "product": {
        "name": "testProduct",
        "public": True
      },
      "user": {
        "name": "testUserOwner2",
        "email": "testEmailOwner2",
        "password": "testPassword"
      }
    },
    # Expected
    "Product with name testProduct already exists"),

    # Input data
    ({
      "product": {
        "name": "testProductMissingUser",
        "public": True
      },
      "user_id": "c34a7368-344a-11eb-adc1-0242ac120002"
    },
    # Expected
    "Failed to create product: Error 1452: Cannot add or update a child row: a foreign key constraint fails (`user_database`.`users_products`, CONSTRAINT `users_products_ibfk_2` FOREIGN KEY (`users_id`) REFERENCES `users` (`id`))") 
]

ids=['No existing product', 'Existing product', 'Missing user']

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
  elif "user_id" in data:
    dataToSend["product"] = data["product"]
    dataToSend["user"] = data["user_id"]
  else:
    dataToSend = data

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
    if r.text != expected:
      pytest.fail(f"Request failed\nStatus code: {r.status_code}\nReturned: {r.text}\nExpected: {expected}")
    return

createTestData = [
    # Input data
    ({
      "product": {
        "name": "testProductGet",
        "public": True
      },
      "user": {
        "name": "testUserOwnerGet",
        "email": "testEmailOwnerGet",
        "password": "testPassword"
      }
    },
    # Expected
    { 
      'Name': 'testProductGet',
      'Public': True,
      'base_asset_path': 'testPath'
    }),

    # Input data
    ({
      "id": "c34a7368-344a-11eb-adc1-0242ac120002"
    },
    # Expected
    "The selected product not found")
]

ids=['Existing product', 'No existing product']

@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetProduct(httpConnection, data, expected):
  userUUID = ""
  productUUID = ""

  if "user" in data:
    try:
      r = httpConnection.POST("/add-user", data["user"])
    except Exception as e:
      pytest.fail(f"Failed to send POST request")
      return

    if r.status_code != 201:
      pytest.fail(f"Failed to add product.\nDetails: {r.text}")
      return

    response = json.loads(r.text)
    userUUID = response["ID"]

  if "product" in data:
    dataToSend = dict()
    dataToSend["product"] = data["product"]
    if userUUID == "":
      pytest.fail(f"Missing user test data")
      return
    dataToSend["user"] = userUUID
    try:
      r = httpConnection.POST("/add-product", dataToSend)
    except Exception as e:
      pytest.fail(f"Failed to send POST request")
      return

    if r.status_code != 201:
      pytest.fail(f"Failed to add product.\nDetails: {r.text}")
      return

    response = json.loads(r.text)
    productUUID = response["ID"]
  else:
    productUUID = data["id"]
  
  try:
    r = httpConnection.GET("/get-product", {"id": productUUID})
  except Exception as e:
    pytest.fail(f"Failed to send GET request")
    return

  if r.status_code == 200:
    response = json.loads(r.text)
    try:
      if response["Name"] != expected["Name"] or \
        response["Public"] != expected["Public"] or \
        "base_asset_path" not in response["Assets"]["DataMap"] or \
        response["Assets"]["DataMap"]["base_asset_path"] != "testPath" or \
        "base_asset_path" not in response["Details"]["DataMap"] or \
        response["Details"]["DataMap"]["base_asset_path"] != "testPath":
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
import pytest
import json
from functionalTest import httpConnection

dataColumns = ("data", "expected")
createTestData = [
    ({
      'name': 'testUser',
      'email': 'testEmail',
      'password': 'testPassword'
    },
    { 
      "Name":"testUser",
      "Email":"testEmail",
      "Password":"dGVzdFBhc3N3b3Jk",
    }),

    ({
      'name': 'testUser',
      'email': 'testEmail',
      'password': 'testPassword'
    },
    "User with this email already exists"),
    
    ({
      'name': 'testUser',
      'email': 'testEmailNew',
      'password': 'testPassword'
    },
    "User with this name already exists")
]

ids=['No existing email', 'Existing email', 'Existing name']

@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_CreateUser(httpConnection, data, expected):
  try:
    r = httpConnection.POST("/add-user", data)
  except Exception as e:
    pytest.fail(f"Failed to send POST request")
    return

  if r.status_code == 201:
    response = json.loads(r.text)
    if response["Name"] != expected["Name"] or \
      response["Email"] != expected["Email"] or \
      response["Password"] != expected["Password"]:
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
      'name': 'testUserGet',
      'email': 'testEmailGet',
      'password': 'testPassword'
    },
    { 
      "Name":"testUserGet",
      "Email":"testEmailGet",
    }),

    ({
      "id": "c34a7368-344a-11eb-adc1-0242ac120002"
    },
    "The selected user not found")
]

ids=['Existing user', 'No existing user']

@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetUser(httpConnection, data, expected):
  uuid = ""
  if "name" in data:
    try:
      r = httpConnection.POST("/add-user", data)
    except Exception as e:
      pytest.fail(f"Failed to send POST request")
      return

    if r.status_code != 201:
      pytest.fail(f"Failed to add user.\nDetails: {r.text}")
      return

    response = json.loads(r.text)
    uuid = response["ID"]
  else:
    uuid = data["id"]
  
  try:
    r = httpConnection.GET("/get-user", {"id": uuid})
  except Exception as e:
    pytest.fail(f"Failed to send GET request")
    return

  if r.status_code == 200:
    response = json.loads(r.text)
    try:
      if response["Name"] != expected["Name"] or \
        response["Email"] != expected["Email"] or \
        response["Settings"]["ID"] == '' or \
        response["Settings"]["ID"] == '00000000-0000-0000-0000-000000000000' or \
        response["Assets"]["ID"] == '' or \
        response["Assets"]["ID"] == '00000000-0000-0000-0000-000000000000':
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

createTestData = [
    # Input data
    ({
      "product": {
        "name": "testProductUsers",
        "public": True
      },
      "user": {
        "name": "testProductUser",
        "email": "testEmailProductUser",
        "password": "testPassword"
      },
      "partner_user": [{ 
        "user" : {
          "name": "testUserPartner",
          "email": "testEmailPartner",
          "password": "testPassword"
        },
        "privilege" : 3
      }
      ]
    },
    # Expected
    "Add product user completed")
]

ids=['Add product users']

@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_AddProductUsers(httpConnection, data, expected):
  userUUID = ""
  productUUID = ""
  partnerUUIDs = list()

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

  if "partner_user" in data:
    for user in data["partner_user"]:
      try:
        r = httpConnection.POST("/add-user", user["user"])
      except Exception as e:
        pytest.fail(f"Failed to send POST request")
        return

      if r.status_code != 201:
        pytest.fail(f"Failed to add product.\nDetails: {r.text}")
        return

      response = json.loads(r.text)
      partnerUUID = dict()
      partnerUUID["id"] = response["ID"]
      partnerUUID["privilege"] = user["privilege"]
      partnerUUIDs.append(partnerUUID)
  

  dataToSend = dict()
  dataToSend["product_id"] = productUUID
  dataToSend["users"] = partnerUUIDs
  try:
    r = httpConnection.POST("/add-product-user", dataToSend)
  except Exception as e:
    pytest.fail(f"Failed to send POST request")
    return

  if r.status_code == 201 or r.status_code == 202:
    if r.text != expected:
      pytest.fail(f"Request failed\nStatus code: {r.status_code}\nReturned: {r.text}\nExpected: {expected}")
  else:
    pytest.fail(f"Request failed\nStatus code: {r.status_code}\nDetails: {r.text}")
    return

createTestData = [
    # Input data
    ({
      "product": {
        "name": "testProductUserDelete",
        "public": True
      },
      "user": {
        "name": "testProductUserDelete",
        "email": "testEmailProductUserDelete",
        "password": "testPassword"
      }
    },
    # Expected
    "Delete product user completed")
]

ids=['Delete product users']

@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_DeleteProductUser(httpConnection, data, expected):
  userUUID = ""
  productUUID = ""

  # Prepare test data
  # Create user
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

  # Create product
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
  
  # Run test
  dataToSend = dict()
  dataToSend["product_id"] = productUUID
  dataToSend["user_id"] = userUUID
  try:
    r = httpConnection.POST("/delete-product-user", dataToSend)
  except Exception as e:
    pytest.fail(f"Failed to send POST request")
    return

  if r.status_code == 200 or r.status_code == 202:
    if r.text != expected:
      pytest.fail(f"Request failed\nStatus code: {r.status_code}\nReturned: {r.text}\nExpected: {expected}")
  else:
    pytest.fail(f"Request failed\nStatus code: {r.status_code}\nDetails: {r.text}")
    return

createTestData = [
    # Input data
    ({
      "user_to_delete": {
        "name": "testUserDelete",
        "email": "testEmailDelete",
        'password': 'testPassword'
      },
      "products_to_delete": [{
        "name": "testProductDelete",
        "public": True,
      }],
      "nominated_users": [{
          "name": "testUserNominated1",
          "email": "testEmailNominated1",
          'password': 'testPassword'
        }
      ]
    },
    # Expected
    "Delete completed"),
    
    # Input data
    ({
      "id": "c34a7368-344a-11eb-adc1-0242ac120002"
    },
    # Expected
    "The selected user not found"),

    ({
      "user_to_delete": {
        "name": "testUserDeleteNoNominee",
        "email": "testEmailDeleteNoNominee",
        'password': 'testPassword'
      },
      "products_to_delete": [{
        "name": "testProductDeleteNoNominee",
        "public": True,
      }]
    },
    # Expected
    "Delete completed")
]

ids=['Existing user', 'Non existing user', 'No nominees']

@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_DeleteUser(httpConnection, data, expected):
  # Prepare data
  userUUID = ""
  # 1. Add the user to be deleted
  if "user_to_delete" in data:
    try:
      r = httpConnection.POST("/add-user", data["user_to_delete"])
    except Exception as e:
      pytest.fail(f"Failed to send POST request")
      return

    if r.status_code != 201:
      pytest.fail(f"Failed to add user.\nDetails: {r.text}")
      return

    response = json.loads(r.text)
    userUUID = response["ID"]
  else:
    userUUID = data["id"]

  productUUIDs = list() 

  # 2. Add the product to be deleted
  if "products_to_delete" in data:
    for product in data["products_to_delete"]:
      dataToSend = dict()
      dataToSend["product"] = product
      dataToSend["user"] = userUUID
      try:
        r = httpConnection.POST("/add-product", dataToSend)
      except Exception as e:
        pytest.fail(f"Failed to send POST request")
        return

      if r.status_code != 201:
        pytest.fail(f"Failed to add user.\nDetails: {r.text}")
        return

      response = json.loads(r.text)
      productUUIDs.append(response["ID"])

  # Add nominated users and their product relationship
  nominatedUUIDs = list()
  if "nominated_users" in data:
    for nominatedUser in data["nominated_users"]:
      try:
        r = httpConnection.POST("/add-user", nominatedUser)
      except Exception as e:
        pytest.fail(f"Failed to send POST request")
        return

      if r.status_code != 201:
        pytest.fail(f"Failed to add nominated user.\nDetails: {r.text}")
        return

      response = json.loads(r.text)
      nominatedUUIDs.append(response["ID"])

      productUsers = list()
      productUser = dict()
      productUser["id"] = response["ID"]
      productUser["privilege"] = 2
      productUsers.append(productUser)
      dataToSend = dict()
      dataToSend["product_id"] = productUUIDs[0]
      dataToSend["users"] = productUsers
      
      try:
        r = httpConnection.POST("/add-product-user", dataToSend)
      except Exception as e:
        pytest.fail(f"Failed to send POST request")
        return

      if r.status_code != 201:
        pytest.fail(f"Failed to add nominated product user.\nDetails: {r.text}")
        return

  if len(nominatedUUIDs) > len(productUUIDs):
    pytest.fail(f"Too many nominated users.")
    return

  nominees = dict()
  for index, nominee in enumerate(nominatedUUIDs):
    nominees[productUUIDs[index]] = nominee
  
  # Run test
  try:
    r = httpConnection.POST("/delete-user", {"id": userUUID, "nominees":nominees})
  except Exception as e:
    pytest.fail(f"Failed to send POST request")
    return

  if r.status_code == 200 or r.status_code == 202:
    if r.text != expected:
      pytest.fail(f"Request failed\nStatus code: {r.status_code}\nReturned: {r.text}\nExpected: {expected}")
  else:
    pytest.fail(f"Request failed\nStatus code: {r.status_code}\nDetails: {r.text}")
    return


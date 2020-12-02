import pytest
import json
from functionalTest import httpConnection
import uuid

dataColumns = ("data", "expected")
creatUserTestData = [
    ({
      'name': 'testUser',
      'email': 'testEmail',
      'password': 'testPassword'
    },
    { 
      "ID":"65ae421d-343f-11eb-be1c-0242ac120003",
      "Name":"testUser",
      "Email":"testEmail",
      "Password":"dGVzdFBhc3N3b3Jk",
      "SettingsID":"65ae4218-343f-11eb-be1c-0242ac120003",
      "AssetsID":"65ae4211-343f-11eb-be1c-0242ac120003"
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

@pytest.mark.parametrize(dataColumns, creatUserTestData, ids=ids)
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

creatUserTestData = [
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

ids=['Existing user', 'Non existing user']

@pytest.mark.parametrize(dataColumns, creatUserTestData, ids=ids)
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

# creatUserTestData = [
#     ({
#       "user_to_delete": {
#         "name": "testUserDelete",
#         "email": "testEmailDelete",
#         'password': 'testPassword'
#       },
#       "nominated_users": [{
#           "name": "testUserNominated1",
#           "email": "testEmailNominated1",
#           'password': 'testPassword'
#         },
#         {
#           "name": "testUserNominated2",
#           "email": "testEmailNominated2",
#           'password': 'testPassword'
#         }
#       ]
#     },
#     "Delete completed"),
#     ({
#       "id": "c34a7368-344a-11eb-adc1-0242ac120002"
#     },
#     "The selected user not found")
# ]

# ids=['Existing user', 'Non existing user']

# @pytest.mark.parametrize(dataColumns, creatUserTestData, ids=ids)
# def test_DeleteUser(httpConnection, data, expected):
#   userUUID = ""
#   if "user_to_delete" in data:
#     try:
#       r = httpConnection.POST("/add-user", data["user_to_delete"])
#     except Exception as e:
#       pytest.fail(f"Failed to send POST request")
#       return

#     if r.status_code != 201:
#       pytest.fail(f"Failed to add user.\nDetails: {r.text}")
#       return

#     response = json.loads(r.text)
#     userUUID = response["ID"]
#   else:
#     userUUID = data["id"]

#   nominatedUUIDs = list()
#   if "nominated_users" in data:
#     for nominatedUser in data["nominated_users"]:
#       try:
#         r = httpConnection.POST("/add-user", nominatedUser)
#       except Exception as e:
#         pytest.fail(f"Failed to send POST request")
#         return

#       if r.status_code != 201:
#         pytest.fail(f"Failed to add nominated user.\nDetails: {r.text}")
#         return

#       response = json.loads(r.text)
#       nominatedUUIDs.append(response["ID"])
  
#   print(nominatedUUIDs)
#   print({"id": userUUID, "nominees":nominatedUUIDs})
#   try:
#     r = httpConnection.GET("/delete-user", {"id": userUUID, "nominees":nominatedUUIDs})
#   except Exception as e:
#     pytest.fail(f"Failed to send GET request")
#     return

#   if r.status_code == 200 or r.status_code == 202:
#     if r.text != expected:
#       pytest.fail(f"Request failed\nStatus code: {r.status_code}\nReturned: {r.text}\nExpected: {expected}")
#   else:
#     pytest.fail(f"Request failed\nStatus code: {r.status_code}\nDetails: {r.text}")
#     return


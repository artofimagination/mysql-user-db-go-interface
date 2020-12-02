import pytest
import json
from functionalTest import httpConnection

dataColumns = ("data", "expected")
creatUserTestData = [
    ({
      'name': 'testUser000',
      'email': 'testEmail000',
      'password': 'testPassword'
    },
    { 
      "ID":"65ae421d-343f-11eb-be1c-0242ac120003",
      "Name":"testUser000",
      "Email":"testEmail000",
      "Password":"dGVzdFBhc3N3b3Jk",
      "SettingsID":"65ae4218-343f-11eb-be1c-0242ac120003",
      "AssetsID":"65ae4211-343f-11eb-be1c-0242ac120003"
    }),
    ({
      'name': 'testUser000',
      'email': 'testEmail000',
      'password': 'testPassword'
    },
    "User with this email already exists"),
    ({
      'name': 'testUser000',
      'email': 'testEmail001',
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
      response = json.dumps(r.text)
      if response != expected:
        pytest.fail(f"Test failed\nReturned: {response}\nExpected: {expected}")
        return
    elif r.status_code == 200:
      if r.text != expected:
        pytest.fail(f"Request failed\nStatus code: {r.status_code}\nReturned: {r.text}\nExpected: {expected}")
      return
    else:
      pytest.fail(f"Request failed\nStatus code: {r.status_code}\nDetails: {r.text}")
      return

creatUserTestData = [
    ({
      'id': '65ae421d-343f-11eb-be1c-0242ac120003',
    },
    { 
      "ID":"65ae421d-343f-11eb-be1c-0242ac120003",
      "Name":"testUser000",
      "Email":"testEmail000",
      "Password":"dGVzdFBhc3N3b3Jk",
      "SettingsID":"65ae4218-343f-11eb-be1c-0242ac120003",
      "AssetsID":"65ae4211-343f-11eb-be1c-0242ac120003"
    }),
    ({
      'id': 'c34a7368-344a-11eb-adc1-0242ac120002',
    },
    "User with this email already exists")
]

ids=['Existing user', 'Non existing user']

@pytest.mark.parametrize(dataColumns, creatUserTestData, ids=ids)
def test_GetUser(httpConnection, data, expected):
    try:
      r = httpConnection.GET("/get-user", data)
    except Exception as e:
      pytest.fail(f"Failed to send GET request")
      return

    if r.status_code == 200:
      response = json.dumps(r.text)
      if response != expected:
        pytest.fail(f"Test failed\nReturned: {response}\nExpected: {expected}")
        return
    else:
      pytest.fail(f"Request failed\nStatus code: {r.status_code}\nDetails: {r.text}")
      return

import pytest
import common

dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            'username': 'testUser',
            'email': 'testEmail',
            'password': 'testPassword'
        },
        # Expected
        {
            "username": "testUser",
            "email": "testEmail",
            "password": "dGVzdFBhc3N3b3Jk"
        }),

    (
        # Input data
        {
            'username': 'testUserEmailExists',
            'email': 'testEmail',
            'password': 'testPassword'
        },
        # Expected
        "User with this email already exists"),

    (
        # Input data
        {
            'username': 'testUser',
            'email': 'testEmailUserExists',
            'password': 'testPassword'
        },
        # Expected
        "User with this name already exists")
]

ids = ['No existing email', 'Existing email', 'Existing name']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_CreateUser(httpConnection, data, expected):
    try:
        r = httpConnection.POST("/add-user", data)
    except Exception:
        pytest.fail("Failed to send POST request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 201:
        if response["username"] != expected["username"] or \
          response["email"] != expected["email"]:
            pytest.fail(
                f"Test failed\nReturned: {response}\nExpected: {expected}")
        return

    if response != expected:
        pytest.fail(
            f"Request failed\nStatus code: \
            {r.status_code}\nReturned: {response}\nExpected: {expected}")


createTestData = [
    (
        # Input data
        {
            'username': 'testUserGet',
            'email': 'testEmailGet',
            'password': 'testPassword'
        },
        # Expected
        {
            "username": "testUserGet",
            "email": "testEmailGet"
        }),

    (
        # Input data
        {
          "id": "c34a7368-344a-11eb-adc1-0242ac120002"
        },
        # Expected
        "The selected user not found")
]

ids = ['Existing user', 'No existing user']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetUser(httpConnection, data, expected):
    uuid = ""
    if "username" in data:
        try:
            r = httpConnection.POST("/add-user", data)
        except Exception:
            pytest.fail("Failed to send POST request")
            return

        if r.status_code != 201:
            pytest.fail(f"Failed to add user.\nDetails: {r.text}")
            return

        response = common.getResponse(r.text, expected)
        if response is None:
            return None
        uuid = response["id"]
    else:
        uuid = data["id"]

    try:
        r = httpConnection.GET("/get-user", {"id": uuid})
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 200:
        try:
            zeroID = '00000000-0000-0000-0000-000000000000'
            if response["username"] != expected["username"] or \
                response["email"] != expected["email"] or \
                response["settings"]["id"] == '' or \
                response["settings"]["id"] == zeroID or \
                response["assets"]["id"] == '' or \
                    response["assets"]["id"] == zeroID:
                pytest.fail(
                  f"Test failed\nReturned: {response}\nExpected: {expected}")
            return
        except Exception as e:
            pytest.fail(f"Failed to compare results.\nDetails: {e}")
            return
    if response != expected:
        pytest.fail(
            f"Request failed\nStatus code: \
            {r.status_code}\nReturned: {response}\nExpected: {expected}")


createTestData = [
    (
        # Input data
        {
            'username': 'testUserGetByEmail',
            'email': 'testEmailGetByEmail',
            'password': 'testPassword'
        },
        # Expected
        {
            "username": "testUserGetByEmail",
            "email": "testEmailGetByEmail",
        }),

    (
      # Input data
      {
          "email": "testEmailGetWrong"
      },
      # Expected
      "The selected user not found")
]

ids = ['Existing user', 'No existing user']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetUserByEmail(httpConnection, data, expected):
    email = ""
    if "username" in data:
        try:
            r = httpConnection.POST("/add-user", data)
        except Exception:
            pytest.fail("Failed to send POST request")
            return

        response = common.getResponse(r.text, expected)
        if response is None:
            return None
        if r.status_code != 201:
            pytest.fail(f"Failed to add user.\nDetails: {response}")
            return

        email = response["email"]
    else:
        email = data["email"]

    try:
        r = httpConnection.GET("/get-user-by-email", {"email": email})
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 200:
        try:
            zeroID = '00000000-0000-0000-0000-000000000000'
            if response["username"] != expected["username"] or \
                response["email"] != expected["email"] or \
                response["settings"]["id"] == '' or \
                response["settings"]["id"] == zeroID or \
                response["assets"]["id"] == '' or \
                    response["assets"]["id"] == zeroID:
                pytest.fail(
                    f"Test failed\nReturned: {response}\nExpected: {expected}")
            return
        except Exception as e:
            pytest.fail(f"Failed to compare results.\nDetails: {e}")
            return

    if response != expected:
        pytest.fail(
            f"Request failed\nStatus code: \
            {r.status_code}\nReturned: {response}\nExpected: {expected}")


createTestData = [
    (
        # Input data
        [{
            'username': 'testUserGetMultiple1',
            'email': 'testEmailGetMultiple1',
            'password': 'testPassword'
        }, {
            'username': 'testUserGetMultiple2',
            'email': 'testEmailGetMultiple2',
            'password': 'testPassword'
        }],
        # Expected
        [{
            "username": "testUserGetMultiple1",
            "email": "testEmailGetMultiple1"
        }, {
            "username": "testUserGetMultiple2",
            "email": "testEmailGetMultiple2"
        }]
    ),

    (
        # Input data
        [{
            'username': 'testUserGetMultipleFail',
            'email': 'testEmailGetMultipleFail',
            'password': 'testPassword'
        }, {
            "id": "c34a7368-344a-11eb-adc1-0242ac120002"
        }],
        # Expected
        [{
            "username": "testUserGetMultipleFail",
            "email": "testEmailGetMultipleFail"
        }]),

    (
        [{
            "id": "c34a7368-344a-11eb-adc1-0242ac120002"
        }],
        "The selected user not found")
]

ids = ['Existing users', 'Missing a user', 'No user']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetUsers(httpConnection, data, expected):
    uuidList = list()
    for element in data:
        if "username" in element:
            try:
                r = httpConnection.POST("/add-user", element)
            except Exception:
                pytest.fail("Failed to send POST request")
                return

            response = common.getResponse(r.text, expected)
            if response is None:
                return None
            if r.status_code != 201:
                pytest.fail(f"Failed to add user.\nDetails: {response}")
                return

            uuidList.append(response["id"])
        else:
            uuidList.append(element["id"])

    try:
        r = httpConnection.GET("/get-users", {"ids": uuidList})
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 200:
        try:
            for index, user in enumerate(response):
                asset = user["assets"]["datamap"]
                settings = user["settings"]["datamap"]
                if user["username"] != expected[index]["username"] or \
                    user["email"] != expected[index]["email"] or \
                    "base_asset_path" not in asset or \
                    asset["base_asset_path"] != "testPath" or \
                    "base_asset_path" not in settings or \
                        settings["base_asset_path"] != "testPath":
                    pytest.fail(
                        f"Test failed\nReturned: \
                        {response}\nExpected: {expected}")
            return
        except Exception as e:
            pytest.fail(f"Failed to compare results.\nDetails: {e}")
            return

    if response != expected:
        pytest.fail(
            f"Request failed\nStatus code: \
            {r.status_code}\nReturned: {response}\nExpected: {expected}")


createTestData = [
    (
        # Input data
        {
            "user": {
                'username': 'testUserGetPassword',
                'email': 'testEmailGetPassword',
                'password': 'testPassword'
            },
            "login": {
              "email": "testEmailGetPassword",
              "password": "testPassword",
            }
        },
        # Expected
        'OK'),

    (
        # Input data
        {
            "user": {
                'username': 'testUserGetPasswordInvalid',
                'email': 'testEmailGetPasswordInvalid',
                'password': 'testPassword'
            },
            "login": {
                "email": "testEmailGetPasswordInvalid",
                "password": "testPasswordWrong"
            }
        },
        # Expected
        'Invalid password'),

    (
        # Input data
        {
          "id": "c34a7368-344a-11eb-adc1-0242ac120002"
        },
        # Expected
        "The selected user not found")
]

ids = ['Valid password', 'Invalid Password', 'No user found']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_Authenticate(httpConnection, data, expected):
    uuid = ""
    email = "empty"
    password = "empty"
    if "user" in data:
        try:
            r = httpConnection.POST("/add-user", data["user"])
        except Exception:
            pytest.fail("Failed to send POST request")
            return

        response = common.getResponse(r.text, expected)
        if response is None:
            return None
        if r.status_code != 201:
            pytest.fail(f"Failed to add user.\nDetails: {response}")
            return

        uuid = response["id"]
        email = data["login"]["email"]
        password = data["login"]["password"]
    else:
        uuid = data["id"]

    try:
        r = httpConnection.GET(
            "/authenticate",
            {"id": uuid, "email": email, "password": password})
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if response != expected:
        pytest.fail(
            f"Request failed\nStatus code: \
            {r.status_code}\nReturned: {response}\nExpected: {expected}")


createTestData = [
    (
      # Input data
      {
          "product": {
              "name": "testProductUsers",
              "public": True
          },
          "user": {
              "username": "testProductUser",
              "email": "testEmailProductUser",
              "password": "testPassword"
          },
          "partner_user": [{
              "user": {
                "username": "testUserPartner",
                "email": "testEmailPartner",
                "password": "testPassword"
              },
              "privilege": 3
          }]
      },
      # Expected
      "OK")
]

ids = ['Add product users']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_AddProductUsers(httpConnection, data, expected):
    partnerUUIDs = list()

    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    if "partner_user" in data:
        for user in data["partner_user"]:
            try:
                r = httpConnection.POST("/add-user", user["user"])
            except Exception:
                pytest.fail("Failed to send POST request")
                return

            response = common.getResponse(r.text, expected)
            if response is None:
                return None
            if r.status_code != 201:
                pytest.fail(f"Failed to add product.\nDetails: {response}")
                return

            partnerUUID = dict()
            partnerUUID["id"] = response["id"]
            partnerUUID["privilege"] = user["privilege"]
            partnerUUIDs.append(partnerUUID)

    dataToSend = dict()
    dataToSend["product_id"] = productUUID
    dataToSend["users"] = partnerUUIDs
    try:
        r = httpConnection.POST("/add-product-user", dataToSend)
    except Exception:
        pytest.fail("Failed to send POST request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if response != expected:
        pytest.fail(
            f"Request failed\nStatus code: \
            {r.status_code}\nReturned: {response}\nExpected: {expected}")


createTestData = [
    (
      # Input data
      {
          "product": {
              "name": "testProductUserDelete"
          },
          "user": {
              "username": "testProductUserDelete",
              "email": "testEmailProductUserDelete",
              "password": "testPassword"
          }
      },
      # Expected
      "OK")
]

ids = ['Delete product users']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_DeleteProductUser(httpConnection, data, expected):
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    dataToSend = dict()
    dataToSend["product_id"] = productUUID
    dataToSend["user_id"] = userUUID
    try:
        r = httpConnection.POST("/delete-product-user", dataToSend)
    except Exception:
        pytest.fail("Failed to send POST request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if response != expected:
        pytest.fail(
            f"Request failed\nStatus code: \
            {r.status_code}\nReturned: {response}\nExpected: {expected}")


createTestData = [
    (
        # Input data
        {
            "user": {
                "username": "testUserDelete",
                "email": "testEmailDelete",
                'password': 'testPassword'
            },
            "products_to_delete": [{
                "name": "testProductDelete"
            }],
            "nominated_users": [{
                "username": "testUserNominated1",
                "email": "testEmailNominated1",
                'password': 'testPassword'
            }]
        },
        # Expected
        "OK"),

    (
        # Input data
        {
            "user_id": "c34a7368-344a-11eb-adc1-0242ac120002"
        },
        # Expected
        "The selected user not found"),

    (
        # Input data
        {
            "user": {
                "username": "testUserDeleteNoNominee",
                "email": "testEmailDeleteNoNominee",
                'password': 'testPassword'
            },
            "products_to_delete": [{
                "name": "testProductDeleteNoNominee"
            }]
        },
        # Expected
        "OK")
]

ids = ['Existing user', 'Non existing user', 'No nominees']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_DeleteUser(httpConnection, data, expected):
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUIDs = list()
    if "products_to_delete" in data:
        for product in data["products_to_delete"]:
            dataToSend = dict()
            dataToSend["product"] = product
            dataToSend["user"] = userUUID
            try:
                r = httpConnection.POST("/add-product", dataToSend)
            except Exception:
                pytest.fail("Failed to send POST request")
                return

            response = common.getResponse(r.text, expected)
            if response is None:
                return None
            if r.status_code != 201:
                pytest.fail(f"Failed to add user.\nDetails: {response}")
                return
            productUUIDs.append(response["id"])

    # Add nominated users and their product relationship
    nominatedUUIDs = list()
    if "nominated_users" in data:
        for nominatedUser in data["nominated_users"]:
            try:
                r = httpConnection.POST("/add-user", nominatedUser)
            except Exception:
                pytest.fail("Failed to send POST request")
                return

            response = common.getResponse(r.text, expected)
            if response is None:
                return None
            if r.status_code != 201:
                pytest.fail(
                    f"Failed to add nominated user.\nDetails: {r.text}")
                return

            nominatedUUIDs.append(response["id"])

            productUsers = list()
            productUser = dict()
            productUser["id"] = response["id"]
            productUser["privilege"] = 2
            productUsers.append(productUser)
            dataToSend = dict()
            dataToSend["product_id"] = productUUIDs[0]
            dataToSend["users"] = productUsers

            try:
                r = httpConnection.POST("/add-product-user", dataToSend)
            except Exception:
                pytest.fail("Failed to send POST request")
                return

            response = common.getResponse(r.text, expected)
            if response is None:
                return None
            if r.status_code != 201:
                pytest.fail(
                    f"Failed to add nominated user.\nDetails: {response}")
                return

    if len(nominatedUUIDs) > len(productUUIDs):
        pytest.fail("Too many nominated users.")
        return

    nominees = dict()
    for index, nominee in enumerate(nominatedUUIDs):
        nominees[productUUIDs[index]] = nominee

    try:
        r = httpConnection.POST(
            "/delete-user",
            {"id": userUUID, "nominees": nominees})
    except Exception:
        pytest.fail("Failed to send POST request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if response != expected:
        pytest.fail(
            f"Request failed\nStatus code: \
            {r.status_code}\nReturned: {response}\nExpected: {expected}")

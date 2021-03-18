import pytest
import common

dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "product": {
                "name": "testProduct"
            },
            "user": {
                "name": "testUserOwner",
                "email": "testEmailOwner",
                "password": "testPassword"
            }
        },
        # Expected
        {
            "Name": "testProduct"
        }),

    (
        # Input data
        {
            "product": {
              "name": "testProduct"
            },
            "user": {
              "name": "testUserOwner2",
              "email": "testEmailOwner2",
              "password": "testPassword"
            }
        },
        # Expected
        "Product with name testProduct already exists"),

    (
        # Input data
        {
            "product": {
              "name": "testProductMissingUser"
            },
            "user_id": "c34a7368-344a-11eb-adc1-0242ac120002"
        },
        # Expected
        "The selected user not found")
]

ids = ['No existing product', 'Existing product', 'Missing user']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_CreateProduct(httpConnection, data, expected):
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    dataToSend = dict()
    dataToSend["product"] = data["product"]
    dataToSend["user"] = userUUID

    try:
        r = httpConnection.POST("/add-product", dataToSend)
    except Exception:
        pytest.fail("Failed to send POST request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 201:
        if response["Name"] != expected["Name"]:
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
            "product": {
              "name": "testProductGet"
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
            'base_asset_path': 'testPath'
        }),

    (
      # Input data
      {
          "user": {
            "name": "testUserOwnerGet1",
            "email": "testEmailOwnerGet1",
            "password": "testPassword"
          },
          "product_id": "c34a7368-344a-11eb-adc1-0242ac120002"
      },
      # Expected
      "The selected product not found")
]

ids = ['Existing product', 'No existing product']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetProduct(httpConnection, data, expected):
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    try:
        r = httpConnection.GET("/get-product", {"id": productUUID})
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 200:
        try:
            asset = response["Assets"]["DataMap"]
            details = response["Details"]["DataMap"]
            if response["Name"] != expected["Name"] or \
                "base_asset_path" not in asset or \
                asset["base_asset_path"] != "testPath" or \
                "base_asset_path" not in details or \
                    details["base_asset_path"] != "testPath":
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
            "product": [{
                    "name": "testProductGetMultiple1"
                },
                {
                    "name": "testProductGetMultiple2"
                }],
            "user": {
                "name": "testUserOwnerGetMultiple",
                "email": "testEmailOwnerGetMultiple",
                "password": "testPassword"
            }
        },
        # Expected
        [
            {
                'Name': 'testProductGetMultiple1',
                'base_asset_path': 'testPath'
            },
            {
                'Name': 'testProductGetMultiple2',
                'base_asset_path': 'testPath'
            }]),

        # Input data
        ({
            "product": [
                {
                    "name": "testProductGetMultipleFail"
                },
                {
                    "product_id": "c34a7368-344a-11eb-adc1-0242ac120002"
                }],
            "user": {
                "name": "testUserOwnerGetMultipleFail",
                "email": "testEmailOwnerGetMultipleFail",
                "password": "testPassword"
            }
        },
        # Expected
        [{
            'Name': 'testProductGetMultipleFail',
            'base_asset_path': 'testPath'
        }]),
    (
        # Input data
        {
            "product": [{
                  "product_id": "c34a7368-344a-11eb-adc1-0242ac120002"
            }],
            "user": {
                "name": "testUserOwnerGetMultipleNoProduct",
                "email": "testEmailOwnerGetMultipleNoProduct",
                "password": "testPassword"
            }
        },
        # Expected
        "The selected product not found")
]

ids = ['Existing products', 'Missing a product', 'No product']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetProducts(httpConnection, data, expected):
    uuidList = list()
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    if "product" in data:
        for element in data["product"]:
            if "name" in element:
                dataToSend = dict()
                dataToSend["product"] = element
                if userUUID == "":
                    pytest.fail("Missing user test data")
                    return
                dataToSend["user"] = userUUID
                try:
                    r = httpConnection.POST("/add-product", dataToSend)
                except Exception:
                    pytest.fail("Failed to send POST request")
                    return

                if r.status_code != 201:
                    pytest.fail(f"Failed to add product.\nDetails: {r.text}")
                    return

                response = common.getResponse(r.text, expected)
                if response is None:
                    return None
                uuidList.append(response["ID"])
            else:
                uuidList.append(element["product_id"])

    try:
        r = httpConnection.GET("/get-products", {"ids": uuidList})
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 200:
        try:
            for index, product in enumerate(response):
                asset = product["Assets"]["DataMap"]
                details = product["Details"]["DataMap"]
                if product["Name"] != expected[index]["Name"] or \
                    "base_asset_path" not in asset or \
                    asset["base_asset_path"] != "testPath" or \
                    "base_asset_path" not in details or \
                        details["base_asset_path"] != "testPath":
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
            "product": {
                "name": "testProductDeleteProduct"
            },
            "user": {
                "name": "testUserOwnerDeleteProduct",
                "email": "testEmailOwnerDeleteProduct",
                "password": "testPassword"
            }
        },
        # Expected
        "OK"
        ),
    (
        # Input data
        {
          "user": {
            "name": "testUserOwnerDeleteProduct1",
            "email": "testEmailOwnerDeleteProduct1",
            "password": "testPassword"
          },
          "product_id": "c34a7368-344a-11eb-adc1-0242ac120002"
        },
        # Expected
        "The selected product not found")]

ids = ['Existing product', 'No existing product']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_DeleteProduct(httpConnection, data, expected):
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    dataToSend = dict()
    dataToSend["product_id"] = productUUID

    try:
        r = httpConnection.POST("/delete-product", dataToSend)
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

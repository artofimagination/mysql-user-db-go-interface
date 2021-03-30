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
                "username": "testUserOwner",
                "email": "testEmailOwner",
                "password": "testPassword"
            }
        },
        # Expected
        {
            "name": "testProduct"
        }),

    (
        # Input data
        {
            "product": {
              "name": "testProduct"
            },
            "user": {
              "username": "testUserOwner2",
              "email": "testEmailOwner2",
              "password": "testPassword"
            }
        },
        # Expected
        {
            "error": "Product with name testProduct already exists"
        }),

    (
        # Input data
        {
            "product": {
              "name": "testProductMissingUser"
            },
            "user_id": "c34a7368-344a-11eb-adc1-0242ac120002"
        },
        # Expected
        {
            "error": "The selected user not found"
        })
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
        if response["name"] != expected["name"]:
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
              "username": "testUserOwnerGet",
              "email": "testEmailOwnerGet",
              "password": "testPassword"
            }
        },
        # Expected
        {
            'name': 'testProductGet',
            'assets': {
                'datamap': {}
            },
            'details': {
                'datamap': {}
            }
        }),

    (
      # Input data
      {
          "user": {
            "username": "testUserOwnerGet1",
            "email": "testEmailOwnerGet1",
            "password": "testPassword"
          },
          "product_id": "c34a7368-344a-11eb-adc1-0242ac120002"
      },
      # Expected
      {
          "error": "The selected product not found"
      })
]

ids = ['Existing product', 'No existing product']


def handleProductIDCheck(response, expected):
    zeroID = '00000000-0000-0000-0000-000000000000'
    if "id" in response and response["id"] != '' and response["id"] != zeroID:
        response.pop('id', None)
    else:
        pytest.fail(
            f"Test failed\nReturned: {response}\nExpected: {expected}")
        return None

    if "id" in response["assets"] and \
            response["assets"]["id"] != '' and \
            response["assets"]["id"] != zeroID:
        response["assets"].pop('id', None)
    else:
        pytest.fail(
            f"Test failed\nReturned: {response}\nExpected: {expected}")
        return None

    if "id" in response["details"] and \
            response["details"]["id"] != '' and \
            response["details"]["id"] != zeroID:
        response["details"].pop('id', None)
    else:
        pytest.fail(
            f"Test failed\nReturned: {response}\nExpected: {expected}")
        return None
    return True


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetProduct(httpConnection, data, expected):
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    try:
        r = httpConnection.GET("/get-product-by-id", {"id": productUUID})
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 200:
        try:
            status = handleProductIDCheck(response, expected)
            if status is None:
                return

            if response != expected:
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
                "username": "testUserOwnerGetMultiple",
                "email": "testEmailOwnerGetMultiple",
                "password": "testPassword"
            }
        },
        # Expected
        [{
            'name': 'testProductGetMultiple1',
            'assets': {
                'datamap': {}
            },
            'details': {
                'datamap': {}
            }
        }, {
            'name': 'testProductGetMultiple2',
            'assets': {
                'datamap': {}
            },
            'details': {
                'datamap': {}
            }
        }]),
    (
        # Input data
        {
            "product": [
                {
                    "name": "testProductGetMultipleFail"
                },
                {
                    "product_id": "c34a7368-344a-11eb-adc1-0242ac120002"
                }],
            "user": {
                "username": "testUserOwnerGetMultipleFail",
                "email": "testEmailOwnerGetMultipleFail",
                "password": "testPassword"
            }
        },
        # Expected
        [{
            'name': 'testProductGetMultipleFail',
            'assets': {
                'datamap': {}
            },
            'details': {
                'datamap': {}
            }
        }]),
    (
        # Input data
        {
            "product": [{
                  "product_id": "c34a7368-344a-11eb-adc1-0242ac120002"
            }],
            "user": {
                "username": "testUserOwnerGetMultipleNoProduct",
                "email": "testEmailOwnerGetMultipleNoProduct",
                "password": "testPassword"
            }
        },
        # Expected
        {
            "error": "The selected product not found"
        })
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
                uuidList.append(response["id"])
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
                status = handleProductIDCheck(product, expected[index])
                if status is None:
                    return
                if response[index] != expected[index]:
                    pytest.fail(
                        f"Test failed\n\
Returned: {response[index]}\nExpected: {expected[index]}")
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
                "username": "testUserOwnerDeleteProduct",
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
            "username": "testUserOwnerDeleteProduct1",
            "email": "testEmailOwnerDeleteProduct1",
            "password": "testPassword"
          },
          "product_id": "c34a7368-344a-11eb-adc1-0242ac120002"
        },
        # Expected
        {
            "error": "The selected product not found"
        })]

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

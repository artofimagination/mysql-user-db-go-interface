import pytest
import common

dataColumns = ("data", "expected")
createTestData = [
    (
      # Input data
      {
          "product": {
              "name": "testProductUpdateDetails"
          },
          "user": {
              "name": "testUserOwnerUpdateDetails",
              "email": "testEmailOwnerUpdateDetails",
              "password": "testPassword"
          },
          "details_entry": {
              "test_entry": "test_data"
          }
      },
      # Expected
      "OK")
]

ids = ['Valid product detail']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_UpdateProductDetail(httpConnection, data, expected):
    dataToSend = dict()
    if "user" in data:
        try:
            r = httpConnection.POST("/add-user", data["user"])
        except Exception:
            pytest.fail("Failed to send POST request")
            return

        if r.status_code != 201:
            pytest.fail(f"Failed to add user.\nDetails: {r.text}")
            return

        response = common.getResponse(r.text, expected)
        if response is None:
            return None
        dataToSend["product"] = data["product"]
        dataToSend["user"] = response["ID"]

    if "product" in data:
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
        dataToSend = dict()
        dataToSend["product"] = response
        for k, v in data["details_entry"].items():
            dataToSend["product"]["Details"]["DataMap"][k] = v

    try:
        r = httpConnection.POST("/update-product-details", dataToSend)
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
                "name": "testProductUpdateAssets"
            },
            "user": {
                "name": "testUserOwnerUpdateAssets",
                "email": "testEmailOwnerUpdateAssets",
                "password": "testPassword"
            },
            "details_entry": {
                "test_entry": "test_data"
            }
        },
        # Expected
        "OK")
]

ids = ['Valid product asset']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_UpdateProductAsset(httpConnection, data, expected):
    dataToSend = dict()
    if "user" in data:
        try:
            r = httpConnection.POST("/add-user", data["user"])
        except Exception:
            pytest.fail("Failed to send POST request")
            return

        if r.status_code != 201:
            pytest.fail("Failed to add user.\nDetails: {r.text}")
            return

        response = common.getResponse(r.text, expected)
        if response is None:
            return None
        dataToSend["product"] = data["product"]
        dataToSend["user"] = response["ID"]

    if "product" in data:
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
        dataToSend = dict()
        dataToSend["product"] = response
        for k, v in data["details_entry"].items():
            dataToSend["product"]["Assets"]["DataMap"][k] = v

    try:
        r = httpConnection.POST("/update-product-assets", dataToSend)
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
                "name": "testProjectUpdateDetails"
            },
            "user": {
                "name": "testUserProjectUpdateDetails",
                "email": "testEmailProjectUpdateDetails",
                "password": "testPassword"
            },
            "project": {
                "name": "testProjectUpdateProjectDetails",
                "visibility": "Public"
            },
            "details_entry": {
              "test_entry": "test_data"
            }
        },
        # Expected
        "OK")
]

ids = ['Valid project detail']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_UpdateProjectDetail(httpConnection, data, expected):
    dataToSend = dict()
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    dataToSend["project"] = data["project"]
    dataToSend["product_id"] = productUUID
    dataToSend["owner_id"] = userUUID
    if "project" in data:
        try:
            r = httpConnection.POST("/add-project", dataToSend)
        except Exception:
            pytest.fail("Failed to send POST request")
            return

        if r.status_code != 201:
            pytest.fail(f"Failed to add project.\nDetails: {r.text}")
            return
        response = common.getResponse(r.text, expected)
        if response is None:
            return None

        dataToSend = dict()
        dataToSend["project"] = response
        for k, v in data["details_entry"].items():
            dataToSend["project"]["Details"]["DataMap"][k] = v

    try:
        r = httpConnection.POST("/update-project-details", dataToSend)
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
                "name": "testProjectUpdateAssets"
            },
            "user": {
                "name": "testUserProjectUpdateAssets",
                "email": "testEmailProjectUpdateAssets",
                "password": "testPassword"
            },
            "project": {
                "name": "testProjectUpdateProjectAssets",
                "visibility": "Public"
            },
            "details_entry": {
                "test_entry": "test_data"
            }
        },
        # Expected
        "OK")
]

ids = ['Valid project asset']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_UpdateProjectAsset(httpConnection, data, expected):
    dataToSend = dict()
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    dataToSend["project"] = data["project"]
    dataToSend["product_id"] = productUUID
    dataToSend["owner_id"] = userUUID
    if "project" in data:
        try:
            r = httpConnection.POST("/add-project", dataToSend)
        except Exception:
            pytest.fail("Failed to send POST request")
            return

        if r.status_code != 201:
            pytest.fail(f"Failed to add project.\nDetails: {r.text}")
            return
        response = common.getResponse(r.text, expected)
        if response is None:
            return None

        dataToSend = dict()
        dataToSend["project"] = response
        for k, v in data["details_entry"].items():
            dataToSend["project"]["Assets"]["DataMap"][k] = v

    try:
        r = httpConnection.POST("/update-project-assets", dataToSend)
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
                "name": "testUserUpdateSettings",
                "email": "testEmailUpdateSettings",
                "password": "testPassword"
            },
            "details_entry": {
                "test_entry": "test_data"
            }
        },
        # Expected
        "OK")
]

ids = ['Valid user settings']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_UpdateUserSettings(httpConnection, data, expected):
    dataToSend = dict()
    if "user" in data:
        try:
            r = httpConnection.POST("/add-user", data["user"])
        except Exception:
            pytest.fail("Failed to send POST request")
            return

        if r.status_code != 201:
            pytest.fail(f"Failed to add user.\nDetails: {r.text}")
            return

        response = common.getResponse(r.text, expected)
        if response is None:
            return None
        dataToSend["user"] = response
        for k, v in data["details_entry"].items():
            dataToSend["user"]["Settings"]["DataMap"][k] = v

    try:
        r = httpConnection.POST("/update-user-settings", dataToSend)
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
              "name": "testUserUpdateUserAssets",
              "email": "testEmailUpdateUserAssets",
              "password": "testPassword"
          },
          "details_entry": {
              "test_entry": "test_data"
          }
      },
      # Expected
      "OK")
]

ids = ['Valid user assets']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_UpdateUserAsset(httpConnection, data, expected):
    dataToSend = dict()
    if "user" in data:
        try:
            r = httpConnection.POST("/add-user", data["user"])
        except Exception:
            pytest.fail("Failed to send POST request")
            return

        if r.status_code != 201:
            pytest.fail(f"Failed to add user.\nDetails: {r.text}")
            return

        response = common.getResponse(r.text, expected)
        if response is None:
            return None
        dataToSend["user"] = response
        for k, v in data["details_entry"].items():
            dataToSend["user"]["Assets"]["DataMap"][k] = v

    try:
        r = httpConnection.POST("/update-user-assets", dataToSend)
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

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
              "username": "testUserOwnerUpdateDetails",
              "email": "testEmailOwnerUpdateDetails",
              "password": common.convertPasswdToBase64("testPassword")
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
        dataToSend["user"] = response["id"]

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
            dataToSend["product"]["details"]["datamap"][k] = v

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
                "username": "testUserOwnerUpdateAssets",
                "email": "testEmailOwnerUpdateAssets",
                "password": common.convertPasswdToBase64("testPassword")
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
        dataToSend["user"] = response["id"]

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
            dataToSend["product"]["assets"]["datamap"][k] = v

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
                "username": "testUserProjectUpdateDetails",
                "email": "testEmailProjectUpdateDetails",
                "password": common.convertPasswdToBase64("testPassword")
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
            dataToSend["project"]["details"]["datamap"][k] = v

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
                "username": "testUserProjectUpdateAssets",
                "email": "testEmailProjectUpdateAssets",
                "password": common.convertPasswdToBase64("testPassword")
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
            dataToSend["project"]["assets"]["datamap"][k] = v

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
                "username": "testUserUpdateSettings",
                "email": "testEmailUpdateSettings",
                "password": common.convertPasswdToBase64("testPassword")
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
        dataToSend["user-id"] = response["id"]
        dataToSend["user-data"] = response["settings"]
        for k, v in data["details_entry"].items():
            dataToSend["user-data"]["datamap"][k] = v

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
              "username": "testUserUpdateUserAssets",
              "email": "testEmailUpdateUserAssets",
              "password": common.convertPasswdToBase64("testPassword")
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
        dataToSend["user-id"] = response["id"]
        dataToSend["user-data"] = response["assets"]
        for k, v in data["details_entry"].items():
            dataToSend["user-data"]["datamap"][k] = v

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

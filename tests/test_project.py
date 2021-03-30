import pytest
import common

dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "product": {
              "name": "testProductAddProject",
            },
            "user": {
              "username": "testUserOwnerAddProject",
              "email": "testEmailOwnerAddProject",
              "password": "testPassword"
            },
            "project": {
              "name": "testProjectAddProject",
              "visibility": "Public"
            }
        },
        # Expected
        {
          "name": "testProjectAddProject",
          "visibility": "Public"
        }),
    (
        # Input data
        {
            "project": {
              "name": "testProjectMissingUser",
              "visibility": "Public"
            },
            "user_id": "c34a7368-344a-11eb-adc1-0242ac120002",
            "product_id": "c34a7368-344a-11eb-adc1-0242ac120002"
        },
        # Expected
        "The selected product not found")

]

ids = ['No existing project', 'Missing product']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_CreateProject(httpConnection, data, expected):
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    dataToSend = dict()
    dataToSend["project"] = data["project"]
    dataToSend["product_id"] = productUUID
    dataToSend["owner_id"] = userUUID

    try:
        r = httpConnection.POST("/add-project", dataToSend)
    except Exception:
        pytest.fail("Failed to send POST request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 201:
        details = response["details"]["datamap"]
        if details["name"] != expected["name"] or \
                details["visibility"] != expected["visibility"]:
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
              "name": "testProductGetProject"
          },
          "user": {
              "username": "testUserOwnerGetProject",
              "email": "testEmailOwnerGetProject",
              "password": "testPassword"
          },
          "project": {
              "name": "testProjectGetProject",
              "visibility": "Public"
          }
      },
      # Expected
      {
          "name": "testProjectGetProject",
          "visibility": "Public"
      }),

    (
        # Input data
        {
            "product": {
                "name": "testProductGetProjectMissing"
            },
            "user": {
                "username": "testUserOwnerGetProjectMissing",
                "email": "testEmailOwnerGetProjectMissing",
                "password": "testPassword"
            },
            "id": "c34a7368-344a-11eb-adc1-0242ac120002"
        },
        # Expected
        "The selected project not found")
]

ids = ['Existing project', 'Missing project']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetProject(httpConnection, data, expected):
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    projectUUID = common.addProject(
      data,
      userUUID,
      productUUID,
      httpConnection)
    if projectUUID is None:
        return

    dataToSend = dict()
    dataToSend["id"] = projectUUID

    try:
        r = httpConnection.GET("/get-project", dataToSend)
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 200:
        details = response["details"]["datamap"]
        if details["name"] != expected["name"] or \
                details["visibility"] != expected["visibility"]:
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
              "name": "testProductGetProjectMultiple"
            },
            "user": {
              "username": "testUserOwnerGetProjectMultiple",
              "email": "testEmailOwnerGetProjectMultiple",
              "password": "testPassword"
            },
            "project": [
                {
                    "name": "testProjectGetProjectMultiple1",
                    "visibility": "Public"
                },
                {
                    "name": "testProjectGetProjectMultiple2",
                    "visibility": "Protected"
                }]
        },
        # Expected
        [
            {
                "name": "testProjectGetProjectMultiple1",
                "visibility": "Public"
            },
            {
                "name": "testProjectGetProjectMultiple2",
                "visibility": "Protected"
            }]
        ),
    (
      # Input data
      {
          "product": {
            "name": "testProductGetProjectMultiple2"
          },
          "user": {
            "username": "testUserOwnerGetProjectMultiple2",
            "email": "testEmailOwnerGetProjectMultiple2",
            "password": "testPassword"
          },
          "project": [
              {
                  "name": "testProjectGetProjectMultiple2",
                  "visibility": "Public"
              },
              {
                  "id": "c34a7368-344a-11eb-adc1-0242ac120002"
              }]
      },
      # Expected
      [{
          "name": "testProjectGetProjectMultiple2",
          "visibility": "Public"
      }]
      ),
    (
        # Input data
        {
          "product": {
              "name": "testProductGetProjectMultiple3"
          },
          "user": {
              "username": "testUserOwnerGetProjectMultiple3",
              "email": "testEmailOwnerGetProjectMultiple3",
              "password": "testPassword"
          },
          "project": [{
              "id": "c34a7368-344a-11eb-adc1-0242ac120002"
          }]
        },
        # Expected
        "The selected project not found"
        )
]

ids = ['Existing projects', 'Missing a project', 'No project']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetProjects(httpConnection, data, expected):
    uuidList = list()
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    uuidList = common.addProjects(data, userUUID, productUUID, httpConnection)
    if uuidList is None:
        return

    try:
        r = httpConnection.GET("/get-projects", {"ids": uuidList})
    except Exception:
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 200:
        for index, product in enumerate(response):
            details = product["details"]["datamap"]
            if details["name"] != expected[index]["name"] or \
                    details["visibility"] != expected[index]["visibility"]:
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
              "name": "testProductGetProductProjects"
            },
            "user": {
              "username": "testUserOwnerGetProductProjects",
              "email": "testEmailOwnerGetProductProjects",
              "password": "testPassword"
            },
            "project": [
                {
                    "name": "testProjectGetProductProjects1",
                    "visibility": "Public"
                },
                {
                    "name": "testProjectGetProductProjects2",
                    "visibility": "Protected"
                }]
        },
        # Expected
        [{
            "name": "testProjectGetProductProjects1",
            "visibility": "Public"
        }, {
            "name": "testProjectGetProductProjects2",
            "visibility": "Protected"
        }]
        ),
    (
        # Input data
        {
            "product": {
                "name": "testProductGetProductProjects2"
            },
            "user": {
              "username": "testUserOwnerGetProductProjects2",
              "email": "testEmailOwnerGetProductProjects2",
              "password": "testPassword"
            },
            "project": [{
                "name": "testProjectGetProductProjects2",
                "visibility": "Public"
            }, {
                "id": "c34a7368-344a-11eb-adc1-0242ac120002"
            }]
        },
        # Expected
        [{
            "name": "testProjectGetProductProjects2",
            "visibility": "Public"
        }]
    ),

    (
      # Input data
      {
          "user": {
              "username": "testUserOwnerGetProductProjects3",
              "email": "testEmailOwnerGetProductProjects3",
              "password": "testPassword"
          },
          "product": {
              "name": "testProductGetProductProjects3"
          },
      },
      # Expected
      "No projects for this product"
    )
]

ids = ['Existing projects', 'Missing a project', 'No project']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetProductProjects(httpConnection, data, expected):
    uuidList = list()
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    uuidList = common.addProjects(data, userUUID, productUUID, httpConnection)
    if uuidList is None:
        return

    try:
        r = httpConnection.GET(
            "/get-product-projects", {"product_id": productUUID})
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None
    if r.status_code == 200:
        for index, product in enumerate(response):
            details = product["details"]["datamap"]
            if details["name"] != expected[index]["name"] or \
                    details["visibility"] != expected[index]["visibility"]:
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
              "name": "testProductDeleteProject"
            },
            "user": {
              "username": "testUserOwnerDeleteProject",
              "email": "testEmailOwnerDeleteProject",
              "password": "testPassword"
            },
            "project": {
              "name": "testProjectDeleteProject",
              "visibility": "Public"
            }
        },
        # Expected
        "OK"),

    (
      # Input data
      {
          "product": {
              "name": "testProductDeleteProjectMissing"
          },
          "user": {
              "username": "testUserOwnerDeleteProjectMissing",
              "email": "testEmailOwnerDeleteProjectMissing",
              "password": "testPassword"
          },
          "id": "c34a7368-344a-11eb-adc1-0242ac120002"
      },
      # Expected
      "The selected project not found")
]

ids = ['Existing project', 'Missing project']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_DeleteProject(httpConnection, data, expected):
    userUUID = common.addUser(data, httpConnection)
    if userUUID is None:
        return

    productUUID = common.addProduct(data, userUUID, httpConnection)
    if productUUID is None:
        return

    projectUUID = common.addProject(
        data,
        userUUID,
        productUUID,
        httpConnection)
    if projectUUID is None:
        return

    dataToSend = dict()
    dataToSend["id"] = projectUUID

    try:
        r = httpConnection.POST("/delete-project", dataToSend)
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

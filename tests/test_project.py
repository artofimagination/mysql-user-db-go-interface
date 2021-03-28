import pytest
import common
import json

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
        {
            "error": "The selected product not found"
        })

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
        {
            "error": "The selected project not found"
        })
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
        {
            "error": "The selected project not found"
        })
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
      {
          "error": "No projects for this product"
      }
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
      {
          "error": "The selected project not found"
      })
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


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "viewer_id": "d66aa5f8-2b83-49b5-bb0c-28b8336f7f34",
            "user": {
              "username": "testUserOwnerAddProjectViewer1",
              "email": "testEmailOwnerAddProjectViewer1",
              "password": "testPassword"
            },
            "product": {
              "name": "testProductAddProjectViewer1",
            },
            "project": {
              "name": "testProjectAddProjectViewer1",
              "visibility": "Public"
            },
            "is_owner": True
        },
        # Expected
        {
          "data": "OK",
          "error": ""
        }),
    (
        # Input data
        {
            "viewer_id": "e514d186-0594-4ee0-badd-56ff712be040",
            "product": {
              "name": "testProductAddProjectViewer2",
            },
            "user": {
              "username": "testUserOwnerAddProjectViewer2",
              "email": "testEmailOwnerAddProjectViewer2",
              "password": "testPassword"
            },
            "id": "c34a7368-344a-11eb-adc1-0242ac120002",
            "is_owner": True
        },
        # Expected
        {
          "data": "",
          "error": "The selected project not found"
        }),
    (
        # Input data
        {
            "viewer_id": "d66aa5f8-2b83-49b5-bb0c-28b8336f7f34",
            "product": {
              "name": "testProductAddProjectViewer3",
            },
            "user": {
              "username": "testUserOwnerAddProjectViewer3",
              "email": "testEmailOwnerAddProjectViewer3",
              "password": "testPassword"
            },
            "project": {
              "name": "testProjectAddProjectViewer3",
              "visibility": "Public"
            },
            "is_owner": True
        },
        # Expected
        {
          "data": "",
          "error": "Viewer already exists"
        })
]

ids = ['Success', 'Failure', 'Duplicate owner']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_CreateProjectViewer(httpConnection, data, expected):
    r = common.addProjectViewer(data, httpConnection)
    if r is None:
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None

    if response != expected["data"]:
        pytest.fail(
          f"Request failed\nStatus code: \
          {r.status_code}\nReturned: {response}\nExpected: {expected}")


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "viewer_id": "eda3e9b4-d011-48c7-af8c-07446b628def",
            "user": {
              "username": "testUserOwnerGetProjectViewer1",
              "email": "testEmailOwnerGetProjectViewer1",
              "password": "testPassword"
            },
            "product": {
              "name": "testProductGetProjectViewer1",
            },
            "project": {
              "name": "testProjectGetProjectViewer1",
              "visibility": "Public"
            },
            "is_owner": True
        },
        # Expected
        {
          "data": [
              {
                  "viewer_id": "eda3e9b4-d011-48c7-af8c-07446b628def",
                  "is_owner": True
              }
          ],
          "error": ""
        }),
    (
        # Input data
        {
            "viewer_id": "ea1724ef-e426-43a4-8030-25c3468ef3a2",
        },
        # Expected
        {
          "data": "",
          "error": "The selected project viewer not found"
        })
]

ids = ['Success', 'Not found']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetProjectViewerByViewerID(httpConnection, data, expected):
    if "user" in data:
        r = common.addProjectViewer(data, httpConnection)
        if r is None:
            return

    try:
        dataToSend = dict()
        dataToSend["viewer_id"] = data["viewer_id"]
        r = httpConnection.GET("/get-project-viewer-by-viewer", dataToSend)
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None

    if response[0]["viewer_id"] != expected["data"][0]["viewer_id"]:
        pytest.fail(
          f"Request failed\nStatus code: \
          {r.status_code}\nReturned: {response}\nExpected: {expected}")


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "viewer_id": "9d5f9412-dd75-4907-82cf-3041de584a30",
            "user": {
              "username": "testUserOwnerGetProjectViewer2",
              "email": "testEmailOwnerGetProjectViewer2",
              "password": "testPassword"
            },
            "product": {
              "name": "testProductGetProjectViewer2",
            },
            "project": {
              "name": "testProjectGetProjectViewer2",
              "visibility": "Public"
            },
            "is_owner": True
        },
        # Expected
        {
          "data": [
              {
                  "viewer_id": "9d5f9412-dd75-4907-82cf-3041de584a30",
                  "is_owner": True
              }
          ],
          "error": ""
        }),
    (
        # Input data
        {
            "user": {
              "username": "testUserOwnerGetProjectViewer3",
              "email": "testEmailOwnerGetProjectViewer3",
              "password": "testPassword"
            },
        },
        # Expected
        {
          "data": "",
          "error": "User is not connected to any viewer"
        })
]

ids = ['Success', 'Not found']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetProjectViewerByUserID(httpConnection, data, expected):
    if "user" in data:
        if "project" in data:
            r = common.addProjectViewer(data, httpConnection)
            if r is None:
                return
        else:
            r = common.addUser(data, httpConnection)
            if r is None:
                return

    try:
        dataToSend = dict()
        dataToSend["email"] = data["user"]["email"]
        r = httpConnection.GET("/get-user-by-email", dataToSend)
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    try:
        response = json.loads(r.text)
    except Exception as e:
        pytest.fail(f"Failed to decode json. Details: {e}")
        return

    try:
        dataToSend = dict()
        dataToSend["user_id"] = response["data"]["id"]
        r = httpConnection.GET("/get-project-viewer-by-user", dataToSend)
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None

    if response[0]["viewer_id"] != expected["data"][0]["viewer_id"]:
        pytest.fail(
          f"Request failed\nStatus code: \
          {r.status_code}\nReturned: {response}\nExpected: {expected}")


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "viewer_id": "56439ce5-b2bb-4278-ba97-247ac6a90d9e",
            "user": {
              "username": "testUserOwnerDeleteProjectViewer1",
              "email": "testEmailOwnerDeleteProjectViewer1",
              "password": "testPassword"
            },
            "product": {
              "name": "testProductDeleteProjectViewer1",
            },
            "project": {
              "name": "testProjectDeleteProjectViewer1",
              "visibility": "Public"
            },
            "is_owner": True
        },
        # Expected
        {
          "data": "OK",
          "error": ""
        }),
    (
        # Input data
        {
            "viewer_id": "c5ec82f4-d6df-4057-8dc1-763709c7810e"
        },
        # Expected
        {
          "data": "",
          "error": "No project viewer was deleted"
        })
]

ids = ['Success', 'Not found']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_DeleteProjectViewerByViewerID(httpConnection, data, expected):
    if "user" in data:
        if "project" in data:
            r = common.addProjectViewer(data, httpConnection)
            if r is None:
                return
        else:
            r = common.addUser(data, httpConnection)
            if r is None:
                return

    try:
        dataToSend = dict()
        dataToSend["viewer_id"] = data["viewer_id"]
        r = httpConnection.POST("/delete-project-viewer-by-viewer", dataToSend)
    except Exception:
        pytest.fail("Failed to send POST request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None

    if response != expected["data"]:
        pytest.fail(
          f"Request failed\nStatus code: \
          {r.status_code}\nReturned: {response}\nExpected: {expected}")


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "viewer_id": "f1049c19-dae6-46dd-926f-791f149b14c8",
            "user": {
              "username": "testUserDeleteProjectViewer2",
              "email": "testEmailDeleteProjectViewer2",
              "password": "testPassword"
            },
            "product": {
              "name": "testProductDeleteProjectViewer2",
            },
            "project": {
              "name": "testProjectDeleteProjectViewer2",
              "visibility": "Public"
            },
            "is_owner": True
        },
        # Expected
        {
          "data": "OK",
          "error": ""
        }),
    (
        # Input data
        {
            "user": {
              "username": "testUserDeleteProjectViewer3",
              "email": "testEmailDeleteProjectViewer3",
              "password": "testPassword"
            },
        },
        # Expected
        {
          "data": "",
          "error": "No project viewer was deleted"
        })
]

ids = ['Success', 'Not found']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_DeleteProjectViewerByUserID(httpConnection, data, expected):
    if "user" in data:
        if "project" in data:
            r = common.addProjectViewer(data, httpConnection)
            if r is None:
                return
        else:
            r = common.addUser(data, httpConnection)
            if r is None:
                return

    try:
        dataToSend = dict()
        dataToSend["email"] = data["user"]["email"]
        r = httpConnection.GET("/get-user-by-email", dataToSend)
    except Exception:
        pytest.fail("Failed to send GET request")
        return

    try:
        response = json.loads(r.text)
    except Exception as e:
        pytest.fail(f"Failed to decode json. Details: {e}")
        return

    try:
        dataToSend = dict()
        dataToSend["user_id"] = response["data"]["id"]
        r = httpConnection.POST("/delete-project-viewer-by-user", dataToSend)
    except Exception:
        pytest.fail("Failed to send POST request")
        return

    response = common.getResponse(r.text, expected)
    if response is None:
        return None

    if response != expected["data"]:
        pytest.fail(
          f"Request failed\nStatus code: \
          {r.status_code}\nReturned: {response}\nExpected: {expected}")

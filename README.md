# ase-3-6-attendance-tracking-web
## Installation
Make sure you add the repository under 
`$GOPATH/src/github.com/johannesrohwer/ase-3-6-attendance-tracking-web`. 
Install the required dependencies via the `go get` tool:
```
go get github.com/google/uuid
go get github.com/dgrijalva/jwt-go
go get github.com/gorilla/mux
go get google.golang.org/appengine
go get google.golang.org/appengine/datastore
go get github.com/dgrijalva/jwt-go
go get golang.org/x/crypto/bcrypt
```

After all dependencies are satisfied a local instance of the application can be launched by using:
`dev_appserver.py api/app.yaml`

The `dev_appserver.py` tool is available as part of the `gcloud` (including `google-cloud-sdk-app-engine-go`)
package that is available [here](https://cloud.google.com/sdk/docs/). In order to run it `python2.7` is required.

## Deployment
To deploy the application to Google App Engine make sure that you initialized the `gcloud` CLI tool and run

`gcloud app deploy api/app.yaml`


## API Documentation
The API is using JSON encoding. Therefore, all requests need to set the `Content-Type` header
to `application/json`.

### Authorization
Certain routes require special authorization. The token that is required to access those routes can be acquired via
`/api/login` and has to be added in the `Authorization` header of every request to the protected routes.

### `/api/login`
In order to log in send an `HTTP POST` request and attach the following payload:
```json
{
  "id": "<student or instructor id>",
  "password": "<password>"
}
```


will answer with a JWT token that is passed via a JSON object.

```json
{
  "token": "<base64 encoded JWT>"
}
```

This token contains access permission information and should
be send with every further request as an `Authorization` header.

### `/api/signup`
In order to create new students the `/api/signup` that is an alias to `/api/students`route is used.
To add a new student, perform a POST request with the following payload:
```json
{
  "id": "<matriculation number>",
  "group_id": "<ID of the group>",
  "name": "<first and last name>",
  "password": "<some password>"
}
```

The server will add a new student to the database and return the created student object with an `HTTP 201` code.

### `/api/students`
Besides the alias to create new students the `/api/students` endpoint enables the API client to fetch data about the registered students.
Using an `HTTP GET` request the server will return an array of registered students:
```json
[
  {
    "id": "32126781",
    "group_id": "01",
    "name": "Sam Student"
  },
  {
    "id": "35218879",
    "group_id": "02",
    "name": "Lucas Learner"
  }
]
```

Besides a list of students the endpoint also provides access to a single student record via the student ID.
An `HTTP GET` to `/api/students/32126781` results in:
```json
{
    "id": "32126781",
    "group_id": "01",
    "name": "Sam Student"
}
```

### `/api/instructors`
The instructor API is very similar to the student one. A new instructor is added via `HTTP POST`
to `/api/instructors` with a payload of
```json
{
  "id": "<matriculation number>",
  "name": "<first and last name>",
  "password": "<some password>"
}
```

A list of instructors can be fetched via `HTTP GET` to `/api/instructors` and a single instructor recorded can be
loaded via an `HTTP GET` request to `/api/instructors/{instructor_id}`.

### `/api/groups`
This endpoint provides information about the different groups.
To add a new group perform a `HTTP POST` request to `/api/groups` with the following payload:
```json
{
  "id": "<group id>",
  "place": "<room, e.g. 01.07.023>",
  "time": "<time in human readable format, e.g. Monday, 10-12>",
  "instructor_id": "<id of the responsible instructor>"
}
```

A list of groups can be accessed by using a `HTTP GET` request to `/api/groups` and a single group record can be fetched
via an `HTTP GET` to `/api/groups/{group_id}`.

### `/api/attendances`
This endpoint provides methods for students to get their personalized attendance token that can
later on be registered by an instructor.
To fetch a new attendance token the user must be logged in as a student and perform an `HTTP GET`
to `/api/attendances/new/{student_id}/{presented}` where the presented option is a boolean (`true or false`)
that indicates if a student presented an exercise in this week.
To fetch all attendances for a particular student an `HTTP GET` request to `/api/attendances/for/{student_id}` is required.
An optional query parameter `?missing_attendances=true` can be added to show dummy entries for weeks where the student was absent.

When an instructor wants to register a students attendance token as valid it has to be send via `HTTP POST`
to `/api/attendances/register` as a payload:
```json
{
  "token": "<attendance token>"
}
```

As usual, the attendance list can be obtained via `HTTP GET` to `/api/attendances` and a single attendance record via 
`/api/attendances/{attendance_id}`

### `/api/week/current`
This `HTTP GET` endpoint return the current week id as a JSON encoded object.

### `/api/version`
Returns a JSON encoded version object that includes authors and version number.


## Miscellaneous
 - __WARNING__: The application is currently using a hardcoded __unsafe encryption key__.
 - The API is not stable. Further changes will most likely appear.

## Raspberry Pi Repository

The repository containing the Raspberry Pi implementation can be found here:
https://github.com/johannesrohwer/ase-3-6-attendance-tracking-pi

## Android App Repository

The repository containing the Android application can be found here:

https://github.com/PSchmiedmayer/ase-3-6-attendance-tracking-android
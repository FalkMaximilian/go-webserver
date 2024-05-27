# go-webserver
A Golang backend for a simple flashcards REST API.

This project makes use of Fiber and Gorm.

Passwords are stored in a Postgresql database with bcrypt.

A loggin middleware is used to log every request that has been handled by the server.

On successful login/registration a jwt token is issued which contains the user_id. On each request this jwt is required to perform database operations.
The user id is extracted from the jwt when calling an endpoint.

## TODO
* Finish CRUD implementation
* Finish loading config in one place (config package)
* Remove jwt from client after user deletion
* Add endpoint for logout
* Remove db operations from handlers and move db operations to model
* Implement 'get all cards for set'


# Local development
## Database setup

I recommend to create a local postgresql instance. 

I run Postgresql in docker.

`sudo docker pull postgres`

Start docker image with the following command through the cli

`sudo docker run -d --name testdb -p 5432:5432 -e POSTGRES_PASSWORD=Start12345@ -v postgres:/var/lib/postgresql/data postgres:latest`

Or in Docker Desktop and set the following environment variables

```
POSTGRES_PASSWORD=Start12345@
POSTGRES_USER=admin
POSTGRES_DB=flashcards
```

`sudo docker ps`

`sudo docker stop <CONTAINER ID>`

## DB setup Flashcards project

A .env file has to be created in the root folder of the flashcards project 

```
PGHOST=localhost
PGPORT=5432
PGUSER=postgres
PGPASSWORD=Start12345@
PGNAME=postgres
```

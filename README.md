# go-webserver
Just trying out webservers with go

# Local development

I recommend to create a local postgresql instance. 

I run Postgresql in docker.

`sudo docker pull postgres`

`sudo docker run -d --name testdb -p 5432:5432 -e POSTGRES_PASSWORD=Tardis1808! -v postgres:/var/lib/postgresql/data postgres:latest`

To connect to it a .env file has to be created in the root folder.

```
PGHOST=localhost
PGPORT=5432
PGUSER=postgres
PGPASSWORD=Tardis1808!
PGNAME=postgres
```

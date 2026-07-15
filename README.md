# Task Manager

## This is golang rest api. Tech stack is below

* __Chi Router__
* __Postgresql__
* __pgx for postgres driver__
* __Json Web Token__
* __Docker Compose__
* __Docker__

### Detailed endpoints

* <http://localhost:PORT/auth>
  * POST /register
  * POST /login

* <http://localhost:PORT/users>
  * GET /me -> to check AuthMiddleware

* <http://localhost:PORT/tasks>
  * POST / -> add single task
  * GET / -> get all tasks by autheticated user id via jwt context
  * GET /{task_id} -> get single task by id with autheticated user id via jwt context
  * PATCH /{task_id} -> update a task
  * DELETE /{task_id} -> delete a task
  * __All tasks route requires valid jwt token with Authorization Header__

### How to run

You just need docker installed. Then create .env file like below:

* *DATABASE_URL=postgresql://user:password@db:5432/tamimsdb?sslmode=disable*
* *DB_USER=username*
* *DB_PASSWORD=password*
* *DB_NAME=database_name*
* *PORT=5000*
* *JWT_SECRET=create secret for jwt secret*

Then run:

```bash
    docker compose up -d
```

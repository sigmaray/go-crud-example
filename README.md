# Example of admin panel with authentication and CRUD implemented with Go, Gin, Gorm

- Sign in
- Sign out
- Viewing, adding, editing, deleting users

## TODO:

- Don't display passwords in admin panel
- Encrypt passwords
- Move db credentials into .env
- Don't expose internal error messages (for example SQL errors) to user

## How to run app

* Install PostgreSQL. Edit db credentials in `main.go`
* `go run .`

User with admin:admin credentials is creating during first run

## How to run app with Docker

`docker-compose up`

## How to run Selenium tests

- Start go app (`go run .`)
- `poetry install`
- `poetry run pytest test.py`

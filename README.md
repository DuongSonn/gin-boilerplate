# GIN BOILERPLATE

## Prerequisites

## CMD

-   Server

    -   `make run`: Run project. Default at port 8080
    -   `make dev`: Run project using air (hot reload). Default at port 8080

-   Test:
    -   `mockery  --output ./mocks/repository --dir ./app/repository --all`: Generate mockfile base on interface
    -   `gotests -all -w app/helper/user.helper.go`: Generate test file

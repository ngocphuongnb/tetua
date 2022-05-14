# Tetua - A simple CMS for blogging

Tetua is a simple CMS for blogging written in Golang. With tetua, you can quickly create a blog with a few simple commands.

Tetua is built on top of awesome libraries like:

- [Fiber](https://github.com/gofiber/fiber)
- [Ent](https://github.com/ent/ent)
- [Rclone](https://github.com/rclone/rclone)
- [Jade](https://github.com/Joker/jade)

## Installation

### Download the binary from the release page:
- https://github.com/ngocphuongnb/tetua/releases

### Create the config file

```sh
./tetua init
```

The config file will be created in the current directory.

```json
{
 "app_env": "production",
 "app_key": "{APP_KEY}",
 "app_port": "3000",
 "db_dsn": "",
 "github_client_id": "",
 "github_client_secret": "",
 "db_query_logging": false
}
```

These fields are required:

- `app_key`: the key to encrypt the data
- `db_dsn`: the database connection string
- `github_client_id`: the client id for github
- `github_client_secret`: the client secret for github

You can skip this initialization step by specifying the environment variables:
- `APP_KEY`
- `DB_DSN`
- `GITHUB_CLIENT_ID`
- `GITHUB_CLIENT_SECRET`

### Create the Admin account
```sh
./tetua setup -u admin -p password
```

### Run the server

```sh
./tetua run
```

## Features

* Posts Management
* Topics Management
* Users Management
* Role and Permission Management
* Site Settings Management
* Comment Management
* File Management
* User profile page
* User posts page
* Local file upload
* S3 file upload
* Sign in with Github

## Documentation
- Documentation is hosted live at [https://tetua.net](https://tetua.net)

## Development

*Requirements:*

- Go 1.18+
- [Goreleaser](https://github.com/goreleaser/goreleaser)

The development requires Go 1.18+ as we use `generic`

### Clone the source code

```sh
git clone https://github.com/ngocphuongnb/tetua.git
cd tetua
go mod tidy
```

### Build the static editor

```sh
make build_editor
```

### Run the test

```sh
make test_alll
```

### Build the local release for testing

```sh
make releaselocal
```

### Public a release

```sh
git tag -a vx.y.z -m "Release note"
git push origin vx.y.z
make release
```

## Road Map
* [ ] Pages cache
* [ ] Sign in with Google
* [ ] Sign up with email (local account)
* [ ] Serial posts
* [ ] Report Abuse
* [ ] Complete the Unit Test

## Screenshots

![image](https://user-images.githubusercontent.com/3405842/167983805-ff26b8dc-27cb-4aa6-ae84-8ebfefae7dc8.png)

![image](https://user-images.githubusercontent.com/3405842/167983866-32b3444e-591f-47e8-8b0a-d3f0cfd03aa3.png)

![image](https://user-images.githubusercontent.com/3405842/167983936-66624f6b-660b-4ccf-a19f-71d35926c405.png)

![image](https://user-images.githubusercontent.com/3405842/167984402-295dc7df-8286-4d8a-975e-ae097d9fa9ad.png)

![image](https://user-images.githubusercontent.com/3405842/167984104-d08c9b3e-8f87-4041-b04a-ae384a1f46aa.png)


## Contribute
If you want to say thank you and/or support the active development of Tetua, please consider some of the following:

1. Add a GitHub Star to the project.
2. Create a pull request.
3. Fire an issue.

## License
Copyright (c) 2022-present @ngocphuongnb and Contributors. Tetua is free and open-source software licensed under the MIT License.

**Third-party libraries:**

- entgo.io/ent
- github.com/Joker/hpp
- github.com/Joker/jade
- github.com/go-sql-driver/mysql
- github.com/gofiber/fiber/v2
- github.com/golang-jwt/jwt/v4
- github.com/google/uuid
- github.com/microcosm-cc/bluemonday
- github.com/valyala/fasthttp
- go.uber.org/zap
- github.com/davecgh/go-spew
- github.com/rclone/rclone
- github.com/urfave/cli/v2
- ariga.io/sqlcomment
- github.com/PuerkitoBio/goquery
- github.com/gofiber/utils
- github.com/gorilla/feeds
- github.com/tdewolff/minify/v2
- github.com/gosimple/slug
- github.com/stretchr/testify
- github.com/yuin/goldmark

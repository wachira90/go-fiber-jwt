# go-fiber-jwt
golang fiber jsonwebttoken

## version

```
go version go1.21.1 windows/amd64
golang.org/x/crypto v0.22.0
github.com/dgrijalva/jwt-go in github.com/dgrijalva/jwt-go v3.2.0+incompatible
github.com/gofiber/fiber/v2 in github.com/gofiber/fiber/v2 v2.52.4
golang.org/x/crypto/bcrypt in golang.org/x/crypto v0.22.0
gorm.io/driver/postgres in gorm.io/driver/postgres v1.5.7
gorm.io/gorm in gorm.io/gorm v1.25.9
golang.org/x/sys v0.19.0
```

## command 

```
git clone 

cd go-fiber-jwt

go mod init go-fiber-jwt

go mod tidy
```

## register 

```
POST http://localhost:3000/register
Content-Type: application/json

{
    "username" : "wachira90",
    "password" : "password1234"
}
```

## result

```
HTTP/1.1 200 OK
Date: Tue, 23 Apr 2024 09:53:45 GMT
Content-Type: application/json
Content-Length: 218
Connection: close

{
  "ID": 1,
  "CreatedAt": "2024-04-23T16:53:45.4929246+07:00",
  "UpdatedAt": "2024-04-23T16:53:45.4929246+07:00",
  "DeletedAt": null,
  "username": "wachira90",
  "password": "$2a$10$kBgM/K7S5DdN12u9OrT1VuXps3/075/ckT8w4jxUaE7/6CVu9Rmja"
}
```


## login 

```
POST http://localhost:3000/login
Content-Type: application/json

{
    "username" : "wachira90",
    "password" : "password1234"
}
```
## result 

```
HTTP/1.1 200 OK
Date: Tue, 23 Apr 2024 09:55:58 GMT
Content-Type: application/json
Content-Length: 131
Connection: close

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM5NTI1NTksInN1YiI6IjEifQ.jm1PY3fcUIYLFfX4-5SWqDOhwxpOLcRw185vQb7-H9M"
}
```

## add data (when missing no header authen)

```
POST http://localhost:3000/books
Content-Type: application/json

{
  "title": "The Great Gatsby",
  "author": "F. Scott Fitzgerald",
  "rating": 5
}
```

## result (when missing no header authen)

```
HTTP/1.1 401 Unauthorized
Date: Tue, 23 Apr 2024 10:04:59 GMT
Content-Type: application/json
Content-Length: 51
Connection: close

{
  "error": true,
  "msg": "Missing Authorization header"
}
```

## add data

```
POST http://localhost:3000/books
Authorization: <TOKEN>
Content-Type: application/json

{
  "title": "The Great Gatsby",
  "author": "F. Scott Fitzgerald",
  "rating": 5
}
```

## get data

```
GET http://localhost:3000/books
Authorization: <TOKEN>
Content-Type: application/json
```
# go-fiber-jwt
golang fiber jsonwebttoken


## command 

git clone 

cd go-fiber-jwt

go mod init go-fiber-jwt

go mod tidy


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

## add data
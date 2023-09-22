# Dinosaur Jump
Made with Ebiten.

<img width="637" alt="Screen Shot 2021-12-12 at 23 55 10" src="https://user-images.githubusercontent.com/19848415/145717352-526ee3aa-c2ff-4fd4-8fe9-c4364324fdc2.png">

# How to start
```shell
$ go run main.go
```

クリエイト
```shell
Invoke-WebRequest -Method POST -Headers @{"Content-Type" = "application/json"} -Body '{"name":"YourUserName"}' -Uri http://localhost/user/create
```
ユーザーゲット
```shell
Invoke-WebRequest -Method POST -Headers @{"Content-Type" = "application/json"} -Body '{"auth_token":"2bd314be-ee78-4d33-926d-68e6894b8c57"}' -Uri http://localhost:8080/user/get
```
ランキング情報取得
```shell
Invoke-WebRequest -Method GET -Uri http://localhost:8080/users/get
```


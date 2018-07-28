Install dep
```
https://golang.github.io/dep/docs/installation.html
```

```bash
env GOCACHE=off go test ./...
```

```bash
curl -v -X POST -H "Content-Type: application/json" -d '{"username": "root", "password": "password"}' http://127.0.0.1:1234/auth/login
```

env can override value from config
```
GO_EXAMPLE_POSTGRESQL.CONNECTSTRING=host=172.24.0.2 user=postgres password=postgresqlPassword dbname=postgres connect_timeout=2 statement_timeout=2000 sslmode=disable
```

```bash
curl -v -X POST -H "Content-Type: application/json" -d '{"username": "nikit.cpp@yandex.ru", "password": "password"}' http://127.0.0.1:1234/auth/register
```

```bash
go get github.com/vektra/mockery/.../
```

```bash
(cd services/; mockery -name=Mailer)
```

```bash
dep ensure -add github.com/gobuffalo/packr@1.12.0
```
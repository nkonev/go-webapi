```bash
env GOCACHE=off go test ./...
```

```bash
curl -v -X POST -H "Content-Type: application/json" -d '{"username": "root", "password": "password"}' http://127.0.0.1:1234/auth2/login
```
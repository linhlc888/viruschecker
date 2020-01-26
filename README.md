## Scan Virus Online Tool
This is a toy application for learning go.

![Image](./mongo2es.png)

## Demo
1. Start clamd daemon in docker container or use running daemon on same host

```
docker run --rm --name clamav -p 3310:3310 quay.io/ukhomeofficedigital/clamav:latest

```

2. Run from source

```
go run *.go
```
3. Open http://localhost:8080

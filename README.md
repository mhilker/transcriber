# Transcriber

Convert audio speech to text.

## Build and run

### On your machine

```bash
$ go build -o build/server ./cmd/server/main.go
$ ./build/server
```

### Via docker

```bash
$ docker build -t mhilker/transcriber:latest -f cmd/server/Dockerfile .
$ docker run -p 8080:8080 mhilker/transcriber:latest
```

### Via docker-compose

```bash
$ docker-compose -f cmd/server/docker-compose.yml build
$ docker-compose -f cmd/server/docker-compose.yml up
```

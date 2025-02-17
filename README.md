
## Tech Stack

**Go:** I use go version ^1.22

**Gin-Gonic:** Because Gin-Gonic is very popular, so there's no reason to not to use this powerful framework

I've implement Clean Architecture design and Factory design pattern. So, you could change my existing Library or framework seamlessly.

## Installation

Install this project using go

```bash
  cp .config.example.json config.json
  go mod download && go mod tidy
```

to run this project

```bash
go run main.go
```

To run, watch, and build this project

```bash
CompileDaemon -exclude-dir=docs -command="./julong-onboarding-be"
```

To migrate the database

```bash
go run ./cmd/migration/main.go
```


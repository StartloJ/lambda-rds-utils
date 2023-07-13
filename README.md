# lambda-rds-utils

Utility tools helper for automate some activity with AWS resources

## Prerequisites

1. Task workflow; installation [here](https://taskfile.dev/installation/)
2. Go binary; version 1.18 and up
3. Pre commit; installation [here](https://pre-commit.com/#installation)

## How to build

Preparing a binary to deploy to AWS Lambda requires that it is compiled for Linux and placed into a .zip file.

```bash
# Remember to build your handler executable for Linux!
# When using the `provided.al2` runtime, the handler executable should be named `bootstrap`
GOOS=linux GOARCH=amd64 go build -o hello main.go
zip lambda-handler.zip hello
```

## Build and push to S3

1. You must accessible to AWS cli
2. Run terraform in folder `./terraform` to create S3 bucket before push your a Go binary.
3.

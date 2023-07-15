# lambda-rds-utils

Utility tools helper for automate some activity with AWS resources

## Prerequisites

1. Task workflow; installation [here](https://taskfile.dev/installation/)
2. Go binary; version 1.18 and up
3. Pre commit; installation [here](https://pre-commit.com/#installation)

## Current feature

Now we limit Event of focus on RDS below. you can see more detail [here](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Events.Messages.html)

1. `RDS-EVENT-0042`: Manual snapshot created.

## Required ENV config

- `OPT_SRC_REGION`: Insert source region of your snapshot created.
- `OPT_TARGET_REGION`: Insert destination region of your snapshot to copy.
- `OPT_BEDUG`: If provided, lambda should more log level.
- `OPT_DB_NAME`: Insert your RDS DB indentifier(Name) to specify this tasks copy.
- `OPT_OPTION_GROUP_NAME`: Insert your `option_group` name in destination region to copy.
- `OPT_KMS_KEY_ID`: Insert KMS key(ARN, ID, or Alias) to encrypted snapshot in destination region.

## How to build

Preparing a binary to deploy to AWS Lambda requires that it is compiled for Linux and placed into a .zip file.

   ```bash
   # Remember to build your handler executable for Linux!
   # When using the `provided.al2` runtime, the handler executable should be named `bootstrap`
   GOOS=linux GOARCH=amd64 go build -o hello main.go
   zip lambda-handler.zip hello
   ```

Or you can use `Taskfile` to automate the previous command.

   ```bash
   task build
   ```

## Build and push to S3

1. You must accessible to AWS cli
2. Run terraform in folder `./terraform` to create S3 bucket before push your a Go binary.
3. You can use `Taskfile` to executed.

    ```bash
    $ task build
    task: [build] go build -o $BINARY_NAME main.go
    task: [build] zip lambda-$BINARY_NAME.zip $BINARY_NAME
    updating: hello (deflated 55%)
    task: [clean] rm $BINARY_NAME

    $ task upload:lambda
    task: [upload:lambda] aws s3api put-object --bucket $S3_NAME --key lambda-$BINARY_NAME.zip --body ./lambda-$BINARY_NAME.zip
    ```

4. You can see the object on S3 bucket.

## How to use

1. If you create a function in AWS Lambda, you can set up a trigger from AWS EventBridge to capture and filter the Event from RDS to created snapshot.
2. Then it will execute in nearly time to run a function and gathering snapshot from source and destination region to compared and return unique snapshot that didn't copy to destination region.
3. The function will try to request AWS API for running the CopyDBSnapshot task.
4. Any activity, you can see the log from function on STDOUT(also Cloudwatch log).

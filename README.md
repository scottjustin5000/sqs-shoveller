# sqs-shoveller
> shovel messages between sqs queues


##### Build
```
go build -o sqs-shoveller
```

##### Usage

cli arguments

* -s address of the source queue (queue to move messages from)
* -d address of the destination queue (queue to move messages to)
* -r default AWS region

The utility assumes that there is either a valid AWS credential file or the following env variables set:

```
$ export AWS_ACCESS_KEY_ID=YOUR_AKID
$ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
```

running the cli:

```
./sqs-shoveller -s=https://sqs.us-east-1.amazonaws.com/1234567/some-source-queue -d=https://sqs.us-east-1.amazonaws.com/1234567/some-destination-queue -r=us-east-1
```



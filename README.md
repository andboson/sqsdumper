# sqsdumper

Dump AWS SQS messages to the output

Example:

```shell
<AWS_PROFILE=specific_profile> sqsdumper -s your-queue-dead-letter-queue 
```


### Help:

```shell
 ./sqsdumper -h                                                                                                 
NAME:
   sqsdumper - sqsdumper -s src_queue

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --deleteMessage               delete received messages (default: false)
   --help, -h                    show help (default: false)
   --jsonPath value, --jp value  json path, like x.y, see https://github.com/sinhashubham95/jsonic for more (default: .)
   --queueName value, -s value   the source queue
   --stopOnTotal                 stop when all messages processed (default: true)
   --version, -v                 print the version (default: false)

```

### Build

```shell
make all
```
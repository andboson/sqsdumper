# sqsdumper

Dump AWS SQS messages to the output

### Usage
```shell
<AWS_PROFILE=specific_profile> sqsdumper -s your-queue-dead-letter-queue 
```
with `jq`

```shell
<AWS_PROFILE=specific_profile> sqsdumper -s your-queue-dead-letter-queue | jq .foo
```

get in json_path (if no `jq` installed)

```shell
<AWS_PROFILE=specific_profile> sqsdumper -s your-queue-dead-letter-queue -jp foo
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
   --raw                         dump entire raw messages (default: false)
   --stopAfter value             stop after N messages processed (default: 0)
   --queueName value, -s value   the source queue
   --stopOnTotal                 stop when all messages processed (default: true)
   --version, -v                 print the version (default: false)

```

### Build

```shell
make all
```

### See also

 * https://github.com/mercury2269/sqsmover
 * https://github.com/prashanthpai/sqscat
 * https://github.com/farbodsalimi/go-sqs-wrapper
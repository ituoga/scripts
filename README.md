# Scripts PoC  - Moving from legacy system to new legacy

## Working on Legacy Systems

Sometimes we need somthing very simple

### Why? 

Imagine you have a system that is 20 years old, built on technologies that nobody wants to work with anymore. For instance, there's a need to enhance a CRM by adding new features, parsers, etc. This "proof of concept" system allows any program to be called over HTTP through MQ (nats.io), thereby enabling the transition of programming tasks from the "legacy system" to modern languages and contemporary technologies. 

### Simple example

Get JSON from `bash.json` subcriber
```
curl -X GET -s localhost:8088/mq?t=bash.json  |jq
```

### How to start


edit list.txt in format (add empty line at the end):
```
containername:mq.topic.here:script_from_scripts_dir.sh

```

run `bash gen.sh` to generate docker-compose.yml file


start system:

```
docker-compose up -d --build
```

execute script:

```
echo "hello world" | docker-compose exec -T cmd /app hello.world
```

should return something like:

```
got message: hello world
```

### other

and... pipelines:
```
echo "input" | docker-compose exec -T cmd /app hello.world 2>&1 | docker-compose exec -T cmd /app pipe2
```

should return something like:
```
got message with 1st message: got message: input
```


via http to nats proxy:
```
curl --data "hello my script" -X POST localhost:8088/mq?t=hello.world
```
# Scripts PoC


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

and... pipelines:
```
echo "input" | docker-compose exec -T cmd /app hello.world 2>&1 | docker-compose exec -T cmd /app pipe2
```

should return something like:
```
got message with 1st message: got message: input
```
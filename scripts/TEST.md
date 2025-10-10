


```bash
 docker run -d --name sentinel -p 8080:8080 -e POOLS_FILEPATH=/app/data.json -v ${PWD}/scripts/:/app/ sentinel:latest
```

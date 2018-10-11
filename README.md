# microservice

Build the image:

```
# docker build -t senarius/microservice:0.3 . -f Dockerfile
```

make sure add hosts:

```
# '127.0.0.1       microservice' into hosts file on respective OS
```

deploy:

```
sudo docker stack deploy microservice --compose-file=./docker-compose.yml
```

POST request sample:

```
curl --data-urlencode "url=http://studynihongo101.blogspot.com/search?updated-max=2012-05-03T05:31:00-07:00&max-results=1&start=1&by-date=false" -H 'Content-Type: application/x-www-form-urlencoded' http://microservice:8080/encode
```
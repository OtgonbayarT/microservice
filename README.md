# microservice

Build the image:

```
# docker build -t senarius/microservice:0.2 . -f Dockerfile
```

make sure add hosts:

```
# '127.0.0.1       microservice' into hosts file on respective OS
```

deploy:

```
sudo docker stack deploy microservice --compose-file=./docker-compose.yml
```
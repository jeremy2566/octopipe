# Docker Login

```shell
aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 858157298152.dkr.ecr.ap-southeast-1.amazonaws.com
```

# Docker 打包

```shell
docker build -t devops/octopipe .
```

# Docker Re Tag

```shell
docker tag devops/octopipe:latest 858157298152.dkr.ecr.ap-southeast-1.amazonaws.com/devops/octopipe:latest
```

# Docker push

```shell
docker push 858157298152.dkr.ecr.ap-southeast-1.amazonaws.com/devops/octopipe:latest
```
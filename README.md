# 部署crontab
## 构建镜像
``` 
docker build -t  192.168.239.161:5000/2410/user-crontab:1.0.0 -f Dockerfile.crontab .
```
## 创建网络
```shell
docker network create -d overlay 2410-user-net
```
## 创建配置文件
``` 
docker config create  2410-user-crontab-conf crontab/dev.config.yaml
```
## 部署服务
``` 
docker service create --name 2410-user-crontab -p 50052:50052 \
--config src=2410-user-crontab-conf,target=/app/config.yaml \
--network 2410-user-net \
--replicas 1 \
--health-cmd "grpc_health_probe -addr=:50052" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
192.168.239.161:5000/2410/user-crontab:1.0.0
```

# user

## 构建镜像
``` 
docker build -t 192.168.239.161:5000/2410/user:1.0.0  .
```
## 创建配置文件
``` 
docker config create 2410-user-conf user/dev.config.yaml
```

## 部署服务
``` 
docker service create --name 2410-user -p 8082:8082 \
--replicas 1 \
--network 2410-user-net \
--config src=2410-user-conf,target=/app/config.yaml \
--with-registry-auth \
192.168.239.161:5000/2410/user:1.0.0
```
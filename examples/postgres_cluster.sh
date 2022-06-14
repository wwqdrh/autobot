#!/bin/bash

# master
go run . postgres \
--ip 172.18.3.9 \
--username root \
--password root \
--mode master \
--name db-master-1 \
--dbport 54032 \
--dbpassword zxpostgres

# slave
go run . postgres \
--ip 172.18.3.9 \
--username root \
--password root \
--mode slave \
--name db-slave-1 \
--dbport 54033 \
--dbpassword zxpostgres \
--master_ip 172.18.3.9 \
--master_port 54032


# 创建postgres:postgres14的镜像 作为已有镜像的测试
docker run --name autobot-postgres14 \
-e POSTGRES_PASSWORD=zx123456 -d -p 54302:5432 \
postgres:14-alpine

docker run --name autobot-postgres14 \
-e POSTGRES_PASSWORD=zx123456 -d -p 54302:5432 \
postgres:12-alpine

# 初始化已有的容器作为master节点
go run . postgres \
--ip 172.18.3.9 \
--username root \
--password root \
--mode master \
--name autobot-postgres14 \
--dbport 54302 \
--dbpassword zxpostgres

# 接下来创建slave节点
go run . postgres \
--ip 172.18.3.9 \
--username root \
--password root \
--mode slave \
--name db-slave-1 \
--dbport 54033 \
--dbpassword zxpostgres \
--master_ip 172.18.3.9 \
--master_port 54302
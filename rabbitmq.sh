#! /bin/bash

#
# 参数1等于"stop"，则停止并清理rabbitmq容器
# 否则就启动一个新的rabbitmq容器
#

docker stop rabbitmq && docker rm rabbitmq

if [ "$1" = "stop" ];then
    exit 0
fi

docker run -d --hostname my-rabbit --name rabbitmq \
    -p 8080:15672 \
    -p 5672:5672  \
    rabbitmq:3.6-management

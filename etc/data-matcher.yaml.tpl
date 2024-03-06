# config file for data-matcher
#

version: 3.0

host: 0.0.0.0
manage_port: 3001

# 内部通道缓存大小，默认一百万，根据实际情况扩大到一千万时，大概需要内存4G
channel_size: 1000000

# 并行处理的工作协程数量，默认8个，实际情况应该根据CPU数量来.如果配置为0，则读取CPU数量然后减一。如果超过CPU数量，则配置为CPU数量然后减一
workers: 8

# 规则配置
rules_file: rules.yaml


# nats 服务配置
nats:
    servers:
        - nats://localhost:4222
    user: user
    password: passwd
    queue_name: data-matcher

http_flow:
    subject: httpTopic


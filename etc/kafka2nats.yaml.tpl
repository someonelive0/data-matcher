# config file for kafka2nats

version: 3.0

host: 0.0.0.0
manage_port: 3002


# kafka 服务配置
kafka:
    brokers:
        - localhost:9092
    # SASL machanism, such as PLAIN, SCRAM
    mechanism: PLAIN
    user: user
    password: passwd
    group: kafka2nats

# nats 服务配置
nats:
    servers:
        - nats://localhost:4222
    user: user
    password: passwd

# 读Kafka消息，写入nats
kakfa2nats:
    topic: httpTopic
    subject: httpTopic
    jetstream: false

# 读nats消息，写入kafka
nats2kafka:
    subject: matchTopic
    jetstream: true
    topic: matchTopic

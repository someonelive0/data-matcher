# config file for kafka2nats

# kafka 服务配置
kafka:
    brokers:
        - localhost:9092
    # SASL machanism, such as PLAIN, SCRAM
    mechanism: PLAIN
    user: user
    password: passwd
    group: kafka2nats

nats:
    servers:
        - nats://localhost:4222
    user: user
    password: passwd

kakfa2nats:
    topic: httpTopic
    subject: httpTopic
    jetstream: false

nats2kafka:
    subject: matchTopic
    jetstream: true
    topic: matchTopic

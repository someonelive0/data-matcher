
# 匹配程序

# 说明

从Nats以队列方式读取消息，然后并行匹配


# 优化

放大本地channel的缓存数量，防止Nats服务端返回"慢消费者"的错误导致的丢包


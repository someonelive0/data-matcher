
# 匹配程序

# 说明

从Nats以队列方式读取消息，然后并行匹配


# 优化

放大本地channel的缓存数量，防止Nats服务端返回"慢消费者"的错误导致的丢包



# Nats 

输入主题： flow.*，不持久化，目的是以最快速度传输完成即可

输出流： match_flow, 输出主题：match_flow.*，输出的匹配结果需要持久化，目的是以后可以更多的使用。

./nats.exe stream add match_flow \
  '--subjects=match_flow.*' \
  --storage=file \
  --replicas=1 \
  --retention=limits \
  --discard=old \
  --max-age=30d \
  --max-bytes=10GiB \
  --max-msg-size=1MiB \
  --max-msgs=-1 \
  --max-msgs-per-subject=-1 \
  --dupe-window=120s \
  --no-allow-rollup \
  --allow-direct \
  --no-deny-delete \
  --no-deny-purge \
  --description='for data-matcher'

# Grafana设置auth.proxy
```
[auth.proxy]
enabled = true
header_name = X-WEBAUTH-USER
header_property = username
auto_sign_up = true
```

# Configure
```
server:
  port: 8080                                     # 代理服务端口

cas:
  url: "http://127.0.0.1/cas"                    # cas地址

grafana:
  url: http://127.0.0.1:3000                     # grafana地址
  username_key: X-WEBAUTH-USER                   # 消息头中使用这个key设置用户名

```
# Usage
```
./main -log_dir=logs
```

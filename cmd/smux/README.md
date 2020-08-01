# Benchmark of HTTP and Smux

证书是我们自签的，没有第三方CA作验证，所以客户端需要关闭校验证书有效性的特性。

1. 生成服务端私钥 `openssl genrsa -out default.key 2048`
2. 生成服务端证书 `openssl req -new -x509 -key default.key -out default.pem -days 3650`

```sh
sh bench.sh DELAY YOUR_CERT_PATH YOUR_KEY_PATH
python plot.py
```

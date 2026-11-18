# 本地开发证书

仓库不再提交任何证书与私钥。开启 TLS 的本地调试前，先在仓库根目录生成一对自签名证书：

```shell
make cert
```

生成的 `config/cert/server.crt` / `config/cert/server.key` 已被 `.gitignore` 忽略，
请勿提交到仓库，也不要复用线上环境的证书。

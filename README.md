# bridge
expose your localhost

# install and run

假设服务器IP为10.0.0.1，内网提供web service的设备是192.168.0.2,端口为8888，如果做如下配置，让10.0.0.1:10001转发到内部网络的192.168.0.2:8888


```bash
#run the command on server
go install github.com/juxuny/bridge/cmd/bridge-server
bridge-server -c=server.json
```

```bash
#run the command on client
go install github.com/juxuny/bridge/cmd/bridge-client
bridge-client -c=client.json
```

server.json

```json
{
  "port": 9090,
  "tokenConf": "token.conf"
}
```

token.conf

```conf
#token  |  aes-key | port
nooF2nXPawgx1LUsKTnrtvdlxY2OZgN3sOuv80kduWpVZOl66uSehv9CGnoL8CH5 uNH8RGSW86rRt85b 10001
```

client.json

```json
{
  "token": "nooF2nXPawgx1LUsKTnrtvdlxY2OZgN3sOuv80kduWpVZOl66uSehv9CGnoL8CH5",
  "key": "uNH8RGSW86rRt85b",
  "host": "10.0.0.1:9090",
  "local": "192.168.0.2:8888"
}
```

现在就可以通过访问 10.0.0.1:10001 来访问内部网络的 192.168.0.2:8888页面。

Note：项目更详细的说明在[这里](https://zhuanlan.zhihu.com/p/67373515)




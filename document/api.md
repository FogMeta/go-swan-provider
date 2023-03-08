# RPC Provider Guideline API

* [版本](#版本)
* [高度](#高度)
* [账户余额](#账户余额)
* [状态信息](#状态信息)
* [设置验证节点](#设置验证节点)
* [查看验证节点](#查看验证节点)
* [抵押](#抵押)


## 版本

描述: 检查运行中 pocket 的当前版本

```shell
curl --url http://127.0.0.1:8088/poktsrv/version 
```

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "version": "RC-0.9.2"
  }
}
```


## 高度

描述: 检查运行中 pocket 的当前区块高度

```shell
curl --url http://127.0.0.1:8088/poktsrv/height
```

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "height": 99156
  }
}
```


## 账户余额

描述:查询指定账号的余额

```shell
curl --request POST --url http://127.0.0.1:8088/poktsrv/balance --header 'Content-Type: application/json' \
--data "{\"height\": 0,\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\"}"
```

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "height": 0,
    "address": "ee60841d9afb70ba893c02965537bc0eec4ef1e4",
    "balance": "39999930000"
  }
}
```


## 状态信息

描述:检查运行中 pocket 的状态信息

```shell
curl --url http://127.0.0.1:8088/poktsrv/status
```

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "version": "RC-0.9.2",
    "height": 99167,
    "address": "ee60841d9afb70ba893c02965537bc0eec4ef1e4",
    "publicKey": "7b1739685dcdc10fcc02bc21dd822ef3458fcf543cc89487af9fe512b573e74d",
    "balance": 39999910000,
    "staking": "20000000000",
    "award": "",
    "jailed": false,
    "jailedBlock": 0,
    "jailedUntil": "0001-01-01T00:00:00Z"
  }
}
```


## 设置验证节点

描述:设置运行中 pocket 的验证节点账户

```shell
curl --request POST --url http://127.0.0.1:8088/poktsrv/set-validator --header 'Content-Type: application/json' \
--data "{\"passwd\": \"123456\",\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\"}"
```

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "result": "spawn sh -c pocket accounts set-validator ee60841d9afb70ba893c02965537bc0eec4ef1e4\r\n\2023/03/03 03:06:37 Initializing Pocket Datadir\r\n2023/03/03 03:06:37 datadir = /home/app/.pocket\r\nEnter the password:\r\n"
  }
}
```


## 查看验证节点

描述:查看当前验证节点账户信息

```shell
curl --url http://127.0.0.1:8088/poktsrv/validator
```

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": "ee60841d9afb70ba893c02965537bc0eec4ef1e4"
}
```


## 抵押

描述:设置节点抵押

```shell
curl --request POST --url http://127.0.0.1:8088/poktsrv/custodial --header 'Content-Type: application/json' \
--data "{\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\",\"amount\": \"20000000000\",\"relay_chain_ids\": \"0001,0021\",\"service_url\": \"http://pokt.storefrontiers.cn:80\",\"network_id\": \"testnet\",\"fee\": \"10000\",\"is_before\": \"false\",\"passwd\": \"123456\"}"
```

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "result": "spawn sh -c pocket nodes stake custodial ee60841d9afb70ba893c02965537bc0eec4ef1e4 20000000000 0001,0021 http://pokt.storefrontiers.cn:80 testnet 10000 false\r\n 2023/03/03 03:15:32 Initializing Pocket Datadir\r\n2023/03/03 03:15:32 datadir = /home/app/.pocket\r\nEnter Passphrase: \r\nhttp://localhost:8081/v1/client/rawtx\r\n{\r\n    \"logs\": null,\r\n    \"txhash\": \"0A025220D33B84525E99AFD5BE7ECA95D6234AFB40CD21901700A7F706DE12E7\"\r\n}\r\n\r\n"
  }
}
```

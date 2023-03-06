# Swan Provider
[![Discord](https://img.shields.io/discord/770382203782692945?label=Discord&logo=Discord)](https://discord.gg/MSXGzVsSYf)
[![Twitter Follow](https://img.shields.io/twitter/follow/0xfilswan)](https://twitter.com/0xfilswan)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg)](https://github.com/RichardLitt/standard-readme)

- 加入FilSwan的[Slack](https://filswan.slack.com)频道，了解新闻、讨论和状态更新。 
- 查看FilSwan的[Medium](https://filswan.medium.com)，获取最新动态和公告。

## 目录

- [特性](#特性)
- [前提条件](#前提条件)
- [安装部署](#安装部署)
- [配置](#配置)
- [命令](#命令)
- [API](#API)
- [许可证](#许可证)

## 特性

Swan Provider Pocket 提供以下功能：

* 在容器中自动部署pocket节点。
* 提供对容器中节点基本设置查询命令。
* 提供对容器中节点维护监控的接口。

## 前提条件
- Docker

### 安装 Docker 
```shell
sudo apt install docker
```

参考: [官方安装文档](https://docs.docker.com/engine/install/)


## 安装部署
### 安装选择:one: **预构建软件包**: 参照 [release assets](https://github.com/filswan/go-swan-provider/releases)
####  构建方法
```shell
mkdir swan-provider
cd swan-provider
wget --no-check-certificate https://raw.githubusercontent.com/filswan/go-swan-provider/release-2.0.0/install.sh
chmod +x ./install.sh
./install.sh
```
### 安装选择:two: 从源代码构建
####  构建指引
```shell
git clone https://github.com/filswan/go-swan-provider.git
cd go-swan-provider
git checkout release-2.0.0
./build_from_source.sh
```

### 配置provider
- 编辑配置文件 **~/.swan/provider/config.toml** 和 **~/.swan/provider/config-pokt.toml**, 参照 [配置](#Config) 部分

### 运行
- 后台运行 `swan-provider`
```
nohup swan-provider pocket start --passwd 123456 >> swan-provider.log 2>&1 & 
```

### 下载快照
- 从最新快照下载将极大地缩短主网同步区块链所需的时间。使用wget进行下载，并在下载后解压缩存档。解压路径 `/root/.pocket` 需要与 `config-pokt.toml` 中 `pokt_data_path` 指定的路径保持一致。
```
wget -qO- https://snapshot.nodes.pokt.network/latest.tar.gz | tar -xz -C /root/.pocket
```

### 配置`chains.json`
- 根据自身需求，配置`config-pokt.toml` 中 `pokt_data_path` 指定的路径下的 `config/chains.json` ，例如：
```
[
    {
      "id": "0001",
      "url": "http://localhost:8081/",
      "basic_auth": {
        "username": "",
        "password": ""
      }
    },
    {
      "id": "0021",
      "url": "https://eth-rpc.gateway.pokt.network/",
      "basic_auth": {
          "username": "",
          "password": ""
      }
    }
]
```

### 充值
- 使用命令或钱包，充值高于最低抵押值的POCK，最低抵押值为15,000 POKT（或15,000,000,000 uPOKT）。
- 如果正在使用测试网络，可以使用[测试网络水龙头](https://faucet.pokt.network)为账户提供资金。

### 设置验证节点
- 通过命令设置验证节点地址
```
pocket accounts set-validator [YOUR_ACCOUNT_ADDRESS]
```

### 重启容器
```
docker restart [CONTAINER_ID]
```

## 配置
### 配置 `config.toml` 文件
- **port:** 默认 `8888`，web api 端口
- **release:** 默认为 `true`, 在release模式下工作时设置为true；否则为false

#### [lotus]
- **client_api_url:** lotus 客户端的web api对应的Url, 比如 `http://[ip]:[port]/rpc/v0`, 通常来说 `[port]` 是 `1234`. 参照 [Lotus API](https://docs.filecoin.io/reference/lotus-api/)
- **market_api_url:** lotus 客户端的web api对应的Url, 比如 `http://[ip]:[port]/rpc/v0`, 通常来说 `[port]` 是 `2345`. 当market和miner没有分离时，这也是miner访问令牌的访问令牌. 参照 [Lotus API](https://docs.filecoin.io/reference/lotus-api/)
- **market_access_token:** lotus market web api的访问令牌. 当market和miner没有分离时，这也是miner访问令牌的访问令牌. 参照 [Obtaining Tokens](https://docs.filecoin.io/build/lotus/api-tokens/#obtaining-tokens)

#### [aria2]
- **aria2_download_dir:** 离线交易文件进行下载以供导入的目录
- **aria2_host:** Aria2 服务器地址
- **aria2_port:** Aria2 服务器端口
- **aria2_secret:** 必须与 `aria2.conf` 的rpc-secret值相同
- **aria2_auto_delete_car_file**: 交易状态变为 Active 或 Error 后, 对应的 CAR 文件会被自动删除， 默认: false
- **aria2_max_downloading_tasks**: Aria2 任务最大下载量, 默认: 10

#### [main]
- **api_url:** Swan API 地址. 对于 Swan production, 地址为 `https://go-swan-server.filswan.com`
- :bangbang:**miner_fid:** Filecoin 矿工ID, 须被添加到 Swan Storage Providers 列表， 添加方式:  [Swan Platform](https://console.filswan.com/#/dashboard) -> "个人信息" -> "作为存储服务商" -> "管理" -> "添加" 。
- **import_interval:** 600秒或10分钟。每笔交易之间的导入间隔。
- **scan_interval:** 600秒或10分钟。在Swan平台上扫描所有正在进行的交易并更新状态的时间间隔。
- **api_key:** api key。可以通过 [Swan Platform](https://console.filswan.com/#/dashboard) -> "个人信息"->"开发人员设置" 获得，也可以访问操作指南。
- **access_token:** 访问令牌。可以通过 [Swan Platform](https://console.filswan.com/#/dashboard) -> "个人信息"->"开发人员设置". 可以访问操作指南查看。
- **api_heartbeat_interval:** 300 秒或 5 分钟. 发送心跳的时间间隔.

#### [bid]
- **bid_mode:** 0: 手动, 1: 自动
- **expected_sealing_time:**  默认: 1920 epoch 或 16 小时. 封装交易的预期时间。过早开始交易将被拒绝。
- **start_epoch:** 默认: 2880 epoch 或 24 小时. 当前 epoch 的相对值。
- **auto_bid_deal_per_day:** 上述配置的矿工每天的自动竞价任务限制。

### 配置 `config-pokt.toml` 文件
#### [pokt]
- **pokt_api_url:** 默认 `8081`，pocket API 端口。
- **pokt_access_token:** 访问令牌.可以通过 [Swan Platform](https://console.filswan.com/#/dashboard) -> "个人信息"->"开发人员设置". 可以访问操作指南查看。
- **pokt_docker_image** Docker 镜像，例如 `filswan/pocket:RC-0.9.2`。
- **pokt_docker_name** 容器名称，可自行定义，例如 `pokt-node-v0.9.2`。
- **pokt_data_path** pocket 数据存储路径，例如 `/root/.pocket`。
- **pokt_scan_interval** 600秒或10分钟。扫描Pocket高度状态的时间间隔。
- **pokt_server_api_url** provider pocket 服务Url，例如 `http://127.0.0.1:8088/`。
- **pokt_server_api_port** provider pocket 服务Port，例如 `8088`。
- **pokt_network_type** pocket网络类型，可以是 MAINNET 和 TESTNET 其中之一。

## 命令
 用 `swan-provider pocket` 命令，与运行中的 pocket 节点进行交互.

### 启动节点
在容器中部署运行pocket节点：
- 拉取 `pokt_docker_image` 指定的image镜像到本地;
- 创建 `pokt_docker_name` 指定的容器，并根据命令参数passwd，创建pocket初始账号;
- 启动 `pokt_docker_name` 指定的容器；
- 等待容器中 pocket node 正常运行，获取pocket版本信息及区块高度。
```
./swan-provider pocket start --passwd "123456"
```

### 版本
检查运行中 pocket 的当前版本
```
./swan-provider pocket version
Pocket Version  : RC-0.9.2
```

### 验证节点
检查运行中 pocket 的当前验证节点账户地址
```
./swan-provider pocket validator
Validator Address       : ee60841d9afb70ba893c02965537bc0eec4ef1e4
```

### 账户余额
检查指定账户的余额
```
./swan-provider pocket balance --addr ee60841d9afb70ba893c02965537bc0eec4ef1e4
Address : ee60841d9afb70ba893c02965537bc0eec4ef1e4
Balance : 39999970000
```

### 状态信息
检查运行中 pocket 节点的状态信息
```
./swan-provider pocket status
Version         : RC-0.9.2
Height          : 99131
Address         : ee60841d9afb70ba893c02965537bc0eec4ef1e4
PublicKey       : 7b1739685dcdc10fcc02bc21dd822ef3458fcf543cc89487af9fe512b573e74d
Balance         : 39999970000
Staking         : 20000000000
Jailed          : false
JailedBlock     : 0
JailedUntil     : 0001-01-01 00:00:00 +0000 UTC
```

### 抵押
设置节点抵押
```
./swan-provider pocket custodial --operatorAddress="ee60841d9afb70ba893c02965537bc0eec4ef1e4" --amount="20000000000" --relayChainIDs="0001,0021" --serviceURI="http://pokt.storefrontiers.cn:80" --networkID="testnet" --fee="10000" --isBefore="false" --passwd="123456"

{Result: spawn sh -c pocket nodes stake custodial ee60841d9afb70ba893c02965537bc0eec4ef1e4 20000000000 0001,0021 http://pokt.storefrontiers.cn:80 testnet 10000 false
v2023/03/02 21:50:02 Initializing Pocket Datadir
2023/03/02 21:50:02 datadir = /home/app/.pocket
Enter Passphrase: 
http://localhost:8081/v1/client/rawtx
{
    "logs": null,
    "txhash": "487F8E6FEFCDB1B8324572B411DC1E4239CEAA915958FB06BA6E6655978ADF43"
}
}
```


## API
用 API 命令，与运行中的 pocket 节点进行交互.

### 版本
检查运行中 pocket 的当前版本
```
curl --url http://127.0.0.1:8088/poktsrv/version 

{
  "status": "success",
  "code": "",
  "data": {
    "version": "RC-0.9.2"
  }
}
```

### 高度
检查运行中 pocket 的当前区块高度
```
curl --url http://127.0.0.1:8088/poktsrv/height

{
  "status": "success",
  "code": "",
  "data": {
    "height": 99156
  }
}
```

### 账户余额
检查运行中 pocket 的当前版本
```
curl --request POST --url http://127.0.0.1:8088/poktsrv/balance --header 'Content-Type: application/json' \
--data "{\"height\": 0,\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\"}"

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

### 状态信息
检查运行中 pocket 的状态信息
```
curl --url http://127.0.0.1:8088/poktsrv/status

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

### 设置验证节点
设置运行中 pocket 的验证节点账户
```
curl --request POST --url http://127.0.0.1:8088/poktsrv/set-validator --header 'Content-Type: application/json' \
--data "{\"passwd\": \"123456\",\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\"}"

{
  "status": "success",
  "code": "",
  "data": {
    "result": "spawn sh -c pocket accounts set-validator ee60841d9afb70ba893c02965537bc0eec4ef1e4\r\n\2023/03/03 03:06:37 Initializing Pocket Datadir\r\n2023/03/03 03:06:37 datadir = /home/app/.pocket\r\nEnter the password:\r\n"
  }
}
```

### 查看验证节点
检查运行中 pocket 的验证节点账户
```
curl --url http://127.0.0.1:8088/poktsrv/validator

{
  "status": "success",
  "code": "",
  "data": "ee60841d9afb70ba893c02965537bc0eec4ef1e4"
}
```

### 抵押
设置节点抵押
```
curl --request POST --url http://127.0.0.1:8088/poktsrv/custodial --header 'Content-Type: application/json' \
--data "{\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\",\"amount\": \"20000000000\",\"relay_chain_ids\": \"0001,0021\",\"service_url\": \"http://pokt.storefrontiers.cn:80\",\"network_id\": \"testnet\",\"fee\": \"10000\",\"is_before\": \"false\",\"passwd\": \"123456\"}"

{
  "status": "success",
  "code": "",
  "data": {
    "result": "spawn sh -c pocket nodes stake custodial ee60841d9afb70ba893c02965537bc0eec4ef1e4 20000000000 0001,0021 http://pokt.storefrontiers.cn:80 testnet 10000 false\r\n 2023/03/03 03:15:32 Initializing Pocket Datadir\r\n2023/03/03 03:15:32 datadir = /home/app/.pocket\r\nEnter Passphrase: \r\nhttp://localhost:8081/v1/client/rawtx\r\n{\r\n    \"logs\": null,\r\n    \"txhash\": \"0A025220D33B84525E99AFD5BE7ECA95D6234AFB40CD21901700A7F706DE12E7\"\r\n}\r\n\r\n"
  }
}
```


## 帮助

如有任何使用问题，请在 [Discord 频道](http://discord.com/invite/KKGhy8ZqzK) 联系 Swan Provider 团队或在Github上创建新的问题.

## 许可证

[Apache](https://github.com/filswan/go-swan-provider/blob/main/LICENSE)

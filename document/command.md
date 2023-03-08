# RPC Provider Guideline 命令

* [启动节点](#启动节点)
* [查看版本](#查看版本)
* [查看验证节点](#查看验证节点)
* [查看余额](#查看余额)
* [状态信息](#状态信息)
* [抵押](#抵押)



## 启动节点

描述: 在容器中部署运行pocket节点：
- 拉取 `pokt_docker_image` 指定的image镜像到本地;
- 创建 `pokt_docker_name` 指定的容器，并根据命令参数passwd，创建pocket初始账号;
- 启动 `pokt_docker_name` 指定的容器；
- 等待容器中 pocket node 正常运行，获取pocket版本信息及区块高度。

```shell
swan-provider pocket start --passwd "123456"
```

输出:

```shell

```


## 查看版本

描述: 检查运行中 pocket 的当前版本

```shell
swan-provider pocket version
```

输出:

```shell
Pocket Version  : RC-0.9.2
```


## 查看验证节点

描述: 检查运行中 pocket 的当前验证节点账户地址

```shell
swan-provider pocket validator
```

输出:

```shell
Validator Address       : ee60841d9afb70ba893c02965537bc0eec4ef1e4
```


## 查看余额

描述: 查看指定账户的余额

```shell
swan-provider pocket balance --addr ee60841d9afb70ba893c02965537bc0eec4ef1e4
```

输出:

```shell
Address : ee60841d9afb70ba893c02965537bc0eec4ef1e4
Balance : 39999970000
```


## 状态信息

描述: 检查运行中 pocket 节点的状态信息

```shell
swan-provider pocket status
```

输出:

```shell
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


## 抵押

描述:设置节点抵押

```shell
swan-provider pocket custodial --operatorAddress="ee60841d9afb70ba893c02965537bc0eec4ef1e4" --amount="20000000000" --relayChainIDs="0001,0021" --serviceURI="http://pokt.storefrontiers.cn:80" --networkID="testnet" --fee="10000" --isBefore="false" --passwd="123456"

```

输出:

```shell
{
    Result: spawn sh -c pocket nodes stake custodial ee60841d9afb70ba893c02965537bc0eec4ef1e4 20000000000 0001,0021 http://pokt.storefrontiers.cn:80 testnet 10000 false
    2023/03/02 21:50:02 Initializing Pocket Datadir
    2023/03/02 21:50:02 datadir = /home/app/.pocket
    Enter Passphrase: 
    http://localhost:8081/v1/client/rawtx
    {
        "logs": null,
        "txhash": "487F8E6FEFCDB1B8324572B411DC1E4239CEAA915958FB06BA6E6655978ADF43"
    }
}
```


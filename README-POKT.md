

# 编译 swan-provider

``` shell
#Build Instructions
git clone https://gitee.com/filswan/go-swan-provider
cd go-swan-provider
git checkout -b rpc-provider origin/rpc-provider
go build -o ./swan-provider main.go

```

# 配置文件
根据实际情况，设置配置文件路径
``` shell
# 配置文件路径由 SWAN_PATH 指定，如未配置此环境变量，默认路径为$HOME/.swan/
cd config
mkdir -p /root/.swan/provider
mv *.toml /root/.swan/provider
cd ..
```

修改config-pokt.toml文件
``` shell
vim /root/.swan/provider/config-pokt.toml

[pokt]
pokt_api_url="http://[ip]:[port]/v1"          # Url of pocket client web api, generally the [port] is 8081
pokt_access_token=""                          # Access token of pocket network api
pokt_docker_image="filswan/pocket:RC-0.9.1.3" # Docker image name and version, this image has push to docker hub
pokt_docker_name="pokt-node-v0.9.1.3"         # Pocket containner name
pokt_config_path="/root/.pocket"              # Pocket chain blocks data path
pokt_scan_interval=600                        # Pocket scan interval
pokt_server_api_url="http://127.0.0.1:8088/"  # Url of pocket server api
pokt_server_api_port=8088                     # Port of pocket server api

```

# 启动pocket 容器
```shell
# 安装docker
sudo wget -qO- https://get.docker.com/ | bash
docker version

# 执行启动命令
# 1. 本机如果没有 pokt_docker_image 指定的image，则从docker hub拉取到本地; 
# 2. 如果 pokt_docker_name 指定的容器不存在，则创建容器，并根据命令参数passwd，创建pocket初始账号;
# 3. 如果 pokt_docker_name 指定的容器状态不是运行状态，则运行此容器；
# 4. 如果 容器中 pocket node 正常运行，则能够获取版本信息及创建的pocket账号信息。
./swan-provider pocket start --passwd "123456"

# 查看pocket 容器是否运行正常
docker ps
CONTAINER ID   IMAGE                       COMMAND                  CREATED          STATUS          PORTS     NAMES
2d4d3e9f991c   filswan/pocket:RC-0.9.1.3   "/usr/bin/expect /ho…"   19 minutes ago   Up 19 minutes             pokt-node-v0.9.1.3

#通过 pocket API查看版本，是否能够正确返回版本
curl --url http://127.0.0.1:8081/v1
"RC-0.9.1.3"

#通过 pocket API查看同步高度，是否不断增高
curl --request POST --url http://127.0.0.1:8081/v1/query/height --header 'Content-Type: application/json'
{"height":857}

```

# Pocket子命令列表

```shell
# 1. 启动pocket节点服务器
swan-provider pocket start

# 2. 获取pocket版本信息
./swan-provider pocket version

Pocket Version is: RC-0.9.1.3 

# 3. 获取当前容器内节点的address信息
./swan-provider pocket nodeaddr

Pocket Node Address is: ffad090789253ad0439c56b7b9c301f90424d5b7 

# 4. 获取指定地址的余额信息
./swan-provider pocket balance --addr=ffad090789253ad0439c56b7b9c301f90424d5b7

Address: ffad090789253ad0439c56b7b9c301f90424d5b7 
Balance: 7454955838 

# 5. 获取节点的当前状态信息
./swan-provider pocket status
Version         : RC-0.9.1.3
Height          : 77438
Address         : ffad090789253ad0439c56b7b9c301f90424d5b7
Balance         : 7454955838
Jailed          : true
JailedBlock     : 36165
JailedUntil     : 2021-10-29 15:11:31.418486526 +0000 UTC


# 6. custodial质押命令
swan-provider pocket custodial
--operatorAddress=0123456789012345678901234567890123456789
--amount=15100000000 
--relayChainIDs=0001,0021
--serviceURI=https://pokt.rocks:443 
--networkID=mainnet 
--fee=10000 
--isBefore=false

# 7. noncustodial质押命令
swan-provider pocket noncustodial
--operatorPublicKey=0123456789012345678901234567890123456789012345678901234567890123 
--outputAddress=0123456789012345678901234567890123456789
--amount=15100000000 
--relayChainIDs=0001,0021
--serviceURI=https://pokt.rocks:443 
--networkID=mainnet 
--fee=10000 
--isBefore=false

```

# Proverer Server API 列表

``` shell
# 1.获取pocket版本信息
http://localhost:8088/api/pocket/v1/version

curl --url http://127.0.0.1:8088/api/pocket/v1/version

{
  "status": "success",
  "code": "",
  "data": {
    "version": "RC-0.9.1.3"
  }
}

# 2. 获取节点的当前高度信息
http://localhost:8088/api/pocket/v1/height

curl --url http://127.0.0.1:8088/api/pocket/v1/height
{
  "status": "success",
  "code": "",
  "data": {
    "height": 77267
  }
}

# 3. 获取指定账号的余额信息
http://localhost:8088/api/pocket/v1/balance

curl --request POST --url http://127.0.0.1:8088/api/pocket/v1/balance \
--header 'Content-Type: application/json' \
--data '{"height": 0,"address":"ffad090789253ad0439c56b7b9c301f90424d5b7"}'

{
  "status": "success",
  "code": "",
  "data": {
    "height": 0,
    "address": "ffad090789253ad0439c56b7b9c301f90424d5b7",
    "balance": "7454955838"
  }
}

#4. 获取节点的当前状态信息
curl --url http://127.0.0.1:8088/api/pocket/v1/status
{
  "status": "success",
  "code": "",
  "data": {
    "version": "RC-0.9.1.3",
    "height": 77267,
    "address": "ffad090789253ad0439c56b7b9c301f90424d5b7",
    "balance": 7454955838,
    "award": "",
    "jailed": true,
    "jailedBlock": 35994,
    "jailedUntil": "2021-10-29T15:11:31.418486526Z"
  }
}


```


# Build Instructions

Refer to [README.md](https://github.com/filswan/go-swan-provider/README.md)

# Configuration
Modify the config-pokt.toml file
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

# Start the pocket container
```shell
# Installing docker
sudo wget -qO- https://get.docker.com/ | bash
docker version

Execute the start command
# 1. If there is no image specified by pokt_docker_image on the local device, it will be pulled from the docker hub to the local device.
# 2. If the container specified by pokt_docker_name does not exist, create a container and create an initial pocket account according to command parameter passwd.
# 3. If the container state specified by pokt_docker_name is not a running state, run the container;
# 4. If pocket node in the container is running properly, the version information and pocket account information can be obtained.
./swan-provider pocket start --passwd "123456"

# Check whether the pocket container is running properly
docker ps
CONTAINER ID   IMAGE                       COMMAND                  CREATED          STATUS          PORTS     NAMES
2d4d3e9f991c   filswan/pocket:RC-0.9.1.3   "/usr/bin/expect /ho…"   19 minutes ago   Up 19 minutes             pokt-node-v0.9.1.3

# Check the version through the pocket API to see if the version can be returned correctly
curl --url http://127.0.0.1:8081/v1
"RC-0.9.1.3"

# Check the synchronization height through pocket API to see whether it keeps increasing
curl --request POST --url http://127.0.0.1:8081/v1/query/height --header 'Content-Type: application/json'
{"height":857}

```

# Pocket subcommand list

```shell
# 1. Start the pocket node server
swan-provider pocket start

# 2. Gets pocket version information
./swan-provider pocket version

Pocket Version is: RC-0.9.1.3 

# 3. Gets the address information of the node in the current container
./swan-provider pocket nodeaddr

Pocket Node Address is: ffad090789253ad0439c56b7b9c301f90424d5b7 

# 4. Gets the balance information for the specified address
./swan-provider pocket balance --addr=ffad090789253ad0439c56b7b9c301f90424d5b7

Address: ffad090789253ad0439c56b7b9c301f90424d5b7 
Balance: 7454955838 

# 5. Gets the current status of the node
./swan-provider pocket status
Version         : RC-0.9.1.3
Height          : 77438
Address         : ffad090789253ad0439c56b7b9c301f90424d5b7
Balance         : 7454955838
Jailed          : true
JailedBlock     : 36165
JailedUntil     : 2021-10-29 15:11:31.418486526 +0000 UTC


# 6. Custodial command for staking
swan-provider pocket custodial
--operatorAddress=0123456789012345678901234567890123456789
--amount=15100000000 
--relayChainIDs=0001,0021
--serviceURI=https://pokt.rocks:443 
--networkID=mainnet 
--fee=10000 
--isBefore=false

# 7. Non-Custodial command for staking
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

# Provider Server API list

``` shell
# 1.Gets pocket version information
http://localhost:8088/api/pocket/v1/version

curl --url http://127.0.0.1:8088/api/pocket/v1/version

{
  "status": "success",
  "code": "",
  "data": {
    "version": "RC-0.9.1.3"
  }
}

# 2. Gets the current height of the node
http://localhost:8088/api/pocket/v1/height

curl --url http://127.0.0.1:8088/api/pocket/v1/height
{
  "status": "success",
  "code": "",
  "data": {
    "height": 77267
  }
}

# 3. Gets the balance information for the specified account
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

#4. Gets the current status of the node
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

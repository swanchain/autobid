# Swan Guide

## Features:

This swan tool listens to the tasks that come from Swan platform. It provides the following functions:

* It selects a suitable auto-bid miner for an auto-bid task. 
* If an auto-bid miner is selected for an auto-bid task, the task status will be set to Assigned, otherwise, ActionRequired.
* Synchronize deal status with Swan platform so that client will know the status changes in realtime.

## Prerequisite
- Database for swan platform.

## Config
* ./config/config.toml.example
```shell
port = "8888"
dev = true
auto_bid_interval_sec = 120 #auto bid interval, unit:second

[database]
db_host = "192.168.88.188"
db_port = "3306"
db_schema_name = "sr2"
db_username = "root"
db_password = ""
db_args = "charset=utf8mb4&parseTime=True&loc=Local"
db_max_idle_conn_num = 10
```
## How to use

### Step 1. Download code
```shell
mkdir go-swan
git clone git@192.168.88.183:NebulaAI-BlockChain/go-swan.git
git checkout dev
```

### Step 2. Compile Provider
```shell
make help    # view how to use make tool
make clean   # remove generated binary file and config file
make dep     # Get dependencies
make test    # Run unit tests
make build   # generate binary file and config file
```

### Step 3. Start Swan
```shell
cd go-swan
vi ./config/config.toml
nohup ./go-swan > go-swan.log &
```

The deal status will be synchronized on the filwan.com, both client and miner will know the status changes in realtime.

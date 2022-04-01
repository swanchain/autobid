# Swan Guide

## Features:

This swan tool listens to the auto-bid tasks that come from Swan platform. It provides the following functions:

* It selects a suitable auto-bid miner for an auto-bid task. 
* If an auto-bid miner is selected for an auto-bid task, the task status will be set to Assigned, otherwise, ActionRequired.
* Synchronize deal status with Swan platform so that client will know the status changes in realtime.

## Prerequisite
- Database for swan platform.

## Config
* ./config/config.toml.example
```shell
port = 8888
auto_bid_interval_sec = 120  #auto bid interval, unit:second

[database]
db_host = ""                 # ip of the host for database instance running on
db_port = 3306               # port of the host for database instance running on
db_schema_name = ""          # database schema name for swan
db_username = "root"         # username to access the database
db_password = ""             # password to access the database
db_args = "charset=utf8mb4&parseTime=True&loc=Local" # other arguments to access database
db_max_idle_conn_num = 10    # maximum number of connections in the idle connection pool

[lotus]
client_api_url = "http://[ip]:[port]/rpc/v0"   # Url of lotus client web api, generally the [port] is 1234
```

## How to use

### Step 1. Download code
```shell
git clone https://github.com/filswan/autobid.git
cd autobid
git checkout <branch-name>
```

### Step 2. Compile Code
```shell
make   # generate binary file and config file to ./build folder
```

### Step 3. Start Swan
```shell
cd build
vi ./config/config.toml
./autobid > autobid.log &
```

#### Note
You can add **nohup** before **./autobid > autobid.log &** to ignore the HUP (hangup) signal and therefore avoid stop when you log out.
```shell
nohup ./autobid > ./autobid.log &
```

The deal status will be synchronized on the filwan.com, both client and miner will know the status changes in realtime.

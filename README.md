# go-swan

# build executable bin file

## for linux

GOOS=linux GOARCH=amd64 go build -v ./

## for mac

env GOOS=darwin GOARCH=amd64 go build -v ./

# put the bin file to destination

# create config folder
in the same directory as the bin file

# put config.toml.example under config folder
# and rename it to config.toml
the source file is:
go-swan-provider/config/config.toml.example

# edit config.toml with right values



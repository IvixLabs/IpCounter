# IpCounter

## Build
make build

## Run
dist/app -path=path_to_ips_file -threads=12

if you want build for specific os or processor you can use similar command:

GOOS=linux GOARCH=amd64 make build
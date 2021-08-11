# UI Server
__This is UI project__


## Requirements:
 - Ubuntu 20.04 (other distros not tested)
   - make=4.2.1-1.2
   - go=1.16.4
   - gcc=4:9.3.0-1 (for "go test")
   - docker=20.10.7
   - docker-compose=1.29.2


## Build
```sh
make build
ls ./output  # check the ui binary at ./output
```

## Unit Test
```sh
make test
```

## Run
```sh
make all  # This would run:
              # UI server (port 9900)
              # API Document server (port 9901)
              # Postgres (port 5432)
```

## Clean
```sh
make clean
```

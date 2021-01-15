# Identity Service

Steps to run project

```
source ./scripts/dev.sh
```

```
go run main.go
```

For windows (wsl)

Start Postgres

```
sudo service postgresql start
```

or

find root ip using:

```
grep nameserver /etc/resolv.conf | awk '{print $2}'
```


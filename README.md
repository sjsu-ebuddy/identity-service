# Identity Service

Steps to run project

```
go run main.go -env=dev
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


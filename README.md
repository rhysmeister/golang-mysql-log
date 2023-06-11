# golang-mysql-log
Simple tool to log terminal output to MYsql/MariaDB

# Docker Tests

Run a MariaDB Docker image for testing

```bash
docker run -p 3306:3306 --rm --detach --name some-mariadb --env MARIADB_USER=go --env MARIADB_PASSWORD=secret --env MARIADB_ROOT_PASSWORD=secret mariadb:latest
```

Create database

```bash
go run main.go -u root -p secret --createDB
```

Test stdin input

```bash
echo "Hello World" | go run main.go -u root -p secret
```

Test log flag input

```bash
go run main.go -u root -p secret -log "Hello World"
```

We should then see something like the following...

`
MariaDB [golang]> select * from log;
+----+-------------+---------------------+
| id | log         | created             |
+----+-------------+---------------------+
|  1 | Hello World | 2023-06-10 16:13:28 |
|  2 | Hello World | 2023-06-10 16:14:16 |
+----+-------------+---------------------+
2 rows in set (0,003 sec)
`
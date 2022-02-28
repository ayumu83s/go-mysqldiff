# Go Mysql Diff
Extract diffs from two MySQL schemas and generate DDL such as `CREATE TABLE`, `ALTER TABLE`.
This tool creates a temporary database. So you need a mysql server.

## Usage
```bash
./go-mysqldiff \
--local-db-host 127.0.0.1 \
--local-db-port 3306 \
--local-db-user root \
--local-db-password  `` \
--src-db-host `your db1 host` \
--src-db-port `your db1 port` \
--src-db-user `your db1 user name` \
--src-db-password `your db1 password` \
--src-db-name `your db1 name` \
--src-file-path `path to sql file` \
--dst-db-host `your db2 host` \
--dst-db-port `your db2 port` \
--dst-db-user `your db2 user name` \
--dst-db-password `your db2 password` \
--dst-db-name `your db2 name` \
--dst-file-path `path to sql file`
```

Can be omitted if `--local-db-*` parameter is the same as the below default value.
```bash
--local-db-host 127.0.0.1
--local-db-port 3306
--local-db-user root
--local-db-password  ""
```

It is convenient to register frequently used parameters in environment variables.
Also, Environment variables can be overwritten with arguments.

```bash
# e.g.
export GO_MYSQL_DIFF_LOCAL_DB_HOST=localhost
export GO_MYSQL_DIFF_LOCAL_DB_PORT=13306
export GO_MYSQL_DIFF_LOCAL_DB_USER=user
export GO_MYSQL_DIFF_LOCAL_DB_PASSWORD=password

export GO_MYSQL_DIFF_SRC_DB_HOST=dev-server.com
export GO_MYSQL_DIFF_SRC_DB_PORT=3306
export GO_MYSQL_DIFF_SRC_DB_USER=dev
export GO_MYSQL_DIFF_SRC_DB_PASSWORD=dev
export GO_MYSQL_DIFF_SRC_DB_NAME=test

export GO_MYSQL_DIFF_DST_DB_HOST=dev-server.com
export GO_MYSQL_DIFF_DST_DB_PORT=3306
export GO_MYSQL_DIFF_DST_DB_USER=dev
export GO_MYSQL_DIFF_DST_DB_PASSWORD=dev
export GO_MYSQL_DIFF_DST_DB_NAME=test2

```

## Synopsis
Take diff between dbname1 and dbname2 
(both of databases on the local MySQL)
```bash
./go-mysqldiff \
--src-db-name `dbname1` \
--dst-db-name `dbname2` 
```

Take diff between dbname1 and dbname2 
(both of databases on remote MySQL)
```bash
./go-mysqldiff \
--src-db-name `dbname1` \
--dst-db-host `your db2 host` \
--dst-db-port `your db2 port` \
--dst-db-user `your db2 user name` \
--dst-db-password `your db2 password` \
--dst-db-name `dbname1` 
```

Take diff between sqlfile and dbname2
(both of databases on remote MySQL)
```bash
./go-mysqldiff \
--src-file-path `<path to file>/schema.sql` \
--dst-db-host `your db2 host` \
--dst-db-port `your db2 port` \
--dst-db-user `your db2 user name` \
--dst-db-password `your db2 password` \
--dst-db-name `dbname2` 
```

## See Also

- [mysqldiff](https://github.com/onishi/mysqldiff)

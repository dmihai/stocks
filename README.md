# stocks

Tool for trading stocks

## Configuring the database

1. Create an empty directory

2. Run the docker container replacing `path/to/local/dir` with the newly created directory:

```sh
docker run --name mysql -v path/to/local/dir:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=mysql -d mysql:8.3.0
```

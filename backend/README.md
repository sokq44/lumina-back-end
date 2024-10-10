# BACKEND

## RUN

```bash
go run .
```

## DATABASE

### Installation

#### LINUX (Ubuntu 22.04 LTS)

```bash
sudo apt update
sudo apt install mysql-server
```

#### WINDOWS

Download mysql installer from [mysql installer](https://dev.mysql.com/downloads/installer/).

### SETUP DATABASE

#### LINUX (Ubuntu 22.04 LTS)

```bash
sudo systectl start mysql
sudo mysql -u root -p # No password for root is set by default (should be changed in the future to whatever you want).
```

#### WINDOWS

By default, MySQL is set to run as a Windows Service, which means it will start automatically when your system boots up.

```bash
mysql -u root -p # no password for root is set by default (should be changed whatever you want)
```

#### DATABASE

```sql
/* Remember to change root's password (good practice for security). */
ALTER USER 'root'@'[server (localhost for now)]' IDENTIFIED WITH mysql_native_password BY '[new_password]';

/* Import the database snapshot. */
source [absolute-path-to-the-db-snapshot.sql-file]

/* Create a new user for managing the articles database. */
CREATE USER '[user]'@'[server (localhost for now)]' IDENTIFIED BY '[password]';

/* Make the new user the 'owner' of the newly created database. */
GRANT ALL PRIVILEGES ON articles.* TO '[user]'@'[server (localhost for now)]';

/* You can now exit the mysql cli. */
FLUSH PRIVILEGES;
EXIT;
```

```bash
mysql -u [new_user] -p # From now on you should log in as the user created in the step above.
```

```
DB_USER="[user]"
DB_PASSWD="[password]"
DB_HOST="[address]"
DB_PORT="[port]"
DB_DBNAME="[database name]"
```

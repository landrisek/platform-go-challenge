#!/bin/bash

set -e

echo "Checking Vault availability..."
while ! curl -s -o /dev/null -w "%{http_code}" $VAULT_ADDR/v1/sys/health >/dev/null 2>&1; do
    echo "Vault is not available yet. Waiting for 30 seconds..."
    sleep 30
done

echo "Checking MySQL availability..."

while ! mysqladmin ping -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" >/dev/null 2>&1; do
    echo "MySQL is not available yet. Waiting for 30 seconds..."
    sleep 30
done

echo "MySQL is now available!"

# Enable and configure the secret engine
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type": "database"}' $VAULT_ADDR/v1/sys/mounts/$VAULT_MOUNT
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"plugin_name\":\"mysql-database-plugin\",\"connection_url\":\"root:$MYSQL_PASSWORD@tcp(mysql:3306)/$MYSQL_DATABASE\"}" $VAULT_ADDR/v1/$VAULT_MOUNT/config/connection
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"plugin_name\":\"mysql-database-plugin\",\"connection_url\":\"root:$MYSQL_PASSWORD@tcp(mysql:3306)/$MYSQL_DATABASE\",\"lease\":\"720h\",\"lease_max\":\"720h\"}" $VAULT_ADDR/v1/$VAULT_MOUNT/config/lease
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"plugin_name\":\"mysql-database-plugin\",\"connection_url\":\"root:$MYSQL_PASSWORD@tcp(mysql:3306)/$MYSQL_DATABASE\",\"allowed_roles\":\"readonly,readwrite,sudo\",\"username\":\"$MYSQL_USER\",\"password\":\"$MYSQL_PASSWORD\"}" $VAULT_ADDR/v1/$VAULT_MOUNT/config/$MYSQL_DATABASE
# readonly role
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{
  \"db_name\": \"$MYSQL_DATABASE\",
  \"sql\": \"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT ON $MYSQL_DATABASE.* TO '{{name}}'@'%';\",
  \"creation_statements\": \"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT ON $MYSQL_DATABASE.* TO '{{name}}'@'%';\"
}" $VAULT_ADDR/v1/$VAULT_MOUNT/roles/readonly
# readwrite role
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{
  \"db_name\": \"$MYSQL_DATABASE\",
  \"sql\": \"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT,INSERT,UPDATE,DELETE ON $MYSQL_DATABASE.* TO '{{name}}'@'%';\",
  \"creation_statements\": \"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT,INSERT,UPDATE,DELETE ON $MYSQL_DATABASE.* TO '{{name}}'@'%';\"
}" $VAULT_ADDR/v1/$VAULT_MOUNT/roles/readwrite
# sudo role
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{
  \"db_name\": \"$MYSQL_DATABASE\",
  \"sql\": \"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT,INSERT,UPDATE,DELETE,CREATE,INDEX,TRIGGER,ALTER,REFERENCES ON $MYSQL_DATABASE.* TO '{{name}}'@'%';\",
  \"creation_statements\": \"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT,INSERT,UPDATE,DELETE,CREATE,INDEX,TRIGGER,ALTER,REFERENCES ON $MYSQL_DATABASE.* TO '{{name}}'@'%';\"
}" $VAULT_ADDR/v1/$VAULT_MOUNT/roles/sudo

echo "Vault provisioning completed successfully."

echo "Checking Redis availability..."

while true; do
    if redis-cli -h redis -p ${REDIS_PORT} ping >/dev/null 2>&1; then
        echo "Redis is available."
        break
    else
        echo "Redis is not available yet. Waiting for 30 seconds..."
        sleep 30
    fi
done

redis-cli -h redis -p ${REDIS_PORT} SET blacklist.master controller
redis-cli -h redis -p ${REDIS_PORT} SET blacklist.slave worker
redis-cli -h redis -p ${REDIS_PORT} SET blacklist.greedy ambitious

echo "Redis provisioning completed successfully."
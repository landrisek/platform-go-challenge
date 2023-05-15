#!/bin/bash

set -e

# Check if Vault is available
echo "Checking Vault availability..."
while ! curl -s -o /dev/null -w "%{http_code}" $VAULT_ADDR/v1/sys/health >/dev/null 2>&1; do
    echo "Vault is not available yet. Waiting for 30 seconds..."
    sleep 30
done

# Enable and configure the secret engine
##curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type":"kv-v2"}' $VAULT_ADDR/v1/sys/mounts/secret
#curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"data\":{\"cryptoKey\":\"$CRYPTO_KEY\"}}" $VAULT_ADDR/v1/secret/data/keys/cryptoKey
#curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"data\":{\"hasherSalt\":\"$HASHER_SALT\"}}" $VAULT_ADDR/v1/secret/data/keys/hasherSalt
#curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"data\":{\"sha1sum\":\"$SHA1SUM\"}}" $VAULT_ADDR/v1/secret/data/keys/sha1sum
##curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"type\":\"kv-v2\"}" $VAULT_ADDR/v1/sys/mounts/database
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type": "database"}' $VAULT_ADDR/v1/sys/mounts/mysql_$MYSQL_DATABASE
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"plugin_name\":\"mysql-database-plugin\",\"connection_url\":\"root:$MYSQL_PASSWORD@tcp(mysql:3306)/$MYSQL_DATABASE\"}" $VAULT_ADDR/v1/mysql_$MYSQL_DATABASE/config/connection
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"plugin_name\":\"mysql-database-plugin\",\"connection_url\":\"root:$MYSQL_PASSWORD@tcp(mysql:3306)/$MYSQL_DATABASE\",\"lease\":\"720h\",\"lease_max\":\"720h\"}" $VAULT_ADDR/v1/mysql_$MYSQL_DATABASE/config/lease
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"plugin_name\":\"mysql-database-plugin\",\"connection_url\":\"root:$MYSQL_PASSWORD@tcp(mysql:3306)/$MYSQL_DATABASE\",\"allowed_roles\":\"readonly,readwrite,sudo\",\"username\":\"$MYSQL_USER\",\"password\":\"$MYSQL_PASSWORD\"}" $VAULT_ADDR/v1/mysql_$MYSQL_DATABASE/config/$MYSQL_DATABASE

curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{
  \"db_name\": \"$MYSQL_DATABASE\",
  \"sql\": \"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT,INSERT,UPDATE,DELETE,CREATE,INDEX,TRIGGER,ALTER,REFERENCES ON $MYSQL_DATABASE.* TO '{{name}}'@'%';\"
}" $VAULT_ADDR/v1/mysql_sandbox/roles/sudo


#curl -s -X GET -H "X-Vault-Token: myroot" http://vault:8200/v1/mysql_sandbox/creds/sudo

echo "Vault provisioning script completed successfully."

#docker exec -it ac39580a5898 vault write mysql_sandbox/roles/readonly db_name=sandbox sql="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT,INSERT,UPDATE,DELETE,CREATE,INDEX,TRIGGER,ALTER,REFERENCES ON sandbox.* TO '{{name}}'@'%';"

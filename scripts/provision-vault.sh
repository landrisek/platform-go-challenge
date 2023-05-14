#!/bin/bash

set -e

# Check if Vault is available
echo "Checking Vault availability..."
while ! curl -s -o /dev/null -w "%{http_code}" $VAULT_ADDR/v1/sys/health >/dev/null 2>&1; do
    echo "Vault is not available yet. Waiting for 30 seconds..."
    sleep 30
done

# Enable and configure the secret engine
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d '{"type":"kv-v2"}' $VAULT_ADDR/v1/sys/mounts/secret
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"data\":{\"cryptoKey\":\"$CRYPTO_KEY\"}}" $VAULT_ADDR/v1/secret/data/keys/cryptoKey
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"data\":{\"hasherSalt\":\"$HASHER_SALT\"}}" $VAULT_ADDR/v1/secret/data/keys/hasherSalt
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"data\":{\"sha1sum\":\"$SHA1SUM\"}}" $VAULT_ADDR/v1/secret/data/keys/sha1sum
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"type\":\"kv-v2\"}" $VAULT_ADDR/v1/sys/mounts/database
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"plugin_name\":\"mysql-database-plugin\",\"connection_url\":\"root:$MYSQL_PASSWORD@tcp(mysql:$MYSQL_PORT)/$MYSQL_DATABASE\",\"allowed_roles\":\"readonly,readwrite,sudo\",\"username\":\"$MYSQL_USER\",\"password\":\"$MYSQL_PASSWORD\"}" $VAULT_ADDR/v1/database/config/$MYSQL_DATABASE
curl -s -X POST -H "X-Vault-Token: $VAULT_TOKEN" -d "{\"username\":\"$MYSQL_USER\",\"password\":\"$MYSQL_PASSWORD\",\"ttl\":\"1h\"}" $VAULT_ADDR/v1/database/creds/sudo

echo "Vault provisioning script completed successfully."

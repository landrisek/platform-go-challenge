#!/bin/bash

set -e

echo "Wait for vault to come up. Wait for 5s."
for i in `seq 1 5`;
do
    printf "."
    sleep 1
done

# grabs cryptoKey
if [ ! -z "$CRYPTO_KEY" ]
then
    echo "Using crypto_key from CRYPTO_KEY env variable"
    cryptoKey=${CRYPTO_KEY}
else
    echo "!!! CRYPTO_KEY env var not found"
    cryptoKey=""
fi

# grabs hasherSalt
if [ ! -z "$HASHER_SALT" ]
then
    echo "Using hasher_salt from HASHER_SALT env variable"
    hasherSalt=${HASHER_SALT}
else
    echo "!!! HASHER_SALT env var not found"
    hasherSalt=""
fi

# MWKs list
#MWK1="43650EC1698D76B87029245FF37D80E6553D57AB19493836A7B9C670E2A35F140DB9ECE2DEF22FBC"
#MWK2="26EF5F5E9541CCB7A4318FF2CD7DB73A7BC14CF81A880C60CC45C3B896DCABCFEC7D136FEB4BAEF1"
#MWK3="577FA33C9326E51855C9687C42FBDABE283A8814251D7DB365670055DE13C1C0B5B9D2A46BB8E14C"
MWK1="131a981c68a15f01c0ce0f851d3a9aa039e2a2bbdbdf75da2656ed50b6e2ecba73b1d9f6ea1c43db"

# SHA1SUMs list
#SHA1SUM="40b2c9bf9c7846fdc9af32caf3be02ec63b18798"
SHA1SUM="854c77f24c186538b4402496dc10ffbbadc9421e"

# run vault provisioning
docker exec -it $(docker ps -qf "name=globalwebindex_vault.*") sh -c "vault secrets disable secret; \
vault secrets enable -version=1 -path=secret kv; \
vault kv put secret/keys/cryptoKey value=$cryptoKey; \
vault kv put secret/keys/hasherSalt value=$hasherSalt; \
vault kv put secret/keys/sha1sum value=$SHA1SUM; \ 
vault write database/config/sandbox \
    plugin_name=\"mysql-database-plugin\" \
    connection_url=\"root:pass@tcp(mysql:3306)/sandbox\" \
    allowed_roles=\"readonly,readwrite,sudo\" \
    username=\"sandbox\" \
    password=\"pass\" "



#vault secrets enable -path=mysql_sandbox mysql; \
#vault kv put mysql_sandbox/config/connection connection_url=\"root:pass@tcp(mysql:3306)/sandbox\"; \
#vault kv put mysql_sandbox/config/lease lease=720h lease_max=720h; \
#vault kv put mysql_sandbox/roles/readonly sql=\"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT ON sandbox.* TO '{{name}}'@'%';\"; \
#vault kv put mysql_sandbox/roles/readwrite sql=\"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT,INSERT,UPDATE,DELETE ON sandbox.* TO '{{name}}'@'%';\"; \
#vault kv put mysql_sandbox/roles/sudo sql=\"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT,INSERT,UPDATE,DELETE,CREATE,INDEX,TRIGGER,ALTER,REFERENCES ON sandbox.* TO '{{name}}'@'%';\" "
# vault.hcl
storage "file" {
  path = "/vault/data"
}

#listener "tcp" {
#  address     = "0.0.0.0:8200"
#  tls_disable = 1
#}

# Enable the kv secrets engine
path "secret/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

# Enable the MySQL secrets engine
path "mysql/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

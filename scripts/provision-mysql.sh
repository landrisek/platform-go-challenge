echo "Mysql provisioning started."
mysql --host="$MYSQL_HOST" --port="$MYSQL_PORT" --user="$MYSQL_USER" --password="$MYSQL_PASSWORD" --database="$MYSQL_DATABASE" --execute="INSERT INTO permissions (\`token\`, \`create\`, \`read\`, \`write\`, \`delete\`) VALUES ('XXX', 1, 0, 0, 0);"
echo "Mysql provisioning completed successfully."

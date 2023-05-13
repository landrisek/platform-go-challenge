#!/bin/bash

set +e

if [ ! -z "$MYSQL_PORT" ]
then
    mysqlPort=${MYSQL_PORT}
else
    mysqlPort="3306"
fi

which mysqladmin

if [ $? -ne 0 ]; then
    echo "Mysqladmin wasn't found, so waiting for 20s"
    for i in `seq 1 20`;
    do
        printf "."
        sleep 1
    done
else
    mysqladmin --port=$mysqlPort status --user=sandbox --password=pass > /dev/null 2>&1
    while [ $? -ne 0 ]
    do
        printf "_"
        mysqladmin --port=$mysqlPort status --user=sandbox --password=pass > /dev/null 2>&1
        sleep 1
    done
    echo "Waiting for 10s more"
    for i in `seq 1 10`;
    do
        printf "."
        sleep 1
    done
    echo ""
fi
#!/bin/bash

IFS=
tpl=$(cat docker-compose.template)

cat docker-compose.base > docker-compose.yml

while IFS=":" read -a line; 
do 
    tple=$(echo $tpl | sed "s/_rm_name/${line[0]}/g" )
    tple=$(echo $tple | sed "s/_rm_topic/${line[1]}/g" )
    tple=$(echo $tple | sed "s/_rm_script/${line[2]}/g" )

    echo $tple >> docker-compose.yml
    echo "" >> docker-compose.yml

done < list.txt
#!/bin/bash
counter=1
total=100
while [ $counter -le $total ]
do
    terraform init > /dev/null
    terraform apply -auto-approve
    ((counter++))
done

echo "Done"
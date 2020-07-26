#!/bin/sh

echo "$TERRAFORM_SOURCE" > ./main.tf
terraform init -no-color > /dev/null
terraform apply -auto-approve -no-color > /dev/null
terraform output contents

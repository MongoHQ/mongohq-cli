#!/bin/sh

source .env

go build mongohq.go
shasum mongohq | awk '{print $1}' > checksum

unamestr=`uname`
if [[ "$unamestr" == 'Linux' ]]; then
   s3cmd -c .s3cmd put mongohq s3://mongohq-cli/linux/
   s3cmd -c .s3cmd put checksum s3://mongohq-cli/linux/
elif [[ "$unamestr" == 'FreeBSD' ]]; then
   s3cmd -c .s3cmd put mongohq s3://mongohq-cli/freebsd/
   s3cmd -c .s3cmd put checksum s3://mongohq-cli/freebsd/
elif [[ "$unamestr" == 'Darwin' ]]; then
   s3cmd -c .s3cmd put mongohq s3://mongohq-cli/macosx/
   s3cmd -c .s3cmd put checksum s3://mongohq-cli/macosx/
fi

s3cmd -c .s3cmd put install.sh s3://mongohq-cli/

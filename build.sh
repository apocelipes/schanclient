#!/bin/bash

#testing
cd config
go test
if [ $? -ne 0 ]
then
	exit 1
fi

cd ../parser
go test
if [ $? -ne 0 ]
then
	exit 1
fi

#build
cd ..
go install schanclient

#!/usr/bin/env bash

cd ~/go/src/github.com/dispatchlabs/disgo
go build

cd ~/go/src/github.com/dispatchlabs/samples/run-nodes-locally/seed
rm -f -r db
rm -f -r config
rm -f disgo.log
cp ~/go/src/github.com/dispatchlabs/disgo/disgo .

cd ~/go/src/github.com/dispatchlabs/samples/run-nodes-locally/delegate-1
rm -f -r db
rm -f -r config
rm -f disgo.log
cp ~/go/src/github.com/dispatchlabs/disgo/disgo .

cd ~/go/src/github.com/dispatchlabs/samples/run-nodes-locally/delegate-2
rm -f -r db
rm -f -r config
rm -f disgo.log
cp ~/go/src/github.com/dispatchlabs/disgo/disgo .

cd ~/go/src/github.com/dispatchlabs/samples/run-nodes-locally/delegate-3
rm -f -r db
rm -f -r config
rm -f disgo.log
cp ~/go/src/github.com/dispatchlabs/disgo/disgo .

cd ~/go/src/github.com/dispatchlabs/samples/run-nodes-locally/delegate-4
rm -f -r db
rm -f -r config
rm -f disgo.log
cp ~/go/src/github.com/dispatchlabs/disgo/disgo .

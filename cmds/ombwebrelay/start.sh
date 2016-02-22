#!/bin/bash
go test github.com/soapboxsys/ombudslib/pubrecdb -run=TestSetupDB
FAKEDBPATH=/home/ubuntu/go/src/github.com/soapboxsys/ombudslib/pubrecdb/test/test.db
DBPATH=/home/ubuntu/.ombudscore/node/data/testnet/pubrecord.db
go build -v . && ./ombwebrelay -pubrecpath=$DBPATH


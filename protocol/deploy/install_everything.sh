#!/bin/bash

# This is an install script for an ahimsa-web node. It assumes debian based system 
# that is configured exactly the way we configure our machines (insecurely). This 
# script installs bitcoin, golang, npm, ahimsad and ahimsa-web. Once the blockchain
# downloads you are good to go. 

BTC_VER=0.9.2.1
GO_VER=1.3.1
NODE_VER=0.10.31

# npm must stay in its ugly $HOME based home
NODE_DIST=node-v$NODE_VER-linux-x64
NPM_BIN_PATH=$HOME/$NODE_DIST/bin/npm

install_bitcoin() {
    cd $HOME
    BTC_DIR=bitcoin-$BTC_VER-linux
    wget https://bitcoin.org/bin/0.9.2.1/$BTC_DIR.tar.gz
    tar -xf $BTC_DIR.tar.gz
    sudo mv $BTC_DIR/bin/64/bitcoind /usr/local/bin/
}

install_golang() {
    cd $HOME
    GO_DIST=go$GO_VER.linux-amd64
    wget https://storage.googleapis.com/golang/$GO_DIST.tar.gz
    tar -xf $GO_DIST.tar.gz
    sudo mv go /usr/local/
}

go_conf() {
    cd $HOME
    GOROOT=/usr/local/go
    GOPATH=$HOME/gocode
    echo "export GOROOT=$GOROOT" >> $HOME/.profile
    echo "export GOPATH=$HOME/gocode" >> $HOME/.profile
    echo "export PATH=$PATH:$GOROOT/bin:$GOPATH/bin" >> $HOME/.profile
    if [ ! $HOME/gocode ]; then
	mkdir $HOME/gocode
    fi
    # for go projects managed under hg
    sudo apt-get install -y mercurial
}

install_ahimsad() {
   go get -u -v github.com/NSkelsey/ahimsad/...
}

install_less() {
   # We need npm for the asset build process
   wget http://nodejs.org/dist/v0.10.31/$NODE_DIST.tar.gz
   tar -xf $NODE_DIST.tar.gz
   sudo $NPM_BIN_PATH install -g less
}

ahimsa_web_deps() {
   install_less
   sudo apt-get -y install python-pip nginx gunicorn git libpython-dev 
}

get_ahimsa_web() {
   cd $HOME
   git clone https://github.com/NSkelsey/ahimsa-web
   cd $HOME/ahimsa-web
   sudo pip install -r requirements.txt    
}

delete_everything() {
   sudo rm -rf /usr/local/go /usr/local/bin/bitcoind $HOME/gocode $HOME/$NODE_DIST $HOME/*.tar.gz ahimsa-web $HOME/bitcoin*

}

ahimsad_deps() {
    install_bitcoin
    install_golang
    go_conf
    install_ahimsad
}

install() {
    sudo apt-get -y update
    # ahimsad install
    ahimsad_deps
    # ahimsa web
    ahimsa_web_deps
    get_ahimsa_web
}


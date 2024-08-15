#!/bin/bash
ROOT_DIR=/workspaces/RLAD

fix_openjdk () {
    echo "fixing openjdk..."
    cp $ROOT_DIR/utils/xllr.openjdk.so /usr/local/metaffi/xllr.openjdk.so
    chmod 777 /usr/local/metaffi/xllr.openjdk.so
}

get_dependencies () {
    echo "getting dependencies..."
    sudo apt-get update && sudo apt-get install -y python3.11-dev
    python3.11 -m pip install beautifulsoup4 requests

    #go get github.com/MetaFFI/lang-plugin-go@v0.1.2
    #go mod tidy
    go get
}

build () {
    echo "building..."
    # -buildvcs=false is used to supress weird git error   
    go build -o $ROOT_DIR/output/large-scale-workshop -buildvcs=false
    mkdir -p $ROOT_DIR/output/logs
}

fix_openjdk
get_dependencies
build

echo "FINISHED BUILD"

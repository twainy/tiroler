#!/bin/sh
export GOPATH=`pwd`
if [ ! -x src/github.com/twainy/tiroler ];then
    mkdir -p src/github.com/twainy
    ln -s `pwd` src/github.com/twainy/tiroler 
fi
go get github.com/twainy/goban
go get github.com/zenazn/goji
go get github.com/PuerkitoBio/goquery
go build

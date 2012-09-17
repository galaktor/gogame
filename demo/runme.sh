#!/bin/bash

# rpath appears to have issues in cgo
# hacky workaround: temp use of LD_LIBRARY_PATH
export LD_LIBRARY_PATH=~/mygo/src/github.com/galaktor/gogre3d/ogrelib

# build c wrapper
#g++ $(go env GOGCCFLAGS) -Wl,-rpath=./ogrelib -L./ogrelib -lOgreMain  -shared -o ./ogrelib/libogrec.so ogrec.cpp -L./ogrelib -lOgreMain

# build example using go wrapper
go run demo.go

unset LD_LIBRARY_PATH
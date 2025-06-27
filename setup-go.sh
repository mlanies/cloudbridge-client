#!/bin/bash

# Setup Go 1.21.13 environment
export PATH=/usr/local/go/bin:$PATH
export PATH=$HOME/go/bin:$PATH

echo "Go environment setup complete:"
echo "Go version: $(go version)"
echo "GOPATH: $(go env GOPATH)"
echo "GOROOT: $(go env GOROOT)" 
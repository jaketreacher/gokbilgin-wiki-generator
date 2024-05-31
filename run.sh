#!/bin/bash
set -e

docker build -q --tag gokbilgin .
docker run --rm -v ./data:/data:rw gokbilgin
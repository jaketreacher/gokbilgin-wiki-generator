#!/bin/bash

set -ex

AUTHOR_DIRECTORY=$1

if [ ! -d "$AUTHOR_DIRECTORY" ]; then
    echo "bad input"
    exit 1
fi

cd "$AUTHOR_DIRECTORY"

rename -v 's/.*(\d{4}\.\d{2}.\d{2}).*/$1/' */

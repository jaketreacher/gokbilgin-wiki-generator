#!/bin/bash

set -ex

FILENAME=letter.yaml
TARGET="$1"

if [ ! -f "$TARGET" ]; then
    echo "Target file does not exist"
    exit 1
fi

WORKDIR=$(dirname "$TARGET")
cd "$WORKDIR"

LETTER_DATE=$(basename "$WORKDIR")

if [ ! -f FILENAME ]; then
    touch $FILENAME
fi

VALUE=$LETTER_DATE yq -i '.date = strenv(VALUE)' $FILENAME
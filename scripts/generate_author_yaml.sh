#!/bin/bash

#--------------------------------------------------------
# Input: data directory
# This will generate an author.yaml file based on the
# name of the author directory.
#--------------------------------------------------------

set -ex

DATA_DIR="$1"

if [ ! -d "$DATA_DIR" ]; then
    echo "Input directory does not exist"
    exit 1
fi

function generate_author_yaml {
    set -ex
    YAML_FILE=author.yaml

    TARGET="$1"
    cd "$TARGET"

    AUTHOR_DIR=$(basename "$TARGET")

    if [ -f "$YAML_FILE" ]; then
        echo "author.yaml already exists"
        exit 1
    fi
    touch $YAML_FILE

    AUTHOR="${AUTHOR_DIR%% (*}"
    VALUE=$AUTHOR yq -i '.name = strenv(VALUE)' $YAML_FILE
}

export -f generate_author_yaml

find "$DATA_DIR" -type d -mindepth 1 -maxdepth 1 -exec bash -c 'generate_author_yaml "$0"' {} \;

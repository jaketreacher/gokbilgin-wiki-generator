#!/bin/bash

#--------------------------------------------------------
# Input: data directory
# Requirements:
#    - Letter directory is in the correct format: YYYY.MM.dd
# This will generate a letter.yaml file that includes
# the names and url for all pdfs in the letter directory.
# 
# If you do not want a file included, don't leave it 
# in the directory!
#--------------------------------------------------------

set -ex

YAML_FILE=letter.yaml
DOWNLOAD_BASE="https://gokbilgin.com/wp-content/uploads/letters"

DATA_DIR="$1"

if [ ! -d "$DATA_DIR" ]; then
    echo "Input directory does not exist"
    exit 1
fi

function generate_letter_yaml {
    local LETTER_DIR="$1"
    cd "$LETTER_DIR"

    if [ ! -f $YAML_FILE ]; then
        touch $YAML_FILE
    fi

    local LETTER_DIR=$(basename "$(pwd)")
    local AUTHOR_DIR=$(basename "$(cd .. && pwd)")

    # TODO: Working theory for multiple letters on the same
    #   date is to have some ID in the date folder. As such, it would
    #   be necessary to remove this ID.
    local LETTER_DATE="$LETTER_DIR"
    VALUE=$LETTER_DATE yq -i '.date = strenv(VALUE)' $YAML_FILE
    yq -i 'del(.downloads)' $YAML_FILE

    for name in *.pdf; do
        local url=$(printf "${DOWNLOAD_BASE}/${AUTHOR_DIR}/${LETTER_DIR}/${name}" | node -p 'encodeURI(require("fs").readFileSync(0))' )
        NAME=$name URL=$url yq -i '.downloads += [{"name": strenv(NAME), "url": strenv(URL)}]' $YAML_FILE
    done
}

for author_dir in "$DATA_DIR"/*; do
    if [ -d "$author_dir" ]; then
        for letter_dir in "$author_dir"/*; do
            if [ -d "$letter_dir" ]; then
                generate_letter_yaml "$letter_dir"
            fi
        done
    fi
done

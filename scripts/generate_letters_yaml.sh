#!/bin/bash

#--------------------------------------------------------
# Input: original.pdf filepath
# Requirements:
#    - Directory is in the correct format: YYYY.MM.dd
# This will generate a letter.yaml file.
#--------------------------------------------------------

set -ex

TARGET="$1"
YAML_FILE=letter.yaml
DOWNLOAD_BASE="https://gokbilgin.com/wp-content/uploads/letters"

if [ ! -f "$TARGET" ]; then
    echo "Target file does not exist"
    exit 1
fi

WORKDIR=$(dirname "$TARGET")
cd "$WORKDIR"

FILES_ORIGINAL=$(basename "$TARGET")
FILES_TRANSLATION=${FILES_ORIGINAL/"original.pdf"/"tercume.pdf"}

LETTER_DIR=$(basename "$(pwd)")
AUTHOR_DIR=$(basename "$(cd .. && pwd)")

URLS_ORIGINAL=$(printf "${DOWNLOAD_BASE}/${AUTHOR_DIR}/${LETTER_DIR}/${FILES_ORIGINAL}" | node -p 'encodeURI(require("fs").readFileSync(0))' )
URLS_TRANSLATION=$(printf "${DOWNLOAD_BASE}/${AUTHOR_DIR}/${LETTER_DIR}/${FILES_TRANSLATION}" | node -p 'encodeURI(require("fs").readFileSync(0))' )

LETTER_DATE=$LETTER_DIR

if [ ! -f $YAML_FILE ]; then
    touch $YAML_FILE
fi

VALUE=$LETTER_DATE yq -i '.date = strenv(VALUE)' $YAML_FILE
yq -i 'del(.downloads)' $YAML_FILE
NAME=$FILES_ORIGINAL URL=$URLS_ORIGINAL yq -i '.downloads += [{"name": strenv(NAME), "url": strenv(URL)}]' $YAML_FILE

if [ -f "$FILES_TRANSLATION" ]; then
    NAME=$FILES_TRANSLATION URL=$URLS_TRANSLATION yq -i '.downloads += [{"name": strenv(NAME), "url": strenv(URL)}]' $YAML_FILE
fi
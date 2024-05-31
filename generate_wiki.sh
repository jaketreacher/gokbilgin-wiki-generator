#!/bin/bash

set -x

COOKIE='Cookie: mw_installer_session=b8e17ba511cfcfcee3ad72f869f7bdda; my_wikiUserName=Admin; my_wiki_session=m9p38hhr9f1jj94l2u8cbusmgkklea4m; my_wikiUserID=1'
TEXT_FILE="translation.tur.txt"
TARGET="$1"

if [ ! -f "$TARGET" ]; then
    echo "Target file does not exist"
    exit 1
fi

WORKDIR=$(dirname "$TARGET")
cd "$WORKDIR"

TOKEN=$(curl 'http://localhost:8080/api.php?action=query&meta=tokens&format=json' \
    -H "${COOKIE}" |
    jq -r '.query.tokens.csrftoken')

TITLE=$(
    FILE_NAME=$(basename "$TARGET")

    TITLE_DATE=$(echo "$FILE_NAME" | grep -oE '[0-9]{4}\.[0-9]{2}\.[0-9]{2}')
    TITLE_AUTHOR=$(echo "$FILE_NAME" | sed "s/ $TITLE_DATE.*//" | awk '{for (i=1; i<=NF; i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')

    FORMATTED_TITLE="${TITLE_AUTHOR}'ten $TITLE_DATE Tarihli Mektup"

    echo $FORMATTED_TITLE
)

if [ -f "$TEXT_FILE" ]; then
    TEXT_CONTENT=$(cat $TEXT_FILE)
else
    TEXT_CONTENT="No translations found"
fi

curl -X POST "http://localhost:8080/api.php?action=edit&format=json" \
    -H "application/x-www-form-urlencoded" \
    -H "${COOKIE}" \
    --data-urlencode "title=${TITLE}" \
    --data-urlencode "text=${TEXT_CONTENT}" \
    --data-urlencode "token=${TOKEN}"

#!/bin/bash

#--------------------------------------------------------
# Input: author directory
# This will check whether the author directory, letter
# directories, and original pdf are all named
# approporiately. 
#--------------------------------------------------------

set -eu

TEMPFILE=$(mktemp)

AUTHOR_DIRECTORY=$1

PATTERN_AUTHOR="^[A-Za-z ]*\([A-Za-z]*\)$"
PATTERN_LETTER="^[0-9]{4}\.[0-9]{2}\.[0-9]{2}$"

if [ ! -d "$AUTHOR_DIRECTORY" ]; then
    echo "bad input"
    exit 1
fi

cd "$AUTHOR_DIRECTORY"

LETTER_MESSAGES=0
for dir in */; do
    DIRNAME=${dir%/}

    FILE_MESSAGES=0
    if [ ! -f "${DIRNAME}"/*original.pdf ]; then
        (( FILE_MESSAGES++ ))
    fi

    if ! [[ $DIRNAME =~ $PATTERN_LETTER ]] || [ "$FILE_MESSAGES" -gt 0 ]; then
        echo $FILE_MESSAGES \- \'$DIRNAME\' >> "$TEMPFILE"
        (( LETTER_MESSAGES++ ))
    fi
done

DIRNAME=$(basename "$AUTHOR_DIRECTORY")
if ! [[ $DIRNAME =~ $PATTERN_AUTHOR ]] || [ "$LETTER_MESSAGES" -gt 0 ]; then
    echo $LETTER_MESSAGES \'$DIRNAME\' >> "$TEMPFILE"
else
    echo PASS \'$DIRNAME\' >> "$TEMPFILE"
fi

tail -r "$TEMPFILE"
rm "$TEMPFILE"
#!/bin/bash

set -ex

PDF_IMAGES_DIR="pdfimages"
TESSERACT_OUTPUT="translation.tur.txt"
TARGET="$1"

if [ ! -f "$TARGET" ]; then
    echo "Target file does not exist"
    exit 1
fi

WORKDIR=$(dirname "$TARGET")
cd "$WORKDIR"

if [ ! -d "$PDF_IMAGES_DIR" ]; then
    mkdir $PDF_IMAGES_DIR
fi

pdfimages -all "$TARGET" $PDF_IMAGES_DIR/out

for i in $PDF_IMAGES_DIR/*; do
    tesseract -l tur $i stdout >> $TESSERACT_OUTPUT
done

awk -v RS= -v ORS='\n\n' '{$1=$1} 1' $TESSERACT_OUTPUT  > tmp.txt && mv tmp.txt $TESSERACT_OUTPUT

rm -rf $PDF_IMAGES_DIR

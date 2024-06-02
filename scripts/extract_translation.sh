#!/bin/bash

#--------------------------------------------------------
# Input: tercume.pdf filepath
# This will extract images from the pdf and use tesseract
# to extract the text from the images. A single txt file
# will be created with the complete translation.
#--------------------------------------------------------

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

FILENAME=$(basename "$TARGET")
if [ -d "$PDF_IMAGES_DIR" ]; then
    rm -rf $PDF_IMAGES_DIR
fi
mkdir $PDF_IMAGES_DIR

pdfimages -all "$FILENAME" $PDF_IMAGES_DIR/out

if [ -f $TESSERACT_OUTPUT ]; then
    rm $TESSERACT_OUTPUT
fi

for i in $PDF_IMAGES_DIR/*; do
    tesseract -l tur $i stdout >> $TESSERACT_OUTPUT
done

awk -v RS= -v ORS='\n\n' '{$1=$1} 1' $TESSERACT_OUTPUT  > tmp.txt && mv tmp.txt $TESSERACT_OUTPUT

rm -rf $PDF_IMAGES_DIR

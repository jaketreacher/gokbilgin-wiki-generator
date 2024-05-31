## Rename

```
rename -vn 's/{Match}/{Replace}/' {Input}
```

Note: `-n` will not execute and provides opportunity to see a preview. Remove this flag when we want to make changes.

### Example

Files in directory:
```
'1954.08.10 SATI EL-HUSRI     '
'1954.08.12 SATI EL-HUSRI  '
```

Command:
```
rename -v 's/ *SATI EL-HUSRI *//' *
```

Output:

```
'1954.08.10 SATI EL-HUSRI     ' renamed to '1954.08.10'
'1954.08.12 SATI EL-HUSRI  ' renamed to '1954.08.12'
```

## Combine Images into PDF

```
magick -density 300 -quality 100 image1.jpg image2.png output.pdf
```

### Install
We need to use the imagemagick tool. Install with:
```
brew install imagemagick
```

### Example

Files in directory:
```
1956.06.16 SATI EL HUSRI  pg 1 jpg.jpg
1956.06.16 SATI EL HUSRI . pg 2 jpg.jpg
1956.06.16.SATI EL HUSRI jpg.jpg
```

Manually rename files we want to convert. This makes it easier to capture when we execute imagemagick. The files now look like this:
```
1956.06.16.SATI EL HUSRI jpg.jpg
original_pg1.jpg
original_pg2.jpg
```

Execute the command:
```
magick -density 300 -quality 100 original*.jpg output.pdf
```

Manually rename `output.pdf`.
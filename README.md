# ced
```
ced #RRGGBBAA #RRGGBBAA ... in.png
replaces colors in png file and writes out.png

(1) single argument: show colors in in.png
$ ced in.png
#000000ff 3090
#ff0000ff 7851
#00000000 14318


(2) two arguments: create black-white from greyscale with threshold
$ ced aa in.png

(3) replace color pairs
# e.g. replace transparent black with white
$ ced '#00000000' '#FFFFFFFF' in.png

(4) catenate images in rows/cols
# append all images in two rows and three columns
$ ced 2 '3#' *.png
```

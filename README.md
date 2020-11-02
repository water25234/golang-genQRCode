# golang-shorturl
- pin code to qr code image

## Prerequisite
- GO v1.14

## Implement
- crate file and write to pin code, just like , 21QQ751672(file name) 751672(pin code)
```
19QQ123456
20QQ234456
21QQ456678
17QQ556730
16QQ445011
16QQ545012 545012
21QQ751672 751672
```
- command line
```
go run main.go -file=pinCodeFile.txt -path=file -fileExt=.png

-file: must be, write file.
-path: Optional, write folder and qr code image inside there.
-fileExt: Optional, default .ping.
```
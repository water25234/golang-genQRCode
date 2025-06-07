# golang-shorturl
- pin code to qr code image

## Prerequisite
- GO v1.14
- go get github.com/skip2/go-qrcode

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
go run main.go

go run main.go -r pinCodeFile.txt -f file -e .png

-r: read file content (default pinCodeFile.txt)
-f: output folder (default file)
-e: file extension (default .png)
```

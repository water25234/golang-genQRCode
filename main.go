package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/skip2/go-qrcode"
)

var (
	filePtr    = flag.String("file", "pinCodeFile.txt", "a string")
	pathPtr    = flag.String("path", "file", "a string")
	fileExtPtr = flag.String("fileExt", ".png", "a string")
)

func main() {

	flag.Parse()

	file, err := os.Open(*filePtr)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	str := string(b)
	fileContentArr := strings.Split(str, "\n")
	fileContentCount := len(fileContentArr)
	var errGenQRCode []string

	os.Mkdir(*pathPtr, os.ModePerm)

	c := make(chan bool, fileContentCount)
	for _, value := range fileContentArr {

		if len(value) == 0 {
			continue
		}

		valueArr := strings.Split(value, " ")
		valueName, valuePinCode, err := pinCodeInfo(valueArr)
		if err != nil {
			continue
		}

		pingCode := *pathPtr + "/" + valueName + *fileExtPtr
		go func(pingCode string) {
			err = qrcode.WriteFile(valuePinCode, qrcode.Medium, 256, pingCode)

			if err != nil {
				errGenQRCode = append(errGenQRCode, pingCode)
				fmt.Println("gen QR Code failure", pingCode)
				c <- true
			}

			size, err := fileSize(pingCode)
			if err != nil {
				errGenQRCode = append(errGenQRCode, pingCode)
				fmt.Println("get file size failure", pingCode)
				c <- true
			}

			fmt.Println(fmt.Sprintf("file: %s, file size: %d", pingCode, size))
			c <- true
		}(pingCode)
		<-c
	}
	if len(errGenQRCode) > 0 {
		fmt.Println("gen qr code list : ", errGenQRCode)
	}
}

func pinCodeInfo(valueArr []string) (valueName string, valuePinCode string, err error) {
	if len(valueArr) == 1 {
		valueName = valueArr[0]
		valuePinCode = valueArr[0]
	} else if len(valueArr) == 2 {
		valueName = valueArr[0]
		valuePinCode = valueArr[1]
	} else {
		fmt.Println("value format is error")
		return "", "", fmt.Errorf("value format is error")
	}
	return valueName, valuePinCode, nil
}

func fileSize(pingCode string) (size int64, err error) {
	fi, err := os.Stat(pingCode)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return fi.Size(), nil
}

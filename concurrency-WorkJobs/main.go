package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/skip2/go-qrcode"
)

var (
	filePtr    = flag.String("file", "pinCodeFile.txt", "a string")
	pathPtr    = flag.String("path", "file", "a string")
	fileExtPtr = flag.String("fileExt", ".png", "a string")
)

func main() {

	flag.Parse()

	fmt.Println("--------------- start work ---------------")

	file, err := os.Open(*filePtr)
	if err != nil {
		log.Fatalf("open %s: %v", *filePtr, err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	str := string(b)
	fileContentArr := strings.Split(str, "\n")
	fileContentCount := len(fileContentArr)
	errGenQRCode := &errLog{}

	os.Mkdir(*pathPtr, os.ModePerm)

	// channel for job
	jobChans := make(chan jobChannel, fileContentCount)

	// start workers
	wg := &sync.WaitGroup{}
	wg.Add(fileContentCount)

	// start workers
	for i := 1; i <= runtime.NumCPU(); i++ {
		go func(i int) {
			for job := range jobChans {
				fmt.Printf("worker - %d: started, working, CPU = %d \n", job.index, i)
				work(job.fileContent, errGenQRCode)
				fmt.Printf("worker - %d: completed !, CPU = %d \n", job.index, i)
				wg.Done()
			}
		}(i)
	}

	// collect job
	for i := 0; i < fileContentCount; i++ {
		name := fmt.Sprintf("job-%d", i)
		duration := time.Duration(rand.Intn(1000)) * time.Millisecond
		fmt.Printf("adding: %s %s\n", name, duration)
		jobChans <- jobChannel{
			index:       i,
			fileContent: fileContentArr[i],
		}
	}

	close(jobChans)

	wg.Wait()

	if len(errGenQRCode.errGenQRCode) > 0 {
		fmt.Println("error gen qr code failure list : ", errGenQRCode.errGenQRCode)
	}

	fmt.Println("--------------- finish work ---------------")
}

type jobChannel struct {
	index       int
	fileContent string
}

type errLog struct {
	errGenQRCode []string
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
		return 0, err
	}
	return fi.Size(), nil
}

func work(fileContent string, errGenQRCode *errLog) {

	if len(fileContent) == 0 {
		return
	}

	valueArr := strings.Split(strings.TrimSpace(fileContent), " ")
	valueName, valuePinCode, err := pinCodeInfo(valueArr)
	if err != nil {
		return
	}

	pingCode := *pathPtr + "/" + valueName + *fileExtPtr

	err = qrcode.WriteFile(valuePinCode, qrcode.Medium, 256, pingCode)

	if err != nil {
		fmt.Println("gen QR Code failure", pingCode)
		errGenQRCode.errGenQRCode = append(errGenQRCode.errGenQRCode, pingCode)
		return
	}

	size, err := fileSize(pingCode)
	if err != nil {
		fmt.Println("get file size failure", pingCode)
		errGenQRCode.errGenQRCode = append(errGenQRCode.errGenQRCode, pingCode)
		return
	}

	fmt.Println(fmt.Sprintf("file: %s, file size: %d", pingCode, size))
	return
}

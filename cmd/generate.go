package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/skip2/go-qrcode"
)

type errLog struct {
	errGenQRCode []string
}

type jobChannel struct {
	index       int
	fileContent string
}

func (fs *flags) generateQRCode() error {
	fmt.Println("--------------- start work ---------------")

	file, err := os.Open(fs.readfile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	errGenQRCode := &errLog{}

	os.Mkdir(fs.folder, os.ModePerm)

	// channel for job, size proportional to CPU
	bufferSize := runtime.NumCPU() * 2
	jobChans := make(chan jobChannel, bufferSize)

	// start workers
	wg := &sync.WaitGroup{}

	for i := 1; i <= runtime.NumCPU(); i++ {
		go func(i int) {
			for job := range jobChans {
				fs.work(job.fileContent, errGenQRCode)
				wg.Done()
			}
		}(i)
	}

	// collect job
	index := 0
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		jobChans <- jobChannel{
			index:       index,
			fileContent: line,
		}
		index++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	close(jobChans)

	wg.Wait()

	if len(errGenQRCode.errGenQRCode) > 0 {
		fmt.Println("error gen qr code failure list : ", errGenQRCode.errGenQRCode)
	}

	fmt.Println("--------------- finish work ---------------")
	return nil
}

func (fs *flags) pinCodeInfo(valueArr []string) (valueName string, valuePinCode string, err error) {
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

func (fs *flags) work(fileContent string, errGenQRCode *errLog) {

	if len(fileContent) == 0 {
		return
	}

	valueArr := strings.Split(strings.TrimSpace(fileContent), " ")
	valueName, valuePinCode, err := fs.pinCodeInfo(valueArr)
	if err != nil {
		return
	}

	pingCode := fs.folder + "/" + valueName + fs.fileExt

	err = qrcode.WriteFile(valuePinCode, qrcode.Medium, 256, pingCode)

	if err != nil {
		fmt.Println("gen QR Code failure", pingCode)
		errGenQRCode.errGenQRCode = append(errGenQRCode.errGenQRCode, pingCode)
		return
	}

	fmt.Println("generated", pingCode)
	return
}

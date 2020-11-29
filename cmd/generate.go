package cmd

import (
	"fmt"
	"io/ioutil"
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
	errGenQRCode := &errLog{}

	os.Mkdir(fs.folder, os.ModePerm)

	// channel for job
	jobChans := make(chan jobChannel, fileContentCount)

	// start workers
	wg := &sync.WaitGroup{}
	wg.Add(fileContentCount)

	// start workers
	for i := 1; i <= runtime.NumCPU(); i++ {
		go func(i int) {
			for job := range jobChans {
				// fmt.Printf("worker - %d: started, working, CPU = %d \n", job.index, i)
				fs.work(job.fileContent, errGenQRCode)
				// fmt.Printf("worker - %d: completed !, CPU = %d \n", job.index, i)
				wg.Done()
			}
		}(i)
	}

	// collect job
	for i := 0; i < fileContentCount; i++ {
		// name := fmt.Sprintf("job-%d", i)
		// duration := time.Duration(rand.Intn(1000)) * time.Millisecond
		// fmt.Printf("adding: %s %s\n", name, duration)
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
	return nil
}

func (fs *flags) fileSize(pingCode string) (size int64, err error) {
	fi, err := os.Stat(pingCode)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return fi.Size(), nil
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

	valueArr := strings.Split(fileContent, " ")
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

	size, err := fs.fileSize(pingCode)
	if err != nil {
		fmt.Println("get file size failure", pingCode)
		errGenQRCode.errGenQRCode = append(errGenQRCode.errGenQRCode, pingCode)
		return
	}

	fmt.Println(fmt.Sprintf("file: %s, file size: %d", pingCode, size))
	return
}

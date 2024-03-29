package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	concurrency int
	resume      bool

	bar *progressbar.ProgressBar
}

func NewDownloader(concurrency int, resume bool) *Downloader {
	return &Downloader{concurrency: concurrency, resume: resume}
}

// 问题一：遇到同名文件不会新建
func (d *Downloader) Download(strURL, filename string) error {
	strURL = "http://192.168.51.206:30805/api/s3download/?time=1711619792789&filename=%E8%8B%B1%E8%AF%AD%E5%9F%BA%E7%A1%80%E9%98%85%E8%AF%BB.mp4&URL=http://www.datrix206.com:7480/aa7135485078d0509863339468732701/2023110810_1f5e142294996c0158a4e5c9ea4796f2_lv0.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=c9ca82da496e4d3aff7eacceb2ddb710%2F20240328%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240328T095632Z&X-Amz-Expires=10800&X-Amz-SignedHeaders=host&response-content-type=application%2Foctet-stream&X-Amz-Signature=3da0b2ceb2c9a6adc9926db029c5886d2abe45727bba6610fac2bb26e780bac6"
	//strURL = "http://192.168.51.206:30805/api/s3download/?time=1711616493322&filename=%E9%98%B3%E5%85%89%E7%94%B5%E5%BD%B1www.ygdy8.com.%E9%95%BF%E6%B4%A5%E6%B9%96%E4%B9%8B%E6%B0%B4%E9%97%A8%E6%A1%A5.2022.BD.1080P.%E5%9B%BD%E8%AF%AD%E4%B8%AD%E5%AD%97_20230811114016.mp4&URL=http://www.datrix206.com:7480/aa7135485078d0509863339468732701/2024031216_btb4pu8jo5hyy_lv0.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=c9ca82da496e4d3aff7eacceb2ddb710%2F20240328%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240328T090133Z&X-Amz-Expires=10800&X-Amz-SignedHeaders=host&response-content-type=application%2Foctet-stream&X-Amz-Signature=3d8c036ebe1d41b9713af5924fa2f1ec34a4c2a97f02f93fc9641b04c1f14c73"
	if filename == "" {
		filename = path.Base(strURL)
	}

	// resp, err := http.Head(strURL)
	// if err != nil {
	// 	return err
	// }

	// if resp.Header.Get("Accept-Ranges") == "bytes" {
	// 	return d.multiDownload(strURL, filename, int(resp.ContentLength))
	// }

	return d.singleDownload(strURL, filename)
}

func (d *Downloader) multiDownload(strURL, filename string, contentLen int) error {
	d.setBar(contentLen)

	partSize := contentLen / d.concurrency

	// 创建部分文件的存放目录
	partDir := d.getPartDir(filename)
	os.Mkdir(partDir, 0777)
	defer os.RemoveAll(partDir)

	var wg sync.WaitGroup
	wg.Add(d.concurrency)

	rangeStart := 0

	for i := 0; i < d.concurrency; i++ {
		go func(i, rangeStart int) {
			defer wg.Done()

			rangeEnd := rangeStart + partSize
			// 最后一部分，总长度不能超过 ContentLength
			if i == d.concurrency-1 {
				rangeEnd = contentLen
			}

			downloaded := 0
			if d.resume {
				partFileName := d.getPartFilename(filename, i)
				content, err := os.ReadFile(partFileName)
				if err == nil {
					downloaded = len(content)
				}
				d.bar.Add(downloaded)
			}

			d.downloadPartial(strURL, filename, rangeStart+downloaded, rangeEnd, i)

		}(i, rangeStart)

		rangeStart += partSize + 1
	}

	wg.Wait()

	d.merge(filename)

	return nil
}

func (d *Downloader) downloadPartial(strURL, filename string, rangeStart, rangeEnd, i int) {
	if rangeStart >= rangeEnd {
		return
	}

	req, err := http.NewRequest("GET", strURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	flags := os.O_CREATE | os.O_WRONLY
	if d.resume {
		flags |= os.O_APPEND
	}

	partFile, err := os.OpenFile(d.getPartFilename(filename, i), flags, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer partFile.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(partFile, d.bar), resp.Body, buf)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}
}

func (d *Downloader) merge(filename string) error {
	destFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for i := 0; i < d.concurrency; i++ {
		partFileName := d.getPartFilename(filename, i)
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}
		io.Copy(destFile, partFile)
		partFile.Close()
		os.Remove(partFileName)
	}

	return nil
}

// getPartDir 部分文件存放的目录
func (d *Downloader) getPartDir(filename string) string {
	return strings.SplitN(filename, ".", 2)[0] // 将 filename 字符串按照.进行分割,最多分成两部分。取第一个部分。
}

// getPartFilename 构造部分文件的名字
func (d *Downloader) getPartFilename(filename string, partNum int) string {
	partDir := d.getPartDir(filename)
	return fmt.Sprintf("%s/%s-%d", partDir, filename, partNum)
}

func (d *Downloader) singleDownload(strURL, filename string) error {
	resp, err := http.Get(strURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d.setBar(int(resp.ContentLength))

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(f, d.bar), resp.Body, buf)
	return err
}

func (d *Downloader) setBar(length int) {
	d.bar = progressbar.NewOptions(
		length,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("downloading..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

var (
	baseURL = "http://192.168.51.206:30805/api/s3download/" // 直接下载链接的头部
	s3URL   = "http://192.168.51.206:30805/api/s3/"         // 下载链接
)

type Down struct {
	concurrency int
	resume      bool

	bar *progressbar.ProgressBar
}

func NewDown(concurrency int, resume bool) *Down {
	return &Down{concurrency: concurrency, resume: resume}
}

func (d *Down) Download(fileID string, single bool) error {
	param := InitDownReq{
		FileID:        fileID,
		Status:        0,
		TranscodeType: 0,
	}
	signatureUrl, name, size, err := InitFileDownloads(param)
	if err != nil {
		return err
	}
	if single {
		fullURL := strAdd(signatureUrl, name)
		return d.singleDownload(fullURL, name)
	}
	return d.multiDownload(signatureUrl, name, int(size))
}

func strAdd(signatureUrl, name string) string {
	params := url.Values{}
	params.Set("time", strconv.Itoa(int(time.Now().Unix())))
	params.Set("filename", name)
	params.Set("URL", signatureUrl)
	queryString := fmt.Sprintf("time=%s&filename=%s&URL=%s", params.Get("time"), params.Get("filename"), params.Get("URL"))
	fullURL := fmt.Sprintf("%s?%s", baseURL, queryString)
	return fullURL
}

func (d *Down) multiDownload(strURL, filename string, contentLen int) error {
	d.setBar(contentLen)
	fmt.Println(contentLen, d.concurrency)
	partSize := contentLen / d.concurrency
	fmt.Println("分片大小", formatBytes(uint64(partSize)))
	// 创建部分文件的存放目录
	partDir := d.getPartDir(filename)
	err := os.Mkdir(partDir, 0777)
	if err != nil {
		return err
	}
	defer os.RemoveAll(partDir)
	var totalDownloadTime time.Duration
	start := time.Now()

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
	end := time.Now()
	downloadTime := end.Sub(start)
	totalDownloadTime += downloadTime
	fmt.Println("下载总时长", totalDownloadTime)
	return nil
}

func (d *Down) downloadPartial(url, filename string, rangeStart, rangeEnd, i int) {
	if rangeStart >= rangeEnd {
		return
	}
	req, err := http.NewRequest(http.MethodGet, s3URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("URL", url)
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

func (d *Down) merge(filename string) error {
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
func (d *Down) getPartDir(filename string) string {
	return strings.SplitN(filename, ".", 2)[0]
}

// getPartFilename 构造部分文件的名字
func (d *Down) getPartFilename(filename string, partNum int) string {
	partDir := d.getPartDir(filename)
	return fmt.Sprintf("%s/%s-%d", partDir, filename, partNum)
}

// go run . -i 2024031216_btb4pu8jo5hyy_lv0.mp4 -s true
// go run . -i 2023110810_1f5e142294996c0158a4e5c9ea4796f2_lv0.mp4 -s true
func (d *Down) singleDownload(strURL, filename string) error {
	var totalDownloadTime time.Duration
	start := time.Now()
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
	end := time.Now()
	downloadTime := end.Sub(start)
	totalDownloadTime += downloadTime
	fmt.Println("下载总时长", totalDownloadTime)
	return err
}

func (d *Down) setBar(length int) {
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

// formatBytes 将字节为单位的数值格式化
func formatBytes(bytes uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	var i int
	for i = 0; bytes >= 1024 && i < len(sizes)-1; i++ {
		bytes /= 1024
	}
	return fmt.Sprintf("%d %s", bytes, sizes[i])
}

package main

import (
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
)

func main() {
	// 默认并发数
	concurrencyN := runtime.NumCPU()

	app := &cli.App{
		Name:  "downloader",
		Usage: "File concurrency downloader",
		Flags: []cli.Flag{ // -h 是自带的，定义时不要用help和h
			&cli.StringFlag{
				Name:    "url",               // 选项名称  --url
				Aliases: []string{"u"},       // 选项别名 -u
				Usage:   "`URL` to download", // 用法
				// Required: true,                // 是否必传
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "output `fileName`", // `` 里面包括的表示选项后的参数
			},
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "input `fileID`",
			},
			&cli.IntFlag{
				Name:    "concurrency",
				Aliases: []string{"n"},
				Value:   concurrencyN, // 默认值
				Usage:   "Concurrency `number`",
			},
			&cli.BoolFlag{
				Name:    "resume",
				Aliases: []string{"r"},
				Value:   true,
				Usage:   "Resume download",
			},
			&cli.BoolFlag{
				Name:    "single",
				Aliases: []string{"s"},
				Value:   false,
				Usage:   "single download", // 整个文件直接下载
			},
		},
		Action: func(c *cli.Context) error {
			strURL := c.String("url")
			fileName := c.String("output")
			fileID := c.String("input")
			concurrency := c.Int("concurrency")
			resume := c.Bool("resume")
			single := c.Bool("single")
			if strURL != "" {
				return NewDownloader(concurrency, resume).Download(strURL, fileName)
			}
			return NewDown(concurrency, resume).Download(fileID, single)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

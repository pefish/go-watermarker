package main

import (
	"flag"
	"fmt"
	"github.com/pefish/go-config"
	"github.com/pefish/go-watermarker/pkg/watermark"
	"github.com/pefish/go-watermarker/version"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const noCoverDirName = version.Name

func main() {
	flagSet := flag.NewFlagSet(version.Name, flag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Printf( "%s 是一款为图片添加文字水印的工具. Enjoy it !!!\n\n", flagSet.Name())
		fmt.Printf("Usage: %s [option] <target file/path>\n", flagSet.Name())
		flagSet.PrintDefaults()
		fmt.Printf("\n")
	}
	configPath := flagSet.String("config", "", "path to config file")
	text := flagSet.String("text", "www.pefish.club", "水印的文本")
	isCover := flagSet.Bool("cover", false, "是否覆盖源文件")

	if configPath != nil && *configPath != "" {
		err := go_config.Config.LoadYamlConfig(go_config.Configuration{
			ConfigFilepath: *configPath,
		})
		if err != nil {
			log.Fatalf("load config file error - %s", err)
		}
		go_config.Config.MergeFlagSet(flagSet)
	}

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	w := watermark.NewWatermark(*text)
	srcFilePath := flagSet.Args()[0]
	fileInfo, err := os.Stat(srcFilePath)
	if err != nil {
		log.Fatalf("获取 %s Stat - %s", srcFilePath, err)
	}
	if fileInfo.IsDir() {
		err = filepath.Walk(srcFilePath, func(path string, info os.FileInfo, err error) error {
			path, err = filepath.Abs(path)
			if err != nil {
				log.Fatalf("获取 %s 绝对地址 - %s", path, err)
			}
			if info.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(srcFilePath, path)
			if err != nil {
				log.Fatalf("获取相对地址 %s %s - %s", srcFilePath, path, err)
			}
			if strings.HasPrefix(rel, ".") {  // 忽略目标目录中的隐藏文件文件夹以及文件
				return nil
			}
			if !strings.HasSuffix(path, ".jpg") && !strings.HasSuffix(path, ".png") {
				return nil
			}
			if strings.Contains(rel, noCoverDirName) {
				return nil
			}

			markSingleFile(w, path, *isCover)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		markSingleFile(w, srcFilePath, *isCover)
	}
}

func markSingleFile(w *watermark.Watermark, srcFilePath string, isCover bool) {
	fmt.Printf("%s ...\n", srcFilePath)
	targetFilePath := srcFilePath
	if !isCover {  // 不覆盖源文件
		dirPath := path.Join(filepath.Dir(srcFilePath), noCoverDirName)
		err := os.Mkdir(dirPath, 0777)
		if err != nil {
			if strings.Contains(err.Error(), "file exists") {
				fileInfo, err := os.Stat(dirPath)
				if err != nil {
					log.Fatalf("获取 %s Stat - %s", dirPath, err)
				}
				if !fileInfo.IsDir() {
					err = os.Remove(dirPath)
					if err != nil {
						log.Fatalf("文件 %s 删除失败 - %s", dirPath, err)
					}
					err = os.Mkdir(dirPath, 0777)
					if err != nil {
						log.Fatalf("文件夹 %s 创建失败 - %s", dirPath, err)
					}
				}
			} else {
				log.Fatalf("文件夹 %s 创建失败 - %s", dirPath, err)
			}
		}
		baseFilePath := filepath.Base(srcFilePath)
		srcAbsPath, err := filepath.Abs(filepath.Dir(srcFilePath))
		if err != nil {
			log.Fatal(err)
		}
		targetFilePath = path.Join(srcAbsPath, version.Name, baseFilePath)
	}
	err := w.MarkFile(srcFilePath, targetFilePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: done!!!\n", srcFilePath)
}

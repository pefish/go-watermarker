# go-watermarker

其他语言版本的README: [English](README.md), [简体中文](README_zh-cn.md)

一款为图片添加文字水印的工具，图片暂时只支持 jpg/png

## 安装

```shell script
go get -u github.com/pefish/go-watermarker/cmd/...
```

## 使用

```shell
go-watermarker --text="this is watermark" --cover ./1.jpg

go-watermarker --text="this is watermark" ./test/
```

## 文档

可以通过 `go-watermarker --help` 查看帮助文档。

```shell script
go-watermarker 是一款为图片添加文字水印的工具. Enjoy it !!!

Usage: go-watermarker [option] <target file/path>
  -config string
    	path to config file
  -cover
    	是否覆盖源文件
  -text string
    	水印的文本 (default "www.pefish.club")
```

如果目标是一个目录，则会遍历目录下的所有图片，都会被盖上水印。

如果指定了 `-cover` ，那么原图片会被直接覆盖，没有任何备份（小心一点）。否则图片所在的目录都会生成一个名字为 `go-watermarker` 的文件夹，盖过水印的图片都会放里面，原图片不动。

## Security Vulnerabilities

If you discover a security vulnerability, please send an e-mail to [pefish@qq.com](mailto:pefish@qq.com). All security vulnerabilities will be promptly addressed.

## License

This project is licensed under the [Apache License](LICENSE).


# go-watermarker

A tool to add text watermark for jpg/png images.

## Install

```shell script
go get -u github.com/pefish/go-watermarker/cmd/...
```

## Quick Start

```shell
go-watermarker --text="this is watermark" --cover ./1.jpg

go-watermarker --text="this is watermark" ./test/
```

## Document

You can use `go-watermarker --help` to find document.

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

If your target is a dictionary, this dictionary will be traversed recursively，all jpg/png images will be watermarked.

If `-cover` is set, source file will be covered directly. Or new dictionary named by `go-watermarker` will be created, all watermarked images are put there.

## Security Vulnerabilities

If you discover a security vulnerability, please send an e-mail to [pefish@qq.com](mailto:pefish@qq.com). All security vulnerabilities will be promptly addressed.

## License

This project is licensed under the [Apache License](LICENSE).

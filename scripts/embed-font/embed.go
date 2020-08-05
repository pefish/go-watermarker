package main

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fs, err := os.Open("./pkg/internal/asset/font.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	newFile, err := os.Create("./pkg/internal/asset/asset.go")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = io.Copy(zw, fs)
	if err != nil {
		log.Fatal(err)
	}
	err = zw.Close()  // 必须先close，从而flush，否则buf中不完整，后面写入的也就不完整，解压就会出错
	if err != nil {
		log.Fatal(err)
	}

	textBytes := []byte(fmt.Sprintf(`
package asset

const (
 FontStr = "%s"
)
`, hex.EncodeToString(buf.Bytes())))

	_, err = newFile.Write(textBytes)
	if err != nil {
		log.Fatal(err)
	}
}

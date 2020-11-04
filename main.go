package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

func download(uri string, f func(length, downLen int64)) error {
	var (
		buf     = make([]byte, 1024)
		fsize   int64
		written int64
	)

	uURL, err := url.ParseRequestURI(uri)
	if err != nil {
		panic("网址错误")
	}
	filename := path.Base(uURL.Path)
	// 1. 获取远程文件数据
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Body == nil {
		return errors.New("body is null")
	}
	// 2. 读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	// 3. 创建本地临时文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, ew := file.Write(buf[0:nr])
			// 数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			// 写入出错
			if ew != nil {
				err = ew
				break
			}
			// 读取的数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		// 没有错误，使用 callback
		f(fsize, written)
	}
	if err != nil {
		panic(err)
	}
	// io.Copy(file, resp.Body)
	return nil
}

func main() {
	download("https://down.sandai.net/thunderx/XunLeiWebSetup10.1.38.890gw.exe", func(length, downLen int64) {
		fmt.Printf("%d - %d\n", length, downLen)
	})
}

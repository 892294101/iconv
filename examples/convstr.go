package main

import (
	"fmt"
	"github.com/892294101/iconv"
)

func main() {
	d := []byte("可能是用的不多可能是用的不多可能是用的不多可能是用的不多可能是用的不多可能是用的不多可能是用的不多可能是用的不多可能是用的不多可能是用的不多可能是用的不多可能是用的不多可能是用的不多可能是用的不多")

	cd, err := iconv.Open("gbk", "utf-8")
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return
	}
	defer cd.Close()
	fmt.Println("go")
	var ind int
	outBuf := make([]byte, len(d))

	res, _, err := cd.Conv(d, outBuf)
	if err != nil {
		ind++
	}
	fmt.Println("res: ", string(res))
	fmt.Println("=========失败： ", ind)
	cd.Close()

}

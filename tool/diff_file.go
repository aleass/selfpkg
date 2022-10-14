package main

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"io"
	"os"
	"sort"
	"strings"
)

func Start() {
	var diffType string
	fmt.Println("Please enter your diff type: \r\n" +
		"1:   data from clipboard\n" +
		"any: data from file")

	fmt.Scanln(&diffType)
	var isSort string
	fmt.Println("press any key without sorting or press yes sorting")
	fmt.Scanln(&isSort)

	var info = make([][]string, 2)
	var r1, r2 *bufio.Reader
	if diffType == "1" {
		fmt.Println("get data from clipboard:Press any key after copy data")
		fmt.Scanln()
		d1, err := clipboard.ReadAll()
		if err != nil {
			panic(err.Error())
		}
		r1 = bufio.NewReader(strings.NewReader(d1))
		clipboard.WriteAll("")
		fmt.Println("get data2 from clipboard:Press any key after copy data2")
		fmt.Scanln()
		d1, err = clipboard.ReadAll()
		if err != nil {
			panic(err.Error())
		}
		r2 = bufio.NewReader(strings.NewReader(d1))
		clipboard.WriteAll("")
	} else {
		var f1, f2 string
		fmt.Println("Please input file1 file2")
		fmt.Scanln(&f1, &f2)
		r1, r2 = GetDataFromFile(f1, f2)
	}
	ExtractData(info, r1, r2)

	if len(info[0]) != len(info[1]) {
		fmt.Println("数量不一样")
	}

	//排序
	if isSort == "yes" {
		sort.Strings(info[0])
		sort.Strings(info[1])
	}
	var i, j int
	for ; len(info[0]) > i && len(info[1]) > i; i, j = i+1, j+1 {
		if info[0][i] != info[1][i] {
			fmt.Println((info[0][i]) + " != " + (info[1][i]))
			return
		}
	}
}

func GetDataFromFile(f1, f2 string) (*bufio.Reader, *bufio.Reader) {
	var f, err = os.Open(f1)
	if err != nil {
		panic(err.Error())
	}
	f, err = os.Open(f2)
	if err != nil {
		panic(err.Error())
	}
	return bufio.NewReader(f), bufio.NewReader(f)
}

func ExtractData(info [][]string, reader, reader2 *bufio.Reader) {
	str, _, err := reader.ReadLine()
	for err != io.EOF {
		info[0] = append(info[0], string(str))
		str, _, err = reader.ReadLine()
	}

	str, _, err = reader2.ReadLine()
	for err != io.EOF {
		info[1] = append(info[1], string(str))
		str, _, err = reader2.ReadLine()
	}

}

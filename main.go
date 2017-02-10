package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var decimalOption int = 0

func main() {
	filename := ""
	if len(os.Args) >= 2 {
		if strings.Compare(os.Args[1], "-d") == 0 {
			filename = os.Args[2]
			decimalOption = 1
		} else if strings.Compare(os.Args[1], "-3d") == 0 {
			filename = os.Args[2]
			decimalOption = 2
		} else if strings.Compare(os.Args[1], "-03d") == 0 {
			filename = os.Args[2]
			decimalOption = 3
		} else {
			if os.Args[1][0] == '-' {
				filename = os.Args[2]
			} else {
				filename = os.Args[1]
			}
		}
	} else {
		fmt.Println("[Usage]")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "filename")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-d", "filename", "\t\t: Decimal value.")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-3d", "filename", "\t\t: 3-digit Decimal value with blank padding.")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-03d", "filename", "\t\t: 3-digit Decimal value with 0 padding.")
		//filename = "test.txt"
		return
	}

	PrintHex(filename)
}

func PrintHex(filename string, countPerLineParam ...int) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	countPerLine := 30
	if len(countPerLineParam) >= 1 {
		countPerLine = countPerLineParam[0]
	}

	//fmt.Print(dat)
	//fmt.Println("\n------------------------------------------")

	for i, v := range dat {
		switch decimalOption {
		case 1:
			if v == '\r' || v == '\n' {
				fmt.Printf("%d.", v)
			} else {
				fmt.Printf("%d ", v)
			}
			break
		case 2:
			if v == '\r' || v == '\n' {
				fmt.Printf("%3d.", v)
			} else {
				fmt.Printf("%3d ", v)
			}
			break
		case 3:
			if v == '\r' || v == '\n' {
				fmt.Printf("%03d.", v)
			} else {
				fmt.Printf("%03d ", v)
			}
			break
		default:
			if v == '\r' || v == '\n' {
				fmt.Printf("%02X.", v)
			} else {
				fmt.Printf("%02X ", v)
			}
			break
		}

		if i%countPerLine == (countPerLine - 1) {
			fmt.Println("")
		}
	}

	fmt.Println("")
	fmt.Println("Size :", len(dat))
}

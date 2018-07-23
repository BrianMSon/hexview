package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var decimalOption int = 0
var lineOption int = 0
var verboseOption int = 0
var ucs2Option int = 0
var noSpaceOption int = 0
var outHexOption int = 0
var inHexOption int = 0
var fromOffset int = 0
var toOffset int = -1

func main() {
	filename := ""
	outfilename := ""
	columnCount := 30

	if len(os.Args) >= 2 {
		//fmt.Println("len : ", len(os.Args))
		filename = os.Args[len(os.Args)-1]
		for iArg := 1; iArg < len(os.Args)-1; iArg++ {
			if strings.Compare(os.Args[iArg], "-d") == 0 {
				decimalOption = 1
			} else if strings.Compare(os.Args[iArg], "-d3") == 0 {
				decimalOption = 2
			} else if strings.Compare(os.Args[iArg], "-d03") == 0 {
				decimalOption = 3
			} else if strings.Compare(os.Args[iArg], "-l") == 0 {
				lineOption = 1
			} else if strings.Compare(os.Args[iArg], "-v") == 0 {
				verboseOption = 1
			} else if strings.Compare(os.Args[iArg], "-u") == 0 {
				ucs2Option = 1
				decimalOption = 0
			} else if strings.Compare(os.Args[iArg], "-n") == 0 {
				noSpaceOption = 1
			} else if strings.Contains(os.Args[iArg], "-c") == true {
				strNum := strings.TrimLeft(os.Args[iArg], "-c")
				columnCount, _ = strconv.Atoi(strNum)
				if columnCount <= 0 {
					columnCount = 0
				}
			} else if strings.Compare(os.Args[iArg], "-o") == 0 {
				outHexOption = 1
				outfilename = os.Args[iArg+1]
			} else if strings.Compare(os.Args[iArg], "-i") == 0 {
				inHexOption = 1
				filename = os.Args[iArg+1]
			} else if strings.Contains(os.Args[iArg], "-f") == true {
				strNum := strings.TrimLeft(os.Args[iArg], "-f")
				fromOffset, _ = strconv.Atoi(strNum)
			} else if strings.Contains(os.Args[iArg], "-t") == true {
				strNum := strings.TrimLeft(os.Args[iArg], "-t")
				toOffset, _ = strconv.Atoi(strNum)
			} else {
				if outHexOption == 0 && inHexOption == 0 {
					fmt.Println("Warning : Invalid Option :", os.Args[iArg])
				}
			}
		}
	} else {
		fmt.Println("[HexView] : Copyright by Brian SMG.")
		fmt.Println("")
		fmt.Println("[Usage] :", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "FILENAME")
		fmt.Println("        :", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-l -v FILENAME")
		fmt.Println("[Option]")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-d  ", "FILENAME", "\t\t: show as Decimal value.")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-d3 ", "FILENAME", "\t\t: 3-digit Decimal value with blank padding.")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-d03", "FILENAME", "\t\t: 3-digit Decimal value with 0 padding.")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-l  ", "FILENAME", "\t\t: new Line at LF(0A).")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-v  ", "FILENAME", "\t\t: Verbose output(size, LF dot).")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-u  ", "FILENAME", "\t\t: \\u mark UCS-2(UTF-16).")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-n  ", "FILENAME", "\t\t: No space.")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-c30", "FILENAME", "\t\t: set Column count(-c0 : No new line).")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-o  ", "OUTFILE 4A 5F...", "\t: HEX string to Output file.")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-i  ", "INFILE  ", "\t\t: Input HEX string file to output file.")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-f30", "FILENAME", "\t\t: show From N-th(zero-based) offset.")
		fmt.Println("  ", os.Args[0][strings.LastIndexAny(os.Args[0], "/\\")+1:], "-t90", "FILENAME", "\t\t: show To N-th(zero-based) offset.")
		//filename = "test.txt" // test file
		return
	}

	if outHexOption == 1 {
		hexString := ""
		if inHexOption == 1 {
			//fmt.Println("Load from HEX string file :", filename)
			hexString = GetStringFromInputFile(filename)
		} else {
			if len(os.Args) < 4 {
				fmt.Println("Fail : There is no HEX string.")
				return
			}
			for i := 3; i < len(os.Args); i++ {
				hexString += os.Args[i]
			}
		}

		if len(hexString)%2 == 1 {
			fmt.Println("Fail : Invalid length of hex string.")
			return
		}
		if strings.ContainsAny(hexString, "GHIJKLMNOPQRSTUVWXYZghijklmnopqrstuvwxyz`~!@#$%^&*()-_=+|\\[]{};:'/?,.<>") == true {
			fmt.Println("Fail : Invalid hex string.")
			return
		}
		SaveHexStringToFile(outfilename, hexString)
		return
	} else if inHexOption == 1 {
		fmt.Println("Fail : There is no output file.")
		return
	}

	PrintHex(filename, columnCount)
}

func PrintHex(filename string, countPerLineParam ...int) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error : Cannot read file.")
		fmt.Println(err)
		return
	}

	countPerLine := 30
	if len(countPerLineParam) >= 1 {
		countPerLine = countPerLineParam[0]
	}

	//fmt.Print(data)
	//fmt.Println("\n------------------------------------------")

	lineMark := " "
	if verboseOption == 1 {
		lineMark = "."
	}

	spaceMark := " "
	if noSpaceOption == 1 {
		spaceMark = ""
		lineMark = ""
	}

	ucs2Mark := ""
	if ucs2Option == 1 {
		ucs2Mark = "\\u"
		spaceMark = ""
		lineMark = ""
	}

	printedCount := 0
	for i, v := range data {
		if fromOffset > 0 {
			if i < fromOffset {
				continue
			}
		}
		if toOffset > -1 {
			if i > toOffset {
				continue
			}
		}
		if ucs2Option == 1 {
			if i%2 == 0 {
				fmt.Printf(ucs2Mark)
			}
		}
		switch decimalOption {
		case 1:
			if v == '\r' || v == '\n' {
				fmt.Printf("%d"+lineMark, v)
			} else {
				fmt.Printf("%d"+spaceMark, v)
			}
			break
		case 2:
			if v == '\r' || v == '\n' {
				fmt.Printf("%3d"+lineMark, v)
			} else {
				fmt.Printf("%3d"+spaceMark, v)
			}
			break
		case 3:
			if v == '\r' || v == '\n' {
				fmt.Printf("%03d"+lineMark, v)
			} else {
				fmt.Printf("%03d"+spaceMark, v)
			}
			break
		default:
			if v == '\r' || v == '\n' {
				fmt.Printf("%02X"+lineMark, v)
			} else {
				fmt.Printf("%02X"+spaceMark, v)
			}
			break
		}

		//////////////////////////////////////////////////////

		if lineOption == 0 {
			if ucs2Option == 0 {
				if countPerLine > 0 && i%countPerLine == (countPerLine-1) {
					fmt.Println("")
				}
			}
		} else {
			if v == '\n' {
				fmt.Println("")
			}
		}

		printedCount++
	}

	if ucs2Option == 1 {
		fmt.Println("")
	}

	if verboseOption == 1 {
		fmt.Println("")
		fmt.Println("Size :", len(data))
		if fromOffset > 0 || toOffset > -1 {
			fmt.Println("Printed count :", printedCount)
		}
	}
}

func GetStringFromInputFile(infilename string) string {
	hexString := ""

	fmt.Println("infilename :", infilename)

	data, err := ioutil.ReadFile(infilename)
	if err != nil {
		fmt.Println("Error : Cannot read file.")
		fmt.Println(err)
		return hexString
	}

	for _, nibble := range data {
		var value byte = '0'
		if nibble >= '0' && nibble <= '9' {
			value = '0' + (nibble - '0')
		} else if nibble >= 'A' && nibble <= 'F' {
			value = 'A' + (nibble - 'A')
		} else if nibble >= 'a' && nibble <= 'f' {
			value = 'A' + (nibble - 'a')
		} else {
			continue
		}
		//fmt.Println(nibble, value, string(value))
		hexString += string(value)
	}

	//fmt.Println(hexString)
	return hexString
}

func SaveHexStringToFile(outfilename, hexString string) {
	if _, err := os.Stat(outfilename); err == nil {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(outfilename + " is exist. Overwrite it? (y/N) ")
		if scanner.Scan() {
			input := scanner.Text()
			//fmt.Println("input :", input)
			if !(strings.Compare(input, "Y") == 0 || strings.Compare(input, "y") == 0) {
				fmt.Println("Fail : cancel.")
				return
			}
		}
	}

	fmt.Println("out file :", outfilename)
	//fmt.Println("hex string :", hexString)

	length := len(hexString) / 2
	data := make([]byte, length)
	//fmt.Println(len(hexString), length)
	for i := 0; i < length; i++ {
		orgi := i * 2
		//fmt.Println(orgi, hexString[orgi], hexString[orgi+1])
		//value := (hexString[orgi]-'0')*16 + (hexString[orgi+1] - '0')
		value := CovertToValueFromHexNibble(hexString[orgi]) * 16
		value += CovertToValueFromHexNibble(hexString[orgi+1])
		data[i] = value
	}
	err := ioutil.WriteFile(outfilename, data, 0666)
	if err != nil {
		panic(err)
	}
}

func CovertToValueFromHexNibble(nibble byte) byte {
	var value byte = 0
	if nibble >= '0' && nibble <= '9' {
		value = nibble - '0'
	} else if nibble >= 'A' && nibble <= 'F' {
		value = nibble - 'A' + 10
	} else if nibble >= 'a' && nibble <= 'f' {
		value = nibble - 'a' + 10
	}
	return value
}

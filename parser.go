package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Person struct {
	uid                  string
	nationalityPrimary   string
	nationalitySecondary string
	ethnicValue          int
}

func ParseRtf(path string) []Person {
	result := []Person{}

	rtfFile, rtfErr := os.Open(path)
	if rtfErr != nil {
		fmt.Println(rtfErr)
		os.Exit(0)
	}
	defer rtfFile.Close()

	UIDRegex := regexp.MustCompile("([0-9]){7,}")

	rtfScanner := bufio.NewScanner(rtfFile)
	for rtfScanner.Scan() {
		rtfLine := rtfScanner.Text()
		uid := UIDRegex.Find([]byte(rtfLine))
		if uid != nil {
			rtfData := strings.Split(rtfLine, "|")

			ethnicValue, ethniceValueErr := strconv.Atoi(strings.Trim(rtfData[7], " "))
			if ethniceValueErr != nil {
				fmt.Println(ethniceValueErr)
				continue
			}

			result = append(result, Person{
				uid:                  string(uid),
				nationalityPrimary:   strings.Trim(rtfData[2], " "),
				nationalitySecondary: strings.Trim(rtfData[3], " "),
				ethnicValue:          ethnicValue,
			})
		}
	}

	if rtfScannerErr := rtfScanner.Err(); rtfScannerErr != nil {
		fmt.Println(rtfScannerErr)
	}

	return result
}

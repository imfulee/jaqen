package main

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ParseRtf(path string) ([]Person, error) {
	result := []Person{}

	rtfFile, rtfErr := os.Open(path)
	if rtfErr != nil {
		return nil, rtfErr
	}
	defer rtfFile.Close()

	UIDRegex := regexp.MustCompile("([0-9]){7,}")

	rtfScanner := bufio.NewScanner(rtfFile)
	for rtfScanner.Scan() {
		rtfLine := rtfScanner.Text()
		uid := UIDRegex.Find([]byte(rtfLine))
		if uid != nil {
			rtfData := strings.Split(rtfLine, "|")
			if len(rtfData) < 8 {
				return nil, errors.New("rtf format wrong")
			}

			ethnicValue, ethniceValueErr := strconv.Atoi(strings.Trim(rtfData[7], " "))
			if ethniceValueErr != nil {
				return nil, ethniceValueErr
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
		return nil, rtfScannerErr
	}

	return result, nil
}
}

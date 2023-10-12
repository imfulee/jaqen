package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ErrBadRTFFormat error = errors.New("bad RTF Format")

type INationEthnicMapper interface {
	Map(string) (string, bool)
}

type RTF struct {
	nationEthnicMapper INationEthnicMapper
}

func (r RTF) Parse(path string) ([]Person, error) {
	result := []Person{}

	rtfFile, rtfErr := os.Open(path)
	if rtfErr != nil {
		return nil, rtfErr
	}
	defer rtfFile.Close()

	UIDRegex, err := regexp.Compile("([0-9]){7,}")
	if err != nil {
		return nil, err
	}

	rtfScanner := bufio.NewScanner(rtfFile)
	for rtfScanner.Scan() {
		rtfLine := rtfScanner.Text()
		uidByte := UIDRegex.Find([]byte(rtfLine))
		if uidByte != nil {
			uid := string(uidByte)

			rtfData := strings.Split(rtfLine, "|")
			if len(rtfData) < 8 {
				return nil, errors.Join(ErrBadRTFFormat, errors.New("not enough lines in RTF line"))
			}

			ethnicValue, ethniceValueErr := strconv.Atoi(strings.Trim(rtfData[7], " "))
			if ethniceValueErr != nil {
				return nil, ethniceValueErr
			}

			if ethnicValue > 10 {
				return nil, errors.Join(ErrBadRTFFormat, errors.New("ethnic value out of bounds"))
			}

			ethnicPrimary, hasEthnic := r.nationEthnicMapper.Map(strings.Trim(rtfData[2], " "))
			if !hasEthnic {
				return nil, errors.Join(ErrBadRTFFormat, fmt.Errorf("does not have ethnic on id %s", uid))
			}

			ethnicSecondary, _ := r.nationEthnicMapper.Map(strings.Trim(rtfData[3], " "))

			result = append(result, Person{
				uid:             uid,
				ethnicPrimary:   ethnicPrimary,
				ethnicSecondary: ethnicSecondary,
				ethnicValue:     ethnicValue,
			})
		}
	}

	if rtfScannerErr := rtfScanner.Err(); rtfScannerErr != nil {
		return nil, rtfScannerErr
	}

	return result, nil
}

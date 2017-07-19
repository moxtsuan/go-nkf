package nkf

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var errDet = errors.New("Couldn't detect")

func charDet(b []byte) (string, error) {
	d := chardet.NewTextDetector()
	res, err := d.DetectBest(b)
	if err != nil {
		return "", err
	}
	switch res.Charset {
	case "Shift_JIS":
		return "ShiftJIS", nil
	case "UTF-8":
		return "UTF8", nil
	case "EUC-JP":
		return "EUCJP", nil
	case "ISO-2022-JP":
		return "ISO2022JP", nil
	default:
		return res.Charset, errDet
	}
}

func toUtf8(str string, in string) (string, error) {
	var u8 []byte
	var err error
	switch in {
	case "ISO2022JP":
		u8, err = ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ISO2022JP.NewDecoder()))
	case "EUCJP":
		u8, err = ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.EUCJP.NewDecoder()))
	case "ShiftJIS":
		u8, err = ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	case "UTF8":
		u8, err = []byte(str), nil
	default:
		u8, err = ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	}
	if err != nil {
		return "", err
	}
	return string(u8), err
}

func toSjis(str string) (string, error) {
	sjis, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return "", err
	}
	return string(sjis), err
}

func toEuc(str string) (string, error) {
	euc, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.EUCJP.NewEncoder()))
	if err != nil {
		return "", err
	}
	return string(euc), err
}

func toJis(str string) (string, error) {
	jis, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ISO2022JP.NewEncoder()))
	if err != nil {
		return "", err
	}
	return string(jis), err
}

func nlRep(str string, nl string) string {
	rep := regexp.MustCompile(`\r\n|\r|\n`)
	switch nl {
	case "UNIX":
		return rep.ReplaceAllString(str, "\n")
	case "WINDOWS":
		return rep.ReplaceAllString(str, "\r\n")
	case "MACINTOSH":
		return rep.ReplaceAllString(str, "\r")
	}
	return str
}

/*
Only detect Character encoding
*/
func Guess(file *os.File) (string, error) {
	input, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	det, err := charDet(input)
	return det, err
}

/*
Convert Character encoding
*/
func Convert(file *os.File, in string, out string, nl string) (string, error) {
	input, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	if in == "" {
		in, err = charDet(input)
		if err != nil {
			return "", err
		}
	}
	u8, err := toUtf8(string(input), in)
	if err != nil {
		return "", err
	}
	if nl != "" {
		u8 = nlRep(u8, nl)
	}
	var output string
	switch out {
	case "ISO2022JP":
		output, err = toJis(u8)
	case "ShiftJIS":
		output, err = toSjis(u8)
	case "EUCJP":
		output, err = toEuc(u8)
	case "UTF8":
		output = u8
	}
	return output, err
}

package routes

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func getFileVersion(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return "", err
	}
	sha1 := crypto.SHA1.New()

	_, err = io.Copy(sha1, file)
	file.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}
	retVal := base64.URLEncoding.EncodeToString(sha1.Sum(nil)) + path.Ext(filePath)

	destinationFilePath := path.Join(path.Dir(filePath), retVal)
	if filePath != destinationFilePath {
		copyFile(destinationFilePath, filePath)
	}
	return retVal, nil
}
func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}
func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
func subStr(value string, separator string) string {
	p := strings.Index(value, separator)
	if p > -1 {
		return value[0:p]
	}
	return value
}

func subStrRight(value string, separator string) string {
	p := strings.LastIndex(value, separator)
	if p > -1 {
		return value[0:p]
	}
	return value
}
func urlencode(s string) (result string) {
	for _, c := range s {
		if c <= 0x7f { // single byte
			result += fmt.Sprintf("%%%X", c)
		} else if c > 0x1fffff { // quaternary byte
			result += fmt.Sprintf("%%%X%%%X%%%X%%%X",
				0xf0+((c&0x1c0000)>>18),
				0x80+((c&0x3f000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else if c > 0x7ff { // triple byte
			result += fmt.Sprintf("%%%X%%%X%%%X",
				0xe0+((c&0xf000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else { // double byte
			result += fmt.Sprintf("%%%X%%%X",
				0xc0+((c&0x7c0)>>6),
				0x80+(c&0x3f),
			)
		}
	}

	return result
}

func lastAfterSlash(byteArray []byte) ([]byte, error) {

	lastIndex := bytes.LastIndex(byteArray, []byte("/"))

	if lastIndex == -1 {
		return nil, errors.New("could not find slash")
	}
	return byteArray[lastIndex+1:], nil
}

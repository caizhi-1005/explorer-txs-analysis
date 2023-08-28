package utils

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func TokenIDConvert(s string) (string, error) {
	s = s[2:]
	s = strings.TrimLeft(s, "0")
	if len(s) == 0 {
		s = "0"
	}

	i, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(i)), nil
}

func PathFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//int64转换string
func Int64String(data int64) string {
	var int_string string
	int_string = strconv.FormatInt(data, 10)
	return int_string
}

//string转换int64
func StringInt64(data string) int64 {
	var string_int64 int64
	if data != "" {
		string_int64, _ = strconv.ParseInt(data, 10, 64)
	} else {
		string_int64 = 0
	}
	return string_int64
}

//float64转换string
func Float64String(data float64) string {
	var float64_strig string
	float64_strig = strconv.FormatFloat(data, 'f', -1, 64)
	return float64_strig
}

func ConvertTokenID(tokenID string) (string, error) {
	b, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		return "", errors.New("parse token id error")
	}
	return common.BigToHash(b).String(), nil
}

package utils

import (
	"io/ioutil"
	"net/http"
	"unsafe"
)

// GetPublicIP Get Public Ip
func GetPublicIP() (string, error) {
	url := "https://api.ipify.org?format=text"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return ByteToString(ip), nil
}

// ByteToString byte to string
func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

package request

import (
"net/http"
"strconv"
)

func GetUrlSize(url string) int64 {
	response, error := http.Head(url)

	if error == nil && response.StatusCode == http.StatusOK {
		length, _ := strconv.Atoi(response.Header.Get("Content-Length"))
        size := int64(length)
        return size
	}

	return 0
}
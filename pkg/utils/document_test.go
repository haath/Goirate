package utils

import "testing"

func TestHTTPGet(t *testing.T) {
	_, err := HTTPGet("http://httpbin.org/ip")

	if err != nil {
		t.Errorf("HTTPGet status code: %s", err.Error())
	}
}

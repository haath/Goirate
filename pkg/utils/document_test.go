package utils

import "testing"

func TestHTTPGet(t *testing.T) {
	_, err := HTTPGet("http://httpbin.org/ip")

	if err != nil {
		t.Errorf("HTTPGet status code: %s", err.Error())
	}
}

func TestHTTPGetError(t *testing.T) {
	_, err := HTTPGet("1.2.3.4")

	if err == nil {
		t.Errorf("HTTPGet status code: %s", err.Error())
	}
}

func TestHTTPGetStatus(t *testing.T) {

	res, err := HTTPGet("https://www.reddit.com/asdf")

	if err == nil {
		t.Errorf("Expected 404, got: %v", res)
	}
}

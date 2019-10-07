package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestRecognize(t *testing.T) {

	buffer := new(bytes.Buffer)
	file := "device-window.wav"

	result, err := recognize(buffer, file)

	if err != nil {
		t.Fatalf("recognize gave error: %s", err)
	}

	expectedResult := []string{"device", "window"}

	if !reflect.DeepEqual(expectedResult, result) {
		t.Errorf("expected %s but got %s", expectedResult, result)
	}
}

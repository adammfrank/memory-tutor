package main

import (
	"bytes"
	"fmt"
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

func TestImportVocab(t *testing.T) {
	file := "vocab-test-file.csv"

	memories, err := importVocab(file)

	if err != nil {
		t.Fatalf("importVocab gave error: %s", err)
	}

	if len(memories) != 39 {
		expectationError("39", fmt.Sprintf("#%v", len(memories)), t)
	}

	firstMemory := memory{
		"kñom",
		"I/me",
		"Mi nombre es",
	}

	middleMemory := memory{
		"nia-le-ga",
		"Clock/watch",
		"Flavor Flav wearing a clock near his neck (ga)",
	}

	lastMemory := memory{
		"baang-uik",
		"window",
		"Stop banging on my window, Ulrik!",
	}

	if !reflect.DeepEqual(memories[0], firstMemory) {
		expectationError(fmt.Sprintf("#%v", firstMemory), fmt.Sprintf("#%v", memories[0]), t)
	}
	if !reflect.DeepEqual(memories[20], middleMemory) {
		expectationError(fmt.Sprintf("#%v", middleMemory), fmt.Sprintf("#%v", memories[20]), t)
	}
	if !reflect.DeepEqual(memories[38], lastMemory) {
		expectationError(fmt.Sprintf("#%v", lastMemory), fmt.Sprintf("#%v", memories[38]), t)
	}
}

func TestFindDevice(t *testing.T) {
	file := "vocab-test-file.csv"

	memories, _ := importVocab(file)

	for _, memory := range memories {
		if actualDevice := findDevice(memories, memory.EnglishWord); actualDevice != memory.Device {
			expectationError(memory.Device, actualDevice, t)
		}
	}

}

func expectationError(expected string, actual string, t *testing.T) {
	t.Errorf("expected %s but got %s", expected, actual)
}

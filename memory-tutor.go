package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func main() {

	memories, err := importVocab("vocab.csv")

	words, err := recognize(os.Stdout, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	if words[0] == "device" {
		device := findDevice(memories, words[1])
		cmd := exec.Command("say", device)
		log.Println(device)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

}

func recognize(w io.Writer, file string) ([]string, error) {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})

	var words []string
	// Print the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			wordSlice := strings.Split(alt.Transcript, " ")
			for _, word := range wordSlice {
				words = append(words, word)
			}
		}
	}

	return words, nil
}

type memory struct {
	KhmerWord   string
	EnglishWord string
	Device      string
}

func importVocab(file string) ([]memory, error) {
	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	// skip first line (column names)
	reader.Read()
	var memories []memory
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		memories = append(memories, memory{
			KhmerWord:   line[0],
			EnglishWord: line[1],
			Device:      line[2],
		})
	}

	return memories, nil
}

func findDevice(memories []memory, englishWord string) string {
	for _, memory := range memories {
		if memory.EnglishWord == englishWord {
			return memory.Device
		}
	}

	return "No device found for " + englishWord
}

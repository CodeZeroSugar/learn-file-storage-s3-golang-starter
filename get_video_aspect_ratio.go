package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type AspectRatio struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ProbeResults struct {
	Streams []AspectRatio `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	newBuff := bytes.NewBuffer([]byte{})
	cmd.Stdout = newBuff

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cmd failed: %w", err)
	}

	fmt.Println("ffprobe output:", newBuff.String())

	var aspect ProbeResults

	if err := json.Unmarshal(newBuff.Bytes(), &aspect); err != nil {
		return "", fmt.Errorf("failed to unmarshal json while getting video aspect ratio: %w", err)
	}
	fmt.Printf("Unmarshaled aspect: %+v\n", aspect)
	fmt.Printf("First stream: %+v\n", aspect.Streams[0])

	width := aspect.Streams[0].Width
	height := aspect.Streams[0].Height

	if width == 16*height/9 {
		return "16:9", nil
	} else if height == 16*width/9 {
		return "9:16", nil
	}

	return "other", nil
}

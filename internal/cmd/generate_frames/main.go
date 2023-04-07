package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"github.com/gabe565/ascii-telnet-go/config"
	"github.com/gabe565/ascii-telnet-go/internal/frame"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//go:embed all.go.tmpl
var allTmpl string

//go:embed frame.go.tmpl
var frameTmpl string

func main() {
	// Remove existing frames
	if err := filepath.Walk(config.OutputDir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" && filepath.Base(path) != "stub.go" {
			return os.Remove(path)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	var frameCap int
	var frames []frame.Frame
	var f frame.Frame
	scanner := bufio.NewScanner(bytes.NewReader(config.Movie))

	// Build part of every frame, excluding progress bar and bottom padding
	for lineNum := 0; scanner.Scan(); lineNum += 1 {
		frameLineNum := lineNum % config.FrameHeight
		if frameLineNum == 0 {
			f = frame.Frame{
				Num:  lineNum / config.FrameHeight,
				Data: strings.Repeat("\n", config.PadTop-1),
			}

			v, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}

			f.Duration = time.Duration(float64(v)*(1000.0/15.0)) * time.Millisecond
		} else {
			f.Data += "\n" + strings.Repeat(" ", config.PadLeft) + scanner.Text()
		}

		if frameLineNum == config.FrameHeight-1 {
			frames = append(frames, f)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Compute the total duration
	var totalDuration time.Duration
	for _, f := range frames {
		totalDuration += f.Duration
	}

	// Build the rest of every frame and write to disk
	var currentPosition time.Duration
	for _, f := range frames {
		f.Data += strings.Repeat("\n", config.PadBottom)
		f.Data += strings.Repeat(" ", config.PadLeft-1)
		f.Data += progressBar(currentPosition+f.Duration/2, totalDuration, config.Width)
		f.Data += strings.Repeat(" ", config.PadLeft-1)
		f.Data += strings.Repeat("\n", config.PadBottom)
		f.Height = strings.Count(f.Data, "\n")
		if frameCap < len(f.Data) {
			frameCap = len(f.Data)
		}
		if err := writeFrame(f); err != nil {
			log.Fatal(err)
		}
		currentPosition += f.Duration
	}

	// Write frame list
	if err := writeFrameList(frames, frameCap); err != nil {
		log.Fatal(err)
	}
}

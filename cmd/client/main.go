package main

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/peterh/liner"
)

var (
	historyFilename = filepath.Join(os.TempDir(), "samerkv_history")
	socketFilename  = filepath.Join(os.TempDir(), "samerkv.sock")
)

func main() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(historyFilename); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	// Establish connection.
	client, err := rpc.DialHTTP("unix", socketFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	for {
		input, err := line.Prompt("> ")
		if err != nil {
			if err == liner.ErrPromptAborted || err == io.EOF {
				break
			}
			log.Printf("error reading prompt: %v", err)
			continue
		}
		line.AppendHistory(input)

		var output string
		if err = client.Call("Server.Execute", &input, &output); err != nil {
			log.Printf("write error: %v", err)
			return
		}
		if len(output) > 0 {
			fmt.Println(output)
		}
	}

	if f, err := os.Create(historyFilename); err != nil {
		log.Printf("Error writing history file: %v", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

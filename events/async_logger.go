package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

//AsyncLogger write json logs to file system in different goroutine
type AsyncLogger struct {
	writer io.WriteCloser
	logCh  chan Fact
}

//Consume event fact and put it to channel
func (al *AsyncLogger) Consume(fact Fact) {
	al.logCh <- fact
}

//Close underlying log file writer
func (al *AsyncLogger) Close() (resultErr error) {
	if err := al.writer.Close(); err != nil {
		return fmt.Errorf("Error closing writer: %v", err)
	}
	return nil
}

//Create AsyncLogger and run goroutine that's read from channel and write to file
func NewAsyncLogger(writer io.WriteCloser) Consumer {
	logger := &AsyncLogger{writer: writer, logCh: make(chan Fact, 20000)}

	go func() {
		for {
			fact := <-logger.logCh
			bts, err := json.Marshal(fact)
			if err != nil {
				log.Printf("Error marshaling event to json: %v", err)
				continue
			}

			buf := bytes.NewBuffer(bts)
			buf.Write([]byte("\n"))

			if _, err := logger.writer.Write(buf.Bytes()); err != nil {
				log.Printf("Error writing event to log file: %v", err)
				continue
			}
		}
	}()

	return logger
}
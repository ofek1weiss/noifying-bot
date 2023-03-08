package bot

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

type MessageDataLoader struct {
	Path           string
	SampleInterval time.Duration
}

func (dl *MessageDataLoader) getOldestFile() (string, error) {
	entries, err := os.ReadDir(dl.Path)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", errors.New("no files")
	}
	sort.Slice(entries, func(i, j int) bool {
		iInfo, _ := entries[i].Info()
		jInfo, _ := entries[j].Info()
		return iInfo.ModTime().Before(jInfo.ModTime())
	})
	return filepath.Join(dl.Path, entries[0].Name()), nil
}

func (dl *MessageDataLoader) loadData() (*MessageData, error) {
	filePath, err := dl.getOldestFile()
	if err != nil {
		return nil, err
	}
	data, err := LoadFile[MessageData](filePath)
	if err != nil {
		return nil, err
	}
	if err := os.Remove(filePath); err != nil {
		return nil, err
	}
	return data, nil
}

func (dl *MessageDataLoader) Listen() chan *MessageData {
	messageDatas := make(chan *MessageData, 1)
	go func() {
		for {
			messageData, err := dl.loadData()
			if err != nil {
				if err.Error() != "no files" {
					logrus.Error("Failed to load message data:", err)
				}
				time.Sleep(dl.SampleInterval)
			} else {
				messageDatas <- messageData
			}
		}
	}()
	return messageDatas
}

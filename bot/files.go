package bot

import (
	"encoding/json"
	"os"
)

func LoadFile[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ret := new(T)
	if err := json.Unmarshal(data, ret); err != nil {
		return ret, err
	}
	return ret, nil
}

func DumpFile(value any, path string) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, os.ModeDevice)
}

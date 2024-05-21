package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type TorrodleConfig struct {
	DataDir      string `json:"DataDir"`
	ResultsLimit int    `json:"ResultsLimit"`
	TorrentPort  int    `json:"TorrentPort"`
	HostPort     int    `json:"HostPort"`
	Debug        bool   `json:"Debug"`
}

func (t TorrodleConfig) String() string {
	return fmt.Sprintf(
		`TorrentDir: %v | ResultsLimit: %d | TorrentPort: %d | HostPort: %d | Debug: %v`,
		t.DataDir, t.ResultsLimit, t.TorrentPort, t.HostPort, t.Debug,
	)
}

func InitConfig(path string) error {
	config := TorrodleConfig{
		DataDir:      "",
		ResultsLimit: 100,
		TorrentPort:  9999,
		HostPort:     8080,
	}
	data, _ := json.MarshalIndent(config, "", "\t")
	err := os.WriteFile(path, data, 0o644)
	return err
}

func LoadConfig(path string) (TorrodleConfig, error) {
	var config TorrodleConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

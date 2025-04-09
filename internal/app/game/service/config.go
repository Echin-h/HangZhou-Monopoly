package service

import (
	"encoding/json"
	"os"
)

type StationData struct {
	Stations []StationConfig `json:"model.modelG.Stations"`
}

type StationConfig struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Line           int    `json:"line"`
	Type           int    `json:"type"`
	Price          int    `json:"price"`
	CurrentStation struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"CurrentStation"`
}

func loadStationConfig(filePath string) StationData {
	data, err := os.ReadFile("config/model.modelG.Stations.json")
	if err != nil {
		panic(err)
	}

	var stationData StationData
	if err = json.Unmarshal(data, &stationData); err != nil {
		panic(err)
	}

	return stationData
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var assetsJSON []struct {
	Area string `json:"area"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}

func main() {
	var assets *os.File
	var assetsErr error
	for assets, assetsErr = os.Open("region.json"); assetsErr != nil; {
		assets, assetsErr = os.Open("region.json")
	}
	defer assets.Close()
	fileBytes, err := ioutil.ReadAll(assets)
	if err != nil {
		return
	}

	err = json.Unmarshal(fileBytes, &assetsJSON)
	if err != nil {
		return
	}
	fmt.Println(assetsJSON)
}

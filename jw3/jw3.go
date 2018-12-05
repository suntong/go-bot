package jw3

import (
	"os"
)

var (
	assets, assetsErr = os.Open("region.json")
)

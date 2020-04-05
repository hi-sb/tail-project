package api

import (
	"path/filepath"
	"strconv"
	"tail-project/core"
	"time"
)

// source path  to sys path
func GetFilePath(suffix string) string {
	id := core.GetID()
	year := strconv.Itoa(time.Now().Year())
	month := time.Now().Format("01")
	day := strconv.Itoa(time.Now().Day())
	return filepath.Join(DataPath, year+month+day, id[0:1], id[2:3], id[4:5], id[6:7], id+"."+suffix)
}

package cache

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sctestcase/counter"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

func findCacheFiles(dir string) (files []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, rerr error) (err error) {

		if rerr != nil {
			return rerr
		}

		if !info.IsDir() && filepath.Ext(path) == ".cache" {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func createCacheMap(file *os.File) (cm map[string]cache.Item, err error) {
	scanner := bufio.NewScanner(file)

	cm = make(map[string]cache.Item)
	for scanner.Scan() {
		jsonEl := &counter.JSONElement{}

		err = json.Unmarshal(scanner.Bytes(), jsonEl)
		if err != nil {
			return nil, err
		}

		item := cache.Item{
			Object:     jsonEl.URL,
			Expiration: jsonEl.Expire,
		}

		cm[uuid.NewString()] = item
	}

	if err = scanner.Err(); err != nil {
		fmt.Printf("[ERROR] scan file: %v", err)
	}

	return cm, nil
}

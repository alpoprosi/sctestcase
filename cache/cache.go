package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sctestcase/counter"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type CacheCounter struct {
	cache    *cache.Cache
	cacheDir string
}

var _ counter.Counter = (*CacheCounter)(nil)

func NewCounter(dir string) (cc *CacheCounter, err error) {
	cc = &CacheCounter{
		cacheDir: dir,
	}

	if err = cc.start(); err != nil {
		return nil, fmt.Errorf("start counter: %w", err)
	}

	return cc, nil
}

func (cc *CacheCounter) Inc(url string) {
	uuid := uuid.NewString()

	cc.cache.SetDefault(uuid, url)
}

func (cc *CacheCounter) Count() (count int) {
	return cc.cache.ItemCount()
}

func (cc *CacheCounter) Shutdown() (err error) {
	fmt.Println("[INFO] Saving cache...")
	elements := cc.cache.Items()
	if len(elements) == 0 {
		fmt.Println("[INFO] Nothing to save...")

		return nil
	}

	var toJson int

	fileName := fmt.Sprintf("%s.cache", uuid.NewString())
	file, err := os.Create(filepath.Join(cc.cacheDir, fileName))
	if err != nil {
		return err
	}

	defer file.Close()

	for _, cacheEl := range elements {
		url, ok := cacheEl.Object.(string)
		if !ok {
			fmt.Println("bad interface: ", cacheEl.Object)

			continue
		}

		el := &counter.JSONElement{
			URL:    url,
			Expire: cacheEl.Expiration,
		}

		json, err := json.Marshal(el)
		if err != nil {
			return err
		}

		json = append(json, '\n')

		_, err = file.Write(json)
		if err != nil {
			return err
		}

		toJson++
	}

	fmt.Printf("[INFO] Saving %d cached elements \n", toJson)

	return err
}

func (cc *CacheCounter) start() (err error) {
	files, err := findCacheFiles(cc.cacheDir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		cc.cache = cache.New(time.Minute, time.Minute)

		return nil
	}

	if len(files) > 1 {
		return fmt.Errorf("could be only one .cache file")
	}

	cacheFile := files[0]
	var file *os.File
	file, err = os.OpenFile(cacheFile, os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}

	defer file.Close()

	cacheMap, err := createCacheMap(file)
	if err != nil {
		return err
	}

	cc.cache = cache.NewFrom(time.Minute, time.Second, cacheMap)

	ext := path.Ext(cacheFile)
	extPos := len(cacheFile) - len(ext)
	os.Rename(cacheFile, cacheFile[:extPos]+".used")

	return nil
}

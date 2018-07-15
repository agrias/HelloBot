package db

import (
	"github.com/peterbourgon/diskv"
	"strings"
	"github.com/sirupsen/logrus"
)


type diskvCache struct {
	cache *diskv.Diskv
}

func (c *diskvCache) Get(key string) ([]byte, error) {
	data, err := c.cache.Read(key)
	if err != nil {
		logrus.Error("Could not find data for key. ", key)
		return nil, err
	}

	return data, nil
}

func (c *diskvCache) Put(key string, value []byte) (error) {
	err := c.cache.Write(key, value)
	if err != nil {
		logrus.Error("Problem writing key and value to cache. ", key, value)
		return err
	}

	return nil
}

func (c *diskvCache) Delete(key string) (error) {
	err := c.Delete(key)
	if err != nil {
		logrus.Error("Problem deleting key. ", key)
		return err
	}

	return nil
}

func (c *diskvCache) Keys() ([]string) {
	var keyset []string

	for data := range c.cache.Keys(nil) {
		keyset = append(keyset, data)
	}

	return keyset
}

func NewDiskvCache(path string) Database {
	d := diskv.New(diskv.Options{
		BasePath:     path,
		AdvancedTransform: AdvancedTransformExample,
		InverseTransform:  InverseTransformExample,
		CacheSizeMax: 1024 * 1024,
	})

	return &diskvCache{d}
}

func AdvancedTransformExample(key string) *diskv.PathKey {
	path := strings.Split(key, "/")
	last := len(path) - 1
	return &diskv.PathKey{
		Path:     path[:last],
		FileName: path[last] + ".data",
	}
}

func InverseTransformExample(pathKey *diskv.PathKey) (key string) {
	txt := pathKey.FileName[len(pathKey.FileName)-4:]
	if txt != ".data" {
		panic("Invalid file found in storage folder!")
	}
	return strings.Join(pathKey.Path, "/") + pathKey.FileName[:len(pathKey.FileName)-4]
}

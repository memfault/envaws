package param_providers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

type fileProvider struct {
	path         string
	hash         string
	wantedParams []string
}

func (file fileProvider) getAndHashData() (string, error) {
	f, err := os.Open(file.path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	hash := md5.New()

	if _, err = io.Copy(hash, f); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes), nil
}

func (file *fileProvider) Init() {
	initialHash, err := file.getAndHashData()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	file.hash = initialHash
}

func (file fileProvider) Changed() bool {
	currentHash, err := file.getAndHashData()
	if err != nil {
		log.Println(err.Error())
		return true
	}
	return currentHash != file.hash
}

func NewFileProvider(path string) *fileProvider {
	return &fileProvider{
		path:         path,
		wantedParams: []string{"wow"},
	}
}

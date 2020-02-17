package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// KeyMap is the type used to store keys map in memory
type KeyMap = map[string]bool

func loadKeyFromLine(line string) (string, uint64, error) {
	infos := strings.Split(line, ",")
	address := infos[0][2:]
	const MaxUint = ^uint64(0)
	var valueEth uint64
	if len(infos[1]) >= len(string(MaxUint)) {
		valueEth = MaxUint
	} else {
		value, err := strconv.ParseUint(infos[1], 10, 64)
		if err != nil {
			return "", 0, err
		}
		valueEth = value
	}
	return string([]byte(address)), valueEth, nil
}

func loadKeyFromFile(path string) (KeyMap, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	keys := make(KeyMap)

	scanner.Scan()
	for lineIdx := 0; scanner.Scan(); lineIdx++ {
		line := scanner.Text()
		if lineIdx == 0 {
			continue
		}
		address, _, err := loadKeyFromLine(line)
		if err != nil {
			return nil, err
		}
		keys[address] = true
	}
	return keys, nil
}

// GetEthKeys load n ethereum keys from file(s)
func GetEthKeys(keydir string) (KeyMap, error) {
	var wg sync.WaitGroup
	keys := KeyMap{}
	fileKeysResult := make(chan KeyMap)
	go func() {
		filepath.Walk(keydir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}
			if filepath.Ext(path) == ".csv" {
				println("Loading ", path)
				wg.Add(1)
				go func() {
					fileKeys, err := loadKeyFromFile(path)
					if err != nil {
						panic(err)
					}
					println("Done: ", path)
					fileKeysResult <- fileKeys
					wg.Done()
				}()
			} else {
				println("Skipping ", path)
			}
			return nil
		})
		wg.Wait()
		close(fileKeysResult)
	}()
	for {
		fileKeys, ok := <-fileKeysResult
		if ok {
			for address, value := range fileKeys {
				keys[address] = value
			}
		} else {
			break
		}
	}
	println("Done loading")
	return keys, nil
}

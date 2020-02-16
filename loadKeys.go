package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func loadKeyFromLine(line string) (*Key, error) {
	infos := strings.Split(line, ",")
	address := infos[0][2:]
	value, err := strconv.ParseUint(infos[1], 10, 64)
	if err != nil {
		return nil, err
	}
	key := Key{address: address, value: value}
	return &key, nil
}

// GetEthKeys load n ethereum keys from file(s)
func GetEthKeys(maxLoad uint) (map[string]Key, error) {
	fileName := "./keys.csv"
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	keys := make(map[string]Key)

	scanner.Scan()
	for lineIdx := 0; scanner.Scan(); lineIdx++ {
		line := scanner.Text()
		if lineIdx == 0 {
			continue
		}
		if maxLoad != 0 && uint(lineIdx) > maxLoad {
			break
		}
		key, err := loadKeyFromLine(line)
		if err != nil {
			return nil, err
		}
		keys[key.address] = *key
	}
	return keys, nil
}

package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
)

func loadKeyFromLine(line string) (*Key, error) {
	infos := strings.Split(line, ",")
	addressString := strings.ReplaceAll(infos[0], " ", "")
	address, err := hex.DecodeString(addressString[2:])
	if err != nil {
		return nil, err
	}
	value, err := strconv.Atoi(infos[1])
	if err != nil {
		return nil, err
	}
	key := Key{address: binary.BigEndian.Uint32(address), value: uint32(value)}
	return &key, nil
}

func getEthKeys() (map[uint32]Key, error) {
	fileName := "./keys.csv"
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	keys := make(map[uint32]Key)

	scanner.Scan()
	for lineIdx := 0; scanner.Scan(); lineIdx++ {
		line := scanner.Text()
		if lineIdx == 0 {
			continue
		}
		if lineIdx == 10 {
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

func compute(map[uint32]Key) {

}

func main() {
	keys, err := getEthKeys()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	compute(keys)
}

package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
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

func getEthKeys(maxLoad uint) (map[string]Key, error) {
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

func try(keys map[string]Key) (bool, *Key) {
	key := GenerateKey()
	keyFound, exist := keys[key.address]
	if exist {
		key.value = keyFound.value
		println("Found one key!")
		key.Debug()
	}
	return exist, key
}

func infiniteWorker(start chan bool, done chan bool, keys map[string]Key) {
	for {
		<-start
		found, _ := try(keys)
		done <- found
	}
}

func compute(keys map[string]Key, nbThread uint) {
	nbTried := 0
	startTime := time.Now()
	start := make(chan bool)
	done := make(chan bool)
	for i := 0; i < int(nbThread); i++ {
		go infiniteWorker(start, done, keys)
		start <- true
	}
	for {
		<-done
		start <- true
		nbTried++
		if nbTried%10000 == 0 {
			println(time.Since(startTime).String())
			startTime = time.Now()
			println("Tried: ", nbTried)
		}
	}
}

func main() {
	nbthread := flag.Uint("thread", 2, "Number of threads to use")
	maxkeyloaded := flag.Uint("maxKeys", 10, "Max keys to load in memory")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	keys, err := getEthKeys(*maxkeyloaded)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	compute(keys, *nbthread)
}

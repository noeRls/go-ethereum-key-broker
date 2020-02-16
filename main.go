package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
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
	}
	return exist, key
}

func infiniteWorker(start chan bool, done chan *Key, keys map[string]Key) {
	for {
		<-start
		found, key := try(keys)
		if found {
			done <- key
		} else {
			done <- nil
		}
	}
}

func getTimeEstimationForOneKey(keyLoaded uint, nbGeneratedPerMinute uint) string {
	timeToBreakOneKeyMin := 1.5e+48 / float64(keyLoaded) / float64(nbGeneratedPerMinute)
	if timeToBreakOneKeyMin > 290*350*24*60 { // check if the duration isn't going to overflow
		fstr := fmt.Sprintf("%e", timeToBreakOneKeyMin)
		return ">290y (" + fstr + "m)"
	}
	t := time.Duration(time.Minute * time.Duration(timeToBreakOneKeyMin))
	return t.String()
}

func compute(keys map[string]Key, nbThread uint, savepath string, debugtime uint) {
	nbTried := 0
	start := make(chan bool)
	done := make(chan *Key)
	for i := 0; i < int(nbThread); i++ {
		go infiniteWorker(start, done, keys)
		start <- true
	}
	timeout := time.After(time.Second * time.Duration(debugtime))
	for {
		select {
		case key := <-done:
			start <- true
			nbTried++
			if key != nil {
				key.Debug()
				if err := key.Save(savepath); err != nil {
					println(err)
					panic(err)
				}
			}
		case <-timeout:
			fmt.Printf("Tested %d keys in %ds\n", nbTried, debugtime)
			estimation := getTimeEstimationForOneKey(uint(len(keys)), uint(nbTried))
			fmt.Printf("Estimated time to crack one key %s\n", estimation)
			nbTried = 0
			timeout = time.After(time.Second * time.Duration(debugtime))
		}
	}
}

func haveWriteAccess(path string) error {
	infos, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !infos.IsDir() {
		return errors.New("Path " + path + ": is not a directory")
	}
	if err := unix.Access(path, unix.W_OK); err != nil {
		return errors.New("Path " + path + ": " + err.Error())
	}
	return nil
}

func main() {
	nbthread := flag.Uint("thread", 2, "Number of threads to use")
	maxkeyloaded := flag.Uint("maxKeys", 0, "Max keys to load in memory")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	savepath := flag.String("savepath", "./keys_found", "Path to save keys")
	debugevery := flag.Uint("debugtime", 10, "time in seconds between two performance debug")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if err := haveWriteAccess(*savepath); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	keys, err := getEthKeys(*maxkeyloaded)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	compute(keys, *nbthread, *savepath, *debugevery)
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"golang.org/x/sys/unix"
)

func getTimeEstimationForOneKey(keyLoaded uint, nbGenerated uint, timeToGenerateThemSec uint) string {
	timeToBreakOneKeySec := 1.5e+48 / float64(keyLoaded) / float64(nbGenerated) * float64(timeToGenerateThemSec)
	if timeToBreakOneKeySec > 290*350*24*60*60 { // check if the duration isn't going to overflow
		fstr := fmt.Sprintf("%0.1e", timeToBreakOneKeySec)
		return ">290y (" + fstr + "s)"
	}
	t := time.Duration(time.Second * time.Duration(timeToBreakOneKeySec))
	return t.String()
}

func compute(keys KeyMap, nbThread uint, savepath string, debugtime uint) {
	nbTried := 0
	start := make(chan bool)
	done := make(chan *Key)
	for i := 0; i < int(nbThread); i++ {
		go InfiniteWorker(start, done, keys)
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
			fmt.Printf("Tested %d keys in %ds (with a database of %.1e keys)\n", nbTried, debugtime, float64(len(keys)))
			estimation := getTimeEstimationForOneKey(uint(len(keys)), uint(nbTried), debugtime)
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
	keydir := flag.String("keydir", "./keys_db", "Specify keys database")
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

	keys, err := GetEthKeys(*keydir)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	compute(keys, *nbthread, *savepath, *debugevery)
}

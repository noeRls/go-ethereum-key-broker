package main

func try(keys KeyMap) (bool, *Key) {
	key := GenerateKey()
	_, exist := keys[key.address]
	if exist {
		println("Found one key!")
	}
	return exist, key
}

// InfiniteWorker try to match a key until it receive a start signal
func InfiniteWorker(start chan bool, done chan *Key, keys KeyMap) {
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

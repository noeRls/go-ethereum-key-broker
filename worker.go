package main

func try(keys map[string]Key) (bool, *Key) {
	key := GenerateKey()
	keyFound, exist := keys[key.address]
	if exist {
		key.value = keyFound.value
		println("Found one key!")
	}
	return exist, key
}

// InfiniteWorker try to match a key until it receive a start signal
func InfiniteWorker(start chan bool, done chan *Key, keys map[string]Key) {
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

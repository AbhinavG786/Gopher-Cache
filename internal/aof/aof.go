package aof

import (
	"bufio"
	"sync"
	"os"
)

type AOF struct{
	file 	*os.File
	mu 		sync.Mutex
}

func NewAOF(path string)(*AOF, error){
	f,err:=os.OpenFile(path,os.O_APPEND | os.O_CREATE | os.O_RDWR,0644)
	if err!=nil{
		return nil,err
	}
	return &AOF{file:f},nil
}

func(a *AOF) Append(cmd string) error{
	a.mu.Lock()
	defer a.mu.Unlock()
	_,err:=a.file.WriteString(cmd + "\n")
	return err
}

func(a *AOF) Close() error{
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.file.Close()
}

func(a *AOF) Replay(fn func(string)) error{
	a.mu.Lock()
	defer a.mu.Unlock()
	a.file.Seek(0,0)

	scanner:=bufio.NewScanner(a.file)
	for scanner.Scan(){
		line:=scanner.Text()
		fn(line)
	}

	return scanner.Err()

}
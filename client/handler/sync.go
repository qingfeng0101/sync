package handler

import (
	"bufio"
	"log"
	"os"
	file2 "sync/file"
)

func create(file file2.File) (error) {
	f,e := os.OpenFile(file.Name,os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if e != nil{
		log.Println("file name: ",file.Name)
		log.Println("OpenFile err: ",e)
		return e
	}
	defer f.Close()
	write := bufio.NewWriter(f)
	_,err := write.Write(file.Date)
	if err != nil{
		log.Println("Write err: ",err)

		return err
	}
	write.Flush()
	return nil
}


func sync(file file2.File,createorappend int) error {
	f,e := os.OpenFile(file.Name,os.O_WRONLY|createorappend, 0666)
	if e != nil{
		log.Println("file name: ",file.Name)
		log.Println("OpenFile err: ",e)
		return e
	}
	defer f.Close()
	write := bufio.NewWriter(f)
	_,err := write.Write(file.Date)
	if err != nil{
		log.Println("Write err: ",err)

		return err
	}
	write.Flush()
	return nil
}

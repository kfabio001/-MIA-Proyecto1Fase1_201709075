package function

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"
	"math/rand"
	"time"
	"strings"
)

type mbr struct{
	Size int64
	Time[25] byte
	Disk_signature int8
	Partitions[4] partition
}

type partition struct{
	Status byte
	Type byte 
	Fit byte
	Start int64
	Size int64
	Name[16] byte
}

func ReadBinaryFile(file_name string) mbr {
	file, err := os.Open(file_name)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		fmt.Println("An error occurred while creating the disc")
	}
	var m mbr
	var mbr_size int = int(unsafe.Sizeof(m)) + 1 
	data := ReadNextBytes(file, mbr_size)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return m
}

func printDisk(m mbr){
	fmt.Println("Disk size", m.Size)
	myString := string(m.Time[:])
	fmt.Println("Disk created at", myString)
	fmt.Println("Disk signature", m.Disk_signature)
	for i := 0; i < len(m.Partitions); i++ {
		par := m.Partitions[i]
		fmt.Println("------------------------------------------------------------------")
		fmt.Println("Partition", i)
		fmt.Println("Partition status", string(par.Status))
		fmt.Println("Partition type", string(par.Type))
		fmt.Println("Partition fit", string(par.Fit))
		fmt.Println("Partition start",par.Start)
		fmt.Println("Partition size",par.Size)
		parName := string(par.Name[:])
		fmt.Println("Partition name",parName)
		fmt.Println("------------------------------------------------------------------")
	}

}

func ReadNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func CreateBinaryFile(file_name string, file_path string, file_size int64) {
	if(!strings.HasSuffix(file_path,"/")){
		file_path = file_path + "/"
	}
	file_route := file_path + file_name
	file, err := os.Create(file_route)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		fmt.Println("Cannot create the file")
	}
	var cero int64 = 0
	s := &cero
	//var bin_buf bytes.Buffer
	//binary.Write(&bin_buf, binary.BigEndian, s)
	//WriteNextBytes(file, bin_buf.Bytes())
	var mbrs mbr
	mbrs.Size = file_size
	mbrs.Disk_signature = generateRandom()
	dt := time.Now()
	tiempo := dt.Format("01-02-2006 15:04:00")
	copy(mbrs.Time[:],tiempo)
	for i:=0;i<len(mbrs.Partitions);i++ {
		mbrs.Partitions[i].Status = '0'
	}
	//mbr_size := unsafe.Sizeof(mbrs)
	ss := &mbrs
	var mbr_buf bytes.Buffer
	binary.Write(&mbr_buf, binary.BigEndian, ss)
	WriteNextBytes(file, mbr_buf.Bytes())

	file.Seek(file_size,0)

	var second_buffer bytes.Buffer
	binary.Write(&second_buffer, binary.BigEndian, s)
	WriteNextBytes(file, second_buffer.Bytes())
	file.Seek(0,0)
	
}

func WriteBFile(file_path string, rec mbr, pos int){
	file, err := os.Create(file_path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		fmt.Println("Cannot write the file")
	}
	file.Seek(0, pos)
	ss := &rec
	var mbr_buf bytes.Buffer
	binary.Write(&mbr_buf, binary.BigEndian, ss)
	WriteNextBytes(file, mbr_buf.Bytes())
	file.Seek(0, 0)
}

func WriteNextBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func generateRandom() int8{
	return int8(rand.Int())
}
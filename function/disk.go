package function

import (
	"fmt"
	"strconv"
	"strings"
	"os"
	"reflect"
    "unsafe"
)

var Equalizer string = "->"
var PartitionError bool = false
var Disks[100] mounted_disk  
var Disks_size int = 0
var mounts[100] mount
var Mounts_size int =0

type mounted_disk struct{
	Size int
	Path string
	Name string
	Identifier string
	Created string
	Mounted_partitions int
}


type binaryFile struct {
	Size int
	Path string
	Name string
	Unit string
}

type partition_config struct{
	Size int64
	Unit byte
	Path string
	Type byte
	Fit[1] byte
	Delete bool
	Name string
	Add bool
}

type mount struct{
	Path string
	Name string
	Identifier string
}

/*MOUNT COMMAND*/

func GetMountIdentifier(path string)string{
	for i, element := range Disks{
		absolute_path := element.Path + element.Name
		if(absolute_path==path){
			id := "vd" + element.Identifier + strconv.Itoa(element.Mounted_partitions)
			element.Mounted_partitions = element.Mounted_partitions + 1
			Disks[i] = element
			return id
		}
	}
	return ""
}


func mount_disk(path string, name string){
	m := ReadBinaryFile(path + name)
	var dsk mounted_disk 
	dsk.Identifier = GetIdentifier(Disks_size)
	dsk.Size = int(m.Size)
	dsk.Path = path
	dsk.Name = name
	tm := string(m.Time[:])
	dsk.Created = tm
	Disks[Disks_size] = dsk
	Disks_size += 1
}

func verifyMountedDisk(path string, name string)bool{
	for _,element := range Disks{
		if(element.Name!=""){
			path_abs := element.Path + element.Name
			if(compareBytes(path+name, path_abs)){
				return true
			}
		}
	}
	return false
}

func Exec_mount(com [] string){
	var new_mount mount
	for _,element := range com {
		spplited_command := strings.Split(element, Equalizer)
		switch  strings.ToLower(spplited_command[0]) {
		case "-path":
			if _, err := os.Stat(spplited_command[1]); !os.IsNotExist(err) {
				new_mount.Path = spplited_command[1]
			}else{
				fmt.Println("Especificated disk doesnt exist")
				return
			}
		case "-name":
			dsik := ReadBinaryFile(new_mount.Path)
			if(!verifyMountedDisk(getPath(new_mount.Path))){
				mount_disk(getPath(new_mount.Path))
			}
			for _,element := range dsik.Partitions {
				name_dsk := strings.TrimRight(string(element.Name[:])," ") 
				if(compareBytes(spplited_command[1],name_dsk)){
					new_mount.Name = spplited_command[1]
				}
			}
			if(new_mount.Name == ""){
				fmt.Println("Partition doesnt exists in disk")
			}
		}
	}
	if(new_mount.Path != "" && new_mount.Name != ""){
		new_mount.Identifier = GetMountIdentifier(new_mount.Path)
		fmt.Println("PARTITION ", new_mount.Identifier, "MOUNTED")
		mounts[Mounts_size] = new_mount
		Mounts_size += 1
		
	}else{
		fmt.Println("Too few arguments")
	}
}

func getPath(p string)(string, string){
	sp := strings.Split(p, "/")
	name := sp[len(sp)-1]
	path := strings.TrimRight(p, name)
	return path, name
}

func compareBytes(str1 string, str2 string)bool{
	for i:=0;i< len(str1); i++{
		if(!(i>=len(str1)) && !(i>=len(str2))){
			if(!(str1[i]==str2[i])){
				return false
			}
		}
	}
	return true
}

func PrintMount(){
	for _,element := range mounts{
		if(element.Identifier!=""){
			fmt.Println("IDENTIFIER:", element.Identifier)
			fmt.Println("DISK:", element.Path)
			fmt.Println("PARTITION", element.Name)
			fmt.Println("----------------------------------")
		}
	}
}

func Exec_fdisk(com []string) {
	var new_partition partition_config
	for _, element := range com {
		spplited_command := strings.Split(element, Equalizer)
		switch strings.ToLower(spplited_command[0]) {
		case "-size":
			i, _ := strconv.Atoi(spplited_command[1])
			if i > 0 {
				new_partition.Size = int64(i)
				//fmt.Println("Partition size ",new_partition.Size)
			} else {
				fmt.Println("Partition size must be positive")
				return
			}
		case "-unit":
			new_partition.Unit = spplited_command[1][0]
			//fmt.Println("Partition unit",new_partition.Unit)
		case "-path":
			if _, err := os.Stat(spplited_command[1]); !os.IsNotExist(err) {
				new_partition.Path = spplited_command[1]
				//fmt.Println("Disk path",new_partition.Path )
			}else{
				fmt.Println("Especificated disk doesnt exist")
				return
			}
		case "-type":
			new_partition.Type =  strings.ToLower(spplited_command[1])[0]
			//fmt.Println("Partition type", new_partition.Type)
		case "-fit":
			var fit_slice[1] byte
			copy(fit_slice[:], strings.ToLower(spplited_command[1])) 
			new_partition.Fit = fit_slice
			//fmt.Println("Partition fit", fit_slice)
		case "-delete":
			new_partition.Delete = true
		case "-name":
			new_partition.Name = spplited_command[1]
			//fmt.Println("Partition name",new_partition.Name)
		case "-add":
			new_partition.Add = true
		default:
			if spplited_command[0] != "fdisk" {
				fmt.Println(spplited_command[0], "command unknow")
			}
		}
	}
	if(new_partition.Size!=0  && new_partition.Path!="" && new_partition.Name!=""){
		if new_partition.Unit == 0 {
			new_partition.Unit = 'k'
			fmt.Println("You didnt specify an unit size")
		}
		if(new_partition.Add && !new_partition.Delete){
			
		}else if(new_partition.Delete && !new_partition.Add){
	
		}else if(!new_partition.Delete && !new_partition.Add){
			record := ReadBinaryFile(new_partition.Path)
			e,_,_ := calcPart(record.Partitions)
			if(e==1){
				if(new_partition.Type=='e'){
					fmt.Println("THERE IS ONE EXTENDED PARTITION ALREADY")
					return
				}
			}
			dks := GetDiskByName(new_partition.Path)
			createPartition(&record, new_partition, dks)
			if(!PartitionError){
				WriteBFile(new_partition.Path, record, 1)
				printDisk(ReadBinaryFile(new_partition.Path))
			}else{
				PartitionError = false
				return
			}
		}else {
			fmt.Println("Incorrect params combination")
		}
	}else {
		fmt.Println("Too few arguments")
	}
}

func BytesToString(b []byte) string {
    bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
    sh := reflect.StringHeader{bh.Data, bh.Len}
    return *(*string)(unsafe.Pointer(&sh))
}

func GetDiskByName(p string)mounted_disk{
	var n mounted_disk
	for _,element := range Disks{
		if(element.Path==p){
			n = element
		}
	}
	return n
}

func createPartition(r *mbr, p partition_config, ds mounted_disk){
	par_unit := string(p.Unit)
	disk_size := r.Size
	part_size := calc_filesize(par_unit,int(p.Size), true)
		for i:=0;i<len(r.Partitions);i++ {
			st := r.Partitions[i].Status
			if(st == '0'){
				r.Partitions[i].Status = 'i'
					r.Partitions[i].Type = p.Type
					r.Partitions[i].Fit = p.Fit[0]
					if(i==0){
						start_first := unsafe.Sizeof(r)
						r.Partitions[i].Start = int64(start_first)
					}else{
						ps := int64(r.Partitions[i-1].Start) + r.Partitions[i-1].Size
						total_size := ps + part_size
						if(total_size<=disk_size){
							r.Partitions[i].Start = ps
						}else{
							fmt.Println("Not enough space on disk")
							PartitionError = true
						}
					}
					r.Partitions[i].Size = part_size
					var parN[16] byte
					copy(parN[:],p.Name)
					r.Partitions[i].Name = parN
				return 
			}	
		}
}

func calcPart(parti [4] partition)(int, int, int){
	primary := 0
	free:=0
	extended := 0
	for i:=0;i<len(parti);i++ {
		if(parti[i].Type == 'p'){
			primary += 1
		}else if(parti[i].Type == 'e'){
			extended +=1
		}else{
			free +=1
		}
	}
	return extended, primary, free
}

func Exec_mrdisk(com []string) {
	splitted_command := strings.Split(com[1], Equalizer)
	if splitted_command[0] == "-path" {
		file_name := splitted_command[1]
		deleteFile(file_name)
	} else {
		fmt.Println(splitted_command[0], "command unknow")
	}
}

func Exec_mkdisk(com []string) {
	var new_disk binaryFile
	for _, element := range com {
		spplited_command := strings.Split(element, Equalizer)
		switch strings.ToLower(spplited_command[0]) {
		case "-size":
			i, _ := strconv.Atoi(spplited_command[1])
			if i > 0 {
				new_disk.Size = i
			} else {
				fmt.Println("Size must be positive! ")
				return
			}
		case "-path":
			if _, err := os.Stat(spplited_command[1]); os.IsNotExist(err) {
				os.MkdirAll(spplited_command[1], os.ModePerm)
			}
			new_disk.Path = spplited_command[1]
		case "-name":
			if strings.HasSuffix(spplited_command[1], ".dsk") {
				new_disk.Name = spplited_command[1]
			} else {
				fmt.Println("Error! Name must have .dsk extension")
				return 
			}
		case "-unit":
			new_disk.Unit = spplited_command[1]
		default:
			if spplited_command[0] != "mkdisk" {
				fmt.Println(spplited_command[0], "command unknow")
			}
		}
	}
	if(new_disk.Path!="" && new_disk.Size != 0 && new_disk.Name!=""){
		CreateBinaryFile(new_disk.Name,new_disk.Path, calc_filesize(new_disk.Unit, new_disk.Size,false))
		filen := new_disk.Path+ new_disk.Name 
		printDisk(ReadBinaryFile(filen))
	}else{
		fmt.Println("Too few arguments")
	}	
}

func calc_filesize(unit string, size int, partition bool)int64{
	if(unit=="" && !partition){
		unit = "m"
	}else if(unit=="" && partition){
		unit = "k"
	}
	switch strings.ToLower(unit) {
	case "k":
		return 1024*int64(size)
	case "m":
		return 1024*1024*int64(size)
	case "b":
		return int64(size)
	default:
		fmt.Println("Invalid unit formmat")
	}
	return 0
}

func deleteFile(path string){
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err := os.Remove(path)
		if err != nil {
			fmt.Println(err)
			return
		}	
		fmt.Println("Removed successfully!")
	}else{
		fmt.Println("Error: File doesnt exists!")
	}
}

func GetIdentifier(elements int) string {
	return string(97+elements)
}
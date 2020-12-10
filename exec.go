package comandos

import (
	"os"
	"reflect"
	"strconv"	
	"fmt"
	"strings"
	"unsafe"
	
)
var inicioLog
var typmount=0
var getTypet=""
var Error_Part bool = false
var Error_Disk=false
var Disks [80]montados
var Disk_Creados=0
var Disks_size int = 0
var mounts [80]mount
var conteoMount=0
var tamano_moun int = 0

func otime()
{
time_t t = time(0);
struct tm *tlocal = localtime(&t);
strftime(data,30,"%d/%m/%y %H:%M",tlocal);
}


type mount struct {
	Path       string
	Name       string
	Identifier string
}
func writeFile() {
	file, err := os.Create("test.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var otro int8 = 0

	s := &otro

	fmt.Println(unsafe.Sizeof(otro))
	//Escribimos un 0 en el inicio del archivo.
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	escribirBytes(file, binario.Bytes())
	//Nos posicionamos en el byte 1023 (primera posicion es 0)	
	file.Seek(1023, 0) // segundo parametro: 0, 1, 2.     0 -> Inicio, 1-> desde donde esta el puntero, 2 -> Del fin para atras
	
	//Escribimos un 0 al final del archivo.
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	escribirBytes(file, binario2.Bytes())

//----------------------------------------------------------------------- //
	//Escribimos nuestro struct en el inicio del archivo

	file.Seek(0, 0) // nos posicionamos en el inicio del archivo.
	
	//Asignamos valores a los atributos del struct.
	disco := mbr{Numero: 5}
	disco.Caracter = 'a'

	// Igualar cadenas a array de bytes (array de chars)
	cadenita := "Hola Amigos"
	copy(disco.Cadena[:], cadenita)

	s1 := &disco
	
	//Escribimos struct.
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())
	
}
func IDmount(path string) string {
	for i, element := range Disks {
		ruta := element.Path + element.Name
		if ruta == path {
			id := "vd" + element.Identifier + strconv.Itoa(element.particiones_montadas)
			element.particiones_montadas = element.particiones_montadas + 1
			Disks[i] = element
			return id
		}
	}
	return ""
}

func disco_mount(path string, name string) {
	m := ArchivoLeido(path + name)
	var dsk montados
	dsk.Identifier = ID_ob(Disks_size)
	dsk.Size = int(m.Size)
	dsk.Path = path
	dsk.Name = name
	tm := string(m.Time[:])
	dsk.Created = tm
	Disks[Disks_size] = dsk
	Disks_size += 1
}

func verificaMount(path string, name string) bool {
	for _, element := range Disks {
		if element.Name != "" {
			path_abs := element.Path + element.Name
			if compareBytes(path+name, path_abs) {
				return true
			}
		}
	}
	return false
}

func comando_mounts(com []string) {
	var new_mount mount

	for _, element := range com {
		spplited_command := strings.Split(element, "->")
		switch strings.ToLower(spplited_command[0]) {
		case "path":
			if _, err := os.Stat(spplited_command[1]); !os.IsNotExist(err) {
				new_mount.Path = spplited_command[1]
			} else {
				fmt.Println("Disco no existente")
				return
			}
		case "name":
			dsik := ArchivoLeido(new_mount.Path)
			if !verificaMount(getPath(new_mount.Path)) {
				disco_mount(getPath(new_mount.Path))
			}
			for _, element := range dsik.Partitions {
				name_dsk := strings.TrimRight(string(element.Name[:]), " ")
				if compareBytes(spplited_command[1], name_dsk) {
					new_mount.Name = spplited_command[1]
				}
			}
			if new_mount.Name == "" {
				fmt.Println("Particion inexistente")
			}
		}
	}
	if new_mount.Path != "" && new_mount.Name != "" {
		new_mount.Identifier = IDmount(new_mount.Path)
		fmt.Println("PARTITION ", new_mount.Identifier, "MOUNTED")
		mounts[tamano_moun] = new_mount
		tamano_moun += 1

	} else {
		fmt.Println("Too few arguments")
	}
}
type propiedadesPart struct {
	Size   int64
	Unit   byte
	Path   string
	Fit    [1]byte
	Delete bool
	Name   string
	Type   byte
	Add    bool
}
func getPath(p string) (string, string) {
	sp := strings.Split(p, "/")
	name := sp[len(sp)-1]
	path := strings.TrimRight(p, name)
	return path, name
}

func compareBytes(str1 string, str2 string) bool {
	for i := 0; i < len(str1); i++ {
		if !(i >= len(str1)) && !(i >= len(str2)) {
			if !(str1[i] == str2[i]) {
				return false
			}
		}
	}
	return true
}

func imp_mounts() {
	for _, element := range mounts {
		if element.Identifier != "" {
			fmt.Println("ID:", element.Identifier)
			fmt.Println("DISK:", element.Path)
			fmt.Println("PARTI", element.Name)
			fmt.Println("#########################3")
		}
	}
}
type propiedadesDisco struct {
	Size int
	Path string
	Name string
	Unit string //aun no funciona :(
}

type montados struct {
	Size int
	Path string
	Name string
	Identifier string
	Created string
	particiones_montadas int
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func obtenerNombrDisc(p string) montados {
	var n montados
	for _, element := range Disks {
		if element.Path == p {
			n = element
		}
	}
	return n
}

func crearPartition(r *mbr, p propiedadesPart, ds montados) {
	par_unit := string(p.Unit)
	disk_size := r.Size
	part_size := calcular_Tamano(par_unit, int(p.Size), true)
	for i := 0; i < len(r.Partitions); i++ {
		st := r.Partitions[i].Status
		if st == '0' {
			r.Partitions[i].Status = 'i'
			r.Partitions[i].Type = p.Type
			r.Partitions[i].Fit = p.Fit[0]
			if i == 0 {
				start_first := unsafe.Sizeof(r)
				r.Partitions[i].Start = int64(start_first)
			} else {
				ps := int64(r.Partitions[i-1].Start) + r.Partitions[i-1].Size
				total_size := ps + part_size
				if total_size <= disk_size {
					r.Partitions[i].Start = ps
				} else {
					fmt.Println("Espacio insuficiente")
					Error_Part = true
				}
			}
			r.Partitions[i].Size = part_size
			var parN [16]byte
			copy(parN[:], p.Name)
			r.Partitions[i].Name = parN
			return
		}
	}
}

func calcPart(parti [4]partition) (int, int, int) {
	primary := 0
	free := 0
	extended := 0
	for i := 0; i < len(parti); i++ {
		if parti[i].Type == 'p' {
			primary += 1
		} else if parti[i].Type == 'e' {
			extended += 1
		} else {
			free += 1
		}
	}
	return extended, primary, free
}


func comando_mkdisks(com []string) {
	var new_disk propiedadesDisco
	//contenido := ""
	for j := 0; j < len(com); j++ {
		com[j] = strings.Replace(com[j], "\"", "", -1)
		com[j] = strings.ToLower(com[j])

	}
	//contenido := string(com)
	//fmt.Println(contenido + "contenido")
	//come := strings.Split(contenido, " ")
	for _, element := range com {
		spplited_command := strings.Split(element, "->")
		switch strings.ToLower(spplited_command[0]) {
		case "size":
			i, _ := strconv.Atoi(spplited_command[1])
			if i > 0 {
				new_disk.Size = i
			} else {
				fmt.Println("Tamano no admitido ")
				return
			}
		case "path":
			if _, err := os.Stat(spplited_command[1]); os.IsNotExist(err) {
				os.MkdirAll(spplited_command[1], os.ModePerm)
			}
			new_disk.Path = spplited_command[1]
		case "name":
			if strings.HasSuffix(spplited_command[1], ".dsk") {
				new_disk.Name = spplited_command[1]
			} else {
				fmt.Println("Error en el comando")
				return
			}
		case "unit":
			new_disk.Unit = spplited_command[1]
		default:
			if spplited_command[0] != "mkdisk" {
				fmt.Println(spplited_command[0], "Comando Desconocido")
			}
		}
	}
	if new_disk.Path != "" && new_disk.Size != 0 && new_disk.Name != "" {
		crearArchivo(new_disk.Name, new_disk.Path, calcular_Tamano(new_disk.Unit, new_disk.Size, false))
		filen := new_disk.Path + new_disk.Name
		consolaDisco(ArchivoLeido(filen))
	} else {
		fmt.Println("Too few arguments")
	}
}

func calcular_Tamano(unit string, size int, partition bool) int64 {
	if unit == "" && !partition {
		unit = "m"
	} else if unit == "" && partition {
		unit = "k"
	}
	switch strings.ToLower(unit) {
	case "k":
		return 1024 * int64(size)
	case "m":
		return 1024 * 1024 * int64(size)
	case "b":
		return int64(size)
	default:
		fmt.Println("Formato no adimitido")
	}
	return 0
}



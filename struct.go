package comandos

import (
	"bytes"
	"strings"
	"time"
	"unsafe"
	"fmt"
	"log"
	"math/rand"
	"os"
)
type struct MBR_particion{
var status; //status
var type;// itpos
char fit;// S M R ajuste
Size size;
byte name[16];
}
type struct T_inodo{
var i_link int64; //numero de enlaces duros;
var i_uid int64; // UID del usuario propietario del archivo o carpeta;
var i_gid int64; // GID del grupo al que pertenece el archivo o carpeta;
var i_size int64; //tamao del fichero en bytes
var i_dia int64;  //fecha que se leyo el inodo sin modificarlos
var i_diac int64;  //fecha que se creo el inodo
var m_dia int64;  //fecha de la ultima modificacion del inodo
var i_block[15] int64; //bloques de inodos 12 Directos y 3 Indirectos
var i_type byte; //inidica si es archivo a, carpeta c o enlace simbolico e;
var i_perm[9] int;
}
func ArchivoLeido(nombre string) mbr {
	file, err := os.Open(nombre)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error en la creacion del disco")
	}
	var m mbr
	var mbr_size int = int(unsafe.Sizeof(m)) + 1
	data := BytesPLeidos(file, mbr_size)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("Error en lectura", err)
	}
	return m
}
type mbr struct {
	Size           int64
	Time           [27]byte
	Disk_signature int8
	Partitions     [4]partition
}
type mbr struct {
	Numero uint8
	Caracter byte
	Cadena [20]byte
}
func consolaDisco(m mbr) {
	fmt.Println("Tamano disco", m.Size)
	myString := string(m.Time[:])
	fmt.Println("Disco creado", myString)
	fmt.Println("Asignado", m.Disk_signature)
	for i := 0; i < len(m.Partitions); i++ {
		par := m.Partitions[i]
		fmt.Println("#########################################################3#")
		fmt.Println("Particion No", i)
		fmt.Println("Particion Estado", string(par.Status))
		fmt.Println("Particion Tipo", string(par.Type))
		fmt.Println("Particion start", par.Start)
		fmt.Println("Particion tamano", par.Size)
		fmt.Println("Particion fit", string(par.Fit))
		parName := string(par.Name[:])
		fmt.Println("Particion nombre", parName)
		fmt.Println("############################################################")
	}

}

func BytesPLeidos(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func crearArchivo(nombre string, file_path string, file_size int64) {
	if !strings.HasSuffix(file_path, "/") {
		file_path = file_path + "/"
	}
	file_route := file_path + nombre
	file, err := os.Create(file_route)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		fmt.Println("No se creo")
	}
	

}
type partition struct {
	Status byte
	Type   byte
	Fit    byte
	Start  int64
	Size   int64
	Name   [16]byte
}
func escribirArchivo(file_path string, rec mbr, pos int) {
	file, err := os.Create(file_path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		fmt.Println("No se puede escribir archivo")
	}
	file.Seek(0, pos)
	ss := &rec
	var mbr_buf bytes.Buffer
	binary.Write(&mbr_buf, binary.BigEndian, ss)
	escribirBytes(file, mbr_buf.Bytes())
	file.Seek(0, 0)
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}


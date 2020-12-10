package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)


var comandoos string = ""

func imprimir_rashos(){
fmt.Println("tamaño superbloque:%i\n",sizeof(struct Super_bloque));
fmt.Println("tamaño inodo:%i\n",sizeof(struct T_inodo));
fmt.Println("tamaño bloque:%i\n",sizeof(struct T_bloque));
fmt.Println("tamaño Journaling:%i\n",sizeof(struct Journaling));
int a = calcular_n_ext2(tam,sizeof(struct Super_bloque),sizeof(struct T_inodo),sizeof(struct T_bloque));
int b = calcular_n_ext3(tam,sizeof(struct Super_bloque),sizeof(struct T_inodo),sizeof(struct T_bloque),sizeof(Journaling));
fmt.Println("tamaño ext2:%i\ntamaño ext3:%i\n",a,b);
}
func main() {
	fmt.Println("Bienvenido a sistema LWH")
	fmt.Println("Ingrese comando para iniciar o escriba cerrar para salir")
	leer := bufio.NewReader(os.Stdin)
	cerrar := false
	for !cerrar {
		input, _ := leer.ReadString('\n')
		input = obtener(input)
		if input != "cerrar" {
			if !strings.HasPrefix(input, "#") {
				mostrar(input)
			}
		} else {
			fmt.Println("Gracias por usar el sistema")
			cerrar = true
		}
	}
}

func mostrar(i string) {
	//fmt.Println("execute")
	if !strings.HasSuffix(i, "/*") {
		comandoos += obtener(i)
		//fmt.Println(strings.Split(comandoos, " -"))
		entro_todo(strings.Split(comandoos, " -"))
		//fmt.Println(splitter(comandoos))
		//entro_todo(splitter(comandoos))
		comandoos = ""
	} else {
		comandoos += strings.TrimRight(i, "/*")
	}
}

/*func ejecutarComando(commandArray []string) {
  data := strings.ToLower(commandArray[0])
	if data == "crear" {
		fmt.Println("Creando un archivo")
	}else {
		fmt.Println("Otro Comando")
	}
}*/
func entro_todo(linea_cm []string) {
	//fmt.Println("reconogize")

	//linea_cm[0] = strings.ToLower(linea_cm[0])
	fmt.Println(linea_cm[0])
	//fmt.Println(linea_cm[0] + " l " + linea_cm[0] + " l ")
	linea_cm[0] = strings.ToLower(linea_cm[0])
	switch strings.ToLower(linea_cm[0]) {
	case "mkdisk":
		fmt.Println(linea_cm)
		comandos.comando_mkdisks(linea_cm)
	case "exec":
		linea_cm2 := strings.Split(linea_cm[1], "->")
		if strings.ToLower(linea_cm2[0]) == "path" {
			readFile(linea_cm2[1])
		} else {
			fmt.Println("No reconocido ruta o comando no encontrado ")

		}
	case "rmdisk":
		comandos.comando_mrdisks(linea_cm)
	case "fdisk":
		comandos.comando_fdisks(linea_cm)
	case "pause":
		fmt.Print("Pausa Teclee una letra para continuar ")
		leer := bufio.NewReader(os.Stdin)
		x, _ := leer.ReadString('\n')
		x += ""
	case "mount":
		fmt.Println(linea_cm)
		if len(linea_cm) >= 2 {
			comandos.comando_mounts(linea_cm)
		} else {
			fmt.Println("Particiones montadas")
			fmt.Println("#########################################")
			comandos.imp_mounts()
		}
	default:
		fmt.Println("Comando no reconozido ")
	}
}
func separar(txt string) []string {
	linea_cm := strings.Split(txt, " ")
	return linea_cm
}
```go
func readFile() {
	//Abrimos/creamos un archivo.
	file, err := os.Open("test.bin")
	defer file.Close() 
	if err != nil { //validar que no sea nulo.
		log.Fatal(err)
	}

	//Declaramos variable de tipo mbr
	m := mbr{}
	//Obtenemos el tamanio del mbr
	var size int = int(unsafe.Sizeof(m))

	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)
	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)
	
	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	
	//Se imprimen los valores guardados en el struct
	fmt.Println(m)
	fmt.Printf("Caracter: %c\nCadena: %s\n", m.Caracter,  m.Cadena)
}
```
func readFile(nombre string) {
	//fmt.Println("entro" + nombre)
	comandoos = "d"

	bytesLeidos, err := ioutil.ReadFile(nombre)
	contenido := string(bytesLeidos)
	//salto := "\\*\r\n"
	contenido = strings.Replace(contenido, "\\*\r\n", "", -1)
	b := []byte(contenido)
	//f, err := os.Open(nombre)
	err = ioutil.WriteFile(nombre, b, 0644)
	//fmt.Println(contenido)
	f, err := os.Open(nombre)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	cadena := bufio.NewScanner(f)

	for cadena.Scan() {
		if cadena.Text() != " " {

			if !strings.HasPrefix(cadena.Text(), "#") {
				fmt.Println("Comando ", cadena.Text(), " ")
				//fmt.Println(strings.TrimRight(cadena.Text(), " "))
				mostrar(strings.TrimRight(cadena.Text(), " "))
			} else {
				fmt.Println(cadena.Text())
			}
		}
	}
	if err := cadena.Err(); err != nil {
		log.Fatal(err)
		return
	}
}



package main

import (
	"os"
	"fmt"
	//"log"
	"time"
	"context"
	"math/rand"
	"strconv"
	"encoding/csv"
	"github.com/432i/T1SisDistribuidos/logistica/chat"
	"google.golang.org/grpc"
)

type Camion struct {
	IdCamion string
	Tipo string
	Paquete1 *chat.Paquete
	Paquete2 *chat.Paquete
}

func crearRegistro(nombreArchivo string){
	archivo, _ := os.Create(nombreArchivo)
	archivo.Close()
}

func guardarPaquete(nombreArchivo string, id string, tipo string, valor string, origen string, destino string, intentos string, fechaEntrega string){
	orden := []string{id, tipo, valor, origen, destino, intentos, fechaEntrega}
	archivo, _ := os.OpenFile(nombreArchivo, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	w := csv.NewWriter(archivo)

	w.Write(orden)

	w.Flush()
	archivo.Close()
}

func getTime() string {
    t := time.Now()
    return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())
}

func Intento(paquete *chat.Paquete) {
	intentos, _ := strconv.Atoi(paquete.Intentos)
	valor, _ := strconv.Atoi(paquete.Valor)
	if paquete.Estado != "Recibido" || paquete.Estado != "No Recibido" {
		if paquete.Tipo == "retail" {
			if intentos < 3 {
				if rand.Float64() <= 0.8 {
					paquete.Estado = "Recibido"
				} else {
					intentos += 1
					paquete.Intentos = strconv.Itoa(intentos)
				}
			} else {
				paquete.Estado = "No Recibido"
			}
		} else {
			if intentos * 10 < valor && intentos < 2 {
				if rand.Float64() <= 0.8 {
					paquete.Estado = "Recibido"
				} else {
					intentos += 1
					paquete.Intentos = strconv.Itoa(intentos)
				}
			} else {
				paquete.Estado = "No Recibido"
			}
		}
	}
}

func Entrega(camion Camion, tEnvio int) bool {
	if camion.Paquete1.Estado == "" && camion.Paquete2.Estado == "" {
		return false;
	} else if camion.Paquete1.Estado == "" && camion.Paquete2.Estado != "" && camion.Paquete2.Estado != "No Recibido" && camion.Paquete2.Estado != "Recibido" {
		Intento(camion.Paquete2)
	} else if camion.Paquete1.Estado != "" && camion.Paquete2.Estado == "" && camion.Paquete1.Estado != "No Recibido" && camion.Paquete1.Estado != "Recibido" {
		Intento(camion.Paquete1)
	} else if camion.Paquete1.Estado != "En Camino" && camion.Paquete2.Estado != "En Camino" {
		return false
	} else if camion.Paquete1.Estado == "Recibido" || camion.Paquete1.Estado == "No Recibido" {
		Intento(camion.Paquete2)
	} else if camion.Paquete2.Estado == "Recibido" || camion.Paquete2.Estado == "No Recibido" {
		Intento(camion.Paquete1)
	} else if camion.Paquete1.Valor > camion.Paquete2.Valor {
		Intento(camion.Paquete1)
		time.Sleep(time.Duration(tEnvio) * time.Second)
		Intento(camion.Paquete2)
	} else {
		Intento(camion.Paquete2)
		time.Sleep(time.Duration(tEnvio) * time.Second)
		Intento(camion.Paquete1)
	}
	return true
}

func Carga(camion Camion, tEspera int, tEnvio int, nombreArchivo string) {
	var conn *grpc.ClientConn
	conn, _ = grpc.Dial("10.6.40.149:9000", grpc.WithInsecure())
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	mensaje := chat.Message{
		Body: camion.Tipo,
	}

	time.Sleep(time.Duration(tEspera) * time.Second)

	paquete1, _ := c.PaqueteQueueToCamion(context.Background(), &mensaje)
	if paquete1.GetId() != "" {
		camion.Paquete1 = paquete1
		camion.Paquete1.Estado = "En Camino"
		msj := chat.Message{
			Body: camion.Paquete1.GetSeguimiento() + ",En Camino",
		}
		respuesta, _ := c.ModificarEstado(context.Background(), &msj)
		fmt.Println(respuesta.GetBody())
		fmt.Printf("Paquete recibido, detalle:\n")
		fmt.Println("     Id: ", camion.Paquete1.Id)
		fmt.Println("     Seguimiento: ", camion.Paquete1.Seguimiento)
		fmt.Println("     Tipo: ", camion.Paquete1.Tipo)
		fmt.Println("     Valor: ", camion.Paquete1.Valor)
		fmt.Println("     Intentos: ", camion.Paquete1.Intentos)
		fmt.Println("     Estado: ", camion.Paquete1.Estado)
		fmt.Println("     Origen: ", camion.Paquete1.Origen)
		fmt.Println("     Destino: ", camion.Paquete1.Destino)
	} else {
		fmt.Println("No hay paquetes en la cola")
		camion.Paquete1 = paquete1
	}

	time.Sleep(time.Duration(tEspera) * time.Second)

	paquete2, _ := c.PaqueteQueueToCamion(context.Background(), &mensaje)
	if paquete2.GetId() != "" {
		camion.Paquete2 = paquete2
		camion.Paquete2.Estado = "En Camino"
		msj := chat.Message{
			Body: camion.Paquete2.GetSeguimiento() + ",En Camino",
		}
		respuesta, _ := c.ModificarEstado(context.Background(), &msj)
		fmt.Println(respuesta.GetBody())
		fmt.Printf("     Paquete recibido, detalle:\n")
		fmt.Println("     Id: ", camion.Paquete2.Id)
		fmt.Println("     Seguimiento: ", camion.Paquete2.Seguimiento)
		fmt.Println("     Tipo: ", camion.Paquete2.Tipo)
		fmt.Println("     Valor: ", camion.Paquete2.Valor)
		fmt.Println("     Intentos: ", camion.Paquete2.Intentos)
		fmt.Println("     Estado: ", camion.Paquete2.Estado)
		fmt.Println("     Origen: ", camion.Paquete2.Origen)
		fmt.Println("     Destino: ", camion.Paquete2.Destino)
	} else {
		fmt.Println("No hay paquetes en la cola")
		camion.Paquete2 = paquete2
	}

	aux := true
	for aux {
		aux = Entrega(camion, tEnvio)
	}
	var msj chat.Message

	respuesta, _ := c.PaqueteCamionToQueue(context.Background(), camion.Paquete1)
	fmt.Println(respuesta)
	msj = chat.Message{
			Body: camion.Paquete1.GetSeguimiento() + "," + camion.Paquete1.GetEstado(),
		}
	respuesta, _ = c.ModificarEstado(context.Background(), &msj)

	if camion.Paquete1.Estado != "No Recibido" {
		guardarPaquete(nombreArchivo, camion.Paquete1.Id, camion.Paquete1.Tipo, camion.Paquete1.Valor, camion.Paquete1.Origen, camion.Paquete1.Destino, camion.Paquete1.Intentos, getTime())
	} else {
		guardarPaquete(nombreArchivo, camion.Paquete1.Id, camion.Paquete1.Tipo, camion.Paquete1.Valor, camion.Paquete1.Origen, camion.Paquete1.Destino, camion.Paquete1.Intentos, "0")
	}

	respuesta, _ = c.PaqueteCamionToQueue(context.Background(), camion.Paquete2)
	fmt.Println(respuesta)
	msj = chat.Message{
			Body: camion.Paquete2.GetSeguimiento() + "," + camion.Paquete2.GetEstado(),
		}
	respuesta, _ = c.ModificarEstado(context.Background(), &msj)
	
	if camion.Paquete2.Estado != "No Recibido" {
		guardarPaquete(nombreArchivo, camion.Paquete2.Id, camion.Paquete2.Tipo, camion.Paquete2.Valor, camion.Paquete2.Origen, camion.Paquete2.Destino, camion.Paquete2.Intentos, getTime())
	} else {
		guardarPaquete(nombreArchivo, camion.Paquete2.Id, camion.Paquete2.Tipo, camion.Paquete2.Valor, camion.Paquete2.Origen, camion.Paquete2.Destino, camion.Paquete2.Intentos, "0")
	}
	fmt.Println(respuesta.GetBody())
}

func main() {
	var tEspera int
	var tEnvio int
	fmt.Printf("Ingrese el tiempo de espera de los camiones\n")
	fmt.Scanln(&tEspera)
	fmt.Printf("El tiempo de espera para tomar el segundo paquete es de %d segundos\n", tEspera)

	fmt.Printf("Ingrese el tiempo de envio de los paquetes\n")
	fmt.Scanln(&tEnvio)
	fmt.Printf("El tiempo de envío entre paquetes es de %d segundos\n", tEnvio)

    CamionR1 := Camion {
    	IdCamion: "R1",
		Tipo: "retail",
	}
	/*
	CamionR2 := Camion {
		IdCamion: "R2",
		Tipo: "retail",
	}
	CamionN := Camion{
		IdCamion: "N"
		Tipo: "normal",
	}*/

	crearRegistro("registroCamion" + CamionR1.IdCamion + ".csv")
	//crearRegistro("registroCamion" + CamionR2.IdCamion + ".csv")
	//crearRegistro("registroCamion" + CamionN.IdCamion + ".csv")

	for {
		Carga(CamionR1, tEspera, tEnvio, "registroCamionR1.csv")
		//fmt.Println(CamionR1.Paquete1.Seguimiento)
		//Carga(CamionR2, tEspera, tEnvio)
		//Carga(CamionN, tEspera, tEnvio)

	}

}
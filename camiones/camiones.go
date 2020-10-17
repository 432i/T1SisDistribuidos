package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"bufio"
	"context"
	"math/rand"
	"strconv"
	"github.com/432i/T1SisDistribuidos/logistica/chat"
	"google.golang.org/grpc"
)

type Camion struct {
	Tipo string
	Paquete1 *chat.Paquete
	Paquete2 *chat.Paquete
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
	fmt.Println("Debug")
	if paquete.Tipo == "retail" {
		fmt.Println("Debug")
		if intentos < 3 {
			fmt.Println("Debug")
			if rand.Float64() <= 0.8 {
				fmt.Println("Debug")
				paquete.Estado = "Recibido"
				fmt.Println("Debug")
			} else {
				fmt.Println("Debug2")
				intentos += 1
				paquete.Intentos = strconv.Itoa(intentos)
				fmt.Println("Debug2")
			}
		} else {
			fmt.Println("Debug3")
			paquete.Estado = "No Recibido"
		}
	} else {
		fmt.Println("Debug4")
		if intentos * 10 < valor && intentos < 2 {
			fmt.Println("Debug4")
			if rand.Float64() <= 0.8 {
				fmt.Println("Debug4")
				paquete.Estado = "Recibido"
			} else {
				fmt.Println("Debug5")
				intentos += 1
				paquete.Intentos = strconv.Itoa(intentos)
				fmt.Println("Debug5")
			}
		} else {
			fmt.Println("Debug6")
			paquete.Estado = "No Recibido"
			fmt.Println("Debug6")
		}
	}
}

func Entrega(camion Camion, tEnvio int) bool {
	fmt.Println("A")
	if camion.Paquete1.Valor > camion.Paquete2.Valor {
		fmt.Println("A")
		time.Sleep(time.Duration(tEnvio) * time.Second)
		fmt.Println("A")
		Intento(camion.Paquete1)
		fmt.Println("A")
	} else {
		fmt.Println("B")
		time.Sleep(time.Duration(tEnvio) * time.Second)
		fmt.Println("B")
		Intento(camion.Paquete2)
		fmt.Println("B")
	}
	fmt.Println("C")
	if camion.Paquete1.Estado != "En Camino" && camion.Paquete2.Estado != "En Camino"{
		fmt.Println("C")
		return false
	}
	fmt.Println("Debug4")
	return true
}

func Carga(camion Camion, tEspera int, tEnvio int) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("10.6.40.149:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)
	if err != nil {
		log.Fatalf("No se pudo generar comunicacion: %s", err)
	}

	mensaje := chat.Message{
		Body: camion.Tipo,
	}

	paquete1, _ := c.PaqueteQueueToCamion(context.Background(), &mensaje)
	if paquete1.GetId() != "" {
		camion.Paquete1 = paquete1
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
	}

	time.Sleep(time.Duration(tEspera) * time.Second)

	paquete2, _ := c.PaqueteQueueToCamion(context.Background(), &mensaje)
	if paquete2.GetId() != "" {
		camion.Paquete2 = paquete2
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
	}

	aux := true
	for aux {
		aux = Entrega(camion, tEnvio)
	}

	//PaqueteCamionToQueue(context.Background(), &camion.paquete1)
	//PaqueteCamionToQueue(context.Background(), &camion.paquete2)

}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Ingrese el tiempo de espera de los camiones\n")
	tEspera1, _ := reader.ReadString('\n')
	fmt.Printf("El tiempo de espera para tomar el segundo paquete es de %s segundos\n", tEspera1)

	fmt.Printf("Ingrese el tiempo de envio de los paquetes\n")
	tEnvio1, _ := reader.ReadString('\n')
	fmt.Printf("El tiempo de envío entre paquetes es de %s segundos\n", tEnvio1)

	tEspera, _ := strconv.Atoi(tEspera1)
	tEnvio, _ := strconv.Atoi(tEnvio1)

    CamionR1 := Camion {
		Tipo: "retail",
	}
	/*
	CamionR2 := Camion {
		Tipo: "retail",
	}
	CamionN := Camion{
		Tipo: "normal",
	}*/

	for {
		Carga(CamionR1, tEspera, tEnvio)
		fmt.Println(CamionR1.Paquete1.Seguimiento)
		//Carga(CamionR2, tEspera, tEnvio)
		//Carga(CamionN, tEspera, tEnvio)

	}

}
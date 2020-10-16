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

func Intento(paquete chat.Paquete) {
	intentos, _ = strconv.Atoi(paquete.Intentos)
	valor, _ = strconv.Atoi(paquete.Valor)
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

func Entrega(camion *Camion, tEnvio int64) bool {
	if camion.Paquete1.Valor > camion.Paquete2.Valor {
		time.Sleep(tEnvio * time.Second)
		Intento(camion.Paquete1)
	} else {
		time.Sleep(tEnvio * time.Second)
		Intento(camion.Paquete2)
	}
	if camion.Paquete1.Estado != "En Camino" && camion.Paquete2.Estado != "En Camino"{
		return false
	}
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
		camion.Paquete1 = chat.Paquete{
			Id:       paquete1.GetId(),
			Track:    paquete1.GetTrack(),
			Tipo:     paquete1.GetTipo(),
			Intentos: paquete1.GetIntentos(),
			Estado:   paquete1.GetEstado(),
		}
		msj = chat.Message{
			Body: camion.Paquete1.GetTrack() + ",En Camino",
		}
		respuesta, _ = c.EstadoPaquete(context.Background(), &msj)
		fmt.Println(respuesta)
	}
	fmt.Println(camion.Paquete1.Id)
	time.Sleep(tEspera * time.Second)

	paquete2, _ := c.PaqueteQueueToCamion(context.Background(), &mensaje)
	if paquete2.GetId() != "" {
		camion.Paquete2 = chat.Paquete{
			Id:       paquete2.GetId(),
			Track:    paquete2.GetTrack(),
			Tipo:     paquete2.GetTipo(),
			Intentos: paquete2.GetIntentos(),
			Estado:   paquete2.GetEstado(),
		}
		msj = chat.Message{
			Body: camion.Paquete2.GetTrack() + ",En Camino",
		}
		respuesta, _ = c.EstadoPaquete(context.Background(), &msj)
		fmt.Println(respuesta)
	}

	aux = true
	for aux {
		aux = Entrega(*camion, tEnvio)
	}

	PaqueteCamionToQueue(context.Background(), &camion.paquete1)
	PaqueteCamionToQueue(context.Background(), &camion.paquete2)

}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Ingrese el tiempo de espera de los camiones\n")
	tEspera, _ := reader.ReadString('\n')
	fmt.Println("El tiempo de espera para tomar el segundo paquete es de %s segundos", tEspera)

	fmt.Println("Ingrese el tiempo de envio de los paquetes\n")
	tEnvio, _ := reader.ReadString('\n')
	fmt.Println("El tiempo de envío entre paquetes es de %s segundos", tEnvio)

	tEspera, _ = strconv.Atoi(tEspera)
	tEnvio, _ = strconv.Atoi(tEnvio)

    CamionR1 := Camion {
		Tipo: "retail",
	}
	CamionR2 := Camion {
		Tipo: "retail",
	}
	CamionN := Camion{
		Tipo: "normal",
	}

	for {
		Carga(CamionR1, tEspera, tEnvio)
		//Carga(CamionR2, tEspera, tEnvio)
		//Carga(CamionN, tEspera, tEnvio)

	}

}
package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"bufio"
	"context"
	"github.com/432i/T1SisDistribuidos/logistica/chat"
	"google.golang.org/grpc"
)

type Camion struct {
	Tipo string
	Paquete1 chat.Paquete
	Paquete2 chat.Paquete
}



func Carga(camion Camion) {
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
	paquete2, _ := c.PaqueteQueueToCamion(context.Background(), &mensaje)

	camion.Paquete1 = chat.Paquete{
		Id:       paquete1.GetId(),
		Track:    paquete1.GetTrack(),
		Tipo:     paquete1.GetTipo(),
		Intentos: paquete1.GetIntentos(),
		Estado:   paquete1.GetEstado(),
	}
	camion.Paquete2 = chat.Paquete{
		Id:       paquete2.GetId(),
		Track:    paquete2.GetTrack(),
		Tipo:     paquete2.GetTipo(),
		Intentos: paquete2.GetIntentos(),
		Estado:   paquete2.GetEstado(),
	}

	//Llamar a funcion que entrega los paquetes

}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Ingrese el tiempo de espera de los camiones\n")
	tEspera, _ := strconv.Atoi(reader.ReadString('\n'))

	fmt.Println("Ingrese el tiempo de envio de los paquetes\n")
	tEnvio, _ := strconv.Atoi(reader.ReadString('\n'))

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
		time.Sleep(tEspera * time.Second)
		go Carga(CamionR1)
		go Carga(CamionR2)
		go Carga(CamionN)

	}

}
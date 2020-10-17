package main

import (
	"log"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"strconv"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
	
type Paquete struct {
	Id string `json:"id"`
	Tipo string `json:"tipo"`
	Valor string `json:"valor"`
	Intentos string `json:"intentos"`
	Estado string `json:"estado"`
}



var gastos int
var ingresos int

func conexion(){
	//Se establece conexion a rabbit con usuario e ip del servidor
	conn, err := amqp.Dial("amqp://finanzas:finanzas@10.6.40.150:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "hello", // name
        false,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare a queue")

    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    failOnError(err, "Failed to register a consumer")

    forever := make(chan bool)

    go func() {
        for d := range msgs {
			pakete := Paquete{}
			json.Unmarshal([]byte(d.Body), &pakete)
			//intentos, _ := strconv.Atoi(pakete.Intentos)
			fmt.Println(pakete)

        }
    }()

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}

func main() {
	var respuesta int
	go conexion()
	for{
		fmt.Println("Ingrese 432 y preione enter para salir del sistema y mostrar balance financiero \n")
		fmt.Scanln(&respuesta)
		if respuesta == 432{
			fmt.Println("\n ---------------------- \n")
			fmt.Println("\n BALANCE FINANCIERO ")
			fmt.Println("\n Ganancias: ")
			fmt.Println(ingresos)
			fmt.Println("\n Gastos: ")
			fmt.Println(gastos)
			fmt.Println("\n Total: ")
			fmt.Println(ingresos-gastos)
			fmt.Println("\n ---------------------- \n")
			break
		}
	}

}

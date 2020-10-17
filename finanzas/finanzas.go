package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
	
type registro struct {
    enviosCompletados int32
    cantIntentosPaquete int32
    paquetesNoEntregados int32
    perdidasOGananciasPaquete float64
}

func main() {
	//Se establece conexion a rabbit con usuario e ip del servidor
	conn, err := amqp.Dial("amqp://finanzas:finanzas@10.6.40.150:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//Se crea un canal para la comunicacion
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		"TestQueue", // queue
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
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

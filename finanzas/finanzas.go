package main

import (
	"log"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"strconv"
	"os"
	"encoding/csv"
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

func crearRegistro(){
	archivo, err := os.Create("paquetes.csv")
	if err != nil{
			log.Println(err)
	}
	archivo.Close()
}

func guardarPaquete(estado string, intentos string, valor int){ 
	//valor es la perdida o la ganancia del pakete
	//estado de Recibido indica que fue entregado (envio completado), No recibido que no se pudo entregar
	
	orden := []string{estado, intentos, strconv.Itoa(valor)}
	archivo, err := os.OpenFile("paquetes.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(archivo)

	w.Write(orden)

	w.Flush()
	archivo.Close()
}

var gastos float64
var ingresos float64

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
			fmt.Println(pakete)
			intentos, _ := strconv.Atoi(pakete.Intentos)
			valor, _ := strconv.Atoi(pakete.Valor)
			if pakete.Tipo == "retail"{
				if pakete.Estado == "No Recibido"{
					ingresos += valor
					gastos += intentos*10
					guardarPaquete(pakete.Estado, pakete.Intentos, valor-intentos*10)
				}
				if pakete.Estado == "Recibido"{
					ingresos += valor
					gastos += intentos*10
					guardarPaquete(pakete.Estado, pakete.Intentos, valor-intentos*10)
				}

			}
			if pakete.Tipo == "normal"{
				if pakete.Estado == "No Recibido"{
					gastos += intentos*10
					guardarPaquete(pakete.Estado, pakete.Intentos, -1*intentos*10)
				}
				if pakete.Estado == "Recibido"{
					ingresos += valor
					gastos += intentos*10
					guardarPaquete(pakete.Estado, pakete.Intentos, valor-intentos*10)
				}
			}	
			if pakete.Tipo == "prioritario"{
				if pakete.Estado == "No Recibido"{
					ingresos += valor*0.3
					gastos += intentos*10
					guardarPaquete(pakete.Estado, pakete.Intentos, valor*0.3-intentos*10)
				}
				if pakete.Estado == "Recibido"{
					ingresos += valor
					gastos += intentos*10
					guardarPaquete(pakete.Estado, pakete.Intentos, valor-intentos*10)
				}
			}
			

        }
    }()

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}

func main() {
	crearRegistro()
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

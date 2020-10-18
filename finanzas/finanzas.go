package main

import (
	"log"
	"encoding/json"
	"fmt"
	"strconv"
	"os"
	"encoding/csv"
	"github.com/streadway/amqp"
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
/*
Funcion: crearRegistro
Parametros:
	- Ninguno
Descripcion:
	- Crea el registro de paquetes paquetes.csv
Retorno:
	- No tiene retorno
*/
func crearRegistro(){
	archivo, err := os.Create("paquetes.csv")
	if err != nil{
			log.Println(err)
	}
	archivo.Close()
}
/*
Funcion: guardarPaquete
Parametros:
	- estado: string que indica el estado del paquete (recibido, no recibido)
	- intentos: string que indica el número de intentos que se realizaron para entregar el paquete
	- valor: el costo o la ganancia que tiene el paquete para la empresa
Descripcion:
	- Guarda paquetes con los campos que recibe en paquetes.csv
Retorno:
	- No tiene retorno
*/
func guardarPaquete(estado string, intentos string, valor float64){ 
	//valor es la perdida o la ganancia del pakete
	//estado de Recibido indica que fue entregado (envio completado), No recibido que no se pudo entregar
	
	orden := []string{estado, intentos, strconv.FormatFloat(valor, 'f', -1, 64)}
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
/*
Funcion: conexion
Parametros:
	- Ninguno
Descripcion:
	- Realiza la conexión con rabbitMQ y saca paquetes que recibe desde la cola para procesarlos
Retorno:
	- No tiene retorno
*/
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
			a, _ := strconv.Atoi(pakete.Intentos)
			intentos := float64(a)
			b, _ := strconv.Atoi(pakete.Valor)
			valor := float64(b)
			if pakete.Tipo == "retail"{
				if pakete.Estado == "No Recibido"{
					ingresos += valor
					gastos += intentos*10.0
					guardarPaquete(pakete.Estado, pakete.Intentos, valor-intentos*10.0)
				}
				if pakete.Estado == "Recibido"{
					ingresos += valor
					gastos += intentos*10.0
					guardarPaquete(pakete.Estado, pakete.Intentos, valor-intentos*10.0)
				}

			}
			if pakete.Tipo == "normal"{
				if pakete.Estado == "No Recibido"{
					gastos += intentos*10.0
					guardarPaquete(pakete.Estado, pakete.Intentos, -1.0*intentos*10.0)
				}
				if pakete.Estado == "Recibido"{
					ingresos += valor
					gastos += intentos*10.0
					guardarPaquete(pakete.Estado, pakete.Intentos, valor-intentos*10.0)
				}
			}	
			if pakete.Tipo == "prioritario"{
				if pakete.Estado == "No Recibido"{
					ingresos += valor*0.3
					gastos += intentos*10.0
					guardarPaquete(pakete.Estado, pakete.Intentos, valor*0.3-intentos*10.0)
				}
				if pakete.Estado == "Recibido"{
					ingresos += valor
					gastos += intentos*10.0
					guardarPaquete(pakete.Estado, pakete.Intentos, valor-intentos*10.0)
				}
			}
			

        }
    }()

    log.Printf(" [*] Esperando paquetes.")
    <-forever
}

func main() {
	crearRegistro()
	var respuesta int
	go conexion()
	for{
		fmt.Println("Ingrese 432 y presione enter para salir del sistema y mostrar balance financiero \n")
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

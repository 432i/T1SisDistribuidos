package chat

import (
        "os"
        "fmt"
        "log"
        "time"
        "golang.org/x/net/context"
        "encoding/csv"
        "strings"
        "strconv"
        "github.com/streadway/amqp"
)

type Server struct {
        todos_paquetes []Paquete //registro en memoria de todos los paquetes, para consultar su estado y demás
        cola_ret_a_camion []Paquete
        cola_prio_a_camion []Paquete
        cola_norm_a_camion []Paquete
        cola_a_server []Paquete
        cola_a_finanzas []Paquete
}

/*
Funcion: guardarOrden
Parametros:
    - Todos los strings: Atributos de un registro
Descripcion:
    - Hace registro de un paquete para almacenarlo en el registro del servidor
Retorno:
    - No tiene retorno
*/
func guardarOrden(id string, producto string, valor string, tienda string, destino string, prioritario, codigo string){
        var tipof string
        tiempoactual := time.Now()
        timestamp := tiempoactual.Format("02-01-2006 15:04")
        if strings.Compare("0", prioritario) == 0{
                tipof = "normal"
        }
        if strings.Compare("1", prioritario) == 0{
                tipof = "prioritario"
        }
        if strings.Compare("2", prioritario) == 0{
                tipof = "retail"
        }

        orden := []string{timestamp, id, tipof, producto, valor, tienda, destino, codigo}
        archivo, err := os.OpenFile("registro.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(archivo)

	w.Write(orden)

	w.Flush()
	archivo.Close()
}

/*
Funcion: ModificarEstado
Parametros:
    - message: Contiene el nuevo estado de un paquete
Descripcion:
    - Modifica el estado de un paquete para mantener al cliente al tanto
Retorno:
    - Retorna un mensaje indicando el exito de este
*/
func (s *Server) ModificarEstado(ctx context.Context, message *Message) (*Message, error){
        estring := message.GetBody()
        l := strings.Split(estring, ",")
        codigoSeguimiento := l[0]
        nuevoEstado := l[1]
        //se modifica el estado del paquete 
        cont := 0 //para saber la posicion de la lista
        for _, pakete := range s.todos_paquetes{
                if strings.Compare(pakete.GetSeguimiento(), codigoSeguimiento) == 0{
                        nuevopakete := Paquete{
                                Id: pakete.GetId(),
                                Seguimiento: pakete.GetSeguimiento(),
                                Tipo: pakete.GetTipo(),
                                Valor: pakete.GetValor(),
                                Intentos: pakete.GetIntentos(),
                                Estado: nuevoEstado,
                                Origen: pakete.GetOrigen(),
                                Destino: pakete.GetDestino(),
                        }
                        s.todos_paquetes = append(s.todos_paquetes[:cont], s.todos_paquetes[cont+1:]...)
                        s.todos_paquetes = append(s.todos_paquetes, nuevopakete)
                        break

                }
                cont = cont +1
        }
        m := "Estado modificado exitosamente"
        msj := Message{
                Body: m,
        }
        return &msj, nil

}

/*
Funcion: SolicitarSeguimiento
Parametros:
    - message: Contiene el codigo de seguimiento
Descripcion:
    - Dado un codigo de seguimiento, busca cual es el estado actual del paquete
Retorno:
    - Retorna el estado del paquete, si este no existe, se indica que la orden no se encuentra en el sistema
*/
func (s *Server) SolicitarSeguimiento(ctx context.Context, message *Message) (*Message, error) {
        codigoSeguimiento := message.GetBody()
        var msj Message
        flag := 1
        for _, pakete := range s.todos_paquetes{
                if strings.Compare(pakete.GetSeguimiento(), codigoSeguimiento) == 0{
                        m := "El estado de su pedido "+codigoSeguimiento+" es "+pakete.GetEstado()
                        //buscar estado del pedido 
                        msj = Message{
                                Body: m,
                        }
                        flag = 0
                        break
                }
        }
        if flag == 1{
                m := "Orden no encontrada en el sistema, fijese bien porfa"
                        //buscar estado del pedido 
                msj = Message{
                        Body: m,
                }
        }
        return &msj, nil
}

/*
Funcion: EnviarOrden
Parametros:
    - orden: Orden realizada por un cliente
Descripcion:
    - Asigna un codigo de seguimiento y genera un paquete dada la orden del cliente y lo ingresa al sistema
Retorno:
    - Retorna el codigo se seguimiento de la orden
*/
var idSeg = 1
func (s *Server) EnviarOrden(ctx context.Context, orden *Orden) (*Message, error) {
        codigoSeguimiento := strconv.Itoa(idSeg)+orden.GetId()
        cuerpo :="Codigo de seguimiento para su producto:   "+codigoSeguimiento

        //mensaje que se envia al cliente
        msj := Message{
                Body: cuerpo,
        }

        guardarOrden(orden.GetId(), orden.GetProducto(), orden.GetValor(), orden.GetTienda(), orden.GetDestino(), orden.GetPrioritario(), codigoSeguimiento)


        //colas segun el tipo, tipo -> retail prioritario normal  
        if strings.Compare(orden.GetPrioritario(), "0") == 0 {
                pakete := Paquete{
                        Id: orden.GetId(),
                        Seguimiento: codigoSeguimiento,
                        Tipo: "normal",
                        Valor: orden.GetValor(),
                        Intentos: "0",
                        Estado: "En bodega",
                        Origen: orden.GetTienda(),
                        Destino: orden.GetDestino(),
                }
                s.cola_norm_a_camion = append(s.cola_norm_a_camion, pakete)
                s.todos_paquetes = append(s.todos_paquetes, pakete)
        }
        if strings.Compare(orden.GetPrioritario(), "1") == 0 {
                pakete := Paquete{
                        Id: orden.GetId(),
                        Seguimiento: codigoSeguimiento,
                        Tipo: "prioritario",
                        Valor: orden.GetValor(),
                        Intentos: "0",
                        Estado: "En bodega",
                        Origen: orden.GetTienda(),
                        Destino: orden.GetDestino(),
                }
                s.cola_prio_a_camion = append(s.cola_prio_a_camion, pakete)
                s.todos_paquetes = append(s.todos_paquetes, pakete)
        }
        if strings.Compare(orden.GetPrioritario(), "2") == 0 {
                pakete := Paquete{
                        Id: orden.GetId(),
                        Seguimiento: codigoSeguimiento,
                        Tipo: "retail",
                        Valor: orden.GetValor(),
                        Intentos: "0",
                        Estado: "En bodega",
                        Origen: orden.GetTienda(),
                        Destino: orden.GetDestino(),
                }
                s.cola_ret_a_camion = append(s.cola_ret_a_camion, pakete)
                s.todos_paquetes = append(s.todos_paquetes, pakete)

        }
        idSeg++
        return &msj, nil
}

/*
Funcion: failOnError
Parametros:
    - err: Tipo de error
    - msg: Tipo de mensaje
Descripcion:
    - Manejo de errores, printea cuando se produce un error de conexion indicando el motivo (No deberia suceder)
Retorno:
    - No tiene retorno
*/
func failOnError(err error, msg string) {
        if err != nil {
            log.Fatalf("%s: %s", msg, err)
        }
}

/*
Funcion: PaquetesAFinanzas
Parametros:
    - pakete: Contiene informacion de paquete
Descripcion:
    - Establece conexion con el servidor de finanzas y envia un string en json con informacion de un paquete
Retorno:
    - No tiene retorno
*/
func PaquetesAFinanzas(pakete string){
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
        body := pakete
        err = ch.Publish(
                "",     // exchange
                q.Name, // routing key
                false,  // mandatory
                false,  // immediate
                amqp.Publishing{
                ContentType: "application/json",
                Body:        []byte(body),
        })
        log.Printf(" [x] Sent %s", body)
        failOnError(err, "Failed to publish a message")

}

/*
Funcion: PaqueteQueueToCamion
Parametros:
    - mensaje: string retail/prioritario/normal
Descripcion:
    - Dependiendo del tipo de mensaje (retail, prioritario, normal), toma un paquete de una de las colas y lo envía a camiones
Retorno:
    - Retorna un paquete de la cola seleccionada
*/
func (s *Server) PaqueteQueueToCamion(ctx context.Context, mensaje *Message) (*Paquete, error) {
        var msj Paquete

        if mensaje.GetBody() == "retail" {
                if len(s.cola_ret_a_camion) > 0 {
                        msj = Paquete {
                                Id: s.cola_ret_a_camion[0].GetId(),
                                Seguimiento: s.cola_ret_a_camion[0].GetSeguimiento(),
                                Tipo: s.cola_ret_a_camion[0].GetTipo(),
                                Valor: s.cola_ret_a_camion[0].GetValor(),
                                Intentos: s.cola_ret_a_camion[0].GetIntentos(),
                                Estado: s.cola_ret_a_camion[0].GetEstado(),
                                Origen: s.cola_ret_a_camion[0].GetOrigen(),
                                Destino: s.cola_ret_a_camion[0].GetDestino(),
                        }
                        if len(s.cola_ret_a_camion) == 1 {
                                s.cola_ret_a_camion = make([]Paquete, 0)
                        } else {
                                s.cola_ret_a_camion = s.cola_ret_a_camion[1:]
                        }
                } else if len(s.cola_prio_a_camion) > 0 {
                        msj = Paquete {
                                Id: s.cola_prio_a_camion[0].GetId(),
                                Seguimiento: s.cola_prio_a_camion[0].GetSeguimiento(),
                                Tipo: s.cola_prio_a_camion[0].GetTipo(),
                                Valor: s.cola_prio_a_camion[0].GetValor(),
                                Intentos: s.cola_prio_a_camion[0].GetIntentos(),
                                Estado: s.cola_prio_a_camion[0].GetEstado(),
                                Origen: s.cola_prio_a_camion[0].GetOrigen(),
                                Destino: s.cola_prio_a_camion[0].GetDestino(),
                        }
                        if len(s.cola_prio_a_camion) == 1 {
                                s.cola_prio_a_camion = make([]Paquete, 0)
                        } else {
                                s.cola_prio_a_camion = s.cola_prio_a_camion[1:]
                        }
                } else {
                        msj = Paquete{
                                Id: "",
                                Seguimiento: "",
                                Tipo: "",
                                Valor: "",
                                Intentos: "",
                                Estado: "",
                        }
                }
        } else {
                if len(s.cola_prio_a_camion) > 0 {
                        msj = Paquete {
                                Id: s.cola_prio_a_camion[0].GetId(),
                                Seguimiento: s.cola_prio_a_camion[0].GetSeguimiento(),
                                Tipo: s.cola_prio_a_camion[0].GetTipo(),
                                Valor: s.cola_prio_a_camion[0].GetValor(),
                                Intentos: s.cola_prio_a_camion[0].GetIntentos(),
                                Estado: s.cola_prio_a_camion[0].GetEstado(),
                                Origen: s.cola_prio_a_camion[0].GetOrigen(),
                                Destino: s.cola_prio_a_camion[0].GetDestino(),
                        }
                        if len(s.cola_prio_a_camion) == 1 {
                                s.cola_prio_a_camion = make([]Paquete, 0)
                        } else {
                                s.cola_prio_a_camion = s.cola_prio_a_camion[1:]
                        }
                } else if len(s.cola_norm_a_camion) > 0 {
                        msj = Paquete {
                                Id: s.cola_norm_a_camion[0].GetId(),
                                Seguimiento: s.cola_norm_a_camion[0].GetSeguimiento(),
                                Tipo: s.cola_norm_a_camion[0].GetTipo(),
                                Valor: s.cola_norm_a_camion[0].GetValor(),
                                Intentos: s.cola_norm_a_camion[0].GetIntentos(),
                                Estado: s.cola_norm_a_camion[0].GetEstado(),
                                Origen: s.cola_norm_a_camion[0].GetOrigen(),
                                Destino: s.cola_norm_a_camion[0].GetDestino(),
                        }
                        if len(s.cola_norm_a_camion) == 1 {
                                s.cola_norm_a_camion = make([]Paquete, 0)
                        } else {
                                s.cola_norm_a_camion = s.cola_norm_a_camion[1:]
                        }
                } else {
                        msj = Paquete{
                                Id: "",
                                Seguimiento: "",
                                Tipo: "",
                                Valor: "",
                                Intentos: "",
                                Estado: "",
                                Origen: "",
                                Destino: "",
                        }
                }
        }
        return &msj, nil
}

/*
Funcion: PaqueteCamionToQueue
Parametros:
    - paquete: Puntero a Paquete
Descripcion:
    - Toma el paquete del parametro y lo deja en las colas de retorno (servidor y finanzas)
Retorno:
    - Retorna un mensaje, el cual indica si el paquete se ingreso en las colas de servidor y finanzas
    - En caso contrario, se indica que no habia mensaje (Con esto se maneja que se haya dejado de enviar paquetes u ordenes)
*/
func (s *Server) PaqueteCamionToQueue(ctx context.Context, paquete *Paquete) (*Message, error) {
        var msj Message
        if paquete.Tipo != "" {
                s.cola_a_server = append(s.cola_a_server, Paquete {
                        Id: paquete.GetId(),
                        Seguimiento: paquete.GetSeguimiento(),
                        Tipo: paquete.GetTipo(),
                        Valor: paquete.GetValor(),
                        Intentos: paquete.GetIntentos(),
                        Estado: paquete.GetEstado(),
                        Origen: paquete.GetOrigen(),
                        Destino: paquete.GetDestino(),
                })
                s.cola_a_finanzas = append(s.cola_a_finanzas, Paquete {
                        Id: paquete.GetId(),
                        Seguimiento: paquete.GetSeguimiento(),
                        Tipo: paquete.GetTipo(),
                        Valor: paquete.GetValor(),
                        Intentos: paquete.GetIntentos(),
                        Estado: paquete.GetEstado(),
                        Origen: paquete.GetOrigen(),
                        Destino: paquete.GetDestino(),
                })
                msj = Message {Body: "El paquete ingreso a las colas de servidor y finanzas"}
                //sacamos el paquete y lo dejamos en json para mandarlo a finanzas
                pakete := s.cola_a_finanzas[0]
                //eliminamos el paquete de la cola
                if len(s.cola_a_finanzas) == 1 {
                        s.cola_a_finanzas = make([]Paquete, 0)
                } else {
                        s.cola_a_finanzas = s.cola_a_finanzas[1:]
                }
                //pasamos a json el pakete
                body := fmt.Sprintf(`{"id":"%s", "tipo":"%s", "valor":"%s", "intentos":"%s", "estado":"%s"}`, pakete.GetId(), pakete.GetTipo(), pakete.GetValor(), pakete.GetIntentos(), pakete.GetEstado())
                PaquetesAFinanzas(body)
        } else {
                msj = Message {Body: "No habia paquete"}
        }
        return &msj, nil
}
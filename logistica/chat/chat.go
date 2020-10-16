package chat

import (
        "os"
 //       "fmt"
        "log"
        "time"
        "golang.org/x/net/context"
        "encoding/csv"
        "strings"
)

type Server struct {
        todos_paquetes []Paquete //registro en memoria de todos los paquetes, para consultar su estado y demás
        cola_ret_a_camion []Paquete
        cola_prio_a_camion []Paquete
        cola_norm_a_camion []Paquete
        cola_ret_a_server []Paquete
        cola_prio_a_server []Paquete
        cola_norm_a_server []Paquete
}
//se guarda la orden en registro.csv
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

func (s *Server) SolicitarSeguimiento(ctx context.Context, message *Message) (*Message, error) {
        codigoSeguimiento := message.GetBody()
        m := "sorry todavia no está implementado ._.XD  "+codigoSeguimiento
        //buscar estado del pedido 
        msj := Message{
                Body: m,
        }


        return &msj, nil
}

func (s *Server) EnviarOrden(ctx context.Context, orden *Orden) (*Message, error) {
        codigoSeguimiento := "432"+orden.GetId()
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
                }
                s.cola_ret_a_camion = append(s.cola_ret_a_camion, pakete)
                s.todos_paquetes = append(s.todos_paquetes, pakete)

        }

        return &msj, nil
}

func (s *Server) PaqueteQueueToCamion(ctx context.Context, mensaje *Message) (*Paquete, error) {
        var msj Paquete

        if mensaje.GetBody() == "retail" {
                if len(s.cola_ret_a_camion) > 0 {
                        msj := Paquete{
                                Id: s.cola_ret_a_camion[0].GetId(),
                                Seguimiento: s.cola_ret_a_camion[0].GetSeguimiento(),
                                Tipo: s.cola_ret_a_camion[0].GetTipo(),
                                Valor: s.cola_ret_a_camion[0].GetValor(),
                                Intentos: s.cola_ret_a_camion[0].GetIntentos(),
                                Estado: s.cola_ret_a_camion[0].GetEstado(),
                        }
                        //se modifica el estado del paquete a En transito
                        cont := 0 //para saber la posicion de la lista
                        for _, pakete := range s.todos_paquetes{
                                if strings.Compare(pakete.GetSeguimiento(), s.cola_ret_a_camion[0].GetSeguimiento()) == 0{
                                        nuevopakete := Paquete{
                                                Id: pakete.GetId(),
                                                Seguimiento: pakete.GetSeguimiento(),
                                                Tipo: pakete.GetTipo(),
                                                Valor: pakete.GetValor(),
                                                Intentos: "0",
                                                Estado: "En Camino",
                                        }
                                        s.todos_paquetes = append(s.todos_paquetes[:cont], s.todos_paquetes[cont+1:]...)
                                        s.todos_paquetes = append(s.todos_paquetes, nuevopakete)
                                        cont = cont +1

                                }
                                cont = cont +1
                        }
                        s.cola_ret_a_camion = s.cola_ret_a_camion[1:]
                }
        } else {
                if len(s.cola_prio_a_camion) > 0 {
                        msj := Paquete{
                                Id: s.cola_prio_a_camion[0].GetId(),
                                Seguimiento: s.cola_prio_a_camion[0].GetSeguimiento(),
                                Tipo: s.cola_prio_a_camion[0].GetTipo(),
                                Valor: s.cola_prio_a_camion[0].GetValor(),
                                Intentos: s.cola_prio_a_camion[0].GetIntentos(),
                                Estado: s.cola_prio_a_camion[0].GetEstado(),
                        }
                        //se modifica el estado del paquete a En transito
                        cont := 0 //para saber la posicion de la lista
                        for _, pakete := range s.todos_paquetes{
                                if strings.Compare(pakete.GetSeguimiento(), s.cola_prio_a_camion[0].GetSeguimiento()) == 0{
                                        nuevopakete := Paquete{
                                                Id: pakete.GetId(),
                                                Seguimiento: pakete.GetSeguimiento(),
                                                Tipo: pakete.GetTipo(),
                                                Valor: pakete.GetValor(),
                                                Intentos: "0",
                                                Estado: "En Camino",
                                        }
                                        s.todos_paquetes = append(todos_paquetes[:cont], todos_paquetes[cont+1:]...)
                                        s.todos_paquetes = append(s.todos_paquetes, nuevopakete)
                                        cont = cont +1

                                }
                                cont = cont +1
                        }
                        s.cola_prio_a_camion = s.cola_prio_a_camion[1:]
                } else if len(s.cola_norm_a_camion) > 0 {
                        msj := Paquete{
                                Id: s.cola_norm_a_camion[0].GetId(),
                                Seguimiento: s.cola_norm_a_camion[0].GetSeguimiento(),
                                Tipo: s.cola_norm_a_camion[0].GetTipo(),
                                Valor: s.cola_norm_a_camion[0].GetValor(),
                                Intentos: s.cola_norm_a_camion[0].GetIntentos(),
                                Estado: s.cola_norm_a_camion[0].GetEstado(),
                        }
                        //se modifica el estado del paquete a En transito
                        cont := 0 //para saber la posicion de la lista
                        for _, pakete := range s.todos_paquetes{
                                if strings.Compare(pakete.GetSeguimiento(), s.cola_norm_a_camion[0].GetSeguimiento()) == 0{
                                        nuevopakete := Paquete{
                                                Id: pakete.GetId(),
                                                Seguimiento: pakete.GetSeguimiento(),
                                                Tipo: pakete.GetTipo(),
                                                Valor: pakete.GetValor(),
                                                Intentos: "0",
                                                Estado: "En Camino",
                                        }
                                        s.todos_paquetes = append(s.todos_paquetes[:cont], s.todos_paquetes[cont+1:]...)
                                        s.todos_paquetes = append(s.todos_paquetes, nuevopakete)
                                        cont = cont +1

                                }
                                cont = cont +1
                        }
                        s.cola_norm_a_camion = s.cola_norm_a_camion[1:]
                } else {
                        return &msj, nil
                }
                return &msj, nil
        }
}

func (s *Server) PaqueteCamionToQueue(ctx context.Context, paquete *Paquete) {
        if paquete.GetTipo() == "retail" {
                s.cola_ret_a_server = append(s.cola_ret_a_server, &paquete)
        } else if paquete.GetTipo() == "prioritario"{
                s.cola_prio_a_server = append(s.cola_prio_a_server, &paquete)
        } else {
                s.cola_norm_a_server = append(s.cola_norm_a_server, &paquete)
        }
}
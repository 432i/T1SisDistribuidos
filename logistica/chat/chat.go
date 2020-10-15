package chat

import (
        "os"
        "fmt"
        "log"
        "time"
        "golang.org/x/net/context"
        "encoding/csv"
        "strings"
)

type Server struct {
        cola_ret_a_camion []Paquete
        cola_prio_a_camion []Paquete
        cola_norm_a_camion []Paquete
}
//se guarda la orden en registro.csv
func guardarOrden(id string, producto string, valor string, tienda string, destino string, codigo string){
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
                Body: m
        }


        return &msj, nil
}

func (s *Server) EnviarOrden(ctx context.Context, orden *Orden) (*Message, error) {
        codigoSeguimiento := "432"+orden.GetId()
        cuerpo :="Codigo de seguimiento para su producto:   "+codigoSeguimiento
        msj := Message{
                Body: cuerpo,
        }

        guardarOrden(orden.GetId(), orden.GetProducto(), orden.GetValor(), orden.GetTienda(), orden.GetDestino(), orden.GetPrioritario(), codigoSeguimiento)
        //tambien se implementan las colas segun el tipo
        return &msj, nil
}

func (s *Server) PaqueteQueueToCamion(ctx context.Context, mensaje *Message) (*Paquete, error) {
        var msj Paquete

        if mensaje.GetBody() == "retail" {
                if len(s.cola_ret_a_camion) > 0 {
                        msj = Paquete{
                                Id: s.cola_ret_a_camion[0].GetId(),
                                Seguimiento: s.cola_ret_a_camion[0].GetSeguimiento(),
                                Tipo: s.cola_ret_a_camion[0].GetTipo(),
                                Valor: s.cola_ret_a_camion[0].GetValor(),
                                Intentos: s.cola_ret_a_camion[0].GetIntentos(),
                                Estado: s.cola_ret_a_camion[0].GetEstado(),
                        }
                        s.cola_ret_a_camion = s.cola_ret_a_camion[1:]
                }
        } else {
                if len(s.cola_prio_a_camion) > 0 {
                        msj = Paquete{
                                Id: s.cola_prio_a_camion[0].GetId(),
                                Seguimiento: s.cola_prio_a_camion[0].GetSeguimiento(),
                                Tipo: s.cola_prio_a_camion[0].GetTipo(),
                                Valor: s.cola_prio_a_camion[0].GetValor(),
                                Intentos: s.cola_prio_a_camion[0].GetIntentos(),
                                Estado: s.cola_prio_a_camion[0].GetEstado(),
                        }
                        s.cola_prio_a_camion = s.cola_prio_a_camion[1:]
                } else if len(s.cola_norm_a_camion) > 0 {
                        msj = Paquete{
                                Id: s.cola_norm_a_camion[0].GetId(),
                                Seguimiento: s.cola_norm_a_camion[0].GetSeguimiento(),
                                Tipo: s.cola_norm_a_camion[0].GetTipo(),
                                Valor: s.cola_norm_a_camion[0].GetValor(),
                                Intentos: s.cola_norm_a_camion[0].GetIntentos(),
                                Estado: s.cola_norm_a_camion[0].GetEstado(),
                        }
                        s.cola_norm_a_camion = s.cola_norm_a_camion[1:]
                } 
                else {
                        return &msj, nil
                }
        }
}
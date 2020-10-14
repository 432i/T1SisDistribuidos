package chat

import (
        "os"
        "fmt"
        "log"
        "time"
        "golang.org/x/net/context"
        "strings"
)

type Server struct {
}
//se guarda la orden en registro.csv
func guardarOrden(id string, producto string, valor string, tienda string, destino string, codigo string){
        tiempoactual := time.Now()
        timestamp := tiempoactual.Format("02-01-2006 15:04")
        if stringsCompare("0", prioritario) == 0{
                tipof := "normal"
        }
        if stringsCompare("1", prioritario) == 0{
                tipof := "prioritario"
        }
        if stringsCompare("2", prioritario) == 0{
                tipof := "retail"
        }

        orden := timestamp+","+ id+","+tipof +","+ producto+","+ valor+","+tienda +","+destino +","+codigo
        archivo, err := os.OpenFile("registro.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(archivo)

	w.Write(orden)

	w.Flush()
	archivo.Close()
}

//func (s *Server) SolicitarSeguimiento(ctx context.Context, codigoSeguimiento *Message) (*Message, error) {
//        log.Printf("Receive message body from client: %s", in.Body)
//        return &Message{Body: "Hello From the Server!"}, nil
//}

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

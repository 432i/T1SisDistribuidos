package chat

import (
        "golang.org/x/net/context"
)

type Server struct {
}

//func (s *Server) SolicitarSeguimiento(ctx context.Context, codigoSeguimiento *Message) (*Message, error) {
//        log.Printf("Receive message body from client: %s", in.Body)
//        return &Message{Body: "Hello From the Server!"}, nil
//}

func (s *Server) EnviarOrden(ctx context.Context, orden *Orden) (*Message, error) {
        codigoSeguimiento := "432"+orden.GetId()
        cuerpo :="Codigo de seguimiento para su producto: "+codigoSeguimiento
        msj := Message{
                Body: cuerpo,
        }
        //aqui se implementa que se guarde la info recibida de cliente en un archivo
        //tambien se implementan las colas segun el tipo
        return &msj, nil
}

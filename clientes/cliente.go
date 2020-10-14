package main
import(
        "os"
        "strings"
        "io"
        "encoding/csv"
        "log"
        "fmt"
        "encoding/json"
        "strconv"
        "golang.org/x/net/context"
        "google.golang.org/grpc"
        "github.com/432i/logistica/chat"
)


type Retail struct {
        tipo string 
        id string 
        producto string
        valor string
        tienda string 
        destino string 
type Pyme struct{
        tipo string 
        id string
        producto string 
        valor string
        tienda string 
        destino string 
        prioritario string 
}
func cargarCsv(){
        fmt.Println("saludos")
        csvpyme, _ := os.Open("pymes.csv")
        csvretail, _ := os.Open("retail.csv")

        readerpyme := csv.NewReader(csvpyme)
        readerretail := csv.NewReader(csvretail)

        var pedidospyme []Pyme
        var pedidosretail []Retail

        for {
                lineapyme, error := readerpyme.Read()
                if error == io.EOF {
                        break
                }else if error != nil{
                        log.Fatal(error)
                }

                pedidospyme = append(pedidospyme, Pyme{
                        tipo: "pyme",
                        id: lineapyme[0],
                        producto: lineapyme[1],
                        valor: lineapyme[2],
                        tienda: lineapyme[3],
                        destino: lineapyme[4],
                        prioritario: lineapyme[5],
                })
        }
        fmt.Println(pedidospyme)

        for {
                linearetail, error := readerretail.Read()
                if error == io.EOF {
                        break
                }else if error != nil{
                        log.Fatal(error)
                }
                pedidosretail = append(pedidosretail, Retail{
                        tipo: "retail",
                        id: linearetail[0],
                        producto: linearetail[1],
                        valor: linearetail[2],
                        tienda: linearetail[3],
                        destino: linearetail[4],
                })
        }
        fmt.Println(pedidosretail)
}

func main(){
        var conn *grpc.ClientConn
        conn, err := grpc.Dial("10.6.40.149:9000", grpc.WithInsecure())
        if err != nil {
                log.Fatalf("did not connect: %s", err)
        }
        defer conn.Close()

        c := chat.NewChatServiceClient(conn)

        //response, err := c.SayHello(context.Background(), &chat.Message{Body: "Hello From Client!"})
        //if err != nil {
        //        log.Fatalf("Error when calling SayHello: %s", err)
        //}
        //log.Printf("Response from server: %s", response.Body)
        for{
                var respuesta string
                fmt.Println("
                ⠀⠀⠀⣠⣶⡾⠏⠉⠙⠳⢦⡀⠀⠀⠀⢠⠞⠉⠙⠲   ⡀⠀
                ⠀⠀⠀⣴⠿⠏⠀⠀⠀⠀⠀⠀⢳⡀⠀⡏⠀⠀⠀⠀⠀     ⢷
                ⠀⠀⢠⣟⣋⡀⢀⣀⣀⡀⠀⣀⡀⣧⠀⢸⠀⠀⠀⠀⠀      ⡇
                ⠀⠀⢸⣯⡭⠁⠸⣛⣟⠆⡴⣻⡲⣿⠀⣸⠀BIENVENIDO ⡇
                ⠀⠀⣟⣿⡭⠀⠀⠀⠀⠀⢱⠀⠀⣿⠀⢹⠀⠀⠀⠀⠀       ⡇
                ⠀⠀⠙⢿⣯⠄⠀⠀⠀⢀⡀⠀⠀⡿⠀⠀⡇⠀⠀⠀⠀    ⡼
                ⠀⠀⠀⠀⠹⣶⠆⠀⠀⠀⠀⠀⡴⠃⠀⠀⠘⠤⣄⣠ ⠞⠀
                ⠀⠀⠀⠀⠀⢸⣷⡦⢤⡤⢤⣞⣁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
                ⠀⠀⢀⣤⣴⣿⣏⠁⠀⠀⠸⣏⢯⣷⣖⣦⡀⠀⠀⠀⠀⠀⠀
                ⢀⣾⣽⣿⣿⣿⣿⠛⢲⣶⣾⢉⡷⣿⣿⠵⣿⠀⠀⠀⠀⠀⠀
                ⣼⣿⠍⠉⣿⡭⠉⠙⢺⣇⣼⡏⠀⠀⠀⣄⢸⠀⠀⠀⠀⠀⠀
                ⣿⣿⣧⣀⣿………⣀⣰⣏⣘⣆⣀
                \n")
                fmt.Println("Ingrese la alternativa que desee: \n")
                fmt.Println("1 Enviar ordenes desde una Pyme")
                fmt.Println("2 Enviar ordenes desde el Retail")
                fmt.Println("3 Realizar seguimiento de un pedido")
                fmt.Println("432 para salir")
                _, err := fmt.Scanln(&respuesta)
                if err != nil {
                        fmt.Fprintln(os.Stderr, err)
                        return
                }
                fmt.Println("Tu respuesta fue:")
                fmt.Println(respuesta)
                if strings.Compare("1", respuesta) == 0{
                        fmt.Println("XD1")
                }
                if strings.Compare("2", respuesta) == 0{
                        fmt.Println("XD2")
                }
                if strings.Compare("3", respuesta) == 0{
                        fmt.Println("X3D")
                }
                if strings.Compare("432", respuesta) == 0{
                        fmt.Println("X432D")
                        break
                }
        }
}
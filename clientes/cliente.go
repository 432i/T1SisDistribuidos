package main
import(
        "os"
        "io"
        "encoding/csv"
        "log"
        "fmt"
        "encoding/json"
        "strconv"
)


type Retail struct {
        tipo string `json:"tipo"`
        id int `json:"id"`
        producto string `json:"producto"`
        valor int `jso:"valor"`
        tienda string `json:"tienda"`
        destino string `json:"destino"`
}
type Pyme struct{
        tipo string `json:"tipo"`
        id int `json:"id"`
        producto string `json:"producto"`
        valor int `json:"valor"`
        tienda string `json:"tienda"`
        destino string `json:"destino"`
        prioritario int `json:"prioritario"`
}


func main(){
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

                idcsv, error := strconv.Atoi(lineapyme[0])
                valorcsv, error := strconv.Atoi(lineapyme[2])
                prioritariocsv, error := strconv.Atoi(lineapyme[5])
                pedidospyme = append(pedidospyme, Pyme{
                        tipo: "pyme",
                        id: idcsv,
                        producto: lineapyme[1],
                        valor: valorcsv,
                        tienda: lineapyme[3],
                        destino: lineapyme[4],
                        prioritario: prioritariocsv,
                })
        }
        fmt.Println(pedidospyme)
        pymejson, _ := json.Marshal(pedidospyme)
        fmt.Println(pymejson)

        for {
                linearetail, error := readerretail.Read()
                if error == io.EOF {
                        break
                }else if error != nil{
                        log.Fatal(error)
                }
                idcsv, error := strconv.Atoi(linearetail[0])
                valorcsv, error := strconv.Atoi(linearetail[2])
                pedidosretail = append(pedidosretail, Retail{
                        tipo: "retail",
                        id: idcsv,
                        producto: linearetail[1],
                        valor: valorcsv,
                        tienda: linearetail[3],
                        destino: linearetail[4],
                })
        }
        fmt.Println(pedidosretail)
        retailjson, _ := json.Marshal(pedidosretail)
        fmt.Println(retailjson)

}
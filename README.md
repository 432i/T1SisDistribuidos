# T1SisDistribuidos
# Tarea 1 del ramo Sistema Distribuidos
## Comenzando 🚀

## 1. Entrar a la máquina correspondiente:
- Máquina 1 (dist09): logistica
- Máquina 2 (dist10): finanzas
- Máquina 3 (dist11): camiones
- Máquina 4 (dist12): clientes
## 2. Entrar a la carpeta T1SisDistribuidos y entrar a la carpeta correspondiente según la entidad de la máquina

- Para logistica: carpeta logistica (no chat)
- Para finanzas: carpeta finanzas
- Para camiones: carpeta camiones
- Para clientes: carpeta clientes


## 3. Escribir make y presionar enter en la consola para ejecutar el código
## Consideraciones (leer antes):
El archivo chat.go (dentro de chat) tiene toda la implementación de logística y logistica.go levanta el server del mismo


Debe ejecutarse el servidor de logistica primero, y antes de hacerlo se deben exportar variables, para esto escribir los siguientes comandos en consola:
- export GOROOT=/usr/local/go
- export GOPATH=$HOME/go
- export GOBIN=$GOPATH/bin
- export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN

Presionar enter y ejecutar el servidor de logística haciendo make

Asegurarse de que el firewall está desactivado o los métodos gRPC no funcionarán:

- service firewalld stop

## Autores ✒️

* **Ignacio Aedo, rol 201773556-2** - *Desarrollo* - [432i](https://github.com/432i)
* **Ethiel Carmona, rol 201773533-3** - *Desarrollo* - [ethielc](https://github.com/ethielc)

## Construido con 🛠️
* Lenguaje Go
* gRPC
* RabbitMQ
* Protocol Buffers

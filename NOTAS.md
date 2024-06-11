# OpenMarket

### Multi-Module workspaces monorepo

- Cada Microservicio dentro del repo quedará dentro de su propia carpeta
- Iniciar un modulo en cada carpeta con go mod init github.com/milluki/nombre-modulo
- Iniciar workspace con el comando go work init ./common ./gateway ./orders ./payment ./production ./stock
- El archivo que se crea go.work debería lucir

go 1.22

use (
./common
./gateway
./orders
./payment
./production
./stock
)

## Instrucciones

- Instalar requerimientos
- Posicionarse desde una terminal en la carpeta raíz
- Asegurarse de que docker se está ejecutando
- Ejecutar el comando docker compose up
- Ingresar a la ui de rabbit mq en localhost:5672 con user-pass guest

#### Referencias

- https://www.youtube.com/watch?v=KdnxzgSNLTU&ab_channel=Tiago

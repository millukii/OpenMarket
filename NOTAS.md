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

#### Referencias

- https://www.youtube.com/watch?v=KdnxzgSNLTU&ab_channel=Tiago

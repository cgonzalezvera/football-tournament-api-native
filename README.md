# REST API de Gestión de Torneo de Fútbol en Go

## 📋 Descripción

REST API construida con Go puro (sin frameworks) para gestionar torneos de fútbol. Implementa Clean Architecture y utiliza PostgreSQL como base de datos, todo dockerizado.

## 🏗️ Arquitectura del Proyecto

```
tournament-api/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point de la aplicación
├── internal/
│   ├── domain/
│   │   ├── player.go              # Entidad Player
│   │   ├── team.go                # Entidad Team
│   │   ├── tournament.go          # Entidad Tournament
│   │   └── match.go               # Entidad Match
│   ├── repository/
│   │   ├── player_repository.go   # Interface y implementación
│   │   ├── team_repository.go
│   │   ├── tournament_repository.go
│   │   └── match_repository.go
│   ├── usecase/
│   │   ├── player_usecase.go      # Lógica de negocio
│   │   ├── team_usecase.go
│   │   ├── tournament_usecase.go
│   │   └── match_usecase.go
│   └── handler/
│       ├── player_handler.go      # HTTP handlers
│       ├── team_handler.go
│       ├── tournament_handler.go
│       └── match_handler.go
├── pkg/
│   └── database/
│       └── postgres.go            # Conexión a PostgreSQL
├── migrations/
│   └── 001_initial_schema.sql     # Schema de BD
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Conceptos Clave para Desarrolladores C#

### Diferencias importantes entre Go y C#

| Concepto | C# | Go |
|----------|----|----|
| **Manejo de dependencias** | NuGet packages | Go modules (go.mod) |
| **Interfaces** | Explícitas (implements) | Implícitas (duck typing) |
| **Herencia** | Clases y herencia | Composición (embedding) |
| **Null** | null, Nullable<T> | nil, zero values |
| **Excepciones** | try/catch | error returns (múltiples valores) |
| **Constructores** | Constructor explícito | Funciones New() por convención |
| **Propiedades** | get/set | Campos públicos o métodos Get/Set |
| **Async/Await** | async/await | goroutines y channels |

### Convenciones Idiomáticas de Go

1. **Nombres de paquetes**: minúsculas, sin guiones bajos
2. **Exportación**: Mayúscula inicial = público, minúscula = privado
3. **Manejo de errores**: Siempre verificar explícitamente
4. **Interfaces pequeñas**: Preferir interfaces con 1-2 métodos
5. **Inicialización**: Usar funciones `New()` en lugar de constructores

## 📦 Requisitos Previos

- **Go 1.23+** (versión más reciente)
- **Docker** y **Docker Compose**
- **PostgreSQL** (gestionado por Docker)

## 🔧 Paso 1: Instalación de Go

### Windows
```bash
# Descargar desde https://go.dev/dl/
# Verificar instalación
go version
```

### Linux/Mac
```bash
wget https://go.dev/dl/go1.23.x.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.x.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

## 🏁 Paso 2: Inicializar el Proyecto

```bash
# Crear directorio del proyecto
mkdir tournament-api
cd tournament-api

# Inicializar módulo de Go (similar a crear un .csproj)
# Reemplaza "github.com/tuusuario" con tu usuario de GitHub
go mod init github.com/tuusuario/tournament-api

# Crear estructura de directorios
mkdir -p cmd/api
mkdir -p internal/domain
mkdir -p internal/repository
mkdir -p internal/usecase
mkdir -p internal/handler
mkdir -p pkg/database
mkdir -p migrations
```

**📝 Nota para C#**: `go mod init` es equivalente a crear un nuevo proyecto en Visual Studio. El archivo `go.mod` es como tu `.csproj`.

## 📥 Paso 3: Instalar Dependencias

```bash
# Driver de PostgreSQL (similar a Entity Framework Core para PostgreSQL)
go get github.com/lib/pq

# Generador de UUIDs
go get github.com/google/uuid

# Estas dependencias se agregarán automáticamente a go.mod
```

**📝 Nota para C#**: `go get` es equivalente a `dotnet add package` o usar NuGet Package Manager.

## 🗄️ Paso 4: Configurar Docker

Los archivos Docker ya están incluidos en los artifacts. Solo necesitas:

```bash
# Construir y levantar los contenedores
docker-compose up --build

# En otra terminal, aplicar migraciones
docker exec -i tournament-postgres psql -U tournament_user -d tournament_db < migrations/001_initial_schema.sql
```

## 🔨 Paso 5: Compilar y Ejecutar

### Desarrollo Local (sin Docker)

```bash
# Compilar el proyecto
go build -o bin/api cmd/api/main.go

# Ejecutar con variables de entorno
DB_HOST=localhost \
DB_PORT=5432 \
DB_USER=tournament_user \
DB_PASSWORD=tournament_pass \
DB_NAME=tournament_db \
API_PORT=8080 \
./bin/api
```

### Con Docker (Recomendado)

```bash
# Construir y ejecutar
docker-compose up --build

# Ver logs
docker-compose logs -f api

# Detener
docker-compose down
```

**📝 Nota para C#**: No hay equivalente directo a `dotnet run` en Go. Siempre debes compilar primero.

## 🧪 Paso 6: Probar la API

### Crear un Jugador (Player)

```bash
curl -X POST http://localhost:8080/api/players \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Lionel Messi",
    "date_birth": "1987-06-24T00:00:00Z"
  }'
```

### Crear un Equipo (Team)

```bash
curl -X POST http://localhost:8080/api/teams \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Argentina"
  }'
```

### Agregar Jugador a Equipo

```bash
curl -X POST http://localhost:8080/api/teams/{team_id}/players/{player_id}
```

### Crear un Torneo (Tournament)

```bash
curl -X POST http://localhost:8080/api/tournaments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Copa América 2024"
  }'
```

### Crear un Partido (Match)

```bash
curl -X POST http://localhost:8080/api/matches \
  -H "Content-Type: application/json" \
  -d '{
    "match_number": 1,
    "date": "2024-06-20T20:00:00Z",
    "team1_id": "uuid-del-equipo-1",
    "team2_id": "uuid-del-equipo-2",
    "goal_scored_team1": 2,
    "goal_scored_team2": 1
  }'
```

### Listar Todos los Jugadores

```bash
curl http://localhost:8080/api/players
```

### Obtener un Jugador por ID

```bash
curl http://localhost:8080/api/players/{player_id}
```

## 🔍 Comandos Útiles de Go

```bash
# Ver dependencias del proyecto
go list -m all

# Actualizar dependencias
go get -u ./...

# Limpiar módulos no utilizados
go mod tidy

# Ver documentación de un paquete
go doc net/http

# Formatear código (equivalente a Prettier/ReSharper)
go fmt ./...

# Analizar código en busca de problemas
go vet ./...

# Ejecutar tests (cuando los crees)
go test ./...

# Compilar para diferentes plataformas
GOOS=linux GOARCH=amd64 go build -o bin/api-linux cmd/api/main.go
GOOS=windows GOARCH=amd64 go build -o bin/api.exe cmd/api/main.go
```

## 📚 Conceptos de Clean Architecture Implementados

### 1. **Domain Layer** (`internal/domain/`)
- Entidades de negocio puras
- Sin dependencias externas
- Equivalente a tus "Domain Entities" en C#

### 2. **Repository Layer** (`internal/repository/`)
- Interfaces que definen contratos de acceso a datos
- Implementaciones concretas para PostgreSQL
- Equivalente a tus "Repositories" en C# con Entity Framework

### 3. **Use Case Layer** (`internal/usecase/`)
- Lógica de negocio
- Orquesta repositories
- Equivalente a tus "Services" o "Application Layer" en C#

### 4. **Handler Layer** (`internal/handler/`)
- Controladores HTTP
- Manejo de request/response
- Equivalente a tus "Controllers" en ASP.NET

## 🛠️ Manejo de Errores en Go

En Go, los errores se manejan como valores de retorno:

```go
// C# equivalente
// var player = repository.GetById(id);

// Go idiomático
player, err := repository.GetByID(id)
if err != nil {
    // Manejar el error
    return nil, err
}
// Usar player
```

**No hay excepciones** en Go. Todo error debe verificarse explícitamente.

## 🔐 Variables de Entorno

Configuración en `.env` o `docker-compose.yml`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=tournament_user
DB_PASSWORD=tournament_pass
DB_NAME=tournament_db
API_PORT=8080
```

## 📖 Recursos de Aprendizaje

1. **Tour Oficial de Go**: https://go.dev/tour/
2. **Effective Go**: https://go.dev/doc/effective_go
3. **Go by Example**: https://gobyexample.com/
4. **Standard Library**: https://pkg.go.dev/std

## 🐛 Debugging

### VS Code
1. Instalar extensión "Go"
2. Agregar configuración en `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch API",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/api",
            "env": {
                "DB_HOST": "localhost",
                "DB_PORT": "5432",
                "DB_USER": "tournament_user",
                "DB_PASSWORD": "tournament_pass",
                "DB_NAME": "tournament_db",
                "API_PORT": "8080"
            }
        }
    ]
}
```

## 🎯 Próximos Pasos

1. **Agregar validaciones**: Usar paquete `validator`
2. **Implementar tests**: `testing` package
3. **Agregar middleware**: Logging, CORS, Authentication
4. **Documentar API**: Swagger/OpenAPI
5. **Implementar paginación**: Para endpoints de listado
6. **Agregar CI/CD**: GitHub Actions, GitLab CI

## ❓ Preguntas Frecuentes (C# → Go)

**P: ¿Dónde está el equivalente a Entity Framework?**  
R: Go no tiene un ORM dominante. Se usa SQL directo con `database/sql` o librerías ligeras como `sqlx`.

**P: ¿Cómo manejo la inyección de dependencias?**  
R: Manualmente, pasando dependencias en constructores (`New()` functions).

**P: ¿Hay algo como LINQ?**  
R: No. Se usan loops explícitos. Es más verboso pero más claro.

**P: ¿Cómo hago async/await?**  
R: Con `goroutines` (go keyword) y `channels` para comunicación.

**P: ¿Hay generics?**  
R: Sí, desde Go 1.18+, pero con sintaxis diferente.

## 📝 Licencia

MIT License - Siéntete libre de usar este código como base para tus proyectos.

---

**¡Feliz Coding en Go! 🎉**
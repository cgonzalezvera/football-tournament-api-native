# 🔧 Guía de Troubleshooting y Tips Avanzados

## 🐛 Problemas Comunes y Soluciones

### 1. Error: "package X is not in GOROOT"

**Problema:**
```bash
package github.com/lib/pq is not in GOROOT (/usr/local/go/src/github.com/lib/pq)
```

**Solución:**
```bash
# Asegúrate de estar en el directorio del proyecto
cd tournament-api

# Descarga las dependencias
go mod download

# O reinstala
go get github.com/lib/pq

# Verifica que go.mod tenga las dependencias
cat go.mod
```

---

### 2. Error: "imported and not used"

**Problema:**
```go
import (
    "fmt"  // imported and not used: "fmt"
)
```

**Solución:**
Go es **muy estricto** con imports no usados. Debes:
- Eliminar el import si no lo usas
- O usa `goimports` para auto-limpiar:

```bash
# Instalar goimports
go install golang.org/x/tools/cmd/goimports@latest

# Auto-limpiar imports
goimports -w .
```

**📝 Nota para C#**: En C# los using no usados son warnings, en Go son errores de compilación.

---

### 3. Error: "cannot use X (type Y) as type Z"

**Problema:**
```go
var id string = "123"
uuid.Parse(id)  // Correcto

var id2 = 123
uuid.Parse(id2)  // Error: cannot use id2 (type int) as type string
```

**Solución:**
Go **NO tiene conversión implícita de tipos**. Debes convertir explícitamente:

```go
id2 := 123
uuid.Parse(strconv.Itoa(id2))  // Convertir int a string
```

---

### 4. Error: "panic: runtime error: invalid memory address"

**Problema:**
```go
var player *Player
fmt.Println(player.Name)  // PANIC: nil pointer dereference
```

**Solución:**
Siempre verifica punteros antes de usarlos:

```go
var player *Player
if player != nil {
    fmt.Println(player.Name)
} else {
    fmt.Println("Player is nil")
}
```

---

### 5. Docker: "Connection refused" al conectar a PostgreSQL

**Problema:**
```
Failed to connect to database: dial tcp 127.0.0.1:5432: connect: connection refused
```

**Soluciones:**

**A) Si usas Docker Compose:**
```bash
# Verifica que los contenedores estén corriendo
docker-compose ps

# Si no están levantados:
docker-compose up -d

# Revisa logs de postgres
docker-compose logs postgres
```

**B) Variables de entorno incorrectas:**
En Docker Compose, usa el nombre del servicio, NO `localhost`:
```yaml
environment:
  DB_HOST: postgres  # ❌ NO usar localhost
  DB_PORT: 5432
```

**C) Orden de inicio:**
Agrega `depends_on` con health check:
```yaml
api:
  depends_on:
    postgres:
      condition: service_healthy
```

---

### 6. Error: "http: multiple response.WriteHeader calls"

**Problema:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.WriteHeader(http.StatusCreated)  // ERROR: segunda llamada
}
```

**Solución:**
Solo puedes llamar `WriteHeader` UNA vez:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Decide el status code primero
    statusCode := http.StatusOK
    if someCondition {
        statusCode = http.StatusCreated
    }
    
    w.WriteHeader(statusCode)  // Solo una llamada
    w.Write([]byte("Response"))
}
```

---

### 7. JSON no se deserializa correctamente

**Problema:**
```go
type Player struct {
    id   uuid.UUID  // ❌ Campo privado (minúscula)
    name string     // ❌ Campo privado
}

// JSON: {"id": "...", "name": "Messi"}
// Resultado: campos vacíos
```

**Solución:**
Los campos deben ser **públicos** (mayúscula inicial):

```go
type Player struct {
    ID   uuid.UUID `json:"id"`    // ✅ Público
    Name string    `json:"name"`  // ✅ Público
}
```

---

### 8. Fechas en formato incorrecto

**Problema:**
```go
// Input JSON: "2024-01-15T10:30:00Z"
time.Parse("2006-01-02", dateStr)  // Error: parsing time
```

**Solución:**
Usa el formato correcto de Go (RFC3339 para ISO 8601):

```go
// ISO 8601 / RFC3339
dateTime, err := time.Parse(time.RFC3339, "2024-01-15T10:30:00Z")

// Solo fecha
date, err := time.Parse("2006-01-02", "2024-01-15")

// Fecha y hora personalizada
custom := "2006-01-02 15:04:05"
dt, err := time.Parse(custom, "2024-01-15 10:30:00")
```

**📝 Nota**: Go usa `2006-01-02 15:04:05` como formato de referencia (no es arbitrario, es mnemotécnico: 1/2 3:4:5 PM '06 -0700).

---

## 🎯 Tips de Performance

### 1. Reutiliza conexiones de DB

**❌ Malo:**
```go
func GetPlayer(id uuid.UUID) (*Player, error) {
    db, _ := sql.Open("postgres", connStr)  // Nueva conexión cada vez
    defer db.Close()
    // ...
}
```

**✅ Bueno:**
```go
// Crear conexión una sola vez en main()
db, _ := database.NewConnection(config)

// Reutilizar en repositorios
repo := NewPostgresPlayerRepository(db)  // Pasa la conexión
```

---

### 2. Usa Context para Timeouts

```go
import "context"

func (r *PostgresPlayerRepository) GetByID(ctx context.Context, id uuid.UUID) (*Player, error) {
    // Timeout de 5 segundos
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    query := `SELECT id, name FROM players WHERE id = $1`
    var player Player
    err := r.db.QueryRowContext(ctx, query, id).Scan(&player.ID, &player.Name)
    
    return &player, err
}
```

---

### 3. Pooling de Conexiones

```go
db.SetMaxOpenConns(25)        // Máximo de conexiones abiertas
db.SetMaxIdleConns(5)         // Conexiones en idle
db.SetConnMaxLifetime(5 * time.Minute)  // Tiempo de vida
```

---

### 4. Prepared Statements para queries repetitivas

```go
// Preparar una vez
stmt, err := db.Prepare(`SELECT id, name FROM players WHERE id = $1`)
defer stmt.Close()

// Reutilizar muchas veces
for _, id := range ids {
    var player Player
    stmt.QueryRow(id).Scan(&player.ID, &player.Name)
}
```

---

## 🔐 Tips de Seguridad

### 1. SQL Injection - Usa Prepared Statements

**❌ PELIGROSO:**
```go
query := fmt.Sprintf("SELECT * FROM players WHERE name = '%s'", name)
db.Query(query)  // Vulnerable a SQL injection
```

**✅ SEGURO:**
```go
query := `SELECT * FROM players WHERE name = $1`
db.Query(query, name)  // Parámetros seguros
```

---

### 2. Validación de UUIDs

```go
func (h *PlayerHandler) GetByID(w http.ResponseWriter, r *http.Request, idStr string) {
    id, err := uuid.Parse(idStr)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid UUID format")
        return
    }
    // Continuar...
}
```

---

### 3. CORS Configurado Correctamente

```go
func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ❌ Producción: NO uses "*"
        w.Header().Set("Access-Control-Allow-Origin", "*")
        
        // ✅ Producción: Especifica dominios
        // allowedOrigins := []string{"https://tuapp.com", "https://admin.tuapp.com"}
        // origin := r.Header.Get("Origin")
        // if contains(allowedOrigins, origin) {
        //     w.Header().Set("Access-Control-Allow-Origin", origin)
        // }
        
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

---

## 🧪 Tips de Testing

### 1. Tests con Subtests

```go
func TestPlayerUseCase(t *testing.T) {
    t.Run("GetPlayer returns player when exists", func(t *testing.T) {
        // Test específico
    })
    
    t.Run("GetPlayer returns error when not found", func(t *testing.T) {
        // Otro test
    })
}
```

### 2. Table-Driven Tests (idiomático en Go)

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "not-an-email", true},
        {"empty email", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## 📝 Tips de Código Limpio

### 1. Nombrado de Errores

```go
// ❌ Malo
var ErrSomethingWentWrong = errors.New("something went wrong")

// ✅ Bueno
var (
    ErrPlayerNotFound     = errors.New("player not found")
    ErrInvalidPlayerData  = errors.New("invalid player data")
    ErrDuplicatePlayer    = errors.New("player already exists")
)
```

### 2. Interfaces Pequeñas

```go
// ❌ Malo: Interfaz grande
type Repository interface {
    Create(entity interface{}) error
    Update(entity interface{}) error
    Delete(id string) error
    GetByID(id string) (interface{}, error)
    GetAll() ([]interface{}, error)
    FindBy(field, value string) ([]interface{}, error)
}

// ✅ Bueno: Interfaces específicas y pequeñas
type PlayerReader interface {
    GetByID(id uuid.UUID) (*Player, error)
    GetAll() ([]Player, error)
}

type PlayerWriter interface {
    Create(player *Player) error
    Update(player *Player) error
    Delete(id uuid.UUID) error
}
```

### 3. Early Returns

```go
// ❌ Malo: Anidación profunda
func ProcessPlayer(id uuid.UUID) error {
    player, err := repo.GetByID(id)
    if err == nil {
        if player != nil {
            if player.IsActive {
                if player.Age >= 18 {
                    // Lógica aquí
                    return nil
                } else {
                    return errors.New("too young")
                }
            } else {
                return errors.New("inactive")
            }
        } else {
            return errors.New("nil player")
        }
    } else {
        return err
    }
}

// ✅ Bueno: Early returns
func ProcessPlayer(id uuid.UUID) error {
    player, err := repo.GetByID(id)
    if err != nil {
        return err
    }
    
    if player == nil {
        return errors.New("nil player")
    }
    
    if !player.IsActive {
        return errors.New("inactive")
    }
    
    if player.Age < 18 {
        return errors.New("too young")
    }
    
    // Lógica aquí
    return nil
}
```

---

## 🔄 Tips de Migración desde C#

### 1. No busques reemplazos 1:1

**C# → Go**
- Entity Framework → SQL directo (es más común)
- LINQ → Loops (está bien, es idiomático)
- async/await → Síncrono (para APIs REST simples)

### 2. Abraza la simplicidad

```go
// En C# harías:
// var adults = players.Where(p => p.Age >= 18).ToList();

// En Go, está bien hacer:
var adults []Player
for _, p := range players {
    if p.Age >= 18 {
        adults = append(adults, p)
    }
}
```

### 3. Errores explícitos > Excepciones

```go
// No intentes hacer try/catch en Go
// Acepta que verificarás errores explícitamente

player, err := repo.GetByID(id)
if err != nil {
    log.Printf("Error getting player: %v", err)
    return nil, err
}
```

---

## 🚀 Comandos Útiles del Día a Día

```bash
# Formatear código (hazlo SIEMPRE antes de commit)
go fmt ./...

# Analizar código en busca de problemas
go vet ./...

# Ejecutar tests
go test ./...

# Tests con cobertura
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # Ver en navegador

# Ejecutar tests de un paquete específico
go test ./internal/usecase

# Ver documentación de un paquete
go doc net/http

# Actualizar dependencias
go get -u ./...
go mod tidy

# Compilar para diferentes plataformas
GOOS=linux GOARCH=amd64 go build -o bin/api-linux cmd/api/main.go
GOOS=windows GOARCH=amd64 go build -o bin/api.exe cmd/api/main.go
GOOS=darwin GOARCH=arm64 go build -o bin/api-mac cmd/api/main.go

# Ver dependencias del proyecto
go list -m all

# Ver información de un módulo
go mod graph

# Verificar módulos
go mod verify
```

---

## 📚 Recursos Recomendados

### Documentación Oficial
- **Tour de Go**: https://go.dev/tour/ (EMPIEZA AQUÍ)
- **Effective Go**: https://go.dev/doc/effective_go
- **Go by Example**: https://gobyexample.com/
- **Standard Library**: https://pkg.go.dev/std

### Libros
- "The Go Programming Language" (Donovan & Kernighan)
- "Let's Go" (Alex Edwards) - Específico para web development

### Blogs y Artículos
- Go Blog oficial: https://go.dev/blog/
- Dave Cheney's Blog: https://dave.cheney.net/
- Ardan Labs Blog: https://www.ardanlabs.com/blog/

### Videos
- JustForFunc (Francesc Campoy): YouTube
- GopherCon talks: YouTube

### Comunidad
- r/golang (Reddit)
- Gophers Slack: https://gophers.slack.com/
- Go Forum: https://forum.golangbridge.org/

---

## ⚡ Shortcuts de VS Code para Go

```json
// settings.json
{
  "go.useLanguageServer": true,
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

**Atajos útiles:**
- `Ctrl+Shift+O`: Ver outline de funciones
- `F12`: Go to definition
- `Shift+F12`: Find all references
- `Ctrl+.`: Quick fix / Import package

---

## 🎓 Ejercicios Recomendados

1. **Agrega Middleware de Logging**
    - Registra cada request (método, path, duración)
    - Tip: Usa `http.HandlerFunc` wrapper

2. **Implementa Paginación**
    - Agrega `?page=1&limit=10` a GetAll
    - Retorna metadatos (total, pages)

3. **Agrega Validación de Datos**
    - Usa tags de validación
    - Librería recomendada: `go-playground/validator`

4. **Implementa Tests Unitarios**
    - Crea mocks de repositorios
    - Usa table-driven tests

5. **Agrega Autenticación JWT**
    - Implementa middleware de auth
    - Librería: `golang-jwt/jwt`

6. **Dockeriza con CI/CD**
    - GitHub Actions para tests y build
    - Deploy a Cloud Run / Railway

---

**¡Buena suerte con tu journey en Go! 🚀**

Recuerda: Go es simple, pero no simplista. Abraza su filosofía y verás cómo tu código se vuelve más mantenible y claro.
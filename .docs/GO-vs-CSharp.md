# Guía Comparativa: Go vs C# para Desarrolladores

Esta guía te ayudará a entender las diferencias clave entre Go y C# mientras construyes APIs REST.

## 🏗️ Estructura de Proyecto

### C# (ASP.NET Core)
```
MyApi/
├── Controllers/
├── Models/
├── Services/
├── Data/
├── Program.cs
├── appsettings.json
└── MyApi.csproj
```

### Go
```
my-api/
├── cmd/api/main.go      # Entry point
├── internal/
│   ├── domain/         # Entities/Models
│   ├── handler/        # Controllers
│   ├── usecase/        # Services
│   └── repository/     # Data Access
├── pkg/                # Código reutilizable
├── go.mod              # Dependencias
└── go.sum              # Lock file
```

**📝 Diferencias clave:**
- Go usa `internal/` para código privado del proyecto
- `pkg/` es para código que puede ser importado externamente
- `cmd/` contiene entry points (pueden ser múltiples)

---

## 📦 Gestión de Dependencias

### C# - NuGet
```bash
dotnet add package Npgsql.EntityFrameworkCore.PostgreSQL
dotnet restore
```

### Go - Módulos
```bash
go get github.com/lib/pq
go mod tidy  # Limpia dependencias no usadas
```

**📝 Equivalencias:**
- `go.mod` = `.csproj`
- `go.sum` = `packages.lock.json`
- `go get` = `dotnet add package`

---

## 🎯 Tipos y Estructuras

### C# - Clases
```csharp
public class Player
{
    public Guid Id { get; set; }
    public string Name { get; set; }
    public DateTime DateBirth { get; set; }
    
    public Player(string name, DateTime dateBirth)
    {
        Id = Guid.NewGuid();
        Name = name;
        DateBirth = dateBirth;
    }
}
```

### Go - Structs
```go
type Player struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    DateBirth time.Time `json:"date_birth"`
}

func NewPlayer(name string, dateBirth time.Time) *Player {
    return &Player{
        ID:        uuid.New(),
        Name:      name,
        DateBirth: dateBirth,
    }
}
```

**📝 Diferencias clave:**
- Go no tiene clases, usa structs
- No hay constructores explícitos, se usan funciones `New*()`
- Tags de JSON se definen con backticks: `` `json:"name"` ``
- Campos en mayúscula son públicos, minúscula son privados

---

## 🔄 Interfaces

### C# - Explícitas
```csharp
public interface IPlayerRepository
{
    Task<Player> GetByIdAsync(Guid id);
    Task CreateAsync(Player player);
}

public class PlayerRepository : IPlayerRepository
{
    // Implementación explícita
    public async Task<Player> GetByIdAsync(Guid id) { }
    public async Task CreateAsync(Player player) { }
}
```

### Go - Implícitas (Duck Typing)
```go
type PlayerRepository interface {
    GetByID(id uuid.UUID) (*Player, error)
    Create(player *Player) error
}

type PostgresPlayerRepository struct {
    db *sql.DB
}

// Implementa la interfaz automáticamente
func (r *PostgresPlayerRepository) GetByID(id uuid.UUID) (*Player, error) { }
func (r *PostgresPlayerRepository) Create(player *Player) error { }
```

**📝 Diferencias clave:**
- Go no requiere declarar explícitamente que implementas una interfaz
- Si un tipo tiene todos los métodos de una interfaz, la implementa automáticamente
- Esto se llama "duck typing": "Si camina como un pato y hace cuac como un pato, es un pato"

---

## ❌ Manejo de Errores

### C# - Excepciones
```csharp
public async Task<Player> GetPlayerAsync(Guid id)
{
    try
    {
        var player = await _repository.GetByIdAsync(id);
        if (player == null)
            throw new NotFoundException("Player not found");
        return player;
    }
    catch (Exception ex)
    {
        _logger.LogError(ex, "Error getting player");
        throw;
    }
}
```

### Go - Valores de Retorno
```go
func (uc *PlayerUseCase) GetPlayer(id uuid.UUID) (*Player, error) {
    player, err := uc.repo.GetByID(id)
    if err != nil {
        // Manejo del error
        return nil, fmt.Errorf("error getting player: %w", err)
    }
    
    if player == nil {
        return nil, fmt.Errorf("player not found")
    }
    
    return player, nil
}
```

**📝 Diferencias clave:**
- Go NO tiene try/catch/finally
- Los errores son valores que se retornan
- Convención: el último valor de retorno es el error
- **SIEMPRE** debes verificar `if err != nil`
- `%w` en `fmt.Errorf` permite wrappear errores (similar a InnerException)

---

## 🔁 Múltiples Valores de Retorno

### C# - Tuplas o Out Parameters
```csharp
// Opción 1: Tuplas (C# 7+)
public (Player player, string error) GetPlayer(Guid id)
{
    if (playerNotFound)
        return (null, "Player not found");
    return (player, null);
}

// Opción 2: Out parameters
public bool TryGetPlayer(Guid id, out Player player)
{
    // ...
}
```

### Go - Nativo
```go
// Esto es idiomático en Go
func (r *Repository) GetByID(id uuid.UUID) (*Player, error) {
    // Retornar múltiples valores es natural en Go
    return player, nil
}

// Uso:
player, err := repo.GetByID(id)
if err != nil {
    // manejar error
}
// usar player
```

---

## 🌐 Controllers vs Handlers

### C# - Controller
```csharp
[ApiController]
[Route("api/[controller]")]
public class PlayersController : ControllerBase
{
    [HttpGet("{id}")]
    public async Task<ActionResult<Player>> GetPlayer(Guid id)
    {
        var player = await _useCase.GetPlayerAsync(id);
        if (player == null)
            return NotFound();
        return Ok(player);
    }
    
    [HttpPost]
    public async Task<ActionResult<Player>> CreatePlayer(CreatePlayerDto dto)
    {
        var player = await _useCase.CreatePlayerAsync(dto);
        return CreatedAtAction(nameof(GetPlayer), new { id = player.Id }, player);
    }
}
```

### Go - Handler con net/http
```go
type PlayerHandler struct {
    useCase *usecase.PlayerUseCase
}

func (h *PlayerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Routing manual
    switch r.Method {
    case http.MethodGet:
        h.GetPlayer(w, r)
    case http.MethodPost:
        h.CreatePlayer(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (h *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
    // Extraer ID de la URL manualmente
    id, err := uuid.Parse(extractID(r.URL.Path))
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid ID")
        return
    }
    
    player, err := h.useCase.GetPlayer(id)
    if err != nil {
        respondWithError(w, http.StatusNotFound, err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusOK, player)
}
```

**📝 Diferencias clave:**
- Go no tiene atributos de routing como `[HttpGet]`
- Debes implementar el routing manualmente
- `http.ResponseWriter` y `*http.Request` son los equivalentes a Response/Request
- No hay ActionResult, retornas void y escribes directamente a ResponseWriter

---

## 💾 Acceso a Datos

### C# - Entity Framework
```csharp
public class PlayerRepository : IPlayerRepository
{
    private readonly AppDbContext _context;
    
    public async Task<Player> GetByIdAsync(Guid id)
    {
        return await _context.Players
            .FirstOrDefaultAsync(p => p.Id == id);
    }
    
    public async Task CreateAsync(Player player)
    {
        _context.Players.Add(player);
        await _context.SaveChangesAsync();
    }
}
```

### Go - SQL Directo
```go
type PostgresPlayerRepository struct {
    db *sql.DB
}

func (r *PostgresPlayerRepository) GetByID(id uuid.UUID) (*Player, error) {
    query := `SELECT id, name, date_birth FROM players WHERE id = $1`
    
    var player Player
    err := r.db.QueryRow(query, id).Scan(
        &player.ID,
        &player.Name,
        &player.DateBirth,
    )
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("player not found")
    }
    if err != nil {
        return nil, err
    }
    
    return &player, nil
}

func (r *PostgresPlayerRepository) Create(player *Player) error {
    query := `INSERT INTO players (id, name, date_birth) VALUES ($1, $2, $3)`
    _, err := r.db.Exec(query, player.ID, player.Name, player.DateBirth)
    return err
}
```

**📝 Diferencias clave:**
- Go no tiene un ORM dominante como EF Core
- Usas SQL directo con `database/sql`
- `$1, $2` son placeholders (PostgreSQL), `?` para MySQL
- `Scan()` mapea los resultados a las variables (similar a AutoMapper manual)
- Alternativas: `sqlx`, `gorm` (más parecido a EF), pero menos comunes

---

## 🚀 Dependency Injection

### C# - Contenedor IoC Built-in
```csharp
// Program.cs
builder.Services.AddScoped<IPlayerRepository, PlayerRepository>();
builder.Services.AddScoped<PlayerService>();
builder.Services.AddControllers();

// Constructor injection automático
public class PlayersController : ControllerBase
{
    private readonly PlayerService _service;
    
    public PlayersController(PlayerService service)
    {
        _service = service;
    }
}
```

### Go - Manual (Constructor Injection)
```go
// main.go
func main() {
    db, _ := database.NewConnection(config)
    
    // Crear dependencias manualmente
    playerRepo := repository.NewPostgresPlayerRepository(db)
    playerUseCase := usecase.NewPlayerUseCase(playerRepo)
    playerHandler := handler.NewPlayerHandler(playerUseCase)
    
    // Registrar handlers
    http.Handle("/api/players", playerHandler)
    http.ListenAndServe(":8080", nil)
}

// Constructores reciben dependencias
func NewPlayerHandler(useCase *usecase.PlayerUseCase) *PlayerHandler {
    return &PlayerHandler{useCase: useCase}
}
```

**📝 Diferencias clave:**
- Go NO tiene DI container built-in
- Inyectas dependencias manualmente en constructores
- Esto hace el código más explícito y fácil de seguir
- Librerías como `wire` existen, pero son menos comunes

---

## 🔒 Null Safety

### C# 8+ - Nullable Reference Types
```csharp
public class Player
{
    public string Name { get; set; }      // Puede ser null (warning)
    public string? Bio { get; set; }      // Explícitamente nullable
    public int Age { get; set; }          // Value type, nunca null
    public int? OptionalAge { get; set; } // Nullable value type
}
```

### Go - Zero Values y Punteros
```go
type Player struct {
    Name string    // Zero value: "" (string vacío)
    Age  int       // Zero value: 0
    Bio  *string   // Puntero, puede ser nil
}

// Verificar nil
if player.Bio != nil {
    fmt.Println(*player.Bio) // Desreferenciar con *
}

// Crear valor nullable
bio := "Some bio"
player.Bio = &bio  // Obtener dirección con &
```

**📝 Diferencias clave:**
- Go no tiene `null`, tiene `nil` (solo para punteros, interfaces, slices, maps, channels)
- Tipos por valor tienen "zero values": `0` para int, `""` para string, `false` para bool
- Si necesitas distinguir "no establecido" vs "valor cero", usa punteros

---

## 🔄 Async/Await vs Goroutines

### C# - async/await
```csharp
public async Task<List<Player>> GetPlayersAsync()
{
    var players = await _repository.GetAllAsync();
    var enriched = await EnrichWithStatsAsync(players);
    return enriched;
}

public async Task ProcessMultipleAsync()
{
    var task1 = GetPlayersAsync();
    var task2 = GetTeamsAsync();
    
    await Task.WhenAll(task1, task2);
}
```

### Go - Goroutines y Channels
```go
func (uc *PlayerUseCase) GetPlayers() ([]Player, error) {
    // No hay async/await, las funciones son síncronas por defecto
    return uc.repo.GetAll()
}

func ProcessMultiple() error {
    // Canal para sincronización
    resultsChan := make(chan Result, 2)
    
    // Lanzar goroutines (concurrencia)
    go func() {
        players, err := getPlayers()
        resultsChan <- Result{players, err}
    }()
    
    go func() {
        teams, err := getTeams()
        resultsChan <- Result{teams, err}
    }()
    
    // Esperar resultados
    result1 := <-resultsChan
    result2 := <-resultsChan
    
    return nil
}
```

**📝 Diferencias clave:**
- Go no tiene palabras clave `async/await`
- Usa `go` para lanzar goroutines (hilos ligeros)
- Channels para comunicación entre goroutines
- Las goroutines son mucho más ligeras que threads (miles de goroutines = OK)
- Para I/O simple, Go es síncrono por defecto (más simple)

---

## 📝 JSON Serialization

### C# - Atributos y System.Text.Json
```csharp
public class Player
{
    [JsonPropertyName("id")]
    public Guid Id { get; set; }
    
    [JsonPropertyName("full_name")]
    public string FullName { get; set; }
    
    [JsonIgnore]
    public string InternalField { get; set; }
}

// Serializar
var json = JsonSerializer.Serialize(player);

// Deserializar
var player = JsonSerializer.Deserialize<Player>(json);
```

### Go - Struct Tags
```go
type Player struct {
    ID           uuid.UUID `json:"id"`
    FullName     string    `json:"full_name"`
    InternalField string    `json:"-"` // Ignorar
    Optional     *string   `json:"optional,omitempty"` // Omitir si nil
}

// Serializar
jsonBytes, err := json.Marshal(player)
if err != nil {
    // manejar error
}

// Deserializar
var player Player
err := json.Unmarshal(jsonBytes, &player)
if err != nil {
    // manejar error
}
```

**📝 Tags importantes:**
- `json:"name"` - Nombre del campo en JSON
- `json:"-"` - Ignorar campo
- `json:",omitempty"` - Omitir si zero value
- Los campos deben empezar con mayúscula para ser exportables

---

## 🎨 LINQ vs Loops

### C# - LINQ
```csharp
var adults = players
    .Where(p => p.Age >= 18)
    .OrderBy(p => p.Name)
    .Select(p => new PlayerDto 
    { 
        Name = p.Name,
        Age = p.Age 
    })
    .ToList();

var total = players.Sum(p => p.Goals);
```

### Go - Loops Explícitos
```go
// Filtrar adultos
var adults []Player
for _, p := range players {
    if p.Age >= 18 {
        adults = append(adults, p)
    }
}

// Ordenar
sort.Slice(adults, func(i, j int) bool {
    return adults[i].Name < adults[j].Name
})

// Mapear a DTOs
dtos := make([]PlayerDto, 0, len(adults))
for _, p := range adults {
    dtos = append(dtos, PlayerDto{
        Name: p.Name,
        Age:  p.Age,
    })
}

// Sumar
total := 0
for _, p := range players {
    total += p.Goals
}
```

**📝 Diferencias clave:**
- Go NO tiene LINQ
- Usas loops `for` explícitos
- `range` itera sobre slices, maps, channels
- Es más verboso pero más claro y performante
- Librerías como `go-funk` existen pero son poco comunes

---

## 🔧 Configuración

### C# - appsettings.json
```csharp
// appsettings.json
{
  "ConnectionStrings": {
    "DefaultConnection": "Host=localhost;Database=mydb"
  },
  "ApiSettings": {
    "Port": 8080
  }
}

// Program.cs
var connectionString = builder.Configuration.GetConnectionString("DefaultConnection");
var port = builder.Configuration["ApiSettings:Port"];
```

### Go - Variables de Entorno
```go
import "os"

func main() {
    dbHost := os.Getenv("DB_HOST")
    if dbHost == "" {
        dbHost = "localhost" // Valor por defecto
    }
    
    port := os.Getenv("API_PORT")
    if port == "" {
        port = "8080"
    }
    
    // O usar una librería como viper para archivos de configuración
}
```

**📝 Opciones:**
- Variables de entorno (lo más común)
- Archivos `.env` con librerías como `godotenv`
- Archivos JSON/YAML con `viper`
- Flags de línea de comandos con `flag` package

---

## 🧪 Testing

### C# - xUnit/NUnit
```csharp
public class PlayerServiceTests
{
    [Fact]
    public async Task GetPlayer_ReturnsPlayer_WhenExists()
    {
        // Arrange
        var mockRepo = new Mock<IPlayerRepository>();
        mockRepo.Setup(r => r.GetByIdAsync(It.IsAny<Guid>()))
                .ReturnsAsync(new Player { Name = "Test" });
        var service = new PlayerService(mockRepo.Object);
        
        // Act
        var result = await service.GetPlayerAsync(Guid.NewGuid());
        
        // Assert
        Assert.NotNull(result);
        Assert.Equal("Test", result.Name);
    }
}
```

### Go - testing Package
```go
func TestGetPlayer_ReturnsPlayer_WhenExists(t *testing.T) {
    // Arrange
    mockRepo := &MockPlayerRepository{
        GetByIDFunc: func(id uuid.UUID) (*Player, error) {
            return &Player{Name: "Test"}, nil
        },
    }
    useCase := NewPlayerUseCase(mockRepo)
    
    // Act
    result, err := useCase.GetPlayer(uuid.New())
    
    // Assert
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if result.Name != "Test" {
        t.Errorf("Expected name 'Test', got '%s'", result.Name)
    }
}

// Mock manual (o usar librería como testify)
type MockPlayerRepository struct {
    GetByIDFunc func(uuid.UUID) (*Player, error)
}

func (m *MockPlayerRepository) GetByID(id uuid.UUID) (*Player, error) {
    return m.GetByIDFunc(id)
}
```

**📝 Comandos:**
```bash
# Ejecutar tests
go test ./...

# Con cobertura
go test -cover ./...

# Verbose
go test -v ./...
```

---

## 📊 Resumen de Diferencias Fundamentales

| Aspecto | C# | Go |
|---------|----|----|
| **Paradigma** | OOP, functional | Procedural, concurrent |
| **Tipo de Tipado** | Estático, fuerte | Estático, fuerte |
| **Null** | `null`, `Nullable<T>` | `nil` (solo punteros), zero values |
| **Errores** | Excepciones (try/catch) | Valores de retorno (if err != nil) |
| **Herencia** | Sí (clases) | No (composición) |
| **Interfaces** | Explícitas | Implícitas |
| **Genéricos** | Sí (desde C# 2.0) | Sí (desde Go 1.18) |
| **Async** | async/await | goroutines + channels |
| **ORM** | Entity Framework | No dominante (SQL directo) |
| **DI** | Contenedor IoC built-in | Manual (constructores) |
| **Colecciones** | LINQ | Loops explícitos |
| **Compilación** | JIT + AOT | AOT (binario nativo) |
| **GC** | Generacional | Concurrent mark-sweep |

---

## 🎯 Conclusiones para Desarrolladores C#

**Lo que amarás de Go:**
- Simplicidad y minimalismo
- Binarios compilados súper rápidos
- Concurrencia fácil con goroutines
- Tooling excelente (go fmt, go vet, go test)
- Deployment trivial (un solo binario)

**Lo que extrañarás de C#:**
- LINQ
- Entity Framework
- async/await
- DI Container
- Generics más potentes
- Tooling de Visual Studio

**Consejos finales:**
1. ✅ Abraza la simplicidad de Go, no intentes replicar C#
2. ✅ Verifica errores explícitamente, siempre
3. ✅ Usa `go fmt` para formatear (no discutas estilo)
4. ✅ Lee "Effective Go" y "Go Proverbs"
5. ✅ Prefiere composición sobre herencia
6. ✅ Mantén tus interfaces pequeñas (1-2 métodos)
7. ✅ No tengas miedo de escribir código "verboso" y claro

**¡Disfruta aprendiendo Go! 🎉**
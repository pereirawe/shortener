# Shortener — Documentação de Contexto

> Arquivo de contexto para o assistente Antigravity.  
> Atualizado em: 2026-03-30

---

## Visão geral

API REST em **Go** para encurtamento de URLs.  
Módulo: `github.com/pereirawe/shortener`  
Entrypoint: `cmd/api/`

### Stack

| Componente | Tecnologia |
|---|---|
| Linguagem | Go 1.22+ |
| HTTP | `net/http` (stdlib, sem framework) |
| ORM / Banco | GORM + PostgreSQL |
| Cache | Redis (via `go-redis`) |
| Parser HTML (SEO) | `golang.org/x/net/html` |
| Infra local | Docker Compose |
| Testes | `testing` stdlib + `httptest` |
| Docs de API | Bruno (`/bruno`) |

---

## Estrutura de diretórios

```
shortener/
├── cmd/api/            # Entrypoint (main.go)
├── internal/
│   ├── api/
│   │   ├── handler.go       # HTTP handlers (ShortenURL, RedirectURL, Health)
│   │   ├── handler_test.go  # Testes unitários com mocks
│   │   └── seo.go           # Fetch e parse de metadados SEO
│   ├── config/         # Leitura de variáveis de ambiente
│   ├── dto/
│   │   └── url_dto.go  # Request / Response structs
│   ├── model/
│   │   └── url.go      # GORM model da tabela de URLs
│   └── service/
│       ├── postgres_service.go  # Implementação de URLStore (GORM)
│       └── redis_service.go     # Implementação de CacheStore (Redis)
├── bruno/              # Coleção Bruno de exemplos de requisição
├── docker-compose.yml
├── .env / .env.example
└── go.mod / go.sum
```

---

## Rotas

| Método | Path | Handler | Descrição |
|---|---|---|---|
| `POST` | `/api/shorten` | `ShortenURL` | Cria URL encurtada |
| `GET` | `/{shortCode}` | `RedirectURL` | Redireciona para a URL original |
| `GET` | `/health` | `Health` | Liveness check |

---

## `POST /api/shorten` — detalhes completos

### Request body

```json
{
  "original_url": "https://example.com/muito-longa",
  "custom_code": "MeuLink"   // opcional
}
```

### Campos do `custom_code`

| Regra | Comportamento |
|---|---|
| Apenas `[a-zA-Z0-9]` | `400 Bad Request` se violar |
| Sem espaços ou tabs | `400 Bad Request` se violar |
| Máximo **12 caracteres** | `400 Bad Request` se exceder |
| Case-sensitive | `"MeuLink"` ≠ `"meulink"` |
| Já em uso | `409 Conflict` |
| Não enviado | Short code aleatório de 7 chars é gerado |

### Response `201 Created`

```json
{
  "short_code":      "MeuLink",
  "short_url":       "http://localhost:8080/MeuLink",
  "original_url":    "https://example.com/muito-longa",
  "url_available":   true,
  "seo_title":       "Example Domain",
  "seo_description": "...",
  "seo_image":       "https://example.com/og.png",
  "warning":         ""   // ausente quando url_available=true
}
```

Quando `url_available: false` (URL inacessível no momento do encurtamento):

```json
{
  "url_available": false,
  "warning": "the destination URL appears to be temporarily unavailable; the short link was created anyway"
}
```

### Fluxo interno (`ShortenURL`)

```
1. Decode + validar body
2. Validar original_url (não vazio, prefixo http/https)
3. Se custom_code fornecido:
   a. Checar len ≤ 12
   b. Checar ausência de espaços
   c. Checar regex ^[a-zA-Z0-9]+$
   d. Checar unicidade no DB → 409 se já existe
4. Senão: gerar código aleatório (charset alfanumérico, 7 chars, até 10 tentativas)
5. fetchSEO(original_url) — timeout 5s, parseia <title>, <meta description>, og:*
6. Persistir no Postgres (com campos SEO)
7. Cachear no Redis (TTL 2 semanas, chave: "shortener:url:<shortCode>")
8. Retornar 201 com resp (+ warning se !available)
```

---

## Model — `model.URL`

```go
type URL struct {
    ID             uint           // PK
    CreatedAt      time.Time
    UpdatedAt      time.Time
    DeletedAt      gorm.DeletedAt // soft delete
    ShortCode      string         // UNIQUE, NOT NULL, size:20
    OriginalURL    string         // text
    Clicks         int64          // default:0
    URLAvailable   bool           // default:true
    SEOTitle       string         // text, omitempty
    SEODescription string         // text, omitempty
    SEOImage       string         // text, omitempty
}
```

> **Migração:** O GORM faz `AutoMigrate` na inicialização — as colunas novas são criadas automaticamente.

---

## Interfaces

### `URLStore`

```go
FindByShortCode(shortCode string) (*model.URL, error)
Create(url *model.URL) error
IncrementClicks(shortCode string)
ExistsByShortCode(shortCode string) (bool, error)
```

Implementação real: `service.PostgresService`

### `CacheStore`

```go
Get(ctx context.Context, key string) (string, error)
Set(ctx context.Context, key, value string, ttl time.Duration) error
```

Implementação real: `service.RedisService`  
Prefixo de chave: `shortener:url:<shortCode>`

---

## SEO (`api/seo.go`)

```go
type SEOData struct {
    Available   bool
    Title       string
    Description string
    Image       string
}

func fetchSEO(rawURL string) SEOData
```

- HTTP GET com timeout de **5 segundos**
- Segue até 5 redirecionamentos
- Parseia somente respostas `Content-Type: text/html`
- Extrai (em ordem de prioridade):
  - `Title`: `<og:title>` → `<title>`
  - `Description`: `<og:description>` → `<meta name="description">`
  - `Image`: `<og:image>`
- `Available = false` se: erro de rede, timeout, status ≥ 400

---

## Testes

```bash
go test ./...                  # roda todos os testes
go test ./internal/api/...     # apenas testes do handler
go test -v ./internal/api/...  # verbose
go test -run TestShortenURL ./internal/api/... # filtro por nome
```

### Cenários cobertos (`handler_test.go`)

| Teste | O que verifica |
|---|---|
| `TestHealth` | 200 no /health |
| `TestShortenURL_Success_AutoCode` | 201 + short_code gerado automaticamente |
| `TestShortenURL_CustomCode_Success` | 201 com short_code = custom_code |
| `TestShortenURL_CustomCode_MaxLength` | 12 chars — aceito |
| `TestShortenURL_CustomCode_TooLong` | 13 chars → 400 |
| `TestShortenURL_CustomCode_WithSpace` | espaço → 400 |
| `TestShortenURL_CustomCode_SpecialChar_Dash` | traço → 400 |
| `TestShortenURL_CustomCode_SpecialChar_At` | @ → 400 |
| `TestShortenURL_CustomCode_SpecialChar_Underscore` | underline → 400 |
| `TestShortenURL_CustomCode_Conflict` | código em uso → 409 |
| `TestShortenURL_EmptyURL` | URL vazia → 400 |
| `TestShortenURL_InvalidURL` | sem prefixo http → 400 |
| `TestShortenURL_InvalidJSON` | JSON inválido → 400 |
| `TestShortenURL_ResponseHasURLAvailableField` | campo `url_available` presente na resposta |
| `TestShortenURL_UnavailableURL_HasWarning` | URL inacessível → `url_available=false` + `warning` |
| `TestRedirectURL_FromDB` | 302 via banco |
| `TestRedirectURL_FromCache` | 302 via cache |
| `TestRedirectURL_NotFound` | 404 para código inexistente |
| `TestRedirectURL_CachePopulatedAfterDBHit` | cache re-populado após hit no banco |

---

## Variáveis de ambiente (`.env`)

```env
BASE_URL=http://localhost:8080

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=shortener
POSTGRES_PASSWORD=shortener
POSTGRES_DB=shortener

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
```

---

## Docker Compose

```bash
docker compose up -d    # sobe Postgres + Redis
docker compose down     # derruba
```

---

## Decisões de design

- **Sem framework HTTP** — `net/http` puro com pattern `METHOD /path` do Go 1.22
- **Interfaces injetadas no Handler** — facilita mocks nos testes unitários sem precisar de banco/redis real
- **SEO no momento do encurtamento** — evita latência no redirect; dado armazenado no banco para uso futuro
- **Warning em vez de erro** — URL inacessível não bloqueia a criação do link curto
- **AutoMigrate** — sem arquivos de migration separados; colunas novas são adicionadas automaticamente
- **Cache TTL de 2 semanas** — balanceia hit-rate com custo de memória Redis

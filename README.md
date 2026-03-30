# 🔗 Shortener

API de encurtamento de URLs com metadados SEO e tokens personalizados.

---

## Sumário

- [Funcionalidades](#funcionalidades)
- [Stack](#stack)
- [Estrutura do projeto](#estrutura-do-projeto)
- [Configuração e execução](#configuração-e-execução)
- [Rotas da API](#rotas-da-api)
- [Tokens personalizados](#tokens-personalizados)
- [Metadados SEO](#metadados-seo)
- [Testes](#testes)
- [Documentação Bruno](#documentação-bruno)

---

## Funcionalidades

- ✂️ **Encurtamento de URLs** com código aleatório (7 caracteres alfanuméricos)
- 🏷️ **Token personalizado** (`custom_code`) — defina seu próprio short code
- 🔍 **Validação de URL** — verifica disponibilidade antes de salvar
- 📊 **Metadados SEO** — coleta `title`, `description` e `og:image` da URL original
- ⚠️ **Warning de indisponibilidade** — cria o link mesmo se a URL estiver offline
- 📈 **Contador de cliques** — incremento atômico a cada redirect
- ⚡ **Cache Redis** — TTL de 2 semanas, fallback automático para o banco
- 🗄️ **Persistência PostgreSQL** via GORM com AutoMigrate

---

## Stack

| Componente | Tecnologia |
|---|---|
| Linguagem | Go 1.22+ |
| HTTP | `net/http` stdlib |
| ORM | GORM |
| Banco de dados | PostgreSQL 15 |
| Cache | Redis 7 |
| Parser HTML | `golang.org/x/net/html` |
| Infra local | Docker Compose |
| Testes | `testing` + `httptest` (sem dependências externas) |
| Docs de API | Bruno |

---

## Estrutura do projeto

```
shortener/
├── cmd/api/
│   └── main.go                  # Entrypoint
├── internal/
│   ├── api/
│   │   ├── handler.go           # HTTP handlers
│   │   ├── handler_test.go      # Testes unitários
│   │   └── seo.go               # Fetch e parse de metadados SEO
│   ├── config/
│   │   └── config.go            # Leitura de variáveis de ambiente
│   ├── dto/
│   │   └── url_dto.go           # Request / Response types
│   ├── model/
│   │   └── url.go               # GORM model
│   └── service/
│       ├── postgres_service.go  # Implementação URLStore
│       └── redis_service.go     # Implementação CacheStore
├── bruno/                       # Coleção de exemplos de requisição
├── .antigravity/
│   └── context.md               # Contexto do projeto para o assistente IA
├── docker-compose.yml
├── .env.example
└── go.mod
```

---

## Configuração e execução

### 1. Pré-requisitos

- Go 1.22+
- Docker e Docker Compose

### 2. Variáveis de ambiente

```bash
cp .env.example .env
```

`.env.example`:

```env
# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5433
POSTGRES_USER=shortener_user
POSTGRES_PASSWORD=shortener_password
POSTGRES_DB=shortener_db

# API
API_PORT=8080
API_BASE_URL=http://localhost:8080
```

### 3. Subir os serviços de infra

```bash
docker compose up -d
```

Isso sobe o PostgreSQL (porta `5433`) e o Redis (porta `6379`).

### 4. Rodar a API

```bash
go run ./cmd/api
```

Ou compilar e executar:

```bash
go build -o bin/shortener ./cmd/api && ./bin/shortener
```

A API estará disponível em `http://localhost:8080`.

---

## Rotas da API

### `GET /health`

Verifica se a API está no ar.

**Response `200`:**
```json
{ "status": "ok" }
```

---

### `POST /api/shorten`

Cria uma URL encurtada.

**Request body:**
```json
{
  "original_url": "https://example.com/pagina-muito-longa",
  "custom_code": "meulink"
}
```

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `original_url` | string | ✅ | URL completa (deve começar com `http://` ou `https://`) |
| `custom_code` | string | ❌ | Token personalizado (veja [Tokens personalizados](#tokens-personalizados)) |

**Response `201 Created`:**
```json
{
  "short_code": "meulink",
  "short_url": "http://localhost:8080/meulink",
  "original_url": "https://example.com/pagina-muito-longa",
  "url_available": true,
  "seo_title": "Example Domain",
  "seo_description": "Descrição da página",
  "seo_image": "https://example.com/og-image.png"
}
```

Se a URL original estiver inacessível, o link é criado mesmo assim e a resposta inclui um warning:

```json
{
  "url_available": false,
  "warning": "the destination URL appears to be temporarily unavailable; the short link was created anyway"
}
```

**Erros possíveis:**

| Status | Motivo |
|---|---|
| `400` | `original_url` ausente, vazia ou sem prefixo `http://`/`https://` |
| `400` | `custom_code` inválido (espaço, caractere especial, mais de 12 chars) |
| `409` | `custom_code` já está em uso |
| `500` | Erro interno |

---

### `GET /{shortCode}`

Redireciona para a URL original.

**Fluxo:**
1. Busca no cache Redis
2. Se não encontrar, busca no PostgreSQL e repopula o cache
3. Incrementa o contador de cliques de forma assíncrona

**Response `302 Found`** → redireciona para a URL original  
**Response `404 Not Found`** → short code não encontrado

---

## Tokens personalizados

Ao criar uma URL encurtada, você pode definir seu próprio código via `custom_code`:

```json
{
  "original_url": "https://github.com/pereirawe",
  "custom_code": "github"
}
```

**Regras de validação:**

| Regra | Exemplo inválido | Erro |
|---|---|---|
| Apenas letras `[a-zA-Z]` e dígitos `[0-9]` | `"my-link"`, `"email@"`, `"my_link"` | `400` |
| Sem espaços | `"my link"` | `400` |
| Máximo **12 caracteres** | `"abcdefghijklm"` (13) | `400` |
| Não pode já estar em uso | código duplicado | `409` |
| **Case-sensitive** | `"GitHub"` ≠ `"github"` | — |

Se `custom_code` não for enviado, um código aleatório de 7 caracteres é gerado automaticamente.

---

## Metadados SEO

Ao encurtar uma URL, a API faz uma requisição HTTP à URL original (timeout: **5 segundos**) para extrair:

| Campo | Fonte (prioridade) |
|---|---|
| `seo_title` | `<meta property="og:title">` → `<title>` |
| `seo_description` | `<meta property="og:description">` → `<meta name="description">` |
| `seo_image` | `<meta property="og:image">` |

Os metadados são armazenados no banco e retornados na resposta de criação.  
Se a URL não responder em 5s ou retornar status `>= 400`, a URL é marcada como `url_available: false` e um `warning` é incluído na resposta — mas o link é criado normalmente.

---

## Testes

Os testes são unitários, sem dependência de banco ou Redis reais. A lógica de SEO é testada com `httptest.Server` inline.

```bash
# Rodar todos os testes
go test ./...

# Apenas o pacote api (com output detalhado)
go test -v ./internal/api/...

# Filtrar por nome
go test -run TestShortenURL_CustomCode ./internal/api/...
```

**Cobertura dos testes:**

| Cenário | Teste |
|---|---|
| Health check | `TestHealth` |
| Criação com código automático + SEO | `TestShortenURL_Success_AutoCode` |
| Criação com custom_code válido | `TestShortenURL_CustomCode_Success` |
| Custom_code com exatamente 12 chars | `TestShortenURL_CustomCode_MaxLength` |
| Custom_code com 13 chars → 400 | `TestShortenURL_CustomCode_TooLong` |
| Custom_code com espaço → 400 | `TestShortenURL_CustomCode_WithSpace` |
| Custom_code com traço → 400 | `TestShortenURL_CustomCode_SpecialChar_Dash` |
| Custom_code com `@` → 400 | `TestShortenURL_CustomCode_SpecialChar_At` |
| Custom_code com `_` → 400 | `TestShortenURL_CustomCode_SpecialChar_Underscore` |
| Custom_code já em uso → 409 | `TestShortenURL_CustomCode_Conflict` |
| URL vazia → 400 | `TestShortenURL_EmptyURL` |
| URL sem prefixo → 400 | `TestShortenURL_InvalidURL` |
| JSON inválido → 400 | `TestShortenURL_InvalidJSON` |
| Campo `url_available` presente | `TestShortenURL_ResponseHasURLAvailableField` |
| URL 503 → warning + url_available false | `TestShortenURL_UnavailableURL_HasWarning` |
| SEO com og:tags | `TestShortenURL_SEO_OGTags` |
| Redirect via banco | `TestRedirectURL_FromDB` |
| Redirect via cache | `TestRedirectURL_FromCache` |
| Short code não encontrado → 404 | `TestRedirectURL_NotFound` |
| Cache repopulado após DB hit | `TestRedirectURL_CachePopulatedAfterDBHit` |

---

## Documentação Bruno

A coleção de exemplos está em `./bruno/`. Para usar:

1. Instale o [Bruno](https://www.usebruno.com/)
2. Abra a pasta `bruno/` como coleção
3. Configure o ambiente com `API_BASE_URL=http://localhost:8080`

Endpoints disponíveis:
- `POST /api/shorten` — Criar URL encurtada
- `GET /{shortCode}` — Redirecionar
- `GET /health` — Health check

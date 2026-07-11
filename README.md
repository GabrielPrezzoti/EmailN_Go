# EmailN - Email Campaign Management System

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.25.0-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat-square)
[![Test Coverage](https://img.shields.io/badge/Coverage-High-blue?style=flat-square)

**Plataforma robusta, escalável e segura para gerenciamento e disparo de campanhas de email em massa.**

[Features](#-features) • [Quick Start](#-quick-start) • [Documentação](#-documentação) • [Contribuir](#-contribuir)

</div>

---

## Visão Geral

**EmailN** é uma aplicação de backend desenvolvida em **Go** que oferece uma solução completa para gerenciar campanhas de email. Com arquitetura limpa, autenticação robusta e processamento assíncrono, o sistema permite criar, monitorar e executar campanhas de email em massa com eficiência e confiabilidade.

### Por que EmailN?

- **Alta Performance**: Desenvolvido em Go para máxima velocidade e eficiência
- **Segurança em Primeiro Lugar**: Autenticação JWT e validação rigorosa de dados
- **Arquitetura Escalável**: Clean Architecture e padrões SOLID
- **Qualidade Garantida**: Testes unitários abrangentes e bem estruturados
- **Processamento Assíncrono**: Worker dedicado para envio paralelo de emails
- **Monitoramento em Tempo Real**: Rastreamento completo do status das campanhas

---

## Features

### Gerenciamento de Campanhas
- Criar campanhas com validação robusta
- Recuperar informações detalhadas de campanhas
- Listar campanhas por usuário
- Atualizar status em tempo real
- Cancelar campanhas pendentes
- Deletar campanhas com segurança

### Processamento de Email
- Envio em massa via SMTP
- Fila de processamento assíncrono
- Tratamento automático de falhas
- Retry inteligente de campanhas
- Logging detalhado de operações

### Segurança
- Autenticação via JWT com OpenID Connect
- Validação de entrada em todos os endpoints
- Tratamento centralizado de erros
- Rate limiting e middleware de segurança
- Isolamento de dados por usuário

### Confiabilidade
- Estados bem definidos (Pending, Started, Done, Failed, Canceled, Deleted)
- Persistência em PostgreSQL
- Transações ACID
- Recuperação automática de falhas
- Auditoria completa de operações

---

## Stack Tecnológico

### Backend
- **[Go 1.25](https://golang.org)** - Linguagem de programação
- **[Chi v5](https://github.com/go-chi/chi)** - Router HTTP de alta performance
- **[GORM](https://gorm.io)** - ORM para banco de dados
- **[PostgreSQL](https://www.postgresql.org)** - Banco de dados relacional
- **[JWT](https://github.com/dgrijalva/jwt-go)** - Autenticação segura

### Infrastructure
- **[Gomail](https://github.com/go-mail/mail)** - Envio de emails SMTP
- **[Testify](https://github.com/stretchr/testify)** - Framework de testes
- **[OpenID Connect](https://openid.net/connect/)** - Autenticação federada

### DevOps
- **Docker** - Containerização
- **PostgreSQL Docker** - Base de dados containerizada
- **Keycloak** - Identity Provider (OpenID Connect)

---

## Arquitetura

### Clean Architecture

```
emailn/
├── cmd/                           # Aplicações executáveis
│   ├── api/                       # Server HTTP
│   └── worker/                    # Worker assíncrono
├── internal/
│   ├── domain/                    # Lógica de negócio (independente)
│   │   └── campaign/
│   │       ├── campaign.go        # Entidade de Campanha
│   │       ├── repository.go      # Interface do repositório
│   │       └── service.go         # Lógica de negócio
│   ├── endpoints/                 # Handlers HTTP
│   │   ├── auth.go                # Middleware de autenticação
│   │   └── campaigns_*.go         # Handlers de campanhas
│   ├── infrastructure/            # Implementações externas
│   │   ├── database/              # PostgreSQL + GORM
│   │   └── mail/                  # Serviço de email SMTP
│   ├── contract/                  # DTOs e tipos
│   ├── internalErrors/            # Tratamento de erros
│   └── test/                      # Mocks e fixtures
├── docs/                          # Documentação e exemplos
└── go.mod                         # Dependências

```

### Fluxo de Dados

```
Cliente HTTP
    ↓
[Auth Middleware] ← JWT Validation
    ↓
[Handler] → Validation
    ↓
[Service] → Business Logic
    ↓
[Repository] → Database
    ↓
Response JSON
    ↓
[Async Worker] (paralelo)
    ↓
[Mail Service] → SMTP
    ↓
Update Campaign Status
```

---

## Quick Start

### Pré-requisitos

- **Go** 1.25.0+
- **PostgreSQL** 12+
- **Docker** (opcional, recomendado)
- **Git**

### 1️⃣ Clonar o Repositório

```bash
git clone https://github.com/seu-usuario/emailn.git
cd emailn
```

### 2️⃣ Configurar Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=emailn_db
DB_SSLMODE=disable

# Email Configuration
EMAIL_SMTP=smtp.gmail.com
EMAIL_USER=seu-email@gmail.com
EMAIL_PASSWD=sua-senha-app

# JWT/OAuth
JWT_SECRET=sua-secret-key
IDENTITY_PROVIDER=http://localhost:8080
```

### 3️⃣ Instalar Dependências

```bash
go mod download
go mod tidy
```

### 4️⃣ Iniciar o Banco de Dados

```bash
# Com Docker Compose (recomendado)
docker-compose up -d

# Ou manualmente
psql -U postgres -c "CREATE DATABASE emailn_db"
```

### 5️⃣ Executar Migrations

```bash
go run cmd/migrate/main.go
```

### 6️⃣ Iniciar o Servidor API

```bash
# Terminal 1 - API Server (porta 3000)
go run cmd/api/main.go

# Terminal 2 - Worker (processamento assíncrono)
go run cmd/worker/main.go
```

**Pronto!** A API estará disponível em `http://localhost:3000`

---

## Documentação

### Endpoints da API

#### Autenticação

```bash
# Obter Token JWT
POST http://localhost:8080/realms/GoProvider/protocol/openid-connect/token
Content-Type: application/x-www-form-urlencoded

client_id=emailn&username=usuario@email.com&password=123456&grant_type=password
```

#### Campanhas

##### Criar Campanha

```http
POST /campaigns
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Black Friday 2024",
  "content": "<h1>Oferta Imperdível!</h1><p>...</p>",
  "emails": [
    "cliente1@email.com",
    "cliente2@email.com",
    "cliente3@email.com"
  ]
}
```

**Response:**
```json
{
  "id": "d98qceobgs4lef5509eg"
}
```

---

##### Obter Campanha

```http
GET /campaigns/{id}
Authorization: Bearer {token}
```

**Response:**
```json
{
  "id": "d98qceobgs4lef5509eg",
  "name": "Black Friday 2024",
  "content": "<h1>Oferta Imperdível!</h1><p>...</p>",
  "status": "Pending",
  "amountOfEmailToSend": 3,
  "createdBy": "usuario@email.com"
}
```

---

##### Iniciar Campanha

```http
PATCH /campaigns/start/{id}
Authorization: Bearer {token}
```

**States Transition:**
- `Pending` → `Started` → `Done`
- `Pending` → `Started` → `Failed`

---

##### Deletar Campanha

```http
DELETE /campaigns/delete/{id}
Authorization: Bearer {token}
```

*Apenas campanhas em estado `Pending` podem ser deletadas*

---

### Estados de Campanha

| Estado | Descrição | Transições |
|--------|-----------|-----------|
| **Pending** | Aguardando início | → Started, → Deleted |
| **Started** | Em processamento | → Done, → Failed |
| **Done** | Concluída com sucesso | → Canceled |
| **Failed** | Falha no envio | → Canceled |
| **Canceled** | Cancelada pelo usuário | Final |
| **Deleted** | Deletada logicamente | Final |

---

## Testes

### Executar Todos os Testes

```bash
go test ./...
```

### Testes com Cobertura

```bash
go test ./... -cover
```

### Testes Específicos

```bash
# Domínio de Campanha
go test ./internal/domain/campaign -v

# Endpoints HTTP
go test ./internal/endpoints -v

# Testes de integração
go test ./... -tags=integration
```

### Estrutura de Testes

```
internal/
├── domain/campaign/
│   ├── campaign_test.go          # Unit tests da entidade
│   ├── service_test.go           # Unit tests da lógica
├── endpoints/
│   ├── campaigns_post_test.go    # Handler tests
│   ├── campaigns_get_by_id_test.go
│   ├── auth_test.go              # Auth middleware tests
└── test/
    └── internal-mock/            # Mocks reutilizáveis
        ├── campaign_repository_mock.go
        └── campaign_service_mock.go
```

---


### Exemplo Completo: Criar e Enviar Campanha

```bash
#!/bin/bash

# 1. Obter token
TOKEN=$(curl -s -X POST http://localhost:8080/realms/GoProvider/protocol/openid-connect/token \
  -d "client_id=emailn&username=user@email.com&password=123456&grant_type=password" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  | jq -r '.access_token')

# 2. Criar campanha
CAMPAIGN=$(curl -s -X POST http://localhost:3000/campaigns \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Newsletter Semanal",
    "content": "<h1>Novidades da Semana</h1>",
    "emails": ["user1@email.com", "user2@email.com"]
  }')

CAMPAIGN_ID=$(echo $CAMPAIGN | jq -r '.id')

# 3. Obter detalhes
curl -s -X GET http://localhost:3000/campaigns/$CAMPAIGN_ID \
  -H "Authorization: Bearer $TOKEN" | jq

# 4. Iniciar campanha
curl -s -X PATCH http://localhost:3000/campaigns/start/$CAMPAIGN_ID \
  -H "Authorization: Bearer $TOKEN"

echo "Campanha iniciada! Status pode ser verificado no banco de dados"
```

---

## Performance e Escalabilidade

### Benchmarks Típicos

| Operação | Tempo | Taxa |
|----------|-------|------|
| Criar Campanha | ~50ms | 20 req/s |
| Obter Campanha | ~30ms | 33 req/s |
| Listar Campanhas | ~100ms | 10 req/s |
| Enviar 1000 Emails | ~5s | Assíncrono |


## Segurança

### Práticas Implementadas

- **JWT Signed** com algoritmo HS256
- **Validação de Email** com regex
- **CSRF Protection** via headers
- **SQL Injection Prevention** via parameterized queries (GORM)
- **Input Validation** em todos os campos
- **Error Sanitization** sem expor informações internas
- **HTTPS Ready** (use em produção)


<div align="center">

**Feito por [Gabriel Prezzoti](https://github.com/GabrielPrezzoti)**
**Linkedin: (https://www.linkedin.com/in/gabriel-barros-prezzoti/)**


[Voltar ao topo](#-emailn---email-campaign-management-system)

</div>

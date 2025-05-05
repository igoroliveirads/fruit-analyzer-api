# Fruit Analyzer

API REST em Go para análise de frutas usando Uber Fx.

## Requisitos

- Go 1.21 ou superior

## Instalação

1. Clone o repositório
2. Execute `go mod download` para baixar as dependências

## Executando o projeto

```bash
go run cmd/api/main.go
```

O servidor será iniciado na porta 8080.

## Endpoints

### POST /analyze

Analisa uma imagem de fruta e retorna seu status.

**Request:**
- Content-Type: multipart/form-data
- Campo: "image" (arquivo de imagem)

**Response:**
```json
{
    "status": "verde|madura|passada|desconhecido"
}
```

## Estrutura do Projeto

```
.
├── cmd/                    # Ponto de entrada da aplicação
│   └── api/               # API REST
│       └── main.go        # Arquivo principal
├── internal/              # Código interno da aplicação
│   ├── analyzer/         # Pacote de análise de frutas
│   │   ├── analyzer.go   # Interface e implementação base
│   │   └── banana.go     # Implementação específica para banana
│   ├── config/           # Configurações da aplicação
│   │   └── config.go     # Configurações e variáveis de ambiente
│   ├── handler/          # Handlers HTTP
│   │   └── fruit.go      # Handler para análise de frutas
│   └── server/           # Configuração do servidor
│       └── server.go     # Configuração e inicialização do servidor
├── pkg/                   # Pacotes que podem ser reutilizados externamente
│   └── logger/           # Configuração do logger
│       └── logger.go     # Configuração do zap logger
├── go.mod                # Gerenciamento de dependências
└── README.md            # Documentação do projeto
```

## Arquitetura

O projeto segue uma arquitetura limpa e modular:

- `cmd/`: Contém os pontos de entrada da aplicação
- `internal/`: Código interno que não deve ser importado por outros projetos
  - `analyzer/`: Implementações específicas para cada tipo de fruta
  - `config/`: Configurações da aplicação
  - `handler/`: Handlers HTTP
  - `server/`: Configuração do servidor
- `pkg/`: Pacotes que podem ser reutilizados por outros projetos

## Extensibilidade

O projeto foi projetado para ser facilmente extensível. Para adicionar suporte a novas frutas:

1. Crie um novo arquivo no diretório `internal/analyzer/` (ex: `apple.go`)
2. Implemente a interface `FruitAnalyzer` para a nova fruta
3. Registre a nova implementação no container de dependências

## Logging

O projeto utiliza o pacote `zap` para logging. Os logs são configurados para produção por padrão. 
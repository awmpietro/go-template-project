catwise/
├── cmd/
│   └── main.go                    # Entry point da aplicação
│
├── internal/
│   ├── domain/                    # Modelos de negócio puros (entidades)
│   │   ├── models.go              # User, Cat, Consultation, etc.
│
│   ├── ports/                     # Interfaces (Contracts / Ports) - UseCase depende disso
│   │   ├── user_repository.go
│   │   ├── cat_repository.go
│   │   ├── vector_store.go
│   │   ├── cache.go
│   │   ├── llm.go
│
│   ├── repository/                # Adapters (implementação das ports) - Ex: GORM, Redis
│   │   ├── postgres/
│   │   │   ├── user_postgres.go
│   │   │   ├── cat_postgres.go
│   │   ├── redis/
│   │   │   └── cache_redis.go
│
│   ├── infra/                     # Configuração e conexão com infra externa
│   │   ├── postgres.go            # Init DB
│   │   ├── redis.go               # Init Redis
│   │   ├── firebase.go            # Init Firebase Admin SDK
│   │   ├── chromadb.go            # Init Vector Store
│
│   ├── usecase/                   # Regras de negócio / Application Layer
│   │   ├── auth_usecase.go        # Login, valida Firebase, gera JWT
│   │   ├── user_usecase.go
│   │   ├── cat_usecase.go
│
│   ├── models/                    # DTOs / Request & Response Models (com json tags)
│   │   ├── auth_dto.go
│   │   ├── user_dto.go
│   │   ├── cat_dto.go
│
│   ├── delivery/                  # Camada de entrega - Handlers + Rotas
│   │   ├── handlers/
│   │   │   ├── auth_handler.go
│   │   │   ├── user_handler.go
│   │   │   ├── cat_handler.go
│   │   ├── routes/
│   │   │   ├── auth_routes.go
│   │   │   ├── user_routes.go
│   │   │   ├── cat_routes.go
│   │   ├── main_routes.go         # Carrega todas as rotas principais
│
│   ├── services/                  # Serviços externos (LLM, Image Analysis, etc.)
│   │   ├── llm_service.go
│   │   ├── rag_service.go
│
├── pkg/                           # Utilitários e middlewares genéricos e reutilizáveis
│   ├── utils/
│   │   ├── bcrypt.go              # Hash de senhas (se necessário)
│   │   ├── jwt.go                 # Geração de JWT (se próprio)
│   ├── middlewares/
│   │   ├── auth_middleware.go     # Ex: checar JWT, API key
│   │   ├── cors.go
│
├── go.mod
├── go.sum

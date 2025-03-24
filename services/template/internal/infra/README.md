# Infrastructure Layer (Drivers / External Clients)

## Purpose:
Holds the configuration and initialization of external resources.

- Database connection (PostgreSQL, Redis)
- Vector Store initialization (ChromaDB, Pinecone)
- LLM API client (OpenAI, Anthropic)

## Important:
- Provides ready-to-use clients for the repositories or services
- Does NOT contain business logic
- Should be imported only by the main entry point or repositories

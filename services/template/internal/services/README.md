# Services Layer (External Integrations / AI / RAG)

## Purpose:
Contains domain service implementations that interact with external systems like:

- LLM providers (OpenAI, Claude, Mistral)
- RAG orchestration logic
- Image analysis integrations (Google Vision, GPT-4 Vision)

## Important:
- Can also implement `ports` interfaces
- Should be reusable and testable
- No business rules here, only orchestration of external services

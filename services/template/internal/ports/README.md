# Ports Layer (Interfaces / Contracts)

## Purpose:
Defines the interfaces that represent the boundaries between the core application and external systems.

- Repository interfaces
- Cache interfaces
- Vector Store or LLM service interfaces

## Important:
- `usecase` depends on `ports`
- `repository`, `infra`, or `services` implement `ports`
- Central for applying the Dependency Inversion Principle (DIP)

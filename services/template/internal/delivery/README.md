# Delivery Layer (Handlers and Routes)

## Purpose:
The delivery layer is responsible for handling incoming requests from the outside world (HTTP, gRPC, etc.) and returning responses.

- Contains all HTTP Handlers
- Defines the API routes
- Handles request validation and response formatting
- Calls the use cases (application layer)

## Important:
- No business logic here
- Should only work with use cases and models (DTOs)
- Can transform request/response models to domain models

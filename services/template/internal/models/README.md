# Models Layer (DTOs / Request-Response Objects)

## Purpose:
Contains the API models used by the delivery layer for requests and responses.

- Structs with JSON tags, validation tags, etc.
- Dedicated to the transport layer (HTTP API)

## Important:
- Models here are NOT the same as domain entities
- Used exclusively by the delivery layer (handlers and routes)
- Should not contain business logic or database logic

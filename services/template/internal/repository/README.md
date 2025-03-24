# Repository Adapters Layer

## Purpose:
Implements the persistence logic required by the use cases.

- Concrete implementations of the repository interfaces (ports)
- Communicates with databases, caches, or other storage mechanisms

## Example:
- PostgreSQL repository
- Redis repository
- Vector Store repository

## Important:
- Implements interfaces defined in the `ports` package
- Contains SQL queries, ORM calls, or external storage access

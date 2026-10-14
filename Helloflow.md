# Hello World Flow — DDD Architecture

> A walkthrough of how an HTTP request travels through all four architectural layers from outside inward, and back out.

---

## Request Flow Diagram

```
Browser/Client
    │
    ▼
infrastructure/http        (Inbound Adapter)
    │
    ▼
application/query          (Use Case / Orchestration)
    │  ┌──────────────────────────────┐
    ├──► application/ports            (Port / Interface)
    │  └──────────────────────────────┘
    │
    ▼
infrastructure/persistence (Outbound Adapter)
    │
    ▼
domain/hello               (Pure Business Logic)
    │
    ▼
Response ──────────────────────────────► Browser/Client
```

---

## Layer-by-Layer Breakdown

### 1. `cmd/main.go` — Wiring Layer

Ties everything together. It initializes adapters and injects dependencies:

- Creates `MemoryHelloRepo` (infrastructure)
- Passes it into `GetHelloHandler` (application)
- Passes the use case into `HelloHandler` (HTTP controller)
- Registers the `/hello` route

### 2. `infrastructure/http` — Inbound Adapter

When a request hits `GET /hello`:

- `HelloHandler` receives it
- Does **not** query the DB or domain directly
- Delegates to the application use case: `getHello.Handle()`

### 3. `application/query` — Use Case / Orchestration

`GetHelloHandler` defines the business workflow:

- Calls `HelloRepo.GetStaticName()` via the **Port interface** (`application/ports`)
- This keeps the application layer decoupled from the actual infrastructure implementation

### 4. `infrastructure/persistence` — Outbound Adapter

`MemoryHelloRepo` implements the `HelloRepo` port:

- Returns the raw string `"World"`
- In production, this would query **PostgreSQL** via `pgx`

### 5. `domain/hello` — Pure Business Logic

The application layer passes `"World"` to `domain.NewMessage(name)`:

- The domain applies business rules (e.g. greeting format)
- Returns `"Hello, World from the Domain layer!"`

---

## Response Path

```
domain/hello
    └──► application/query
              └──► infrastructure/http
                        └──► Browser/Client
```

---

## Key Rule: Dependency Direction

> **Dependencies always point inward.**

| Layer             | Can import                               |
| ----------------- | ---------------------------------------- |
| `domain/`         | Go standard library only                 |
| `application/`    | `domain/` only                           |
| `infrastructure/` | `application/`, `domain/`, external libs |
| `cmd/main.go`     | All layers (wiring only)                 |

The **Domain** knows nothing about HTTP or databases. This is the core guarantee of Clean Architecture.

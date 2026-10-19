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

## Dependency Rule Audit

> **Dependencies always point inward — verified against actual source imports.**

| File                                       | Layer             | Imports                                                                             | Status                            |
| ------------------------------------------ | ----------------- | ----------------------------------------------------------------------------------- | --------------------------------- |
| `domain/hello/model.go`                    | `domain/`         | `"errors"`, `"strings"`, `"time"` (all stdlib)                                      | ✅ stdlib only                    |
| `application/ports/hello_repo.go`          | `application/`    | `"context"` (stdlib)                                                                | ✅ stdlib only                    |
| `application/query/get_hello.go`           | `application/`    | `"context"` + `application/ports` + `domain/hello`                                  | ✅ `domain/` + sibling ports only |
| `infrastructure/persistence/hello_repo.go` | `infrastructure/` | `"context"` (stdlib)                                                                | ✅ no upward violations           |
| `infrastructure/http/hello_handler.go`     | `infrastructure/` | `"net/http"` + `application/query`                                                  | ✅ imports `application/`         |
| `cmd/main.go`                              | `cmd/`            | stdlib + `application/query` + `infrastructure/http` + `infrastructure/persistence` | ✅ wires all layers               |

### Import Rules Summary

| Layer             | May import                                  |
| ----------------- | ------------------------------------------- |
| `domain/`         | Go standard library only                    |
| `application/`    | `domain/` only                              |
| `infrastructure/` | `application/`, `domain/`, external drivers |
| `cmd/main.go`     | All layers (wiring only)                    |

> **Note on implicit interface satisfaction:** `infrastructure/persistence/hello_repo.go` does not explicitly import `application/ports`. Instead, `cmd/main.go` passes `*MemoryHelloRepo` into `query.NewGetHelloHandler()`, which accepts the `ports.HelloRepo` interface. Go satisfies interfaces **implicitly**, so the compiler enforces the contract without a direct import — this is idiomatic Go and not a violation.

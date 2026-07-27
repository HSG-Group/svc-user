# user-svc

The **User Service** is the bounded context responsible for managing **user profiles**, **presence** (online status), **friendships**, and **user preferences** in the Harmony platform. It is written in Go and follows a strict **Domain-Driven Design (DDD)** layered architecture with the Hexagonal (Ports & Adapters) pattern.

---

## Table of contents

- [Overview](#overview)
- [Bounded context](#bounded-context)
- [Architecture](#architecture)
- [Folder structure](#folder-structure)
- [Domain layer](#domain-layer)
- [Application layer](#application-layer)
- [Infrastructure layer](#infrastructure-layer)
- [API reference](#api-reference)
- [Domain events](#domain-events)
- [Database](#database)
- [Local development](#local-development)
- [Testing](#testing)
- [Environment variables](#environment-variables)
- [Related services](#related-services)

---

## Overview

| Property     | Value                                                     |
| ------------ | --------------------------------------------------------- |
| Language     | Go 1.22+                                                  |
| Architecture | Domain-Driven Design (DDD) + Hexagonal (Ports & Adapters) |
| Database     | PostgreSQL (primary) · Redis (presence cache)             |
| Messaging    | Apache Kafka (event publishing + consuming)               |
| Transport    | HTTP/REST · gRPC (internal service-to-service)            |
| HTTP Port    | `8083`                                                    |
| gRPC Port    | `9083`                                                    |

---

## Bounded context

The User Service owns everything related to **who a user is** and **how they relate to other users**. It does **not** own authentication (that is `auth-svc`), community membership (that is `community-svc`), or messaging history (that is `messaging-svc`).

```
auth-svc ──► identity.user.registered ──► user-svc (creates profile)
                                       ──► user-svc (sets default preferences)

user-svc ──► user.presence.changed    ──► ws-gateway (fan-out to friends)
         ──► user.profile.updated     ──► search-svc (re-index)
         ──► user.friend.added        ──► notification-svc
```

> **Rule:** No service reads the `user-svc` database directly.
> Cross-context user data (e.g. display name in messages) is denormalised at write time via events, not joined at read time.

---

## Architecture

```
┌──────────────────────────────────────────────────────┐
│                  Infrastructure layer                 │
│                                                      │
│  Inbound adapters          Outbound adapters         │
│  ─────────────────         ─────────────────         │
│  HTTP handlers             PostgreSQL repo           │
│  gRPC server               Redis presence store      │
│  Kafka consumer            Kafka publisher           │
│  Auth middleware           S3 avatar uploader        │
│                                                      │
│   ┌──────────────────────────────────────────────┐   │
│   │              Application layer               │   │
│   │                                              │   │
│   │   Commands            Queries                │   │
│   │   ────────            ───────                │   │
│   │   UpdateProfile       GetProfile             │   │
│   │   SetPresence         GetPresence            │   │
│   │   SendFriendRequest   ListFriends            │   │
│   │   AcceptFriend        SearchUsers            │   │
│   │   BlockUser                                  │   │
│   │                                              │   │
│   │   Ports (interfaces) · DTOs                  │   │
│   │                                              │   │
│   │   ┌──────────────────────────────────────┐   │   │
│   │   │           Domain layer               │   │   │
│   │   │                                      │   │   │
│   │   │   User aggregate root                │   │   │
│   │   │   Friendship entity                  │   │   │
│   │   │   Presence entity                    │   │   │
│   │   │   Value objects: Username, Avatar    │   │   │
│   │   │   Domain events: ProfileUpdated …    │   │   │
│   │   └──────────────────────────────────────┘   │   │
│   └──────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────┘
              ▲ dependencies always point inward
```

**Dependency rule:**

| Layer             | May import                                  |
| ----------------- | ------------------------------------------- |
| `domain/`         | Go standard library only                    |
| `application/`    | `domain/` only                              |
| `infrastructure/` | `application/`, `domain/`, external drivers |
| `cmd/main.go`     | All layers (wiring only)                    |

---

## Folder structure

```
user-svc/
├── domain/                              # Pure business logic — no framework imports
│   ├── user/
│   │   ├── user.go                      # User aggregate root
│   │   ├── user_events.go               # ProfileUpdated, AvatarChanged, UserBlocked …
│   │   └── user_repository.go           # IUserRepo interface
│   ├── friendship/
│   │   ├── friendship.go                # Friendship entity (Pending → Accepted | Declined)
│   │   ├── friendship_events.go         # FriendRequestSent, FriendAccepted, FriendRemoved
│   │   └── friendship_repository.go     # IFriendshipRepo interface
│   ├── presence/
│   │   ├── presence.go                  # Presence entity (Online/Idle/DND/Invisible/Offline)
│   │   └── presence_repository.go       # IPresenceRepo interface
│   └── valueobjects/
│       ├── username.go                  # Username: 2–32 chars, alphanumeric + underscore
│       ├── avatar_url.go                # AvatarURL: must be HTTPS, max 2MB
│       ├── custom_status.go             # CustomStatus: max 128 chars
│       └── discriminator.go            # Discriminator: 4-digit tag e.g. #0042
│
├── application/                         # Orchestration — no HTTP/DB/Kafka knowledge
│   ├── command/
│   │   ├── update_profile.go            # UpdateProfile use case
│   │   ├── set_presence.go              # SetPresence use case
│   │   ├── upload_avatar.go             # UploadAvatar use case
│   │   ├── send_friend_request.go       # SendFriendRequest use case
│   │   ├── accept_friend_request.go     # AcceptFriendRequest use case
│   │   ├── decline_friend_request.go    # DeclineFriendRequest use case
│   │   ├── remove_friend.go             # RemoveFriend use case
│   │   └── block_user.go               # BlockUser use case
│   ├── query/
│   │   ├── get_profile.go               # GetProfile query (cache-first)
│   │   ├── get_presence.go              # GetPresence query (Redis-first)
│   │   ├── list_friends.go              # ListFriends query
│   │   ├── list_pending_requests.go     # ListPendingRequests query
│   │   └── search_users.go             # SearchUsers query (username prefix search)
│   ├── ports/
│   │   ├── user_repo_port.go            # IUserRepo
│   │   ├── friendship_repo_port.go      # IFriendshipRepo
│   │   ├── presence_repo_port.go        # IPresenceRepo
│   │   ├── event_publisher_port.go      # IEventPublisher
│   │   ├── cache_port.go                # ICache
│   │   └── avatar_store_port.go         # IAvatarStore (upload to S3)
│   └── dto/
│       ├── update_profile_dto.go        # Input: displayName, bio, customStatus
│       ├── user_profile_dto.go          # Output: full profile response
│       ├── presence_dto.go              # Output: status + lastSeenAt
│       └── friend_dto.go               # Output: friend list item
│
├── infrastructure/                      # Adapters — implements all ports
│   ├── persistence/
│   │   ├── pg_user_repo.go              # PostgreSQL: implements IUserRepo
│   │   ├── pg_friendship_repo.go        # PostgreSQL: implements IFriendshipRepo
│   │   ├── redis_presence_repo.go       # Redis: implements IPresenceRepo (TTL-based)
│   │   ├── redis_cache.go               # Redis: implements ICache (profile cache)
│   │   ├── outbox.go                    # Writes domain events to outbox table atomically
│   │   └── outbox_relay.go              # Polls outbox → publishes to Kafka
│   ├── http/
│   │   ├── profile_handler.go           # GET/PATCH /users/:id
│   │   ├── avatar_handler.go            # POST /users/:id/avatar (multipart)
│   │   ├── presence_handler.go          # GET/PUT /users/:id/presence
│   │   ├── friends_handler.go           # GET/POST/DELETE /users/:id/friends
│   │   ├── search_handler.go            # GET /users/search?q=
│   │   ├── middleware.go                # JWT auth, rate limit, request logging
│   │   └── router.go                    # Route registration
│   ├── grpc/
│   │   ├── user_grpc_server.go          # gRPC: GetProfile, GetPresence (internal calls)
│   │   └── user.proto                   # Protobuf schema
│   ├── messaging/
│   │   ├── kafka_publisher.go           # Implements IEventPublisher
│   │   └── kafka_consumer.go            # Consumes: identity.user.registered
│   └── storage/
│       └── s3_avatar_store.go           # Implements IAvatarStore (S3 + CloudFront URL)
│
├── cmd/
│   └── main.go                          # DI wiring: new(repo) → use case → handler
│
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_friendships.sql
│   ├── 003_create_preferences.sql
│   └── 004_create_outbox.sql
│
├── Dockerfile
├── go.mod
└── go.sum
```

---

## Domain layer

### User aggregate

`domain/user/user.go`

The `User` is the **aggregate root** for the User context. All profile mutations flow through this struct. The aggregate raises domain events when state changes — these are collected and flushed by the application layer after each command.

```go
type User struct {
    id           UserID
    username     Username        // value object — enforces naming rules
    displayName  string
    bio          string
    avatarURL    AvatarURL       // value object — HTTPS only
    customStatus CustomStatus    // value object — max 128 chars
    discriminator Discriminator  // value object — 4-digit tag
    blockedIDs   []UserID
    preferences  Preferences
    events       []DomainEvent   // uncommitted events
    createdAt    time.Time
    updatedAt    time.Time
}

// NewUser is the factory called when identity.user.registered is consumed
func NewUser(id UserID, username Username) (*User, error)

// UpdateProfile validates and applies profile changes, raises ProfileUpdated
func (u *User) UpdateProfile(displayName, bio string, status CustomStatus) error

// ChangeAvatar sets a new avatar URL, raises AvatarChanged
func (u *User) ChangeAvatar(url AvatarURL) error

// Block prevents a user from DMing or seeing this user, raises UserBlocked
func (u *User) Block(targetID UserID) error

// Unblock removes a block
func (u *User) Unblock(targetID UserID) error
```

**Invariants enforced by the aggregate:**

- Username: 2–32 chars, alphanumeric + underscores, no spaces, unique (checked at repo level)
- Display name: 1–80 chars after trim
- Bio: max 190 chars
- Custom status: max 128 chars
- Avatar URL: must be HTTPS, served from the Harmony CDN domain only
- A user cannot block themselves
- A user cannot have more than 1,000 blocked users

### Friendship entity

`domain/friendship/friendship.go`

```go
type Friendship struct {
    id         FriendshipID
    requesterID UserID
    addresseeID UserID
    status     FriendshipStatus  // Pending | Accepted | Declined | Blocked
    createdAt  time.Time
    updatedAt  time.Time
}

// Send creates a new pending friend request
func Send(requesterID, addresseeID UserID) (*Friendship, error)

// Accept transitions Pending → Accepted, raises FriendAccepted
func (f *Friendship) Accept() error

// Decline transitions Pending → Declined
func (f *Friendship) Decline() error

// Remove removes an accepted friendship
func (f *Friendship) Remove() error
```

**Invariants:**

- Cannot send a friend request to yourself
- Cannot send a duplicate request if one is already Pending or Accepted
- Cannot accept a request you sent (requester cannot self-accept)
- Max 1,000 friends per user (enforced by repo before `Send`)

### Presence entity

`domain/presence/presence.go`

Presence is intentionally **not** part of the User aggregate because it changes constantly (every few seconds) and has different persistence characteristics (ephemeral, TTL-based, stored in Redis rather than PostgreSQL).

```go
type PresenceStatus string

const (
    Online    PresenceStatus = "online"
    Idle      PresenceStatus = "idle"
    DND       PresenceStatus = "dnd"       // Do Not Disturb
    Invisible PresenceStatus = "invisible"
    Offline   PresenceStatus = "offline"
)

type Presence struct {
    userID     UserID
    status     PresenceStatus
    clientID   string    // which device/tab set this
    lastSeenAt time.Time
}
```

### Value objects

| Value object    | File                            | Rule enforced                                      |
| --------------- | ------------------------------- | -------------------------------------------------- |
| `Username`      | `valueobjects/username.go`      | 2–32 chars, `[a-zA-Z0-9_]` only, no reserved words |
| `AvatarURL`     | `valueobjects/avatar_url.go`    | HTTPS, Harmony CDN domain only, max 2MB            |
| `CustomStatus`  | `valueobjects/custom_status.go` | Max 128 chars, optional emoji prefix               |
| `Discriminator` | `valueobjects/discriminator.go` | 4-digit zero-padded number `0001`–`9999`           |

### Domain events

```go
type UserCreated struct {
    UserID    string
    Username  string
    OccurredAt time.Time
}

type ProfileUpdated struct {
    UserID      string
    DisplayName string
    Bio         string
    OccurredAt  time.Time
}

type AvatarChanged struct {
    UserID    string
    AvatarURL string
    OccurredAt time.Time
}

type PresenceChanged struct {
    UserID    string
    Status    string
    OccurredAt time.Time
}

type FriendRequestSent struct {
    RequesterID string
    AddresseeID string
    OccurredAt  time.Time
}

type FriendAccepted struct {
    RequesterID string
    AddresseeID string
    OccurredAt  time.Time
}

type UserBlocked struct {
    BlockerID  string
    BlockedID  string
    OccurredAt time.Time
}
```

---

## Application layer

### Commands

#### `UpdateProfile`

```
1. Load User aggregate via IUserRepo.FindByID
2. Call user.UpdateProfile(displayName, bio, status)   ← invariants checked here
3. Call IUserRepo.Save(user)                           ← writes row + outbox atomically
4. Call ICache.Del("profile:" + userID)                ← invalidate cache
5. Return UserProfileDTO
```

#### `SetPresence`

```
1. Construct Presence entity from DTO
2. Call IPresenceRepo.Set(presence)    ← writes to Redis with TTL 65s
3. Publish PresenceChanged via IEventPublisher (async via outbox)
4. Return PresenceDTO
```

> Presence uses a 65-second TTL. Clients send a heartbeat every 30 seconds.
> If the heartbeat stops, Redis expires the key and the user appears Offline.

#### `UploadAvatar`

```
1. Validate file: JPEG/PNG/GIF/WEBP only, max 2MB
2. Call IAvatarStore.Upload(userID, file)   ← uploads to S3, returns CDN URL
3. Construct AvatarURL value object         ← validates CDN domain
4. Load User aggregate
5. Call user.ChangeAvatar(avatarURL)
6. Call IUserRepo.Save(user)
7. Call ICache.Del("profile:" + userID)
8. Return new avatar URL
```

#### `SendFriendRequest`

```
1. Check IFriendshipRepo.Exists(requesterID, addresseeID) — reject if duplicate
2. Check IUserRepo.IsFriendLimitReached(requesterID)      — max 1,000
3. Check if addresseeID has blocked requesterID           — reject silently
4. Call domain.Friendship.Send(requesterID, addresseeID)
5. Call IFriendshipRepo.Save(friendship)
6. Publish FriendRequestSent event
```

#### `AcceptFriendRequest`

```
1. Load Friendship via IFriendshipRepo.FindByID
2. Verify caller is the addressee (not the requester)
3. Call friendship.Accept()
4. Call IFriendshipRepo.Save(friendship)
5. Publish FriendAccepted event
6. Invalidate friend-list cache for both users
```

#### `BlockUser`

```
1. Load User aggregate (blocker)
2. Call user.Block(targetID)            ← checks not self, not over limit
3. Call IUserRepo.Save(user)
4. Remove any existing friendship between the two users
5. Publish UserBlocked event
```

### Queries

Queries bypass the domain layer entirely and read from the read model.

| Query                 | Cache key                 | TTL | Fallback                                |
| --------------------- | ------------------------- | --- | --------------------------------------- |
| `GetProfile`          | `profile:{userId}`        | 60s | PostgreSQL `users` table                |
| `GetPresence`         | Redis key (IPresenceRepo) | 65s | Returns `offline` if key missing        |
| `ListFriends`         | `friend-list:{userId}`    | 30s | PostgreSQL `friendships` table          |
| `ListPendingRequests` | —                         | —   | PostgreSQL (always fresh)               |
| `SearchUsers`         | —                         | —   | PostgreSQL `users` table (ILIKE prefix) |

### Ports

```go
// ports/user_repo_port.go
type IUserRepo interface {
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id string) (*domain.User, error)
    FindByUsername(ctx context.Context, username string) (*domain.User, error)
    IsFriendLimitReached(ctx context.Context, userID string) (bool, error)
}

// ports/presence_repo_port.go
type IPresenceRepo interface {
    Set(ctx context.Context, presence *domain.Presence) error
    Get(ctx context.Context, userID string) (*domain.Presence, error)
    GetBulk(ctx context.Context, userIDs []string) ([]*domain.Presence, error)
}

// ports/avatar_store_port.go
type IAvatarStore interface {
    Upload(ctx context.Context, userID string, file io.Reader, contentType string) (string, error)
    Delete(ctx context.Context, userID string) error
}

// ports/event_publisher_port.go
type IEventPublisher interface {
    Publish(ctx context.Context, topic string, event any) error
}
```

---

## Infrastructure layer

### PostgreSQL repository (`persistence/pg_user_repo.go`)

Implements `IUserRepo`. Uses `pgx/v5`. The `Save` method opens a single transaction that writes the user row and any uncommitted domain events to the `outbox` table atomically — preventing the dual-write problem.

```go
func (r *PGUserRepo) Save(ctx context.Context, user *domain.User) error {
    return r.db.BeginTx(ctx, func(tx pgx.Tx) error {
        // 1. UPSERT into users table
        // 2. INSERT uncommitted domain events into outbox
        // 3. user.ClearEvents()
    })
}
```

### Redis presence store (`persistence/redis_presence_repo.go`)

Implements `IPresenceRepo`. Each user's presence is stored under key `presence:{userID}` as a JSON hash with a **65-second TTL**. The client heartbeat (every 30s) resets the TTL. If the TTL expires, `GetPresence` returns `offline` by default.

```go
// Key format
"presence:{userID}" → JSON { status, clientID, lastSeenAt }  TTL: 65s
```

### S3 avatar store (`storage/s3_avatar_store.go`)

Implements `IAvatarStore`. Uploads to the `harmony-avatars` S3 bucket under the key `avatars/{userID}/{timestamp}.{ext}`. Returns the CloudFront CDN URL. Old avatars are deleted on upload.

### Kafka consumer (`messaging/kafka_consumer.go`)

Subscribes to `identity.user.registered`. On each event:

1. Creates a new `User` domain object via `domain.NewUser`
2. Assigns a unique `Discriminator` (retry if collision)
3. Saves via `IUserRepo.Save`
4. Publishes `user.created` event

### gRPC server (`grpc/user_grpc_server.go`)

Internal-only. Used by `messaging-svc` and `community-svc` to look up display names and avatars without going through the public HTTP API.

```protobuf
service UserService {
    rpc GetProfile(GetProfileRequest) returns (ProfileResponse);
    rpc GetPresenceBulk(GetPresenceBulkRequest) returns (PresenceBulkResponse);
}
```

### HTTP routes

| Method   | Path                           | Use case              | Auth                |
| -------- | ------------------------------ | --------------------- | ------------------- |
| `GET`    | `/users/:id`                   | `GetProfile`          | ✅                  |
| `PATCH`  | `/users/:id`                   | `UpdateProfile`       | ✅ own profile only |
| `POST`   | `/users/:id/avatar`            | `UploadAvatar`        | ✅ own profile only |
| `GET`    | `/users/:id/presence`          | `GetPresence`         | ✅ friend or self   |
| `PUT`    | `/users/:id/presence`          | `SetPresence`         | ✅ own only         |
| `GET`    | `/users/:id/friends`           | `ListFriends`         | ✅                  |
| `GET`    | `/users/:id/friends/pending`   | `ListPendingRequests` | ✅ own only         |
| `POST`   | `/users/:id/friends`           | `SendFriendRequest`   | ✅                  |
| `PATCH`  | `/users/:id/friends/:friendId` | `AcceptFriendRequest` | ✅ addressee only   |
| `DELETE` | `/users/:id/friends/:friendId` | `RemoveFriend`        | ✅                  |
| `POST`   | `/users/:id/blocks`            | `BlockUser`           | ✅ own only         |
| `DELETE` | `/users/:id/blocks/:targetId`  | `UnblockUser`         | ✅ own only         |
| `GET`    | `/users/search`                | `SearchUsers`         | ✅                  |
| `GET`    | `/health`                      | —                     | ❌                  |

---

## Domain events

### Published by this service

| Event               | Kafka topic             | Consumed by                                    |
| ------------------- | ----------------------- | ---------------------------------------------- |
| `UserCreated`       | `user.created`          | search-svc, community-svc                      |
| `ProfileUpdated`    | `user.profile.updated`  | search-svc, ws-gateway (broadcast to friends)  |
| `AvatarChanged`     | `user.avatar.changed`   | messaging-svc (denormalise in message history) |
| `PresenceChanged`   | `user.presence.changed` | ws-gateway (fan-out to online friends)         |
| `FriendRequestSent` | `user.friend.request`   | notification-svc                               |
| `FriendAccepted`    | `user.friend.accepted`  | notification-svc, ws-gateway                   |
| `UserBlocked`       | `user.blocked`          | messaging-svc (hide DMs)                       |

### Consumed by this service

| Event            | Kafka topic                | Action                            |
| ---------------- | -------------------------- | --------------------------------- |
| `UserRegistered` | `identity.user.registered` | Create User profile with defaults |

---

## Database

### Schema

```sql
-- users
CREATE TABLE users (
    id             UUID PRIMARY KEY,
    username       VARCHAR(32) NOT NULL UNIQUE,
    display_name   VARCHAR(80) NOT NULL,
    bio            VARCHAR(190) NOT NULL DEFAULT '',
    avatar_url     TEXT,
    custom_status  VARCHAR(128) NOT NULL DEFAULT '',
    discriminator  CHAR(4) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (username, discriminator)
);

-- friendships
CREATE TABLE friendships (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requester_id  UUID NOT NULL REFERENCES users(id),
    addressee_id  UUID NOT NULL REFERENCES users(id),
    status        VARCHAR(10) NOT NULL
                  CHECK (status IN ('pending','accepted','declined','blocked')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (requester_id, addressee_id),
    INDEX idx_friendships_addressee (addressee_id, status),
    INDEX idx_friendships_requester (requester_id, status)
);

-- user_preferences
CREATE TABLE user_preferences (
    user_id              UUID PRIMARY KEY REFERENCES users(id),
    theme                VARCHAR(10) NOT NULL DEFAULT 'dark',
    message_display      VARCHAR(10) NOT NULL DEFAULT 'cozy',
    notifications_dm     BOOL NOT NULL DEFAULT TRUE,
    notifications_mention BOOL NOT NULL DEFAULT TRUE,
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- outbox (guaranteed event delivery)
CREATE TABLE outbox (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic       TEXT NOT NULL,
    payload     JSONB NOT NULL,
    processed   BOOL NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    INDEX idx_outbox_unprocessed (processed, created_at)
    WHERE processed = FALSE
);
```

### Redis key patterns

| Key                    | Value                | TTL | Purpose           |
| ---------------------- | -------------------- | --- | ----------------- |
| `presence:{userId}`    | JSON presence object | 65s | Online status     |
| `profile:{userId}`     | JSON profile object  | 60s | Profile cache     |
| `friend-list:{userId}` | JSON array           | 30s | Friend list cache |

---

## Local development

### Prerequisites

- Go 1.22+
- Docker + Docker Compose
- `migrate` CLI

### Start dependencies

From the monorepo root:

```bash
docker-compose up -d postgres redis kafka
```

### Run migrations

```bash
migrate -path ./migrations \
        -database "postgres://harmony:harmony@localhost:5432/user_svc?sslmode=disable" \
        up
```

### Start the service

```bash
cd services/user-svc
cp .env.example .env
go run ./cmd/main.go
```

Service starts on `http://localhost:8083` and `grpc://localhost:9083`.

### Quick smoke test

```bash
# Get a user profile
curl http://localhost:8083/users/<user_id> \
  -H "Authorization: Bearer <token>"

# Update profile
curl -X PATCH http://localhost:8083/users/<user_id> \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"displayName": "Harmony User", "bio": "Hello world"}'

# Set presence
curl -X PUT http://localhost:8083/users/<user_id>/presence \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"status": "online"}'

# Send a friend request
curl -X POST http://localhost:8083/users/<user_id>/friends \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"addresseeId": "<target_user_id>"}'
```

---

## Testing

```bash
# Unit tests — domain + application (no Docker needed)
go test ./domain/... ./application/...

# Integration tests — requires PostgreSQL + Redis
go test ./infrastructure/...

# All tests with coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test strategy

| Layer                         | Type        | Mocks / Dependencies                                                           |
| ----------------------------- | ----------- | ------------------------------------------------------------------------------ |
| `domain/`                     | Unit        | None — pure Go                                                                 |
| `application/`                | Unit        | Mock `IUserRepo`, `IPresenceRepo`, `IEventPublisher`, `ICache`, `IAvatarStore` |
| `infrastructure/persistence/` | Integration | Real PostgreSQL via Testcontainers                                             |
| `infrastructure/http/`        | Integration | Full stack via `httptest.Server`                                               |
| `infrastructure/grpc/`        | Integration | gRPC in-process server                                                         |

---

## Environment variables

| Variable           | Default                  | Description                              |
| ------------------ | ------------------------ | ---------------------------------------- |
| `PORT`             | `8083`                   | HTTP server port                         |
| `GRPC_PORT`        | `9083`                   | gRPC server port                         |
| `DATABASE_URL`     | —                        | PostgreSQL connection string             |
| `REDIS_URL`        | `redis://localhost:6379` | Redis connection string                  |
| `KAFKA_BROKERS`    | `localhost:9092`         | Comma-separated Kafka brokers            |
| `KAFKA_GROUP_ID`   | `user-svc`               | Kafka consumer group ID                  |
| `AUTH_SERVICE_URL` | —                        | auth-svc URL for JWT validation          |
| `S3_BUCKET`        | `harmony-avatars`        | S3 bucket for avatar uploads             |
| `CDN_BASE_URL`     | —                        | CloudFront base URL for avatar serving   |
| `OUTBOX_POLL_MS`   | `500`                    | Outbox relay poll interval in ms         |
| `PRESENCE_TTL_S`   | `65`                     | Redis presence TTL in seconds            |
| `LOG_LEVEL`        | `info`                   | `debug` / `info` / `warn` / `error`      |
| `ENV`              | `development`            | `development` / `staging` / `production` |

---

## Related services

| Service            | Relationship                                                                     |
| ------------------ | -------------------------------------------------------------------------------- |
| `auth-svc`         | Publishes `identity.user.registered` → triggers `UserCreated` in user-svc        |
| `auth-svc`         | Validates JWT tokens on every inbound HTTP request                               |
| `community-svc`    | Calls user-svc gRPC to get display names for member lists                        |
| `messaging-svc`    | Calls user-svc gRPC to get avatar + display name for message author              |
| `ws-gateway`       | Consumes `user.presence.changed` to fan-out to online friends                    |
| `notification-svc` | Consumes `user.friend.request` and `user.friend.accepted`                        |
| `search-svc`       | Consumes `user.created` and `user.profile.updated` to maintain user search index |

---

_Part of the Harmony platform monorepo — see the root [README](../../README.md) for full architecture overview._

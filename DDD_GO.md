# Domain-Driven Design (DDD) in Go

Go's type system and philosophy map well to DDD concepts. This document outlines idiomatic implementations of core DDD patterns.

## 1. Aggregates
An Aggregate is a cluster of objects treated as a single unit for data changes.

### Go Implementation
- **Struct**: The Aggregate Root.
- **Unexported Fields**: Enforce invariants by preventing direct external modification.
- **Factory Function**: Ensures the aggregate is created in a valid state.
- **Methods**: Expose behavior (business logic), not data (getters/setters).

```go
package domain

import (
    "errors"
    "github.com/google/uuid"
)

// Customer is the Aggregate Root
type Customer struct {
    id      uuid.UUID
    name    string
    balance int     // Private: Logic must go through methods
    changes []Event // Stored domain events
}

// Factory: Guarantees valid initial state
func NewCustomer(name string) (*Customer, error) {
    if name == "" {
        return nil, errors.New("name required")
    }
    c := &Customer{
        id:   uuid.New(),
        name: name,
    }
    c.raise(CustomerCreated{ID: c.id})
    return c, nil
}

// Business Method: Enforces invariants
func (c *Customer) Charge(amount int) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    if c.balance < amount {
        return errors.New("insufficient funds")
    }
    c.balance -= amount
    c.raise(BalanceChanged{ID: c.id, NewBalance: c.balance})
    return nil
}
```

## 2. Value Objects
Objects defined by their attributes, not their identity. They should be immutable.

### Go Implementation
- **Struct**: With exported fields (if simple) or unexported (if validation needed).
- **No Pointer Receivers**: Methods use value receivers `(m Money)` to ensure immutability.

```go
type Money struct {
    Amount   int
    Currency string
}

// Value Receiver: Returns a NEW Money, doesn't modify 'm'
func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("currency mismatch")
    }
    return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}, nil
}
```

## 3. Domain Events
Something that happened in the domain.

### Go Implementation
- **Struct**: Immutable data carrier.
- **Marker Interface**: Optional, for type safety.

```go
type Event interface { isEvent() }

type BalanceChanged struct {
    ID         uuid.UUID
    NewBalance int
}
func (BalanceChanged) isEvent() {}

// Internal helper to track events
func (c *Customer) raise(e Event) {
    c.changes = append(c.changes, e)
}
```

## 4. Repositories
Abstraction for retrieving and persisting aggregates.

### Go Implementation
- **Interface**: Defined in the **Domain** layer.
- **Implementation**: Defined in the **Infrastructure** layer.
- **Returns Aggregates**: Not DTOs or DB models.

```go
// Defined in domain/customer.go
type CustomerRepository interface {
    Get(ctx context.Context, id uuid.UUID) (*Customer, error)
    Save(ctx context.Context, c *Customer) error
}

// INFRASTRUCTURE: internal/infra/postgres_repo.go
type PostgresCustomerRepo struct {
    db        *sql.DB
    publisher EventPublisher // Interface for message bus (e.g. Kafka/RabbitMQ)
}

func (r *PostgresCustomerRepo) Save(ctx context.Context, c *Customer) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()

    // 1. Save Aggregate State
    _, err = tx.ExecContext(ctx, "UPDATE customers SET balance = $1 WHERE id = $2", c.balance, c.id)
    if err != nil { return err }

    // 2. Dispatch Domain Events (Outbox Pattern or Direct)
    for _, event := range c.changes {
        if err := r.publisher.Publish(ctx, event); err != nil {
            return err
        }
    }

    // 3. Commit Transaction
    if err := tx.Commit(); err != nil { return err }

    // 4. Clear events to prevent re-dispatching if reused
    c.changes = nil 
    return nil
}
```

## 5. Application Services (Command Handlers)
Orchestrate the flow of data. They do not contain business rules.

### Go Implementation
- **Struct**: Holds dependencies (Repositories, Gateways).
- **Method**: Accepts a Command (DTO), loads Aggregate, invokes method, saves.

```go
type ChargeHandler struct {
    repo CustomerRepository
}

type ChargeCommand struct {
    ID     uuid.UUID
    Amount int
}

func (h *ChargeHandler) Handle(ctx context.Context, cmd ChargeCommand) error {
    // 1. Load
    cust, err := h.repo.Get(ctx, cmd.ID)
    if err != nil { return err }

    // 2. Execute Domain Logic
    if err := cust.Charge(cmd.Amount); err != nil {
        return err
    }

    // 3. Persist (Repository handles event dispatching)
    return h.repo.Save(ctx, cust)
}
```

## 6. Anti-Corruption Layer (ACL)
Isolates your domain from external systems.

### Go Implementation
- **Interface (Port)**: Defined in **Domain**. Describes what the domain *needs*.
- **Adapter**: Defined in **Infrastructure**. Translates between external API and Domain interface.

```go
// DOMAIN: internal/domain/shipping.go
type ShippingService interface {
    Ship(itemID uuid.UUID, address Address) (tracking string, err error)
}

// INFRA: internal/infra/fedex.go
type FedExAdapter struct {
    client *fedex.Client
}

func (f *FedExAdapter) Ship(itemID uuid.UUID, address Address) (string, error) {
    // Translate Domain Address -> FedEx Address
    req := f.toFedExRequest(itemID, address)
    
    resp, err := f.client.CreateShipment(req)
    if err != nil {
        return "", fmt.Errorf("fedex error: %w", err)
    }
    return resp.TrackingNumber, nil
}
```

## 7. Onion Architecture (Project Structure)
Also known as Hexagonal or Clean Architecture. The core principle is that **dependencies point INWARD**.

### Layers
1.  **Domain (Core)**: Aggregates, Value Objects, Repository Interfaces.
    *   *Dependencies*: None. Pure Go.
2.  **Application**: Command Handlers, Use Cases.
    *   *Dependencies*: Domain.
3.  **Infrastructure**: Database, External APIs, File System.
    *   *Dependencies*: Domain (implements interfaces), Application.
4.  **Interfaces (Presentation)**: HTTP Handlers, CLI, gRPC.
    *   *Dependencies*: Application.

### Standard Layout
```
/cmd
  /api          # Main entry point
/internal
  /domain       # Core business logic (Aggregates, Repo Interfaces)
    /customer
      entity.go
      repository.go
  /application  # Orchestration (Command Handlers)
    /customer
      service.go
  /infrastructure # Implementation details
    /postgres
      customer_repo.go
    /stripe
      adapter.go
  /interfaces   # Entry points (HTTP, gRPC)
    /http
      handler.go
```

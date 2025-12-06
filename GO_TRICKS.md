# Go Tricks and Compiler Magic

A collection of advanced Go features, compiler directives, and "magic" behaviors.

## 1. `comparable` Constraint
The `comparable` constraint is a **predeclared identifier** (not a normal interface) that denotes the set of all types that can be compared using `==` and `!=`.
- **Defined**: In the Go Language Specification.
- **Enforced**: Directly by the compiler.
- **Usage**: Allows types to be used as map keys in generic functions.
```go
func Keys[K comparable, V any](m map[K]V) []K { ... }
```

## 2. `any` Alias
`any` is simply a type alias for `interface{}` defined in `builtin`.
```go
type any = interface{}
```

## 3. `unsafe` Package
Allows bypassing Go's type safety and memory safety.
- **`unsafe.Pointer`**: A bridge to convert between arbitrary pointer types.
- **`unsafe.Sizeof` / `Offsetof`**: These are **compile-time constants**, not runtime function calls.

## 4. `internal` Packages
A build-system rule enforcing strict module boundaries.
- Packages named `internal` can only be imported by packages rooted in the parent directory.
- Prevents external users from depending on your private implementation details.

## 5. `go:embed`
Compiles static files directly into the binary.
```go
import _ "embed"

//go:embed input.txt
var input string
```

## 6. `go:linkname` (The Backdoor)
A compiler directive that allows linking a local symbol to a private symbol in another package.
```go
//go:linkname time_now time.now
func time_now() (sec int64, nsec int32, mono int64)
```
*Note: Highly unsafe and unstable; mostly for standard library use.*

## 7. Build Tags
Conditional compilation based on OS, architecture, or custom tags.
```go
//go:build linux && amd64
```

## 8. `init()` Functions
- Run automatically before `main()`.
- Can have multiple `init()` functions per package/file.
- Used for registering drivers or initializing global state.

## 9. Compile-Time Interface Checks
A static assertion to ensure a type implements an interface.
```go
var _ io.Reader = (*MyType)(nil)
```
Does nothing at runtime, but fails compilation if `*MyType` stops satisfying `io.Reader`.

## 10. Zero-Size Types (`struct{}`)
- Occupies **0 bytes** of memory.
- Useful for sets (`map[string]struct{}`) or signal channels (`chan struct{}`).

## 11. `iota`
A predeclared identifier for creating enumerated constants. It resets to 0 at the `const` keyword and increments for each specification.

## 12. `runtime.KeepAlive`
Prevents the Garbage Collector from collecting an object prematurely, specifically when interacting with C code or finalizers.

## 13. Type Aliases
```go
type MyInt = int // Alias: MyInt IS int
type MyInt int   // Definition: MyInt is a NEW type
```
Aliases are useful for refactoring and gradual code migration.

## 14. `go:generate`
A standard way to define code generation commands directly in the source code.
```go
//go:generate stringer -type=Pill
```
Running `go generate ./...` will execute these commands.

## 15. Modifying Named Returns in `defer`
`defer` functions run *after* the return statement has evaluated its arguments but *before* the function actually returns. This allows you to modify named return values.
```go
func c() (i int) {
    defer func() { i++ }()
    return 1
}
// Returns 2
```

## 16. The "Nil Interface" Trap
An interface holding a nil concrete value is **not nil**.
```go
var p *int = nil
var i interface{} = p
fmt.Println(i == nil) // False!
```
`i` contains `(type=*int, value=nil)`. An interface is only `nil` if both type and value are nil.

## 17. `go:noinline`
A compiler directive to prevent function inlining. Useful for debugging, benchmarking, or ensuring stack frames exist for `runtime.Callers`.
```go
//go:noinline
func expensive() { ... }
```

## 18. `sync.Pool`
A "magic" pool for reusing objects to reduce Garbage Collector pressure. The runtime may clear the pool at any time (usually during GC), so it's only suitable for caches, not persistent state.

## 19. Struct Padding & Alignment
The order of fields in a struct matters for memory usage due to alignment.
```go
// Uses 24 bytes (on 64-bit)
type Bad struct {
    Flag    bool  // 1 byte + 7 padding
    Count   int64 // 8 bytes
    Enabled bool  // 1 byte + 7 padding
}

// Uses 16 bytes
type Good struct {
    Count   int64 // 8 bytes
    Flag    bool  // 1 byte
    Enabled bool  // 1 byte + 6 padding
}
```

## 20. `runtime.Gosched()`
Explicitly yields the processor, allowing other goroutines to run. Useful in tight loops that don't perform I/O or channel operations to prevent starving the scheduler.

## 21. Slice Header Manipulation
You can access the underlying structure of a slice using `unsafe` and `reflect.SliceHeader` (or `unsafe.Slice` in newer Go).
```go
func StringToBytes(s string) []byte {
    return unsafe.Slice(unsafe.StringData(s), len(s))
}
```
*Note: This creates a zero-copy byte slice from a string. Extremely dangerous if modified.*

## 22. `singleflight`
(From `golang.org/x/sync/singleflight`)
A powerful pattern to prevent "cache stampedes". It ensures that if multiple goroutines ask for the same key at the same time, the expensive function is only called **once**, and the result is shared with all callers.

## 23. Context Values are Linked Lists
`context.WithValue` wraps the parent context. Searching for a value is a linear scan up the chain.
**Tip**: Don't use context for optional parameters or passing lots of data. It's O(N) depth.

## 24. Idiomatic "Head | Tail" Pattern
Go does not have pattern matching or Tail Call Optimization (TCO), so recursion is dangerous for large lists. The idiomatic way to process a list like `head | tail` is an iterative loop using slice re-slicing.
```go
// Functional style (Dangerous in Go: Stack Overflow)
// func process(list []int) {
//     head, tail := list[0], list[1:]
//     process(tail)
// }

// Idiomatic Go (Safe & Fast)
for len(list) > 0 {
    head := list[0]
    list = list[1:] // "Tail" becomes the new list (Zero allocation)
    
    // Process head...
}
```

## 25. Functional Patterns (FlatMap / SelectMany)
Go 1.18+ Generics allow implementing functional patterns like `FlatMap` (C# `SelectMany`).

**Eager Implementation (Slice-based):**
```go
func FlatMap[T, U any](input []T, transform func(T) []U) []U {
    result := make([]U, 0, len(input))
    for _, item := range input {
        result = append(result, transform(item)...)
    }
    return result
}
```

**Lazy Implementation (Go 1.23+ Iterators):**
```go
import "iter"

func SelectMany[T, U any](seq iter.Seq[T], transform func(T) iter.Seq[U]) iter.Seq[U] {
    return func(yield func(U) bool) {
        for item := range seq {
            for innerItem := range transform(item) {
                if !yield(innerItem) {
                    return
                }
            }
        }
    }
}
```

## 26. DSLs in Go
Go's rigid syntax prevents "natural language" DSLs (like Ruby/Kotlin), but supports:
- **Fluent Interfaces**: Method chaining (`db.Where(...).Find(...)`).
- **Functional Options**: Declarative configuration (`NewServer(WithPort(80))`).
- **Struct Literals**: Declarative tree structures (`Elem("div", Attrs{...})`).

## 27. Functional Options Pattern (The "Keyword List" Alternative)
Go lacks Elixir-style keyword lists (`[key: val]`). The idiomatic alternative for optional/variadic configuration is the **Functional Options Pattern**.

```go
// Instead of: Configure(Opt{"port", 8080})
// Use:
func WithPort(p int) Option { ... }

server := NewServer(
    WithPort(8080),
    WithTimeout(30 * time.Second),
)
```
- **Pros**: Type-safe, discoverable API, clean call site.
- **Cons**: Requires boilerplate to define the Option functions.

## 28. Builder Pattern (Fluent API)
Used for constructing complex objects where configuration steps might fail or be optional. To handle errors in a chain, use the "Sticky Error" pattern.

```go
type ServerBuilder struct {
    server *Server
    err    error
}

func NewBuilder() *ServerBuilder {
    return &ServerBuilder{server: &Server{}}
}

func (b *ServerBuilder) Port(p int) *ServerBuilder {
    if b.err != nil { return b } // Skip if already failed
    if p < 0 {
        b.err = errors.New("invalid port")
        return b
    }
    b.server.Port = p
    return b
}

func (b *ServerBuilder) Build() (*Server, error) {
    if b.err != nil { return nil, b.err }
    return b.server, nil
}

// Usage
srv, err := NewBuilder().Port(8080).Timeout(5).Build()
```

## 29. Observer / Pub-Sub (Channels)
Go's concurrency primitives (Channels + Goroutines) make Pub-Sub trivial without heavy frameworks.

```go
type Broker[T any] struct {
    subscribers []chan T
}

func (b *Broker[T]) Subscribe() chan T {
    ch := make(chan T, 1) // Buffer to prevent blocking
    b.subscribers = append(b.subscribers, ch)
    return ch
}

func (b *Broker[T]) Publish(msg T) {
    for _, sub := range b.subscribers {
        // Non-blocking send to avoid stalling if a sub is slow
        select {
        case sub <- msg:
        default: // Drop message or handle overflow
        }
    }
}
```

## 30. Aggregate Pattern (Composition)
Go doesn't have inheritance. Aggregates are built using **Composition** and **Embedding**.

```go
type Entity struct {
    ID uuid.UUID
}

type User struct {
    Entity // Embedding: User "is-a" Entity (sort of)
    Name   string
}

// User now has ID automatically:
u := User{Name: "Alice"}
u.ID = uuid.New() // Access promoted field directly
```
**Trick**: Embedding interfaces allows a struct to "inherit" method signatures (and panic if called without implementation), useful for mocking.

### Ambiguous Selectors (Collision)
If two embedded structs have the same method name, the compiler cannot decide which one to use (Ambiguous Selector).
- **Direct Call Error**: `u.Start()` fails if both embedded types have `Start()`.
- **Resolution**: You must be explicit: `u.Entity.Start()` or `u.Worker.Start()`.
- **Shadowing**: If `User` implements `Start()`, it overrides both embedded methods.

## 31. Interface Embedding (Intersection)
Interfaces behave differently than structs when embedded.

```go
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }

// ReadWriter is the UNION of methods (Intersection of types)
type ReadWriter interface {
    Reader
    Writer
}
```

**Conflict Resolution (Go 1.14+)**:
- **Identical Signatures**: Allowed. They merge into a single method entry.
- **Different Signatures**: **Compile-time Error**. You cannot create the interface.
    - Unlike structs (which allow it but make it ambiguous to call), interfaces simply refuse to compile if there is a signature mismatch.

## 32. Discriminated Unions (Sum Types)
Go does not have native sum types (like Rust `enum` or TypeScript `|`).

**Idiomatic Emulation: Sealed Interfaces**
Define an interface with a private method so only local types can implement it.

```go
type Result interface {
    isResult() // Private: Seals the interface
}

type Success struct { Value int }
func (Success) isResult() {}

type Failure struct { Error error }
func (Failure) isResult() {}

func Handle(r Result) {
    switch v := r.(type) {
    case Success:
        // ...
    case Failure:
        // ...
    }
}
```
**Limitation**: No compile-time exhaustiveness check. If you add a new variant, the compiler won't warn you to update your switch statements.

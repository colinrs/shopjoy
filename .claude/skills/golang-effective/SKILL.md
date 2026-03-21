---
name: effective-go
description: Use when writing Go code to follow idiomatic patterns, avoid common mistakes, and apply industry best practices from Effective Go and Uber Style Guide
---

# Effective Go Skill

Write idiomatic, maintainable, and performant Go code by applying patterns from Effective Go, Uber Go Style Guide, and common mistake prevention.

## Quick Reference

### Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Package | lowercase, single word, no underscores | `bufio`, `strconv` |
| Interface (1 method) | method name + `-er` suffix | `Reader`, `Writer`, `Formatter` |
| Getter | field name (no `Get` prefix) | `obj.Owner()` not `obj.GetOwner()` |
| Setter | `Set` + field name | `obj.SetOwner(user)` |
| Acronyms | all caps or all lower | `URL`, `urlParser`, `HTTPServer` |
| Unexported constants | camelCase | `maxRetries`, `defaultTimeout` |

### Error Naming

| Type | Pattern | Example |
|------|---------|---------|
| Error types | `FooError` | `type ExitError struct {}` |
| Error values | `ErrFoo` / `ErrorFoo` | `var ErrorOrderNotFound = errs.NewCodeMsg(...)` |

### Error Handling Rules (Project-Specific)

| Rule | Description                                                                 |
|------|-----------------------------------------------------------------------------|
| Wrap with stack trace | Use `code.WithStack(err)` from `github.com/colinrs/shopjoy/pkg/code/`               |
| Business errors | Must be predefined constants in `github.com/colinrs/shopjoy/pkg/code`       |
| No inline errors | **PROHIBITED**: `fmt.Errorf(...)` or `errors.New(...)` in return statements |

## Core Patterns
ag
### 1. Error Handling

**Handle errors exactly once. Either log or return, never both.**

```go
// BAD: Handles error twice
func doThing(ctx context.Context) error {
    err := something()
    if err != nil {
        logx.WithContext(ctx).Errorf("something failed: %v", err) // logged here
        return err // AND returned here - caller may log again
    }
    return nil
}

// GOOD: Return with stack trace, let caller decide
func doThing(ctx context.Context) error {
    if err := something(); err != nil {
        return errs.WithStack(err)
    }
    return nil
}
```

**Use predefined error constants, not inline error creation.**

```go
// BAD: Inline error creation (PROHIBITED)
func GetOrder(ctx context.Context, db *gorm.DB, id int64) (*Order, error) {
    order, err := repo.GetByID(ctx, db, id)
    if err != nil {
        return nil, fmt.Errorf("order %d not found", id)  // ✗ NEVER do this
    }
    return order, nil
}

// BAD: errors.New inline (PROHIBITED)
func ValidateOrder(order *Order) error {
    if order.Amount.IsZero() {
        return errors.New("order amount cannot be zero")  // ✗ NEVER do this
    }
    return nil
}

// GOOD: Use predefined error constants from pkg/errorx/code.go
func GetOrder(ctx context.Context, db *gorm.DB, id int64) (*Order, error) {
    order, err := repo.GetByID(ctx, db, id)
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errorx.ErrorOrderNotFound  // Predefined constant
    }
    if err != nil {
        return nil, errs.WithStack(err)  // Wrap with stack trace
    }
    return order, nil
}

// GOOD: Predefined validation errors
func ValidateOrder(order *Order) error {
    if order.Amount.IsZero() {
        return errorx.ErrorOrderAmountInvalid  // Predefined constant
    }
    return nil
}
```

**Error wrapping for lower-level errors.**

```go
// BAD: String matching
if strings.Contains(err.Error(), "not found") {
    // fragile
}

// GOOD: Wrap with stack trace, check with errors.Is
if err := doSomething(ctx, db); err != nil {
    return errs.WithStack(err)  // Preserves stack trace
}

// Caller can check specific errors
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, errorx.ErrorNotFound
}
```

### 2. Reduce Nesting

**Handle errors first, keep happy path unindented.**

```go
// BAD: Happy path nested
func process(ctx context.Context, db *gorm.DB, data []byte) error {
    if len(data) > 0 {
        if err := validate(data); err == nil {
            if result, err := transform(data); err == nil {
                return save(ctx, db, result)
            } else {
                return errs.WithStack(err)
            }
        } else {
            return errs.WithStack(err)
        }
    }
    return errorx.ErrorEmptyData  // Predefined error constant
}

// GOOD: Early returns, flat structure
func process(ctx context.Context, db *gorm.DB, data []byte) error {
    if len(data) == 0 {
        return errorx.ErrorEmptyData  // Predefined error constant
    }
    if err := validate(data); err != nil {
        return errs.WithStack(err)
    }
    result, err := transform(data)
    if err != nil {
        return errs.WithStack(err)
    }
    return save(ctx, db, result)
}
```

### 3. Unnecessary Else

**If both branches set a variable, drop the else.**

```go
// BAD
var status string
if success {
    status = "ok"
} else {
    status = "failed"
}

// GOOD
status := "failed"
if success {
    status = "ok"
}
```

### 4. Interface Compliance

**Verify interface compliance at compile time.**

```go
// GOOD: Compile-time verification
type Handler struct{}

var _ http.Handler = (*Handler)(nil)       // *Handler implements http.Handler
var _ http.Handler = Handler{}             // Handler implements http.Handler (if value receiver)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
```

### 5. Nil Slices vs Empty Slices

**Prefer nil slices. Only use empty slices for JSON `[]` output.**

```go
// BAD: Unnecessary allocation
func filter(items []int) []int {
    result := []int{}  // allocates
    for _, item := range items {
        if item > 0 {
            result = append(result, item)
        }
    }
    return result
}

// GOOD: Nil slice, append handles it
func filter(items []int) []int {
    var result []int  // nil, no allocation
    for _, item := range items {
        if item > 0 {
            result = append(result, item)
        }
    }
    return result
}

// EXCEPTION: JSON encoding - use empty slice for []
type Response struct {
    Items []string `json:"items"`
}
resp := Response{Items: []string{}}  // JSON: {"items": []}
```

### 6. Copy Slices and Maps at Boundaries

**Slices and maps are references. Copy them to prevent unintended mutation.**

```go
// BAD: Retains reference to caller's slice
type Store struct {
    data []int
}

func (s *Store) SetData(data []int) {
    s.data = data  // caller can mutate s.data
}

// GOOD: Defensive copy
func (s *Store) SetData(data []int) {
    s.data = make([]int, len(data))
    copy(s.data, data)
}

// GOOD: Return copy to prevent mutation
func (s *Store) GetData() []int {
    result := make([]int, len(s.data))
    copy(result, s.data)
    return result
}
```

**Same applies to maps:**

```go
// BAD
func (c *Config) SetOptions(opts map[string]string) {
    c.options = opts
}

// GOOD
func (c *Config) SetOptions(opts map[string]string) {
    c.options = make(map[string]string, len(opts))
    for k, v := range opts {
        c.options[k] = v
    }
}
```

### 7. Specify Container Capacity

**Pre-allocate when size is known.**

```go
// BAD: Multiple reallocations
func collect(n int) []int {
    var result []int
    for i := 0; i < n; i++ {
        result = append(result, i)
    }
    return result
}

// GOOD: Single allocation
func collect(n int) []int {
    result := make([]int, 0, n)
    for i := 0; i < n; i++ {
        result = append(result, i)
    }
    return result
}

// Maps too
m := make(map[string]int, len(keys))
```

### 8. Prefer strconv Over fmt

**strconv is faster for primitive conversions.**

```go
// BAD: Slower
s := fmt.Sprintf("%d", 42)

// GOOD: ~3x faster
s := strconv.Itoa(42)

// Same for parsing
// BAD
var i int
fmt.Sscanf("42", "%d", &i)

// GOOD
i, _ := strconv.Atoi("42")
```

### 9. Avoid Pointers to Interfaces

**Interfaces are already pointer-like. Don't double-indirect.**

```go
// BAD: Pointer to interface
func process(r *io.Reader) {}

// GOOD: Interface directly
func process(r io.Reader) {}
```

### 10. Receiver Type Consistency

**Don't mix value and pointer receivers on a type.**

```go
// BAD: Mixed receivers
func (s *MyStruct) Method1() {}
func (s MyStruct) Method2() {}   // inconsistent

// GOOD: All pointer receivers (if any mutation or large struct)
func (s *MyStruct) Method1() {}
func (s *MyStruct) Method2() {}

// OR all value receivers (small, immutable)
func (p Point) Distance(q Point) float64 {}
func (p Point) String() string {}
```

### 11. Goroutines and Channel Ownership

**Let the producer own the channel. Producer creates, writes, closes.**

```go
// GOOD: Producer owns channel lifecycle
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)  // producer closes
        for _, n := range nums {
            out <- n
        }
    }()
    return out  // return receive-only channel
}

// Consumer just reads
func consume(in <-chan int) {
    for n := range in {
        fmt.Println(n)
    }
}
```

### 12. Goroutine Exit Conditions

**Always ensure goroutines can exit. Use context or done channels.**

```go
// BAD: Goroutine may leak
func watch(updates <-chan Update) {
    go func() {
        for update := range updates {
            process(update)
        }
    }()
}

// GOOD: Controllable exit
func watch(ctx context.Context, updates <-chan Update) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case update, ok := <-updates:
                if !ok {
                    return
                }
                process(update)
            }
        }
    }()
}
```

### 13. Wait for Goroutines

**Use sync.WaitGroup to wait for goroutine completion.**

```go
// GOOD: Proper synchronization
func processAll(items []Item) {
    var wg sync.WaitGroup
    for _, item := range items {
        wg.Add(1)
        go func(item Item) {
            defer wg.Done()
            process(item)
        }(item)
    }
    wg.Wait()
}
```

### 14. Functional Options

**Use functional options for flexible, extensible APIs.**

```go
type Server struct {
    addr    string
    timeout time.Duration
    logger  *log.Logger
}

type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithLogger(l *log.Logger) Option {
    return func(s *Server) {
        s.logger = l
    }
}

func NewServer(addr string, opts ...Option) *Server {
    s := &Server{
        addr:    addr,
        timeout: 30 * time.Second,  // default
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
srv := NewServer(":8080",
    WithTimeout(60*time.Second),
    WithLogger(logger),
)
```

### 15. Table-Driven Tests

**Use test tables for comprehensive, maintainable tests.**

```go
func TestSplit(t *testing.T) {
    tests := []struct {
        name  string
        input string
        sep   string
        want  []string
    }{
        {name: "simple", input: "a,b,c", sep: ",", want: []string{"a", "b", "c"}},
        {name: "empty", input: "", sep: ",", want: []string{""}},
        {name: "no sep", input: "abc", sep: ",", want: []string{"abc"}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := strings.Split(tt.input, tt.sep)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Split(%q, %q) = %v, want %v",
                    tt.input, tt.sep, got, tt.want)
            }
        })
    }
}
```

## Common Mistakes

### Loop Iterator Variables (Go < 1.22)

**In Go versions before 1.22, loop variables are reused. Capture them explicitly.**

```go
// BAD (Go < 1.22): All goroutines see final value
for _, item := range items {
    go func() {
        process(item)  // item is the same variable, captures final value
    }()
}

// GOOD (Go < 1.22): Capture in closure parameter
for _, item := range items {
    go func(item Item) {
        process(item)
    }(item)
}

// GOOD (Go < 1.22): Shadow with local variable
for _, item := range items {
    item := item  // shadows loop variable
    go func() {
        process(item)
    }()
}

// Go 1.22+: Loop variables are per-iteration, original code is safe
```

### Defer in Loops

**Defer runs at function exit, not loop iteration end.**

```go
// BAD: Files stay open until function returns
func processFiles(paths []string) error {
    for _, path := range paths {
        f, err := os.Open(path)
        if err != nil {
            return err
        }
        defer f.Close()  // won't close until function returns!
        // process f
    }
    return nil
}

// GOOD: Use closure to scope defer
func processFiles(paths []string) error {
    for _, path := range paths {
        if err := processFile(path); err != nil {
            return err
        }
    }
    return nil
}

func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()
    // process f
    return nil
}
```

### Nil Map Panic

**Writing to a nil map panics. Always initialize.**

```go
// BAD: Panics
var m map[string]int
m["key"] = 1  // panic!

// GOOD: Initialize
m := make(map[string]int)
m["key"] = 1

// GOOD: Or use composite literal
m := map[string]int{}
```

### Modifying Slice While Iterating

**Don't append/delete while ranging. Use index loop or build new slice.**

```go
// BAD: Undefined behavior
for i, v := range slice {
    if v < 0 {
        slice = append(slice[:i], slice[i+1:]...)
    }
}

// GOOD: Build new slice
result := slice[:0]  // reuse backing array
for _, v := range slice {
    if v >= 0 {
        result = append(result, v)
    }
}
```

## Checklist

Before submitting Go code:

- [ ] Package names are lowercase, single word
- [ ] Interfaces with one method end in `-er`
- [ ] No `Get` prefix on getters
- [ ] Errors handled exactly once (return OR log, not both)
- [ ] Errors wrapped with `errs.WithStack(err)` (NOT `fmt.Errorf`)
- [ ] Business errors use predefined constants from `pkg/errorx/code.go`
- [ ] No `fmt.Errorf(...)` or `errors.New(...)` in return statements
- [ ] Happy path is unindented (early returns for errors)
- [ ] No unnecessary `else` after `if-return`
- [ ] Interface compliance verified with `var _ Interface = (*Type)(nil)`
- [ ] Nil slices used (not `[]T{}`) unless JSON empty array needed
- [ ] Slices/maps copied at API boundaries
- [ ] Container capacity specified when size known
- [ ] `strconv` used instead of `fmt` for conversions
- [ ] No pointers to interfaces
- [ ] Receiver types consistent (all pointer or all value)
- [ ] Goroutines have clear exit conditions
- [ ] `sync.WaitGroup` used to wait for goroutines
- [ ] Loop variables captured correctly (Go < 1.22)
- [ ] Defer not used in loops (use helper function)
- [ ] Maps initialized before use
- [ ] Table-driven tests for multiple test cases

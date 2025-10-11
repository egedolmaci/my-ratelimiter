# Production-Grade Rate Limiter in Go

A flexible, extensible rate limiting system built with **Test-Driven Development**, demonstrating enterprise design patterns and concurrent programming expertise.

## 🎯 Overview

Production-ready rate limiter implementing multiple algorithms through the Strategy Pattern. Built entirely through Outside-In TDD with 25+ comprehensive tests covering acceptance, integration, and unit levels.

**Tech Stack**: Go 1.25 | Strategy Pattern | TDD | Concurrent Programming | HTTP Middleware

## 🏗️ Architecture

```
my-ratelimiter/
├── internal/ratelimiter/   # Core rate limiting engine with strategy pattern
├── internal/strategies/    # Pluggable algorithms (Fixed Window, Sliding Window Log)
├── pkg/middleware/         # HTTP middleware with dependency injection
└── examples/test-server/   # Working HTTP server demonstration
```

### Design Patterns Implemented
- **Strategy Pattern**: Pluggable rate limiting algorithms with clean interface
- **Interface Segregation**: Minimal, focused contracts for extensibility
- **Dependency Injection**: Testable components with mock time support
- **Factory Pattern**: Type-safe configuration and object creation

## 🔧 Implemented Features

### Rate Limiting Algorithms

**Fixed Window Strategy**
- Time-based window counting with automatic cleanup
- Optimized for high throughput (10,000+ concurrent requests)
- Scheduled background cleanup with goroutine lifecycle management

**Sliding Window Log Strategy**
- Timestamp-based precise rate limiting
- Eliminates boundary burst issues
- Hybrid cleanup (per-request + scheduled background)

**Configuration System**
- Type-safe config struct with strategy selection
- Factory methods for multiple initialization patterns
- Backward-compatible API design

### Production-Ready Features
- Thread-safe concurrent access with `sync.RWMutex`
- Memory management with automatic cleanup goroutines
- Graceful shutdown with channel-based lifecycle
- HTTP middleware with client IP extraction
- Mock time provider for deterministic testing

## 🧪 Test-Driven Development

Built using rigorous **Outside-In TDD** methodology:

```
Acceptance Tests (Public API) → Integration Tests → Unit Tests
```

**Test Coverage**:
- 25+ tests across 4 test suites (all passing)
- Acceptance tests through public API
- Integration tests for HTTP middleware
- Unit tests for strategy algorithms
- Concurrency tests with 10,000 goroutines
- Full end-to-end server integration

**TDD Process Applied**:
- RED: Acceptance test defines behavior through public API
- GREEN: Minimal implementation makes test pass
- REFACTOR: Clean up code, add unit tests, optimize

## 💻 Code Examples

### Strategy Pattern Interface
```go
type RateLimitStrategy interface {
    IsRequestAllowed(identifier string) (bool, int)
    Stop()
}
```

### Usage with Configuration
```go
config := ratelimiter.Config{
    Strategy:   "sliding_window_log",
    Limit:      100,
    WindowSize: time.Minute,
}
rl := ratelimiter.NewRatelimiterWithConfig(config)
defer rl.Stop()
```

### HTTP Middleware Integration
```go
middleware := middleware.Middleware{Ratelimiter: rl}
mux.HandleFunc("/api", middleware.RateLimitMiddleware(handler))
```

## 🚀 Running the Project

```bash
# Run the demo server
go run examples/test-server/main.go

# Test rate limiting
curl http://localhost:8080/ratelimited

# Run all tests
go test ./...

# Run with coverage
go test ./... -cover
```

## 🛠️ Technical Challenges Solved

**Concurrency Control**
- Debugged race conditions in concurrent map access
- Implemented proper lock granularity with `RWMutex`
- Designed goroutine lifecycle management with channels

**Network Programming**
- Handled TCP ephemeral port issues in client identification
- Implemented proper IP extraction from `RemoteAddr`
- Managed middleware state lifecycle

**Test Architecture**
- Eliminated test duplication through proper layering
- Implemented mock time provider for deterministic tests
- Solved HTTP route conflicts in integration tests

**Major Refactoring**
- Extracted fixed window algorithm into strategy pattern
- Maintained 100% backward compatibility during refactoring
- Achieved clean separation of concerns across packages

## 📊 Current Status

**Completed**:
- ✅ Strategy Pattern architecture with clean interfaces
- ✅ Fixed Window algorithm with thread safety
- ✅ Sliding Window Log algorithm with precision tracking
- ✅ Type-safe configuration system
- ✅ HTTP middleware with dependency injection
- ✅ Comprehensive test suite (25+ tests, all passing)
- ✅ Memory management with cleanup goroutines

**In Progress**:
- 🔄 Sliding Window Counter (hybrid algorithm)

**Planned**:
- ⏳ Token Bucket (burst traffic support)
- ⏳ Leaky Bucket (traffic smoothing)
- ⏳ Redis-backed distributed storage
- ⏳ Prometheus metrics integration

## 🎓 Skills Demonstrated

**Go Programming**
- Concurrent programming with goroutines and channels
- Interface design and composition
- Generics-free extensibility through interfaces
- Proper error handling and resource management

**Software Engineering**
- Test-Driven Development (Outside-In TDD)
- Design patterns (Strategy, Factory, Dependency Injection)
- Clean Architecture principles
- SOLID principles (especially Interface Segregation, Open/Closed)

**System Design**
- Rate limiting algorithms and trade-offs
- Memory vs accuracy optimization
- Concurrency control strategies
- HTTP middleware architecture

**Development Practices**
- Red-Green-Refactor workflow
- Acceptance-first testing
- Backward compatibility during refactoring
- Production-ready code with proper lifecycle management

---

*Built with TDD principles to demonstrate production-grade Go development and system design expertise.*

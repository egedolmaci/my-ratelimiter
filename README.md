# Enterprise Rate Limiter with Strategy Pattern

A production-ready rate limiting system built in Go, demonstrating **Strategy Pattern**, **Interface Segregation**, and **Test-Driven Development** - the same patterns used by GitHub, Netflix, and Cloudflare.

## 🎯 Project Overview

This project showcases a **flexible, extensible rate limiter** using Go, built entirely through Test-Driven Development methodology. Features a **Strategy Pattern architecture** that supports multiple rate limiting algorithms while maintaining clean interfaces and backward compatibility.

## 🏗️ Architecture & Design Patterns

### Clean Architecture Implementation
```
my-ratelimiter/
├── cmd/                    # Application entry points
├── internal/ratelimiter/   # Core rate limiting logic (Fixed Window algorithm)
├── pkg/middleware/         # HTTP middleware with dependency injection
└── examples/test-server/   # Integration examples and demos
```

### Key Design Patterns Applied

- **Strategy Pattern**: Pluggable rate limiting algorithms (Fixed Window, future: Sliding Window, Token Bucket)
- **Interface Segregation**: Clean separation between core functionality (`RateLimitStrategy`) and optional features (`StorageChecker`)
- **Dependency Injection**: Middleware accepts rate limiter interface for maximum testability
- **Factory Pattern**: Constructor functions ensure proper initialization and encapsulation
- **Single Responsibility**: Each strategy handles one algorithm, each package has focused purpose

## 🧪 Test-Driven Development (TDD)

**22 comprehensive tests** across 3 test suites demonstrate rigorous TDD approach:

### Test Coverage
- **Unit Tests**: Core rate limiting algorithm (`internal/ratelimiter/`)
- **Concurrency Tests**: Thread-safety validation with 10,000 goroutines (`concurenncy_test.go`)
- **Integration Tests**: HTTP middleware functionality (`pkg/middleware/`)
- **End-to-End Tests**: Full server integration (`examples/test-server/`)

### TDD Learning Journey
1. **Red**: Write failing tests first to define expected behavior
2. **Green**: Implement minimal code to pass tests
3. **Refactor**: Improve design while maintaining test coverage

### Latest TDD Achievement: Strategy Pattern Refactoring
Successfully refactored monolithic rate limiter into flexible strategy pattern through TDD:
- **Failing Test**: `TestRateLimiterWithFixedWindowStrategy` - strategy switching capability
- **Red Phase**: Tests failed - no strategy interface existed
- **Green Phase**: Implemented `RateLimitStrategy` interface and `FixedWindowStrategy`
- **Refactor Phase**: Extracted all algorithm logic into strategy while maintaining 100% backward compatibility
- **Result**: Extensible architecture ready for multiple algorithms (Sliding Window, Token Bucket, etc.)

Example TDD progression:
```go
// Test First: Define expected behavior
func TestIsRequestAllowed(t *testing.T) {
    t.Run("single request", func(t *testing.T) {
        rt := NewRateLimiter(1, time.Minute)
        allowed, _ := rt.IsRequestAllowed("user123")
        if !allowed {
            t.Error("First request should be allowed")
        }
    })
}
```

## 🔧 Technical Implementation

### Strategy Pattern Architecture
```go
// Core interface - minimal, focused contract
type RateLimitStrategy interface {
    IsRequestAllowed(identifier string) (bool, int)
    Stop()
}

// Extensible - easy to add new algorithms
type FixedWindowStrategy struct { /* implementation */ }
type SlidingWindowStrategy struct { /* future */ }
type TokenBucketStrategy struct { /* future */ }
```

### Current Implementation: Fixed Window Strategy
- **Algorithm**: Fixed time windows with request counting
- **Storage**: In-memory map with automatic cleanup
- **Thread Safety**: Full concurrent support with `sync.RWMutex` protection
- **Performance**: Handles 10,000+ concurrent requests safely
- **Encapsulation**: Strategy manages own lifecycle, storage, and cleanup

### HTTP Middleware Integration
```go
type Middleware struct {
    Ratelimiter Limiter // Interface for testability
}

func (m *Middleware) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract client IP (handling ephemeral ports)
        host, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            panic("Error while splitting the identifier address")
        }

        if allowed, remaining := m.Ratelimiter.IsRequestAllowed(host); allowed {
            w.Write([]byte(fmt.Sprintf("Remaining limit = %d", remaining)))
            next.ServeHTTP(w, r)
        } else {
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
        }
    })
}
```

## 📚 Key Learning Outcomes

### Software Engineering Practices
- **Interface Design**: Creating testable, loosely-coupled components
- **Error Handling**: Proper error propagation and panic recovery strategies
- **Network Programming**: Understanding TCP connections and ephemeral ports
- **HTTP Middleware**: Building reusable, configurable middleware components

### Testing Strategies
- **Test Isolation**: Each test creates fresh instances to avoid state pollution
- **Behavior Verification**: Testing outcomes rather than implementation details
- **Edge Case Coverage**: Window boundaries, rate limit exhaustion, time-based scenarios
- **Concurrency Testing**: 10,000 goroutines validating thread-safety under extreme load
- **Integration Testing**: Full HTTP request/response cycle validation

### Debugging & Problem Solving
- **Route Conflicts**: Learned to handle HTTP mux pattern conflicts in tests
- **Client Identification**: Solved ephemeral port issues in rate limiting
- **Concurrency Issues**: Debugged "concurrent map writes" and implemented thread-safe solutions
- **Lock Design**: Chose appropriate synchronization primitives (`sync.RWMutex`)
- **State Management**: Understanding shared vs. per-request state in middleware

## 🚀 Getting Started

### Running the Server
```bash
go run examples/test-server/main.go
```

### Testing Rate Limits
```bash
# Test the rate-limited endpoint
curl http://localhost:8080/ratelimited

# Monitor remaining requests
for i in {1..12}; do
  echo "Request $i:"
  curl http://localhost:8080/ratelimited
  echo -e "\n"
done
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run specific test suite
go test ./internal/ratelimiter -v

# Run concurrency test specifically
go test ./internal/ratelimiter -v -run TestConcurrentAccess
```

## 🛠️ Technical Challenges Overcome

1. **Middleware State Management**:
   - Problem: Creating fresh rate limiter per request
   - Solution: Dependency injection with server-scoped instances

2. **Client Identification**:
   - Problem: Ephemeral ports creating unique identifiers per request
   - Solution: Extract IP address only using `net.SplitHostPort()`

3. **Test Isolation**:
   - Problem: HTTP route conflicts in integration tests
   - Solution: Fresh server instances and unique route paths

## 📈 Next Phase: Multiple Strategies (In Progress)

**Roadmap**: Currently implementing additional rate limiting algorithms through TDD:

- **✅ Fixed Window**: Complete with cleanup and concurrency support
- **🔄 Sliding Window**: More accurate rate limiting (next implementation)
- **⏳ Token Bucket**: Burst traffic handling with sustained rates
- **⏳ Leaky Bucket**: Traffic smoothing for consistent output rates

**Future Enterprise Features**:
- **Distributed Storage**: Redis-based storage for multi-instance deployments
- **Metrics & Observability**: Prometheus integration with detailed rate limiting statistics
- **Dynamic Configuration**: Hot-reloadable YAML/JSON configuration

## 🎓 Skills Demonstrated

**Design Patterns & Architecture**:
- **Strategy Pattern**: Pluggable algorithms for flexible system design
- **Interface Segregation**: Clean separation of concerns and focused contracts
- **Factory Pattern**: Proper object initialization and encapsulation
- **Clean Architecture**: Dependency inversion and testable components

**Go Programming Excellence**:
- **Concurrent Programming**: Thread-safe implementations with proper synchronization
- **Interface Design**: Creating minimal, focused, and extensible interfaces
- **Error Handling**: Robust error propagation and recovery strategies
- **Testing**: Comprehensive test coverage with behavior-driven testing

**Software Engineering Practices**:
- **Test-Driven Development**: Red-Green-Refactor methodology with 100% backward compatibility
- **Refactoring**: Major architectural changes while maintaining functionality
- **Problem Solving**: Complex concurrency issues and network-level challenges

---

*This project demonstrates production-ready Go development practices through the lens of rate limiting - a critical component in modern distributed systems.*
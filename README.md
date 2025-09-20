# Rate Limiter Implementation

A comprehensive rate limiting system built in Go, demonstrating clean architecture, Test-Driven Development (TDD), and production-ready patterns.

## 🎯 Project Overview

This project showcases the implementation of a **Fixed Window Rate Limiter** using Go, built entirely through Test-Driven Development methodology. It demonstrates enterprise-level software engineering practices including dependency injection, interface design, and comprehensive testing strategies.

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

- **Dependency Injection**: Middleware accepts rate limiter interface for testability
- **Interface Segregation**: `Limiter` interface defines minimal contract
- **Single Responsibility**: Each package has a focused purpose
- **Factory Pattern**: Constructor functions ensure proper initialization

## 🧪 Test-Driven Development (TDD)

**21 comprehensive tests** across 3 test suites demonstrate rigorous TDD approach:

### Test Coverage
- **Unit Tests**: Core rate limiting algorithm (`internal/ratelimiter/`)
- **Integration Tests**: HTTP middleware functionality (`pkg/middleware/`)
- **End-to-End Tests**: Full server integration (`examples/test-server/`)

### TDD Learning Journey
1. **Red**: Write failing tests first to define expected behavior
2. **Green**: Implement minimal code to pass tests
3. **Refactor**: Improve design while maintaining test coverage

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

### Fixed Window Rate Limiter
- **Algorithm**: Fixed time windows with request counting
- **Storage**: In-memory map for request tracking
- **Thread Safety**: Single-threaded design with clear concurrency considerations
- **Time Handling**: Proper window expiration and reset logic

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
- **Integration Testing**: Full HTTP request/response cycle validation

### Debugging & Problem Solving
- **Route Conflicts**: Learned to handle HTTP mux pattern conflicts in tests
- **Client Identification**: Solved ephemeral port issues in rate limiting
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

## 📈 Future Enhancements

- **Distributed Rate Limiting**: Redis-based storage for multi-instance deployments
- **Algorithm Variations**: Sliding window, token bucket, leaky bucket implementations
- **Metrics & Monitoring**: Prometheus integration for rate limiting statistics
- **Configuration**: YAML/JSON configuration files for flexible rate limit rules

## 🎓 Skills Demonstrated

- **Go Programming**: Idiomatic Go code with proper error handling
- **Test-Driven Development**: Comprehensive test coverage with TDD methodology
- **Clean Architecture**: Separation of concerns and dependency management
- **HTTP Programming**: Middleware patterns and request/response handling
- **Problem Solving**: Debugging network-level issues and architectural challenges
- **Documentation**: Clear technical communication for future developers

---

*This project demonstrates production-ready Go development practices through the lens of rate limiting - a critical component in modern distributed systems.*
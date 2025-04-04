# TraceLite

TraceLite is a lightweight tracing library for Go applications that helps you track and analyze execution time and key points in your program. It provides a simple yet powerful way to create multiple trace points, mark events, and collect timing data.

## Features

- Lightweight and easy to integrate
- Support for multiple sub-traces
- Custom tags for traces
- Thread-safe operations
- Flexible data collection and formatting
- Millisecond precision timing
- Zero external dependencies

## Installation

```bash
go get github.com/OC0311/tracelite
```

## Usage Examples

### Basic HTTP Server Tracing

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Create and start trace
    trace := tracelite.NewTrace("http-request")
    trace.TraceOn()
    defer trace.TraceOff()
    
    // Add request information as tags
    trace.SetTags(map[string]interface{}{
        "path": r.URL.Path,
        "method": r.Method,
        "client_ip": r.RemoteAddr,
    })
    
    // Database operation trace
    trace.BeginTrace("database", map[string]interface{}{
        "operation": "query",
    })
    trace.Mark("database", "start", "Starting database query")
    // ... perform database operations ...
    trace.Mark("database", "end", "Database query completed")
    
    // Response trace
    trace.BeginTrace("response", nil)
    trace.Mark("response", "prepare", "Preparing response")
    // ... prepare response ...
    trace.Mark("response", "send", "Sending response")
    
    // Log trace results
    result := trace.Collect()
    log.Printf("Request processed in %dms\n", result.TotalCost)
}
```

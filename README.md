# File logger

```go
package go_logger

type Config struct {
	Dir string // default: ./_log/
	Level string // LevelDebug = 0 | LevelError = 1 | LevelInfo = 2
}

```
```go

package main

logger = NewLogger(Config{})

```

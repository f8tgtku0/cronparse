# cronparse

> Human-readable cron expression parser and validator with next-run prediction

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/cronparse.svg)](https://pkg.go.dev/github.com/yourusername/cronparse)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## Installation

```bash
go get github.com/yourusername/cronparse
```

---

## Usage

```go
package main

import (
    "fmt"
    "github.com/yourusername/cronparse"
)

func main() {
    expr, err := cronparse.Parse("*/5 9-17 * * MON-FRI")
    if err != nil {
        panic(err)
    }

    // Human-readable description
    fmt.Println(expr.Describe())
    // Output: Every 5 minutes, between 09:00 and 17:00, Monday through Friday

    // Validate an expression
    if err := cronparse.Validate("0 12 * * *"); err != nil {
        fmt.Println("Invalid:", err)
    }

    // Predict the next run time
    next := expr.Next(time.Now())
    fmt.Println("Next run:", next)
}
```

### Supported Syntax

| Field  | Allowed Values | Special Characters |
|--------|----------------|--------------------|
| Minute | 0–59           | `*` `/` `,` `-`    |
| Hour   | 0–23           | `*` `/` `,` `-`    |
| Day    | 1–31           | `*` `/` `,` `-`    |
| Month  | 1–12           | `*` `/` `,` `-`    |
| Weekday| 0–7, SUN–SAT  | `*` `/` `,` `-`    |

---

## License

This project is licensed under the [MIT License](LICENSE).
# rat-formula-parser
generate (*math/big).Rat from string formula.

```go
package main

import (
    "fmt"
)

func main() {
    rt, err := NewRatFromFormula(Formula{Num:"0.01 * 100"})
    if err != nil {
        panic(err)
    }

    fmt.Println(rt.RatString()) // 1
}
```
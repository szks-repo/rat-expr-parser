# rat-expr-parser
generate (*math/big).Rat from string expression.

```go
package main

import (
    "fmt"
)

func main() {
    {
        rt, err := NewRatFromExpr(Expr{Num: "0.01 * 100"})
        if err != nil {
            panic(err)
        }

        fmt.Println(rt.RatString()) // "1"
    }
    {
        rt, err := NewRatFromExpr(Expr{Num: "10 ** 2", Denom: "10 ** 3"})
        if err != nil {
            panic(err)
        }

        fmt.Println(rt.RatString()) // "1/10"
    }
}
```
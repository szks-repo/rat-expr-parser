package parser

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestRatFromExpr(t *testing.T) {
	t.Parallel()

	for i, tt := range []struct {
		expr                 Expr
		wantErr              bool
		wantRatStringFactory func() string
	}{
		{
			expr:    Expr{Num: "1", Denom: "100"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "1/100"
			},
		},
		{
			expr:    Expr{Num: "(10 + 100) * 3", Denom: "30"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "11"
			},
		},
		{
			expr:    Expr{Num: "((10 + 100)) * (3)", Denom: "30"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "11"
			},
		},
		{
			expr:    Expr{Num: "(10 + 100) * 100", Denom: "100+10"},
			wantErr: false,
			wantRatStringFactory: func() string {
				a := new(big.Rat).SetInt64((10 + 100) * 100)
				b := new(big.Rat).SetInt64(100 + 10)
				c := new(big.Rat).Quo(a, b)
				return c.RatString()
			},
		},
		{
			expr:    Expr{Num: "10 * (2 + 3) - 5 / 1", Denom: "5 * (1 + 1)"},
			wantErr: false,
			wantRatStringFactory: func() string {
				a := new(big.Rat).SetInt64(10*(2+3) - 5/1)
				b := new(big.Rat).SetInt64(5 * (1 + 1))
				c := new(big.Rat).Quo(a, b)
				return c.RatString()
			},
		},
		{
			expr:    Expr{Num: "10 * 5 + 2"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "52"
			},
		},
		{
			expr:    Expr{Num: "0.1 + 0.2"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "3/10"
			},
		},
		{
			expr:    Expr{Num: "10.5", Denom: "0.5"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "21"
			},
		},
		{
			expr:    Expr{Num: "100", Denom: "0.01"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "10000"
			},
		},
		{
			expr:    Expr{Num: "0.01 * 100"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "1"
			},
		},
		{
			expr:    Expr{Num: ".01*-0.01"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "-1/10000"
			},
		},
		{
			expr:    Expr{Num: "1.0 / 3.0"},
			wantErr: false,
			wantRatStringFactory: func() string {
				a, _ := new(big.Rat).SetString("1.0")
				b, _ := new(big.Rat).SetString("3.0")
				c := new(big.Rat).Quo(a, b)
				return c.RatString()
			},
		},
		{
			expr:    Expr{Num: "5 % 5"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "0"
			},
		},
		{
			expr:    Expr{Num: "2 * 5 % 3"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "1"
			},
		},
		{
			expr:    Expr{Num: "2 ** 3 % 5", Denom: "1"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "3"
			},
		},
		{
			expr:    Expr{Num: "2 * 3 ** 2"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "18"
			},
		},
		{
			expr:    Expr{Num: "5 ** 0"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "1"
			},
		},
		{
			expr:    Expr{Num: "0 ** 0"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "1"
			},
		},
		{
			expr:    Expr{Num: "2 ** -2"},
			wantErr: false,
			wantRatStringFactory: func() string {
				return "1/4"
			},
		},
	} {
		t.Run(fmt.Sprintf("%d[%s]", i+1, tt.expr.String()), func(t *testing.T) {
			got, err := NewRatFromExpr(tt.expr)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RatFromExpr() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if tt.wantErr {
					t.Errorf("RatFromExpr() = %v, wantErr %v", got, tt.wantErr)
				} else {

					if got.RatString() != tt.wantRatStringFactory() {
						t.Errorf("RatFromExpr() = %v, want %v", got.RatString(), tt.wantRatStringFactory())
					}
					wantRat, ok := new(big.Rat).SetString(tt.wantRatStringFactory())
					if !ok {
						t.Fatal("SetString() error")
					}
					{
						a, _ := got.Float64()
						b, _ := wantRat.Float64()
						a = math.Ceil(a)
						b = math.Ceil(b)
						if a != b {
							t.Errorf("math.Ceil(RatFromExpr()) = %f, want %f", a, b)
						}
					}
					{
						a, _ := got.Float64()
						b, _ := wantRat.Float64()
						a = math.Round(a)
						b = math.Round(b)
						if a != b {
							t.Errorf("math.Round(RatFromExpr()) = %f, want %f", a, b)
						}
					}
					{
						a, _ := got.Float64()
						b, _ := wantRat.Float64()
						a = math.Floor(a)
						b = math.Floor(b)
						if a != b {
							t.Errorf("math.Floor(RatFromExpr()) = %f, want %f", a, b)
						}
					}
				}
			}
		})
	}
}

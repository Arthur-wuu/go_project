package common

import (
	"fmt"
	"math/big"
)

type Calc struct {
	num *big.Int
}

func NewCalc(x *big.Int) *Calc {
	return &Calc{x}
}

func (c *Calc) Add(xs ...*big.Int) *Calc {
	for _, x := range xs {
		c.num = big.NewInt(0).Add(c.num, x)
	}
	return c
}

func (c *Calc) Sub(xs ...*big.Int) *Calc {
	for _, x := range xs {
		c.num = big.NewInt(0).Sub(c.num, x)
	}
	return c
}

func (c *Calc) Mul(xs ...*big.Int) *Calc {
	for _, x := range xs {
		c.num = big.NewInt(0).Mul(c.num, x)
	}
	return c
}

func (c *Calc) Div(xs ...*big.Int) *Calc {
	for _, x := range xs {
		c.num = big.NewInt(0).Div(c.num, x)
	}
	return c
}

func (c *Calc) Pow(n int) *Calc {
	ret := big.NewInt(1)
	for n != 0 {
		if n%2 != 0 {
			ret = big.NewInt(0).Mul(ret, c.num)
		}
		n /= 2
		c.num = big.NewInt(0).Mul(c.num, c.num)
	}

	c.num = ret
	return c
}

func (c *Calc) Get() *big.Int {
	return c.num
}

func fc() {
	fmt.Printf("%#v", NewCalc(big.NewInt(10)).Pow(8).Get())
}

//
//func Mul(x, y *big.Int) *big.Int {
//	21     return big.NewInt(0).Mul(x, y)
//	22 }
//23 func Sub(x, y *big.Int) *big.Int {
//	24     return big.NewInt(0).Sub(x, y)
//	25 }
//26 func Add(x, y *big.Int) *big.Int {
//	27     return big.NewInt(0).Add(x, y)
//	28 }
//29 func Div(x, y *big.Int) *big.Int {
//	30     return big.NewInt(0).Div(x, y)

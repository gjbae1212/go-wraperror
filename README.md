# go-wraperror

<p align="center">
<a href="https://hits.seeyoufarm.com"><img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fgjbae1212%2Fgo-wraperror&count_bg=%2379C83D&title_bg=%23555555&icon=go.svg&icon_color=%2308BEB8&title=hits&edge_flat=false"/></a>
<a href="/LICENSE"><img src="https://img.shields.io/badge/license-MIT-GREEN.svg" alt="license"/></a> 
</p>

WrapError is custom error struct implemented error interface and supporting errors.As, errors.Is, Unwrap ... and so on.

## Getting Started
```go
package main
import (
    "os"
   "errors"
   wraperror "github.com/gjbae1212/go-wraperror"
)


func main() {
  sample1 := errors.New("[err] tests")
  sample2 := &os.PathError{}
  sample3 := &os.SyscallError{}
  
  wrap := wraperror.Error(sample)
  wrap = wrap.Error(sample2)
  
  // true
  errors.Is(wrap, sample1)
  errors.As(wrap, &sample1)
  // true
  errors.Is(wrap, sample2)
  errors.As(wrap, &sample2)
  // false
  errors.Is(wrap, sample3)
  errors.As(wrap, &sample3)
}
``` 

## LICENSE
This project is following The MIT.

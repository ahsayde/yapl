# Getting Started

## Installation

```bash
go get github.com/ahsayde/yapl/yapl
```


## Usage

```go
import (
    "fmt"
	"io/ioutil"
	"github.com/ahsayde/yapl/yapl"
)

func main() {
    input := map[string]interface{}{
        "method": "GET",
        "endpoint": "/users"
    }

    params := map[string]interface{}{
        "method": "GET"
    }

    raw, err := ioutil.ReadFile("policy.yaml")
	if err != nil {
		panic(err)
	}

    policy, err := yapl.Parse(raw)
    if err != nil {
        panic(err)
    }

    result, err := policy.Eval(input, params)
    if err != nil {
        panic(err)
    }

    fmt.Println(result)
}

```


## Writing Policies

See [Policy Syntax](./syntax/README.md) documentation to learn how to write `yapl` policies.

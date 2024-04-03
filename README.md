# nazhard/cli

---

Cute and easy-to-use CLI Framework! Made as simple as possible with a few features...

I'm not sure you'll actually use this simple framework.

Because this framework is just too cute. It's like an adorable loli!

### Example

```go
package main

import (
    "fmt"
    "os"

    "github.com/nazhard/cli"
)

func main() {
    app := cli.App{}
    app.Name = "moe"

    cmd := &cli.Command{
        Name:        "nii-chan",
        Usage:       "[message]",
        Description: "print cute message to onii-chan!",
        Action: func(ctx cli.Context) {
            fmt.Println("Helllo Oni-chan!!")
        },
    }

    app.AddCommand(cmd)
    app.Run(os.Args)
}
```

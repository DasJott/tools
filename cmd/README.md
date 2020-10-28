# CMD - start different parts with a cli parameter
Command handler for managing start parameter commands.

## Usage
Implement the `Command` interface in your package:
```go
type Command interface {
    // Init is called first
    Init()
    // Start is called right after and returns an error, to be returned to caller of cmd.Start()
    Start() error
    // Clean is called last
    Clean()
}
```
Then register that struct along with a certain parameter name.<br>
**Make sure, you provide a pointer!**<br>
You can also use empty command, if you wish.
```go
package main

import (
    "fmt"
    "cleverreach.com/crtools/cmd"
    "mydomain/myproject/mypkg"
)

type Command struct {}

func main() {
    // register something in this file
    cmd.RegisterCmd("", &Command{})

    // register a command in another package
    cmd.RegisterCmd("mine", new(mypkg.Command))

    err := cmd.Start()
    fmt.Println(err) // passed error
}

// functions for the interface to be implemented
func (c *Command) Init() {}
func (c *Command) Start() error {
    // do stuff here
    fmt.Println("I just got started, yay!")
    return nil
}
func (c *Command) Clean() {}
```

## Self registering packages
You can also put the register call into the package itself.<br>
We'll use the `init()` function for this.

```go
package app

import "cleverreach.com/crtools/cmd"

type Command struct {}

func init() {
    cmd.RegisterCmd("app", new(Command))
}

func (c *Command) Init() {}
func (c *Command) Start() (err error) {
    // do stuff here
    return err
}
func (c *Command) Clean() {}
```

You still have to start this thingamabob from within main
```go
package main

import "cleverreach.com/crtools/cmd"

func main() {
    err := cmd.Start()
    // err handling here
}

```

## Additional
- `cmd.PanicEmptyCommand = true`<br>
  panics if no matching cmd was found

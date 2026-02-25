# quake3-rcon for Golang

Tiny Quake3 RCON library for Go.

It was inspired by the excellent Node.js library:
https://github.com/thbaumbach/node-quake3-rcon

The goal of this project is to provide the same simplicity and tiny footprint, but in pure Go.

---

## Features

* tiny and clean implementation
* no dependencies
* supports multi-packet responses
* configurable timeout and quiet window
* debug logging support
* interactive CLI example included

---

## Installation

Clone the repository:

```bash
git clone https://github.com/mambuzrrr/golang-quake3-rcon.git
cd golang-quake3-rcon
```

Initialize module if needed:

```bash
go mod tidy
```

---

## Configuration

Create a `config.json` in the project root:

```json
{
  "address": "127.0.0.1:28960",
  "password": "your_rcon_password",
  "debug": false,
  "timeout_ms": 2000,
  "quiet_ms": 180
}
```

---

## Example Usage

Run the status example:

```bash
go run ./examples/status
```

Example output:

```
map: mp_railyard
num score ping name            lastmsg address
--- ----- ---- --------------- ------- ---------------------
0     8   50 player           0       127.0.0.1:28960
```

---

## CLI Example

Run the CLI:

```bash
go run ./examples/cli
```

Then type commands:

```
> status
> map_restart
> say hello noobs
> exit (exit is to leave the CLI, nothing todo with ingame :D)
```

---

## Library Usage

Example:

```go
package main

import (
    "fmt"
    "q3rcon/q3rcon"
)

func main() {
    client, _ := q3rcon.New("127.0.0.1:28960", "rconpassword")

    resp, _ := client.Send("status") // Change Command here... like map_restart blabla...

    fmt.Println(resp)
}
```

---

## Why this exists

This library exists to provide a minimal, clean, and reliable Go implementation of Quake3 RCON, inspired by the Node.js version but without runtime dependencies.

---

## License

MIT License

---

## Credits

Inspired by:
https://github.com/thbaumbach/node-quake3-rcon

Thanks to the original authors.

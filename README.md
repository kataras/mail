# Mail __/ email sender written in Go

This simple package supports rich email messages, unix `sendmail` command and more.

The `mail` package provides an amazing API to work with.

<a href="https://travis-ci.org/kataras/mail"><img src="https://img.shields.io/travis/kataras/mail.svg?style=flat-square" alt="Build Status"></a>
<a href="https://github.com/kataras/mail/blob/master/LICENSE"><img src="https://img.shields.io/badge/%20license-MIT%20%20License%20-E91E63.svg?style=flat-square" alt="License"></a>
<a href="https://github.com/kataras/mail/releases"><img src="https://img.shields.io/badge/%20release%20-%20v0.0.1-blue.svg?style=flat-square" alt="Releases"></a>
<a href="https://godoc.org/github.com/kataras/mail"><img src="https://img.shields.io/badge/%20docs-reference-5272B4.svg?style=flat-square" alt="Godocs"></a>
<a href="https://kataras.rocket.chat/channel/mail"><img src="https://img.shields.io/badge/%20community-chat-00BCD4.svg?style=flat-square" alt="Build Status"></a>
<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
<a href="#"><img src="https://img.shields.io/badge/platform-Any--OS-yellow.svg?style=flat-square" alt="Platforms"></a>

## Installation

The only requirement is the [Go Programming Language](https://golang.org/dl),  at least version 1.9.

```bash
$ go get -u -v github.com/kataras/mail
```

> Stable release installation by `go get gopkg.in/kataras/mail.v0`

## Getting Started

- `New` returns a new, e-mail sender service
- `Mail#Send` sends a (rich) email message
- `Mail#Subject("...")` returns a `Builder` for the mail message which ends up with `Send() error`, `SendUNIX() error`
- `SendUNIX` make use of the `sendmail` program of *nix OS, no `Mail` sender is needed

```go
// New returns a new *Mail, which contains the `Send(...) error`
// and `Subject(...) *Builder` functions.
New(c Credentials) *Mail
```

### Example

```sh
$ cat send-mail.go
```

```go
package main

import "github.com/kataras/mail"

func main() {
    c := mail.Credentials{
        Addr:     "smtp.sendgrid.net:587",
        Username: "apikey",
        Password: "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",
    }

    message, err := mail.New(c)
    if err != nil {
        panic(err)
    }

    message.
        Subject("Subject").
        BodyString("Body").
        To("receipt@example.com", "receipt2@example.com").
        From("FromName", "from@example.com").
        Send()

    // Tip #1
    //
    // Alternative ways to set Body inside a Builder:
    // Body([]byte)
    // BodyReader(io.Reader)
    // BodyReadCloser(io.ReadCloser)

    // Tip #2
    //
    // If you want to re-use a Builder(= `message.Subject(...)`'s result) after the `Send`
    // then you have to call the Builder's `MarkSingleton()` before its `Send` execution.

    // Small:
    //
    // [ init time, once after the `mail.New(...)` ]
    // message.DefaultFrom = &mail.Address{"FromName", "from@example.com"}
    // [[ run time, many ]]
    // message.Subject("Subject").BodyString("Body").To("receipt@example.com").Send()
}
```

> For the stable release use `import gopkg.in/kataras/mail.v0` instead

```sh
$ go run send-mail.go
```

### More examples

- [Simple](_examples/simple/main.go)
- [Using the Builder](_examples/builder/main.go)
- [Using the Unix `sendmail` program](_examples/unix/main.go)

## FAQ

Explore [these questions](https://github.com/kataras/mail/issues?mail=label%3Aquestion) or navigate to the [community chat](https://kataras.rocket.chat/channel/mail).

## Versioning

Current: **v0.0.1**

Read more about Semantic Versioning 2.0.0

 - http://semver.org/
 - https://en.wikipedia.org/wiki/Software_versioning
 - https://wiki.debian.org/UpstreamGuide#Releases_and_Versions

## People

The author of the `mail` is [@kataras](https://github.com/kataras).

## Contributing

If you are interested in contributing to the `mail` project, please make a PR.

### TODO

- [ ] Add a simple CLI tool for sending emails
- [ ] Read the specification for the email attachment and implement that.

## License

This project is licensed under the MIT License. License file can be found [here](LICENSE).

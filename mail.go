// Package mail is a simple email sender written in Go.
// Please refer to the https://github.com/kataras/mail/tree/master/_examples folder for more.
package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/mail"
	"net/smtp"
	"os/exec"
	"strings"
	"sync"

	"github.com/valyala/bytebufferpool"
)

const (
	// Version current semantic version of the `mail` package.
	Version = "0.0.1"
)

// Credentials are the SMTP credentials.
type Credentials struct {
	// Addr is the server mail server's full host, IP or host:port,
	//
	// if port is missing then `:smtp` will be used.
	//
	// Required.
	Addr string
	// Username is the auth username@domain.com for the sender.
	//
	/// Required.
	Username string
	// Password is the auth password for the sender.
	//
	// Required.
	Password string
}

// Mail is the main structure of this package, it's the mail sender.
type Mail struct {
	Addr string
	// DefaultFrom is being used if from Address is missing from the
	// Send functions, it's the "Username <Username@host>".
	DefaultFrom *Address
	// Auth is exported so caller can change it to another.
	//
	// Defaults to smtp.PlainAuth based on the Credential's Addr, Username and Password.
	Auth smtp.Auth
}

// New returns a new Mail instance based on the given `Credentials`.
func New(c Credentials) (*Mail, error) {
	if !strings.Contains(c.Addr, ":") {
		c.Addr += ":smtp"
	}

	host, _, err := net.SplitHostPort(c.Addr)
	if err != nil {
		return nil, err
	}

	sender := &Mail{
		Addr: c.Addr,
		DefaultFrom: &Address{
			Name:    c.Username,
			Address: c.Username + "@" + host,
		},
		Auth: smtp.PlainAuth("", c.Username, c.Password, host),
	}

	return sender, nil
}

const (
	contentTypeHTML         = `text/html; charset=\"utf-8\"`
	mimeVer                 = "1.0"
	contentTransferEncoding = "base64"
)

var bufPool bytebufferpool.Pool

// Address is an alias of `net/mail#Address`.
type Address = mail.Address

// ParseAddress parses a string to an `Address`.
//
// Usage:
// addr, err := mail.ParseAddress("Gerasimos <gerasimos@example.com>")
// addr.Name is "Gerasimos" and addr.Address is "gerasimos@example.com".
//
// ParseAddress is just a shortcut of `net/mail#ParseAddress`.
var ParseAddress = mail.ParseAddress

type stringWriter interface {
	WriteString(string) (int, error)
}

func writeHeaders(w stringWriter, subject string, body []byte, to []string) {
	w.WriteString(fmt.Sprintf("%s: %s\r\n", "To", strings.Join(to, ",")))
	w.WriteString(fmt.Sprintf("%s: %s\r\n", "Subject", subject))
	w.WriteString(fmt.Sprintf("%s: %s\r\n", "MIME-Version", mimeVer))
	w.WriteString(fmt.Sprintf("%s: %s\r\n", "Content-Type", contentTypeHTML))
	w.WriteString(fmt.Sprintf("%s: %s\r\n", "Content-Transfer-Encoding", contentTransferEncoding))
	w.WriteString(fmt.Sprintf("\r\n%s", base64.StdEncoding.EncodeToString(body)))
}

// Send sends an email to the recipient(s)
// the body can be in HTML format as well.
func (m *Mail) Send(from *Address, subject string, body []byte, to ...string) error {
	if from.Address == "" {
		// from.Name can be empty but Address not.
		return fmt.Errorf("from address is required")
	}

	buffer := bufPool.Get()
	defer bufPool.Put(buffer)
	buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "From", from.String()))
	writeHeaders(buffer, subject, body, to)

	return smtp.SendMail(
		m.Addr,
		m.Auth,
		from.Name,
		to,
		buffer.Bytes(),
	)
}

// Subject returns thae `Builder` with a filled mail subject, from there
// the caller can build the e-mail and Send the e-mail from the `Builder` instance.
func (m *Mail) Subject(subject string) *Builder {
	return acquireBuilder(m).Subject(subject)
}

// Builder is the builder of the e-mail headers and message body,
// it resets the current Builder instance each time `Send` is called to avoid memory allocations.
//
// See `Mail#Subject`.
type Builder struct {
	// from is not resetable so it can be re-used even if it's missing,
	// if never setted then the DefaultFrom from the *Mail will be used.
	from *Address

	subject    string
	body       []byte
	recipients []string

	singleton bool

	sender *Mail
}

var builderPool sync.Pool

func acquireBuilder(m *Mail) *Builder {
	v := builderPool.Get()
	if v != nil {
		b := v.(*Builder)

		return b
	}
	return &Builder{sender: m}
}

func releaseBuilder(b *Builder) {
	// we don't reset yet, so a single builder can be used many times,
	// we reset on the acquireBuilder.
	b.subject = ""
	b.body = b.body[0:0]
	b.recipients = b.recipients[0:0]
	b.from = nil // reset that as well.
	b.singleton = false

	builderPool.Put(b)
}

// MarkSingleton will make this Builder re-usable, even after the `Send` or `SendUNIX` functions.
func (b *Builder) MarkSingleton() *Builder {
	b.singleton = true
	return b
}

// From is from address header, it's not required.
// If not setted then the `Mail.DefaultFrom` will be used instead.
//
// Accepts two input arguments:
// name is the proper name; may be empty.
// address is the full address; user@domain.
func (b *Builder) From(name, address string) *Builder {
	b.from = &mail.Address{Name: name, Address: address}
	return b
}

// Subject ses the subject of the mail header.
func (b *Builder) Subject(subject string) *Builder {
	b.subject = subject
	return b
}

// Body sets the body of the mail.
func (b *Builder) Body(body []byte) *Builder {
	b.body = body
	return b
}

// AppendBody adds more body to the body of the mail.
func (b *Builder) AppendBody(bodyData []byte) *Builder {
	return b.Body(append(b.body, bodyData...))
}

// BodyReader same as `Body` but it accepts
// an io.Reader to read the actual body from and set.
func (b *Builder) BodyReader(r io.Reader) *Builder {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return b
	}
	return b.Body(body)
}

// BodyReadCloser same as `BodyReader` but it closes the reader at the end.
func (b *Builder) BodyReadCloser(r io.ReadCloser) *Builder {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return b
	}

	if err = r.Close(); err != nil {
		return b
	}

	return b.Body(body)
}

// BodyString sets a body like `Body` but it accepts a string instead of []byte.
func (b *Builder) BodyString(body string) *Builder {
	return b.Body([]byte(body))
}

// To adds a recipient to the mail header, it can be called multiple times.
func (b *Builder) To(recipients ...string) *Builder {
	b.recipients = append(b.recipients, recipients...)
	return b
}

// Send sends the e-mail based on the subject, body and recipients.
// It resets the current Builder on finish.
func (b *Builder) Send() error {
	if b.from == nil {
		b.from = b.sender.DefaultFrom
	}
	// Good idea but unexpected so not: If there is an error the Builder is not released, so the caller can retry.
	err := b.sender.Send(b.from, b.subject, b.body, b.recipients...)

	if !b.singleton {
		// if it's not singleton (default behavior)
		// then it will put back to the pool.
		releaseBuilder(b)
	}

	return err
}

// SendUNIX will call the "sendmail" unix command and send the e-mail
// based on the subject, body and recipients.
// If the "sendmail" command is not part of the host's(current machine)
// programs then it will return an error.
// It resets the current Builder on finish.
//
// Note that this function can be used from a *nux operating system host only,
// windows operating system doesn't support this.
// Don't forget to make sure to configure the machine's mail client before make use of this feature.
func (b *Builder) SendUNIX() error {
	if b.from == nil {
		b.from = b.sender.DefaultFrom
	}
	err := SendUNIX(b.from, b.subject, b.body, b.recipients...)
	if !b.singleton {
		// if it's not singleton (default behavior)
		// then it will put back to the pool.
		releaseBuilder(b)
	}

	return err
}

// SendUNIX will call the "sendmail" unix command and send the e-mail
// based on the given subject, body and recipients (to).
// If the "sendmail" command is not part of the host's(current machine)
// programs then it will return an error.
// It resets the current Builder on finish.
//
// Note that this function can be used from a *nux operating system host only,
// windows operating system doesn't support this.
// Don't forget to make sure to configure the machine's mail client before make use of this feature.
func SendUNIX(from *Address, subject string, body []byte, to ...string) error {
	buffer := new(bytes.Buffer)
	writeHeaders(buffer, subject, body, to)

	cmd := exec.Command("sendmail", "-F", from.Name, "-f", from.Address, "-t")
	cmd.Stdin = buffer
	_, err := cmd.CombinedOutput()
	return err
}

package main

import "github.com/kataras/mail"

func main() {
	// SendUNIX will call the "sendmail" unix command and send the e-mail
	// based on the given subject, body and recipients (to).
	// If the "sendmail" command is not part of the host's(current machine)
	// programs then it will return an error.
	// It resets the current Builder on finish.
	//
	// Note that this function can be used from a *nux operating system host only,
	// windows operating system doesn't support this.
	// Don't forget to make sure to configure the machine's mail client before make use of this feature.
	err := mail.SendUNIX(
		&mail.Address{
			Name:    "Example",
			Address: "Example@hotmail.com",
		},
		// message subject.
		"Hello subject",
		// message body.
		[]byte(`<h1>Hello</h1> <br/><br/> <span style="color:red">This is the rich message body </span>`),
		// message receipts.
		"kataras2006@hotmail.com", "kataras2007@hotmail.com")

	if err != nil {
		println("error while sending the e-mail: " + err.Error())
	}
}

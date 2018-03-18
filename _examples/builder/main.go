package main

import "github.com/kataras/mail"

func main() {
	// sender credentials.
	c := mail.Credentials{
		Addr:     "smtp.sendgrid.net:587",
		Username: "apikey",
		Password: "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",
	}

	sender, err := mail.New(c)
	if err != nil {
		panic(err)
	}

	// An, optional, default from address if `Form` is missing.
	// sender.DefaultFrom = &mail.Address{
	// 	Name:    "FromName",
	// 	Address: "from@example.com",
	// }

	// the subject/title of the e-mail.
	err = sender.
		Subject("Hello subject").
		BodyString(`<h1>Hello</h1> <br/><br/> <span style="color:red">
					This is the rich message body </span>`).
		To("kataras2006@hotmail.com", "kataras2007@hotmail.com").
		// the from address, if nil then the `sender.DefaultFrom` will be used.
		From("FromName", "from@example.com").
		Send()

	if err != nil {
		println("error while sending the e-mail: " + err.Error())
	}

	// Tip #1
	//
	// Alternative ways to set Body inside a Builder:
	// Body([]byte)
	// BodyReader(io.Reader)
	// BodyReadCloser(io.ReadCloser)

	// Tip #2
	//
	// If you want to re-use a Builder(= `sender.Subject(...)`'s result) after the `Send`
	// then you have to call the Builder's `.MarkSingleton()` before its `Send` execution.

	// Small:
	//
	// [ init time, once after the `mail.New(...)` ]
	// sender.DefaultFrom = &mail.Address{"FromName", "from@example.com"}
	// [[ run time, many ]]
	// sender.Subject("Subject").BodyString("Body").To("receipt@example.com").Send()
}

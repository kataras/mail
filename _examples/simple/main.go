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

	// An, optional, default from address if `form` input argument of the `Send` functions is missing.
	// sender.DefaultFrom = &mail.Address{
	// 	Name:    "Example",
	// 	Address: "example@example.com",
	// }

	err = sender.Send(
		// the from address, if nil then the `sender.DefaultFrom` will be used.
		&mail.Address{
			Name:    "FromName",
			Address: "from@example.com",
		},
		// message subject.
		"Hello subject",
		// message body, only []byte but Builder method has shortcuts for `string`, `io.Reader` and `io.ReadCloser`.
		[]byte(`<h1>Hello</h1> <br/><br/> <span style="color:red">This is the rich message body </span>`),
		// message receipts.
		"kataras2006@hotmail.com", "kataras2007@hotmail.com", "other@example.com")

	if err != nil {
		println("error while sending the e-mail: " + err.Error())
	}

	// Small:
	//
	// [ init time, once after the `mail.New(...)` ]
	// sender.DefaultFrom = &mail.Address{"FromName", "from@example.com"}
	// [[ runtime, many ]]
	// sender.Send(nil, "Subject", []byte("Body"), "receipt@example.com")
	//
	// [ or/and, many ]
	// sender.Send(&mail.Address{"FromName", "from@example.com"}, "Subject", []byte("Body"), "receipt@example.com")
}

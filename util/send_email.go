package util

import (
	"gopkg.in/gomail.v2"
)

func SendMail() {
	m := gomail.NewMessage()
	m.SetHeader("From", "patharetush@gmail.com")
	m.SetHeader("To", "patharetush@gmail.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	d := gomail.NewPlainDialer("smtp.gmail.com", 587, "patharetush@gmail.com", "Password")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

package util

import (
	"gopkg.in/gomail.v2"
	"os/exec"
	"fmt"
)

func SendMail(from, to, subject, body, attachementPath, password string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if attachementPath != "" {
		m.Attach(attachementPath)
	}
	d := gomail.NewPlainDialer("smtp.gmail.com", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	
	// delete the file
	if attachementPath != "" {
		cmd := exec.Command("rm", "-f", attachementPath)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}


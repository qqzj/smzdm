package main

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

func send() {
	smtpHost := []byte(conf.EmailFromSMTP)[0:strings.LastIndex(conf.EmailFromSMTP, ":")]
	auth := smtp.PlainAuth("", conf.EmailFrom, conf.EmailFromPassword, string(smtpHost))
	mail := fmt.Sprintf("From:Smzdm-Auto-Sign<%s>\nTo:%s\nSubject:%s\nContent-Type: text/plain; charset=utf-8\n\n%s", conf.EmailFrom, conf.EmailTo, "Smzdm Notify", "Hello SMZDM!!"+time.Now().Format(time.Kitchen))
	err := smtp.SendMail(conf.EmailFromSMTP, auth, conf.EmailFrom, conf.EmailTo, []byte(mail))
	if err != nil {
		panic(err)
	}
}

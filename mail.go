package main

import (
	"bytes"
	"log"
	"net/smtp"
	"strings"
	"text/template"
	"time"
)

func send() {
	mailConten := parseMail()
	// 解析邮件服务器主机和端口
	smtpHost := []byte(conf.EmailFromSMTP)[0:strings.LastIndex(conf.EmailFromSMTP, ":")]
	// 认证
	auth := smtp.PlainAuth("", conf.EmailFrom, conf.EmailFromPassword, string(smtpHost))
	// 发送
	err := smtp.SendMail(conf.EmailFromSMTP, auth, conf.EmailFrom, conf.EmailTo, mailConten)
	if err != nil {
		checkError(err)
	} else {
		log.Printf("邮件发送成功%+v", conf.EmailTo)
	}
}
func parseMail() []byte {
	// 接收模板编译
	var mailContent bytes.Buffer
	// 模板静态编译
	tpl, err := template.New("mail.ghtml").Funcs(template.FuncMap{
		"join":     strings.Join,
		"dateTime": time.Now().Format,
	}).ParseFiles("mail.ghtml")
	checkError(err)
	// 传入参数
	err = tpl.Execute(&mailContent, struct {
		Conf          config
		Content       string
		StartAt       time.Time
		EndAt         time.Time
		SignResult    []signJson
		CommentResult []commentJson
	}{
		Conf:          *conf,
		Content:       "Hello SMZDM!!" + time.Now().Format(time.Kitchen),
		StartAt:       startAt,
		EndAt:         time.Now(),
		SignResult:    signResult,
		CommentResult: commentResult,
	})
	checkError(err)
	// 返回
	return mailContent.Bytes()
}

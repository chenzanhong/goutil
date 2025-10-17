package email

import (
	"crypto/tls"
	"errors"
	"log"
	"strings"

	"gopkg.in/gomail.v2"
)

// SendEmail 是一个通用的邮件发送函数
func SendEmail(fromEmail, fromPwd, toEmail, smtpServerHost string, smtpServerPort int, subject, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	d := gomail.NewDialer(smtpServerHost, smtpServerPort, fromEmail, fromPwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // 跳过证书验证，生产环境中应谨慎使用

	if err := d.DialAndSend(m); err != nil {
		log.Printf("发送邮件失败: %v", err)
		if strings.Contains(err.Error(), "535") { // 检查错误消息中是否包含 SMTP 身份验证失败的代码
			log.Printf("可能是 SMTP 身份验证错误")
			return errors.New("发送邮件失败可能是 SMTP 身份验证错误")
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Printf("SMTP 服务器连接被拒绝")
			return errors.New("发送邮件失败：SMTP 服务器连接被拒绝")
		}
		return err
	}
	log.Println("邮件发送成功")
	return nil
}

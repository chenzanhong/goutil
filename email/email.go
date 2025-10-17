package email

import (
	"crypto/tls"
	"fmt"
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
		if strings.Contains(err.Error(), "535") { // 检查错误消息中是否包含 SMTP 身份验证失败的代码
			return fmt.Errorf("发送邮件失败，可能是 SMTP 身份验证错误: %v", err)
		} else if strings.Contains(err.Error(), "connection refused") {
			return fmt.Errorf("发送邮件失败，SMTP 服务器连接被拒绝: %v", err)
		}
		return err
	}
	return nil
}

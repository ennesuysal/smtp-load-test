package main

import (
	"flag"
	"time"

	server "github.com/ennesuysal/smtp-load-test/smtp"
)

func main() {

	workerSize := flag.Int("workerSize", 10, "-wokerSize 10")
	batchSize := flag.Int("batchSize", 10, "-batchSize 10")
	jobCount := flag.Int("jobCount", 10, "-jobCount 10")

	smtpServer := flag.String("smtpServer", "test.yartu.io", "-smtpServer test.yartu.io")
	smtpPort := flag.String("smtpPort", "587", "-smtpPort 587")
	helo := flag.String("helo", "mail.example.com", "-helo mail.example.com")
	smtpPassword := flag.String("smtpPassword", "", "-smtpPassword")
	senderName := flag.String("senderName", "John Doe", "-senderName")
	senderMail := flag.String("senderMail", "sender@example.com", "-senderMail")
	receiverMail := flag.String("receiverMail", "receiver@example.com", "-receiverMail")
	sslEnabled := flag.Bool("sslEnabled", false, "-sslEnabled true")
	startTLS := flag.Bool("startTLS", true, "-startTLS false")
	authEnabled := flag.Bool("authEnabled", true, "-authEnabled false")

	flag.Parse()

	s, _ := server.New(*smtpServer, *smtpPort, *sslEnabled, *startTLS, *authEnabled, *helo, *smtpPassword, *senderName, *senderMail, *receiverMail, *workerSize, *batchSize, *jobCount)

	start := time.Now()
	s.SendTestMails()
	s.St.ProcessDuration = time.Since(start).Seconds()
	s.St.Report()
}

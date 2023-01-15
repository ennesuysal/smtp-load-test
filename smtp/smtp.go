package smtp

import (
	"fmt"
	"net/smtp"
	"strconv"
	"time"

	pool "github.com/ennesuysal/go-thread-pooling"
	statistics "github.com/ennesuysal/smtp-load-test/statistics"
	task "github.com/ennesuysal/smtp-load-test/task"
)

type ServerTest struct {
	server       string
	port         string
	password     string
	senderName   string
	senderMail   string
	receiverMail string
	workerSize   int
	batchSize    int
	jobCount     int
	St           statistics.Statistics
}

func (s *ServerTest) sendMail(order interface{}) error {
	mailTemplate := "To: %s\r\n" +
		"From: %s <%s>" + "\r\n" +
		"Subject: %s\r\n" +
		"\r\n" +
		"%s"

	subject := "SMTP Load Test - Mail " + strconv.Itoa(order.(int))
	body := "Mail " + strconv.Itoa(order.(int)) + "\r\n\r\n...RANDOM DATA..."

	mail := fmt.Sprintf(mailTemplate, s.receiverMail, s.senderName, s.senderMail, subject, body)
	msg := []byte(mail)

	auth := smtp.PlainAuth("", s.senderMail, s.password, s.server)

	to := []string{s.receiverMail}
	start := time.Now()
	err := smtp.SendMail(s.server+":"+s.port, auth, s.senderMail, to, msg)
	duration := time.Since(start).Seconds()

	if err != nil {
		s.St.AddStatistic(statistics.Statistic{
			Duration: duration,
			Success:  false,
		})
		return err
	}

	s.St.AddStatistic(statistics.Statistic{
		Duration: duration,
		Success:  true,
	})

	return nil
}

func (s *ServerTest) SendTestMails() {
	p, _ := pool.NewPool(s.workerSize, s.batchSize)

	p.Start()
	for i := 1; i < s.jobCount+1; i++ {
		p.AddWork(task.NewTask(s.sendMail, func(e error) {
			fmt.Printf("%v\n", e)
		}, i))
	}
	p.Stop()

}

func New(server string, port string, password string, senderName string, senderMail string, receiverMail string, workerSize int, batchSize int, jobCount int) (*ServerTest, error) {
	return &ServerTest{
		server:       server,
		port:         port,
		password:     password,
		senderName:   senderName,
		senderMail:   senderMail,
		receiverMail: receiverMail,
		workerSize:   workerSize,
		batchSize:    batchSize,
		jobCount:     jobCount,
		St:           statistics.Statistics{},
	}, nil
}

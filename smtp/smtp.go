package smtp

import (
	"crypto/tls"
	"fmt"
	"net"
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
	helo         string
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

	auth := smtp.PlainAuth("", s.senderMail, s.password, s.server)

	stat := statistics.NewStatistic(0)

	start := time.Now()
	conn, err := net.Dial("tcp", s.server+":"+s.port)
	if err != nil {
		stat["DIAL"] = time.Since(start).Seconds()
		s.St.AddSuccess(false)
		s.St.AddStatistic(stat)
		return err
	}
	s.St.RemoteIp = conn.RemoteAddr().String()
	stat["DIAL"] = time.Since(start).Seconds()

	start = time.Now()
	c, err := smtp.NewClient(conn, s.server)
	if err != nil {
		stat["TOUCH"] = time.Since(start).Seconds()
		s.St.AddSuccess(false)
		s.St.AddStatistic(stat)
		return err
	}
	defer c.Close()
	stat["TOUCH"] = time.Since(start).Seconds()

	start = time.Now()
	err = c.Hello(s.helo)
	if err != nil {
		stat["HELO"] = time.Since(start).Seconds()
		s.St.AddSuccess(false)
		s.St.AddStatistic(stat)
		return err
	}
	stat["HELO"] = time.Since(start).Seconds()

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{InsecureSkipVerify: true}
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}

	err = c.Auth(auth)
	if err != nil {
		return err
	}

	start = time.Now()
	err = c.Mail(s.senderMail)
	if err != nil {
		stat["MAIL"] = time.Since(start).Seconds()
		s.St.AddSuccess(false)
		s.St.AddStatistic(stat)
		return err
	}
	stat["MAIL"] = time.Since(start).Seconds()

	start = time.Now()
	err = c.Rcpt(s.receiverMail)
	if err != nil {
		stat["RCPT"] = time.Since(start).Seconds()
		s.St.AddSuccess(false)
		s.St.AddStatistic(stat)
		return err
	}
	stat["RCPT"] = time.Since(start).Seconds()

	start = time.Now()
	wc, err := c.Data()
	if err != nil {
		stat["DATA"] = time.Since(start).Seconds()
		s.St.AddSuccess(false)
		s.St.AddStatistic(stat)
		return err
	}
	_, err = fmt.Fprint(wc, mail)
	if err != nil {
		stat["DATA"] = time.Since(start).Seconds()
		s.St.AddSuccess(false)
		s.St.AddStatistic(stat)
		return err
	}

	err = wc.Close()
	if err != nil {
		stat["DATA"] = time.Since(start).Seconds()
		s.St.AddSuccess(false)
		s.St.AddStatistic(stat)
		return err
	}
	stat["DATA"] = time.Since(start).Seconds()

	start = time.Now()
	err = c.Quit()
	if err != nil {
		stat["QUIT"] = time.Since(start).Seconds()
		s.St.AddSuccess(false)
		s.St.AddStatistic(stat)
		return err
	}
	stat["QUIT"] = time.Since(start).Seconds()

	s.St.AddSuccess(true)
	s.St.AddStatistic(stat)

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

func New(server string, port string, helo string, password string, senderName string, senderMail string, receiverMail string, workerSize int, batchSize int, jobCount int) (*ServerTest, error) {
	return &ServerTest{
		server:       server,
		port:         port,
		helo:         helo,
		password:     password,
		senderName:   senderName,
		senderMail:   senderMail,
		receiverMail: receiverMail,
		workerSize:   workerSize,
		batchSize:    batchSize,
		jobCount:     jobCount,
		St:           *statistics.NewStatistics(),
	}, nil
}

package fakes

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	MaxTestMails           = 50
	SmtpDefaultMaxLineSize = 77
)

type FakeSmtpServer struct {
	Mailbox  chan MailParams
	port     string
	listener net.Listener
	logger   *zap.SugaredLogger
}

type MailParams struct {
	From    string
	To      string
	Subject string
	Body    string
}

func NewFakeSmtpServer(port string) *FakeSmtpServer {
	return &FakeSmtpServer{
		port: port,
	}
}

func (n *FakeSmtpServer) StartReceive(logger *zap.SugaredLogger) error {
	var err error
	n.logger = logger
	serverAddr := fmt.Sprintf("localhost:%s", n.port)
	n.listener, err = net.Listen("tcp", serverAddr)
	if err != nil {
		return err
	}

	n.Mailbox = make(chan MailParams, MaxTestMails)
	go func() {
		n.logger.Info("Fake email server is ready to receive")
		for {
			conn, err := n.listener.Accept()
			if err != nil {
				n.logger.Errorf("Failed to accept connection: %s\n", err.Error())
				return
			}

			go n.clientHandler(conn)
		}
	}()
	return nil
}

func (n *FakeSmtpServer) Close() {
	if err := n.listener.Close(); err != nil {
		n.logger.Errorf("Error closing SMTP listener: %s\n", err.Error())
		return
	}
	close(n.Mailbox)
}

func (n *FakeSmtpServer) clientHandler(conn net.Conn) {
	receiveData := false
	readContent := false
	fullEmailContent := ""
	email := MailParams{}
	bufout := bufio.NewWriter(conn)
	bufin := bufio.NewReader(conn)
	_, err := bufout.WriteString("220 welcome\r\n")
	if err != nil {
		n.logger.Errorf("Error writing initial SMTP message: %s\n", err.Error())
		return
	}

	if err := bufout.Flush(); err != nil {
		n.logger.Errorf("Error sending SMTP message: %s\n", err.Error())
		return
	}

	timeout := time.NewTimer(15 * time.Second)
	for {
		select {
		case <-timeout.C:
			return
		default:
			content, err := bufin.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					n.logger.Errorf("Error reading from client:  %s\n", err.Error())
					return
				}
			}
			if strings.Contains(content, "QUIT") {
				_, err := bufout.WriteString("221 Closing connection\r\n")
				if err != nil {
					n.logger.Errorf("Error sending close connection SMTP message: %s\n", err.Error())
					return
				}
				if err := bufout.Flush(); err != nil {
					n.logger.Errorf("Error sending SMTP message: %s\n", err.Error())
					return
				}
				email.Body = fullEmailContent
				n.Mailbox <- email
				return
			} else if strings.Contains(content, "DATA") {
				_, err := bufout.WriteString("354 Ready for receive message. End data with <CR><LF>.<CR><LF>\r\n")
				if err != nil {
					n.logger.Errorf("Error sending ready SMTP message: %s\n", err.Error())
					return
				}
				if err := bufout.Flush(); err != nil {
					n.logger.Errorf("Error sending SMTP message: %s\n", err.Error())
					return
				}
				receiveData = true
			} else {
				if receiveData {
					if readContent {
						fullEmailContent += processContent(content)
					} else if strings.Contains(content, "From: ") {
						email.From = strings.TrimSpace(strings.ReplaceAll(content, "From: ", ""))
					} else if strings.Contains(content, "To: ") {
						email.To = strings.TrimSpace(strings.ReplaceAll(content, "To: ", ""))
					} else if strings.Contains(content, "Subject: ") {
						email.Subject = strings.TrimSpace(strings.ReplaceAll(content, "Subject: ", ""))
					} else if strings.Contains(content, "<html>") {
						fullEmailContent += processContent(content)
						readContent = true
					}
				}

				_, err := bufout.WriteString("250 Received\r\n")
				if err != nil {
					n.logger.Errorf("Error sending acknowledge SMTP message: %s\n", err.Error())
					return
				}

				if err := bufout.Flush(); err != nil {
					n.logger.Errorf("Error sending SMTP message: %s\n", err.Error())
					return
				}
			}

			n.logger.Infof(content)
		}
	}
}

func processContent(content string) string {
	if len(content) == 3 && content == ".\r\n" { // End of msg
		return ""
	}

	if len(content) > SmtpDefaultMaxLineSize {
		content = strings.ReplaceAll(content, "=\r\n", "")
	}
	content = strings.ReplaceAll(content, "=3D", "=")
	return strings.ReplaceAll(content, "\r", "")
}

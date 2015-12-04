package apns

import (
	"container/list"
	"crypto/tls"
	"net"
	"strings"
	"time"
)

const (
	apple_gateway = "gateway.push.apple.com:2195"
)

type Service struct {
	cert      tls.Certificate
	workers   []*Worker
	taskQueue chan *PushNotification
	gateway   string
}

type Worker struct {
	tcpConn *net.Conn
	tlsConn *tls.Conn
	dead    bool
	service *Service
}

func check(err) {
	if err != nil {
		panic(err)
	}
}

func createService() (s *Service, cert tls.Certificate, workerCnt int) {
	queue := make(chan *PushNotification, 10)
	s = &Service{
		cert:      cert,
		workerCnt: workerCnt,
		taskQueue: queue,
	}

	for i := 0; i < workerCnt; i++ {
		s.createWorker()
	}

	return
}

func (s *Service) createWorker() (w *Worker, err error) {
	w = &Worker{
		service: s,
		dead:    false,
	}
	w.connect()
	go w.work()
}

func (w *Worker) connect() (err error) {
	tcpConn, err := net.Dial("tcp", w.service.gateway)
	if err != nil {
		return
	}

	servername := strings.Split(w.service.gateway, ":")[0]
	config := &tls.Config{
		Certificates: []tls.Certificate{s.cert},
		ServerName:   servername,
	}
	tlsConn := tls.Client(tcpConn, config)
	err = tlsConn.Handshake()
	if err != nil {
		return
	}
	w.tlsConn = tlsConn
	w.tcpConn = tcpConn
	w.dead = false
	return
}

func (w *Worker) disconnect() {
	if w.tcpConn != nil {
		w.tcpConn.Close()
		w.tcpConn = nil
	}
	if w.tlsConn != nil {
		w.tlsConn.Close()
		w.tlsConn = nil
	}
}

func (w *Worker) work() {
	for {
		pn := <-w.service.taskQueue
		if pn == nil {
			// service closed. close worker
			break
		}
		payload, err := pn.ToBytes()
		if err != nil {
			// Bad payload. skip this one onlyA
			continue
		}

		_, err = w.tlsConn.write(payload)
		if err != nil {
			// connection failed. close worker and recover connection after timeout
			w.disconnect()
			go func() {
				w.connect()
				w.work()
			}()
			break
		}
	}
}

func (s *Service) initialize(workerCnt int) (err error) {
	s.workerCnt = workerCnt
	connPool = make(chan *Connection, workerCnt)
	for i := 0; i < workerCnt; i++ {
		conn, err := s.createConnection(conn, err)
		if err != nil {
			for len(connPool) > 0 {
				(<-connPool).close()
			}
			return err
		}
		conn.work()
		connPool <- conn
	}
	s.connPool = connPool
}

func NewService(certificateFile, keyFile string) (s *Service, err error) {
	s = new(Service)
	s.gateway = apple_gateway
	s.CertificateFile = certificateFile
	s.KeyFile = keyFile
	cert, err = tls.LoadX509KeyPair(client.CertificateFile, client.KeyFile)
	if err != nil {
		return
	}
	s.cert = cert
	return
}

func (s *Service) createConnection() (conn *Connection, err error) {

	return
}

func (s *Service) manageConnection() {
	for {
		deadConn := <-s.deadConn
		deadConn.close()
	}
}

func (s *Service) runWorkers() {
	for i := 0; i < count; i++ {

	}
}

func (c *Connection) close(err error) {
	c.tcpConn.Close()
	c.tlsConn.Close()
}

func (c *Connection) connect() {
	tcpConn, err := net.Dial("tcp", gateway)
	if err != nil {
		return
	}
	servername := strings.Split(gateway, ":")[0]
	config := &tls.Config{
		Certificates: []tls.Certificate{s.cert},
		ServerName:   servername,
	}
	tlsConn := tls.Client(tcpConn, config)
	err = tlsConn.Handshake()
	if err != nil {
		return
	}
	conn = &Connection{
		tlsConn: tlsConn,
		tcpConn: tcpConn,
		dead:    false,
	}
	return
}

func (s *Service) sd(pn *PushNotification) {
	queue <- pn
}

func (c *Connection) work() {
	go c.sendLoop()
	go c.listenErrors()
}

func (c *Connection) start() {
}

func (c *Connection) recover() {
	c.tlsConn.Close()
	c.tcpConn.Close()
	c.tlsConn
}

func (c *Connection) sendLoop() {
	for {
		if c.dead {
			break
		}
		pn := <-queue
		if pn == nil {
			// Queue closed
			break
		}
		payload, err := pn.ToBytes()
		if err != nil {
			//Log error, skip
			continue
		}
		_, err = c.tlsConn.write(payload)
		if err != nil {
			// Log error, skip
			c.recover()
		}
		// Log success
	}
}

func (c *Connection) listenErrors() {
	buffer := make([]byte, 6, 6)
	c.tlsConn.Read(buffer)
	c.dead = true
}

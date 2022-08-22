package tcp_proxy

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"
)

// Proxy - Manages a Proxy connection, piping data between local and remote.
type Proxy struct {
	lSockFile string
	raddr     string
	stopChan  chan struct{}
}

// New - Create a new Proxy instance. Takes a remote address and creates
// a unix file socket. This allows you to manipulate the connection forwarding
// however you please
// and closes it when finished.
// func New(raddr string) (*Proxy, error) {
func New(raddr string) *Proxy {
	p := &Proxy{
		lSockFile: generateSockFile(),
		raddr:     raddr,
	}
	return p
}
func (p *Proxy) Addr() string {
	return p.lSockFile
}
func (p *Proxy) logError(msg string, err error) {
	fmt.Println("[Connection Proxy][Error]", msg, err.Error())
}
func (p *Proxy) logMsg(msg ...interface{}) {
	fmt.Print("[Connection Proxy] ")
	fmt.Println(msg...)
}
func (p *Proxy) Stop() {
	p.logMsg("Intercepting and stopping connection")
	close(p.stopChan)
}
func (p *Proxy) Start() {
	go p.start()
}
func (p *Proxy) start() {
	laddr, err := net.ResolveUnixAddr("unix", p.lSockFile)
	if err != nil {
		p.logError("Resolve unix socket file error", err)
		return
	}
	raddr, err := net.ResolveTCPAddr("tcp", p.raddr)
	if err != nil {
		p.logError("Resolve remote address error", err)
		return
	}
	p.logMsg("Waiting for new connection")
	listener, err := net.ListenUnix("unix", laddr)
	if err != nil {
		p.logError("Listen Unix Socket Error:", err)
		return
	}
	listener.SetUnlinkOnClose(true)
	p.stopChan = make(chan struct{})
	for {
		if listener == nil {
			break
		}
		lconn, err := listener.AcceptUnix()
		if err != nil {
			p.logError("Accept unix connection error", err)
			return
		}
		p.logMsg("New connection: ", lconn.RemoteAddr())
		go func() {
			defer func() {
				lconn.Close()
				listener.Close()
				listener = nil
			}()
			// connect to remote
			rconn, err := net.DialTCP("tcp", nil, raddr)
			if err != nil {
				p.logError("Remote connection failed: ", err)
				return
			}
			defer rconn.Close()
			// display both ends
			p.logMsg("Opened:: Remote:", raddr.String(), "Local:", laddr.String())
			closer := make(chan struct{})
			// bidirectional copy
			go p.pipe(closer, lconn, rconn)
			go p.pipe(closer, rconn, lconn)
			select {
			case <-closer:
				p.logMsg("Closing connection")
			case <-p.stopChan:
				p.logMsg("Stopping connection")
			}
			p.logMsg("Connection complete")
		}()
	}
	p.logMsg("Proxy shutting down")
}

func (p *Proxy) pipe(closer chan struct{}, src io.Reader, dst io.Writer) {
	_, _ = io.Copy(dst, src)
	closer <- struct{}{}
}

func generateSockFile() string {
	prefix := rand.NewSource(time.Now().UnixNano()).Int63()
	return filepath.Join(os.TempDir(), fmt.Sprintf("%d-conn-proxy.sock", prefix))
}

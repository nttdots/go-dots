package connection

import (
	"net"
	dtls "github.com/nttdots/go-dtls"
)

// ---server socket
/*
ネットワーク待ち受けリスナのインターフェース
 */
type DotsNetworkListener interface {
	Close()
}

/*
ネットワーク待ち受けリスナのファクトリインターフェース
 */
type ListenerFactory interface {
	CreateListener(address string, workerCh chan net.Conn, errorCh chan error) (DotsNetworkListener, error)
}

/*
DTLSを使うリスナのためのファクトリ
 */
type DTLSListenerFactory struct {
	caCertFile     string
	crlFile        string
	serverCertFile string
	serverKeyFile  string
}

func NewDTLSListenerFactory(caCertFile, crlFile, serverCertFile, serverKeyFile string) *DTLSListenerFactory {
	return &DTLSListenerFactory{
		caCertFile,
		crlFile,
		serverCertFile,
		serverKeyFile,
	}
}

func (d *DTLSListenerFactory) CreateListener(address string, workerCh chan net.Conn, errorCh chan error) (listener DotsNetworkListener, err error) {

	context, err := dtls.NewDTLSServerContext(
		d.caCertFile,
		d.crlFile,
		d.serverCertFile,
		d.serverKeyFile)

	if err != nil {
		return
	}

	return context.Listen(address, workerCh, errorCh)
}

// --- client socket

type ClientConnectionFactory interface {
	Connect(address string) (net.Conn, error)
	Close()
}

type DTLSConnectionFactory struct {
	caCertFile     string
	clientCertFile string
	clientKeyFile  string

	ctx *dtls.DTLSCTX
}

func NewDTLSConnectionFactory(caCertFile, clientCertFile, clientKeyFile string) (*DTLSConnectionFactory, error) {
	ctx, err := dtls.NewDTLSClientContext(caCertFile, clientCertFile, clientKeyFile)
	if err != nil {
		return nil, err
	}

	return &DTLSConnectionFactory{
		caCertFile,
		clientCertFile,
		clientKeyFile,
		ctx,
	}, nil
}

func (d *DTLSConnectionFactory) Connect(address string) (net.Conn, error) {
	return d.ctx.Connect(address)
}

func (d *DTLSConnectionFactory) Close() {
	d.ctx.Close()
}

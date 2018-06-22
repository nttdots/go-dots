package main

import (
  "crypto/tls"
  "crypto/x509"
  "errors"
  "fmt"
  "io/ioutil"
  "net"
  "net/http"
  "strconv"

  log "github.com/sirupsen/logrus"
  "github.com/nttdots/go-dots/dots_server/router"
)

func listenData(address string, port uint16, caFile string, certFile string, keyFile string) error {

  rootCerts := x509.NewCertPool()
  pem, err := ioutil.ReadFile(caFile)
  if err != nil {
    return fmt.Errorf("listenData: io.ReadFile() failed: %v", err)
  }
  ok := rootCerts.AppendCertsFromPEM(pem)
  if !ok {
    return errors.New("listenData: CertPool#AppendCertsFromPEM() failed.")
  }

  config := tls.Config {
    ClientAuth: tls.RequireAndVerifyClientCert,
    ClientCAs: rootCerts,
  }

  server := &http.Server {
    Addr: net.JoinHostPort(address, strconv.Itoa(int(port))),
    TLSConfig: &config,
    Handler: router.CreateRouter(),
  }
  go func() {
    err := server.ListenAndServeTLS(certFile, keyFile)
    if err != nil {
      log.Fatal(err)
    }
  }()
  return nil
}

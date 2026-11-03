package mail

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cogikyo/hgmx/credential"
)

func TestCheckHidesPassword(t *testing.T) {
	const good = "right-password-123"
	fakeSMTP(t, good)
	s := SMTP{Host: "127.0.0.1", Port: 18520, User: "relay", From: "site@example.com"}

	s.Password = good
	if err := s.Check(context.Background()); err != nil {
		t.Fatalf("good password: %v", err)
	}
	s.Password = "wrong-password-456"
	err := s.Check(context.Background())
	f, ok := err.(*credential.Failure)
	if !ok || f.Code != 535 {
		t.Fatalf("wrong password: got %#v, want a 535 Failure", err)
	}
	if strings.Contains(err.Error(), s.Password) {
		t.Fatalf("failure echoes the password: %q", err.Error())
	}
}

func fakeSMTP(t *testing.T, password string) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	certFile := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SSL_CERT_FILE", certFile)

	ln, err := tls.Listen("tcp", "127.0.0.1:18520", &tls.Config{Certificates: []tls.Certificate{{Certificate: [][]byte{der}, PrivateKey: key}}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go serveSMTP(conn, password)
		}
	}()
}

func serveSMTP(conn net.Conn, password string) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	reply := func(s string) { conn.Write([]byte(s + "\r\n")) }
	reply("220 fake ESMTP")
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		cmd, arg, _ := strings.Cut(strings.TrimSpace(line), " ")
		switch strings.ToUpper(cmd) {
		case "EHLO":
			reply("250-fake\r\n250 AUTH PLAIN")
		case "AUTH":
			raw, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(arg, "PLAIN "))
			if parts := strings.Split(string(raw), "\x00"); len(parts) == 3 && parts[2] == password {
				reply("235 2.7.0 accepted")
			} else {
				reply("535 5.7.8 rejected " + parts[len(parts)-1])
			}
		case "QUIT":
			reply("221 bye")
			return
		default:
			reply("502 unsupported")
		}
	}
}

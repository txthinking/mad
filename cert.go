package mad

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

type Cert struct {
	CaPEM    []byte
	CaKeyPEM []byte
	C        *x509.Certificate
	CertPEM  []byte
	KeyPEM   []byte
}

func NewCert(caPEM, caKeyPEM []byte, Organization, OrganizationalUnit string) *Cert {
	c := &x509.Certificate{
		Subject: pkix.Name{
			Organization:       []string{Organization},
			OrganizationalUnit: []string{OrganizationalUnit},
		},
		NotBefore:             time.Date(2019, time.June, 1, 0, 0, 0, 0, time.UTC),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	return &Cert{
		CaPEM:    caPEM,
		CaKeyPEM: caKeyPEM,
		C:        c,
	}
}

func (c *Cert) SetIPAddresses(ips []net.IP) {
	c.C.IPAddresses = ips
	if len(ips) > 0 {
		c.C.Subject.CommonName = ips[0].String()
	}
}

func (c *Cert) SetDNSNames(domains []string) {
	c.C.DNSNames = domains
	if len(domains) > 0 {
		c.C.Subject.CommonName = domains[0]
	}
}

func (c *Cert) Create() error {
	tc, err := tls.X509KeyPair(c.CaPEM, c.CaKeyPEM)
	if err != nil {
		return err
	}
	ca, err := x509.ParseCertificate(tc.Certificate[0])
	if err != nil {
		return err
	}

	p, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	pub := p.Public()

	b, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return err
	}
	var spki struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}
	if _, err := asn1.Unmarshal(b, &spki); err != nil {
		return err
	}
	skid := sha1.Sum(spki.SubjectPublicKey.Bytes)
	c.C.SubjectKeyId = skid[:]

	sn, err := rand.Int(rand.Reader, big.NewInt(0).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}
	c.C.SerialNumber = sn

	b, err = x509.CreateCertificate(rand.Reader, c.C, ca, pub, tc.PrivateKey)
	if err != nil {
		return err
	}
	c.CertPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: b})
	// c.KeyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(p)})
	b, err = x509.MarshalPKCS8PrivateKey(p)
	if err != nil {
		return err
	}
	c.KeyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b})
	return nil
}

func (c *Cert) Cert() []byte {
	return c.CertPEM
}

func (c *Cert) Key() []byte {
	return c.KeyPEM
}

func (c *Cert) SaveToFile(cert, key string) error {
	f, err := os.Create(cert)
	if err != nil {
		return err
	}
	if _, err := f.Write(c.CertPEM); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	f, err = os.Create(key)
	if err != nil {
		return err
	}
	if _, err := f.Write(c.KeyPEM); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

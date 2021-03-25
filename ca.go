package mad

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

type Ca struct {
	C      *x509.Certificate
	CaPEM  []byte
	KeyPEM []byte
}

func NewCa(Organization, OrganizationalUnit, CommonName string) *Ca {
	c := &x509.Certificate{
		Subject: pkix.Name{
			Organization:       []string{Organization},
			OrganizationalUnit: []string{OrganizationalUnit},
			CommonName:         CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		MaxPathLenZero:        true,
	}
	return &Ca{
		C: c,
	}
}

func (c *Ca) Create() error {
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

	b, err = x509.CreateCertificate(rand.Reader, c.C, c.C, pub, p)
	if err != nil {
		return err
	}
	c.CaPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: b})
	// c.KeyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(p)})
	b, err = x509.MarshalPKCS8PrivateKey(p)
	if err != nil {
		return err
	}
	c.KeyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b})
	return nil
}

func (c *Ca) Ca() []byte {
	return c.CaPEM
}

func (c *Ca) Key() []byte {
	return c.KeyPEM
}

func (c *Ca) SaveToFile(ca, key string) error {
	f, err := os.Create(ca)
	if err != nil {
		return err
	}
	if _, err := f.Write(c.CaPEM); err != nil {
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

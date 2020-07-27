package rs

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
)

// Secure command manages credentials in /etc/credentials
func Secure(cmd *Cmd) error {
	flags := flag.NewFlagSet("secure", flag.ContinueOnError)
	name := flags.String("a", "", "account name")
	secret := flags.String("s", "", "secret")
	check := flags.Bool("c", false, "check if secret is valid")
	flags.SetOutput(cmd.Out)
	if err := flags.Parse(cmd.Args); err != nil {
		return err
	}
	var acc Account
	if err := cmd.Sys.LoadAccount(&acc, *name); err != nil {
		return err
	}

	cred := NewCredentials()
	cmd.Sys.Load(&cred, "/etc/credentials")

	if *check {
		if err := cred.Check(acc.UID, *secret); err != nil {
			return err
		}
	}
	hash := sha256.New()
	cred.AddSecret(&Secret{
		UID:       acc.UID,
		Encrypted: hash.Sum([]byte(*secret)),
	})
	return cmd.Sys.Save("/etc/credentials", &cred)
}

func NewCredentials() *Credentials {
	return &Credentials{
		Secrets: make([]*Secret, 0),
	}
}

type Credentials struct {
	Secrets []*Secret
}

// AddSecret
func (me *Credentials) AddSecret(s *Secret) {
	me.Secrets = append(me.Secrets, s)
}

// Check
func (me *Credentials) Check(uid int, secret string) error {
	for _, s := range me.Secrets {
		if s.UID != uid {
			continue
		}
		encrypted := sha256.New().Sum([]byte(secret))
		if bytes.Equal(encrypted, s.Encrypted) {
			return nil
		}
	}
	return fmt.Errorf("invalid")
}

type Secret struct {
	UID       int
	Encrypted []byte
}

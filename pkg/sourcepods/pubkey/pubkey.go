package pubkey

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// pk="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDHLLxzamVxBlhRjU3wQ4zzfeCSmpHShQF65csazVsZptNbQCktII/nxaJIh9/WcVzmhdFcSi+mPqFx6Hi5PCS6/8XEjTa+E0FfXyd4Yy+bA8M7a2IgHU//AbF99Itj43VY6Zmruzvli1IOdMzFzHis3AWDjxO3zSARMR+UfctE8DgrUKMo7xvCHpm/JxS4+m3s0dET9mFlxUQ4ZUmyVEpFqh/LMYeoCt2TzbgyGF3cJ3DliivrRUTdzgN+1GZYv69hoqkV7qy4tM7bzZZsEgU+d9TLLvyf7g3cxU/iy22Nj+YssYPYd2YMpN8fFkkHtsir9cWnmgkaGFtbLKsuqlS1\n"
// fp="SHA256:bdUUBV8GPzfOFvSjRmS/gvFVllvTVDJcI37e+3wkfOQ"

type (
	// PubKey holds a pubkey...
	PubKey struct {
		ID          string
		Name        string
		Fingerprint string
		Content     string
		Created     time.Time
		Updated     time.Time
	}
)

func (pk *PubKey) fillFingerprint() error {
	gpk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pk.Content))
	if err != nil {
		return errors.Wrap(err, "could not parse public key")
	}
	pk.Fingerprint = ssh.FingerprintSHA256(gpk)
	return nil
}

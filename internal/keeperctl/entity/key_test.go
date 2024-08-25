package entity_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"

	"github.com/a-x-a/gophkeeper/internal/keeperctl/entity"
	"github.com/a-x-a/gophkeeper/internal/util/creds"
	"github.com/a-x-a/gophkeeper/internal/util/gophtest"
)

func TestKeyToHash(t *testing.T) {
	sat := entity.NewKey(gophtest.Username, gophtest.Password)

	snaps.MatchSnapshot(t, sat.Hash())
}

func TestEncryptDecrypt(t *testing.T) {
	tt := []struct {
		name     string
		username string
		password creds.Password
		msg      []byte
	}{
		{
			name:     "Basic key",
			username: gophtest.Username,
			password: gophtest.Password,
			msg:      []byte("TestEncryptDecrypt"),
		},
		{
			name:     "Short key",
			username: "a",
			password: "b",
			msg:      []byte("TestEncryptDecrypt"),
		},
		{
			name:     "Empty message is noop",
			username: gophtest.Username,
			password: gophtest.Password,
			msg:      []byte{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sat := entity.NewKey(gophtest.Username, gophtest.Password)

			encrypted, err := sat.Encrypt(tc.msg)
			require.NoError(t, err)

			decrypted, err := sat.Decrypt(encrypted)
			require.NoError(t, err)
			require.Equal(t, tc.msg, decrypted)
		})
	}
}

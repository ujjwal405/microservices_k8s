package user

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	usr := User{
		Username: "abcdef",
		User_id:  "12345",
	}
	tok, err := GenerateAlltoken(usr)
	require.NoError(t, err)
	log.Println(tok)
}

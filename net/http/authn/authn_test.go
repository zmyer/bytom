package authn

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"

	dbm "github.com/tendermint/tmlibs/db"

	"github.com/bytom/accesstoken"
	"github.com/bytom/errors"
)

func TestAuthenticate(t *testing.T) {
	ctx := context.Background()

	tokenDB := dbm.NewDB("testdb", "leveldb", "temp")
	defer os.RemoveAll("temp")
	tokenStore := accesstoken.NewStore(tokenDB)
	token, err := tokenStore.Create(ctx, "alice", "test")
	if err != nil {
		t.Errorf("create token error")
	}

	cases := []struct {
		id, tok string
		want    error
	}{
		{"alice", token.Token, nil},
		{"alice", "alice:abcsdsdfassdfsefsfsfesfesfefsefa", ErrInvalidToken},
	}

	api := NewAPI(tokenStore)

	for _, c := range cases {
		var username, password string
		toks := strings.SplitN(c.tok, ":", 2)
		if len(toks) > 0 {
			username = toks[0]
		}
		if len(toks) > 1 {
			password = toks[1]
		}

		req, _ := http.NewRequest("GET", "/", nil)
		req.SetBasicAuth(username, password)

		_, err := api.Authenticate(req)
		if errors.Root(err) != c.want {
			t.Errorf("Authenticate(%s) error = %s want %s", c.id, err, c.want)
		}
	}
}

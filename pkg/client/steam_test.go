package client

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountList(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("../../testdata/fixture_account_list.json")
		if err != nil {
			assert.NoError(t, err)
		}
		w.Write(f)
	})
	server := httptest.NewServer(fn)
	defer server.Close()

	service := newService("abc")
	service.url = server.URL
	steam := SteamService{}
	steam.service = service

	accounts, err := steam.GetAccountList()
	assert.NoError(t, err)
	assert.Equal(t, "76561197900265728", accounts.Actor)
	assert.True(t, accounts.IsBanned)
	assert.Equal(t, 1, len(accounts.Accounts))

	account := accounts.Accounts[0]
	assert.Equal(t, "1111", account.SteamID)
	assert.Equal(t, "abc", account.LoginToken)
	assert.Equal(t, "test", account.Memo)
	assert.Equal(t, 730, account.AppID)
	assert.Equal(t, 1589448113, account.LastLogon)
	assert.True(t, account.IsExpired)
	assert.False(t, account.IsDeleted)
}
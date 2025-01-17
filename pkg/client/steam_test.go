package client

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fakeServer(fn http.HandlerFunc) (SteamService, func()) {
	server := httptest.NewServer(fn)
	service := newService("abc")
	service.url = server.URL
	steam := SteamService{}
	steam.service = service
	return steam, server.Close
}

func TestGetAccountList(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("../../testdata/fixture_account_list.json")
		if err != nil {
			assert.NoError(t, err)
		}
		w.Write(f)
	})
	steam, close := fakeServer(fn)
	defer close()

	accounts, err := steam.GetAccountList()
	assert.NoError(t, err)
	assert.Equal(t, "76561197900265728", accounts.Actor)
	assert.True(t, accounts.IsBanned)
	assert.Equal(t, 1, len(accounts.Servers))

	account := accounts.Servers[0]
	assert.Equal(t, "1111", account.SteamID)
	assert.Equal(t, "abc", account.LoginToken)
	assert.Equal(t, "test", account.Memo)
	assert.Equal(t, 730, account.AppID)
	assert.Equal(t, 1589448113, account.LastLogon)
	assert.True(t, account.IsExpired)
	assert.False(t, account.IsDeleted)
}

func TestCreateAccount(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("../../testdata/fixture_create_account.json")
		if err != nil {
			assert.NoError(t, err)
		}
		w.Write(f)
	})
	steam, close := fakeServer(fn)
	defer close()

	account, err := steam.CreateAccount(730, "hello world")
	assert.NoError(t, err)
	assert.Equal(t, 730, account.AppID)
	assert.Equal(t, "hello world", account.Memo)
	assert.Equal(t, "80068392925402169", account.SteamID)
	assert.Equal(t, "D212EAB4B33A0005CA4CD483AAAA4C9E", account.LoginToken)
}

func TestSetMemo(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("../../testdata/fixture_empty_response.json")
		if err != nil {
			assert.NoError(t, err)
		}
		w.Write(f)
	})
	steam, close := fakeServer(fn)
	defer close()

	err := steam.SetMemo(80068392925402169, "hello world")
	assert.NoError(t, err)
}

func TestResetLoginToken(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("../../testdata/fixture_reset_token.json")
		if err != nil {
			assert.NoError(t, err)
		}
		w.Write(f)
	})
	steam, close := fakeServer(fn)
	defer close()

	token, err := steam.ResetLoginToken(80068392925402169)
	assert.NoError(t, err)
	assert.Equal(t, "A612F82D000F93F800D737B624040080", token.LoginToken)
	assert.Equal(t, "80068392925402169", token.SteamID)
	assert.Equal(t, 0, token.AppID)
	assert.False(t, token.IsDeleted)
	assert.False(t, token.IsExpired)
}

func TestDeleteAccount(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("../../testdata/fixture_empty_response.json")
		if err != nil {
			assert.NoError(t, err)
		}
		w.Write(f)
	})
	steam, close := fakeServer(fn)
	defer close()

	err := steam.DeleteAccount(80068392925402169)
	assert.NoError(t, err)
}

func TestQueryLoginToken(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("../../testdata/fixture_query_login_token.json")
		if err != nil {
			assert.NoError(t, err)
		}
		w.Write(f)
	})
	steam, close := fakeServer(fn)
	defer close()

	account, err := steam.QueryLoginToken("A612F82D000F93F800D737B624040080")
	assert.NoError(t, err)
	assert.False(t, account.IsBanned)
	assert.Equal(t, 0, account.Expires)
	assert.Equal(t, 1, len(account.Servers))

	token := account.Servers[0]
	assert.Equal(t, "A612F82D000F93F800D737B624040080", token.LoginToken)
	assert.Equal(t, "85008392003198734", token.SteamID)
	assert.Equal(t, 0, token.AppID)
	assert.False(t, token.IsDeleted)
	assert.False(t, token.IsExpired)
}

func TestQueryLoginTokenNotValid(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("../../testdata/fixture_empty_response.json")
		if err != nil {
			assert.NoError(t, err)
		}
		w.Write(f)
	})
	steam, close := fakeServer(fn)
	defer close()

	account, err := steam.QueryLoginToken("abcd")
	assert.Error(t, err)
	assert.Nil(t, account)
}

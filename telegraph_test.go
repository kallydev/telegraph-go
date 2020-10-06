package telegraph_test

import (
	"github.com/kallydev/telegraph-go"
	"testing"
	"time"
)

const (
	shortName  = "telegraph-go"
	authorName = "TelegraphGo"
	authorURL  = "https://github.com/kallydev/telegraph-go"
)

var (
	err     error
	client  *telegraph.Client
	account *telegraph.Account
	page    *telegraph.Page
)

func TestCreateClient(t *testing.T) {
	client, err = telegraph.NewClient("", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateAccount(t *testing.T) {
	account, err = client.CreateAccount(shortName, &telegraph.CreateAccountOption{
		AuthorName: authorName,
		AuthorURL:  authorURL,
	})
	if err != nil {
		t.Fatal(err)
	}
	client.AccessToken = account.AccessToken
}

func TestClient_EditAccountInfo(t *testing.T) {
	account, err = client.EditAccountInfo(&telegraph.EditAccountInfoOption{
		ShortName:  shortName,
		AuthorName: authorName,
		AuthorURL:  authorURL,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_RevokeAccessToken(t *testing.T) {
	account, err = client.RevokeAccessToken()
	if err != nil {
		t.Fatal(err)
	}
	client.AccessToken = account.AccessToken
}

func TestClient_CreatePage(t *testing.T) {
	page, err = client.CreatePage("title", []telegraph.Node{
		telegraph.NodeElement{
			Tag: "p",
			Children: []telegraph.Node{
				"hello world",
			},
		},
	}, &telegraph.CreatePageOption{
		AuthorName:    authorName,
		AuthorURL:     authorURL,
		ReturnContent: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_EditPage(t *testing.T) {
	page, err = client.EditPage(page.Path, "title", []telegraph.Node{
		telegraph.NodeElement{
			Tag: "p",
			Children: []telegraph.Node{
				"hello world",
			},
		},
	}, &telegraph.EditPageOption{
		AuthorName:    authorName,
		AuthorURL:     authorURL,
		ReturnContent: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetAccountInfo(t *testing.T) {
	account, err = client.GetAccountInfo(&telegraph.GetAccountInfoOption{
		Fields: []string{
			telegraph.FieldShortName, telegraph.FieldAuthorName, telegraph.FieldAuthorURL,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetPage(t *testing.T) {
	page, err = client.GetPage(page.Path, &telegraph.GetPageOption{
		ReturnContent: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetPageList(t *testing.T) {
	page, err = client.GetPage(page.Path, &telegraph.GetPageOption{
		ReturnContent: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetViews(t *testing.T) {
	now := time.Now()
	_, err = client.GetViews(page.Path, now.Year(), int(now.Month()), now.Day(), &telegraph.GetViewsOption{
		Hour: now.Hour(),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Upload(t *testing.T) {
	if _, err = client.Upload([]string{
		"public/banner.png",
	}); err != nil {
		t.Fatal(err)
	}
}

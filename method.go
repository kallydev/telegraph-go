package telegraph

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	methodCreateAccount     = "createAccount"
	methodEditAccountInfo   = "editAccountInfo"
	methodGetAccountInfo    = "getAccountInfo"
	methodRevokeAccessToken = "revokeAccessToken"
	methodCreatePage        = "createPage"
	methodEditPage          = "editPage"
	methodGetPage           = "getPage"
	methodGetPageList       = "getPageList"
	methodGetViews          = "getViews"
	methodUpload            = "upload"
)

type (
	response struct {
		OK    bool   `json:"ok"`
		Error string `json:"error,omitempty"`
	}

	responseAccount struct {
		response
		Result *Account `json:"result,omitempty"`
	}

	responsePage struct {
		response
		Result *Page `json:"result,omitempty"`
	}

	responsePageList struct {
		response
		Result *PageList `json:"result,omitempty"`
	}

	responsePageViews struct {
		response
		Result *PageViews `json:"result,omitempty"`
	}

	responseUpload struct {
		Path string `json:"src"`
	}
)

type CreateAccountOption struct {
	AuthorName string
	AuthorURL  string
}

func (client *Client) CreateAccount(shortName string, option *CreateAccountOption) (account *Account, err error) {
	params := url.Values{}
	params.Add("short_name", shortName)
	if option != nil {
		if len(option.AuthorName) > 0 {
			params.Add("author_name", option.AuthorName)
		}
		if len(option.AuthorURL) > 0 {
			params.Add("author_url", option.AuthorURL)
		}
	}
	httpResponse, err := client.post(methodCreateAccount, params)
	if err != nil {
		return nil, err
	}
	responseAccountModel := new(responseAccount)
	if err := json.Unmarshal(httpResponse, responseAccountModel); err != nil {
		return nil, err
	}
	if !responseAccountModel.OK {
		return nil, errors.New(responseAccountModel.Error)
	}
	return responseAccountModel.Result, nil
}

type CreatePageOption struct {
	AuthorName    string
	AuthorURL     string
	ReturnContent bool
}

func (client *Client) CreatePage(title string, content []Node, option *CreatePageOption) (page *Page, err error) {
	params := url.Values{}
	params.Add("access_token", client.AccessToken)
	params.Add("title", title)
	contentData, err := json.Marshal(&content)
	if err != nil {
		return nil, err
	}
	params.Add("content", string(contentData))
	if option != nil {
		if len(option.AuthorName) > 0 {
			params.Add("author_name", option.AuthorName)
		}
		if len(option.AuthorURL) > 0 {
			params.Add("author_url", option.AuthorURL)
		}
		if option.ReturnContent {
			params.Add("return_content", strconv.FormatBool(option.ReturnContent))
		}
	}
	httpResponse, err := client.post(methodCreatePage, params)
	if err != nil {
		return nil, err
	}
	responsePageModel := new(responsePage)
	if err := json.Unmarshal(httpResponse, responsePageModel); err != nil {
		return nil, err
	}
	if !responsePageModel.OK {
		return nil, errors.New(responsePageModel.Error)
	}
	return responsePageModel.Result, nil
}

type EditAccountInfoOption struct {
	ShortName  string
	AuthorName string
	AuthorURL  string
}

func (client *Client) EditAccountInfo(option *EditAccountInfoOption) (account *Account, err error) {
	params := url.Values{}
	params.Add("access_token", client.AccessToken)
	if option != nil {
		if len(option.ShortName) > 0 {
			params.Add("short_name", option.ShortName)
		}
		if len(option.AuthorName) > 0 {
			params.Add("author_name", option.AuthorName)
		}
		if len(option.AuthorURL) > 0 {
			params.Add("author_url", option.AuthorURL)
		}
	}
	httpResponse, err := client.post(methodEditAccountInfo, params)
	if err != nil {
		return nil, err
	}
	responseAccountModel := new(responseAccount)
	if err := json.Unmarshal(httpResponse, responseAccountModel); err != nil {
		return nil, err
	}
	if !responseAccountModel.OK {
		return nil, errors.New(responseAccountModel.Error)
	}
	return responseAccountModel.Result, nil
}

type EditPageOption struct {
	AuthorName    string
	AuthorURL     string
	ReturnContent bool
}

func (client *Client) EditPage(path, title string, content []Node, option *EditPageOption) (page *Page, err error) {
	params := url.Values{}
	params.Add("access_token", client.AccessToken)
	params.Add("path", path)
	params.Add("title", title)
	contentData, err := json.Marshal(&content)
	if err != nil {
		return nil, err
	}
	params.Add("content", string(contentData))
	if option != nil {
		if len(option.AuthorName) > 0 {
			params.Add("author_name", option.AuthorName)
		}
		if len(option.AuthorURL) > 0 {
			params.Add("author_url", option.AuthorURL)
		}
		if option.ReturnContent {
			params.Add("return_content", strconv.FormatBool(option.ReturnContent))
		}
	}
	httpResponse, err := client.post(methodEditPage, params)
	if err != nil {
		return nil, err
	}
	responsePageModel := new(responsePage)
	if err := json.Unmarshal(httpResponse, responsePageModel); err != nil {
		return nil, err
	}
	if !responsePageModel.OK {
		return nil, errors.New(responsePageModel.Error)
	}
	return responsePageModel.Result, nil
}

const (
	FieldShortName  = "short_name"
	FieldAuthorName = "author_name"
	FieldAuthorURL  = "author_url"
	FieldAuthURL    = "auth_url"
	FieldPageCount  = "page_count"
)

type GetAccountInfoOption struct {
	Fields []string
}

func (client *Client) GetAccountInfo(option *GetAccountInfoOption) (account *Account, err error) {
	params := url.Values{}
	params.Add("access_token", client.AccessToken)
	if option != nil {
		if len(option.Fields) == 0 {
			option.Fields = []string{FieldShortName, FieldAuthorName, FieldAuthorURL}
		}
		builder := new(strings.Builder)
		for i, value := range option.Fields {
			if i == 0 {
				builder.Write([]byte{'[', 0x22})
			}
			builder.WriteString(value)
			if i == len(option.Fields)-1 {
				builder.Write([]byte{0x22, ']'})
			} else {
				builder.Write([]byte{0x22, ',', 0x22})
			}
		}
		params.Add("fields", builder.String())
	}
	httpResponse, err := client.post(methodGetAccountInfo, params)
	if err != nil {
		return nil, err
	}
	responseAccountModel := new(responseAccount)
	if err := json.Unmarshal(httpResponse, responseAccountModel); err != nil {
		return nil, err
	}
	if !responseAccountModel.OK {
		return nil, errors.New(responseAccountModel.Error)
	}
	return responseAccountModel.Result, nil
}

type GetPageOption struct {
	ReturnContent bool
}

func (client *Client) GetPage(path string, option *GetPageOption) (page *Page, err error) {
	params := url.Values{}
	params.Add("path", path)
	if option != nil {
		if option.ReturnContent {
			params.Add("return_content", strconv.FormatBool(option.ReturnContent))
		}
	}
	httpResponse, err := client.post(methodGetPage, params)
	if err != nil {
		return nil, err
	}
	responsePageModel := new(responsePage)
	if err := json.Unmarshal(httpResponse, responsePageModel); err != nil {
		return nil, err
	}
	if !responsePageModel.OK {
		return nil, errors.New(responsePageModel.Error)
	}
	return responsePageModel.Result, nil
}

type GetPageListOption struct {
	Offset int
	Limit  int
}

func (client *Client) GetPageList(option *GetPageListOption) (pageList *PageList, err error) {
	params := url.Values{}
	params.Add("access_token", client.AccessToken)
	if option != nil {
		if option.Offset > 1 {
			params.Add("offset", strconv.Itoa(option.Offset))
		}
		if option.Limit > 1 {
			params.Add("limit", strconv.Itoa(option.Limit))
		}
	}
	httpResponse, err := client.post(methodGetPageList, params)
	if err != nil {
		return nil, err
	}
	responsePageListModel := new(responsePageList)
	if err := json.Unmarshal(httpResponse, responsePageListModel); err != nil {
		return nil, err
	}
	if !responsePageListModel.OK {
		return nil, errors.New(responsePageListModel.Error)
	}
	return responsePageListModel.Result, nil
}

type GetViewsOption struct {
	Hour int
}

func (client *Client) GetViews(path string, year, month, day int, option *GetViewsOption) (pageViews *PageViews, err error) {
	params := url.Values{}
	params.Add("path", path)
	params.Add("year", strconv.Itoa(year))
	params.Add("month", strconv.Itoa(month))
	params.Add("day", strconv.Itoa(day))
	if option != nil {
		if option.Hour > 1 {
			params.Add("hour", strconv.Itoa(option.Hour))
		}
	}
	httpResponse, err := client.post(methodGetViews, params)
	if err != nil {
		return nil, err
	}
	responsePageViewsModel := new(responsePageViews)
	if err := json.Unmarshal(httpResponse, responsePageViewsModel); err != nil {
		return nil, err
	}
	if !responsePageViewsModel.OK {
		return nil, errors.New(responsePageViewsModel.Error)
	}
	return responsePageViewsModel.Result, nil
}

func (client *Client) RevokeAccessToken() (account *Account, err error) {
	params := url.Values{}
	params.Add("access_token", client.AccessToken)
	httpResponse, err := client.post(methodRevokeAccessToken, params)
	if err != nil {
		return nil, err
	}
	responseAccountModel := new(responseAccount)
	if err := json.Unmarshal(httpResponse, responseAccountModel); err != nil {
		return nil, err
	}
	if !responseAccountModel.OK {
		return nil, errors.New(responseAccountModel.Error)
	}
	return responseAccountModel.Result, nil
}

func (client *Client) Upload(filenames []string) (paths []string, err error) {
	var files []*os.File

	// Close the file handle finished processing.
	defer func() {
		for _, file := range files {
			file.Close()
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		
		files = append(files, file)

		part, err := writer.CreateFormFile(uuid.New().String(), filename)
		if err != nil {
			return nil, err
		}

		if _, err = io.Copy(part, file); err != nil {
			return nil, err
		}
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf(baseURL, methodUpload), body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	httpResponse, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	data, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	responseUploadModels := make([]responseUpload, 0)
	if err := json.Unmarshal(data, &responseUploadModels); err != nil {
		m := map[string]string{}
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, err
		}

		return nil, errors.New(strings.ToLower(m["error"]))
	}

	for _, upload := range responseUploadModels {
		paths = append(paths, upload.Path)
	}

	return paths, nil
}

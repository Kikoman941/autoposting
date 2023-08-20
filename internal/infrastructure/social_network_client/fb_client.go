package social_network_client

import (
	"encoding/json"
	"fmt"
	"github.com/ztrue/tracerr"
	"io/ioutil"
	"net/http"
	"net/url"
)

type fbCreatePostResponse struct {
	PostID string `json:"id"`
}

type FBCredentials struct {
	AppID        string `json:"app_id"`
	AccessToken  string `json:"access_token"`
	ClientSecret string `json:"client_secret"`
}

type fbClient struct {
	httpClient  *http.Client
	authApiUrl  string
	workApiUrl  string
	redirectUrl string
}

type fbAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type fbGetAccountPagesResponse struct {
	Data []struct {
		AccessToken string `json:"access_token"`
		Name        string `json:"name"`
		Id          string `json:"id"`
	} `json:"data"`
	Paging struct {
		Cursors struct {
			Before string `json:"before"`
			After  string `json:"after"`
		} `json:"cursors"`
	} `json:"paging"`
}

func (f *fbClient) GetAuthURL(credentials string) (string, error) {
	fbCredentials, err := f.stringToFBCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v16.0/dialog/oauth", f.authApiUrl), nil)
	if err != nil {
		return "", tracerr.Errorf("cannot create auth url request")
	}

	q := url.Values{
		"client_id":     []string{fbCredentials.AppID},
		"redirect_uri":  []string{f.redirectUrl},
		"response_type": []string{"code"},
		"scope":         []string{"pages_show_list,pages_read_engagement,pages_manage_posts"},
		"display":       []string{"popup"},
	}
	req.URL.RawQuery = q.Encode()

	return req.URL.String(), nil
}

func (f *fbClient) GetAccessToken(credentials string, queryParams map[string][]string) (string, error) {
	var (
		data fbAccessTokenResponse
	)

	fbCredentials, err := f.stringToFBCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", f.workApiUrl, "oauth/access_token"), nil)
	if err != nil {
		return "", tracerr.Errorf("cannot create access token request:\n%s", err)
	}
	q := url.Values{
		"client_id":     []string{fbCredentials.AppID},
		"client_secret": []string{fbCredentials.ClientSecret},
		"redirect_uri":  []string{f.redirectUrl},
		"code":          []string{queryParams["code"][0]},
	}
	req.URL.RawQuery = q.Encode()
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", tracerr.Errorf("cannot get access token:\n%s", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", tracerr.Errorf("cannot read access token response:\n%s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", tracerr.Errorf(
			"access token response status is %d\ntokenResponse:%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return "", tracerr.Errorf("cannot unmarshal access token body:\n%s", err)
	}

	return data.AccessToken, nil
}

func (f *fbClient) GetAccountPages(_, accessToken string) ([]SocialNetworkPage, error) {
	var (
		pages []SocialNetworkPage
		data  fbGetAccountPagesResponse
	)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v16.0/me/accounts", f.workApiUrl), nil)
	if err != nil {
		return nil, tracerr.Errorf("cannot create getting account pages request:\n%s", err)
	}
	q := url.Values{
		"admin`_only":  []string{"true"},
		"access_token": []string{accessToken},
	}
	req.URL.RawQuery = q.Encode()
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, tracerr.Errorf("cannot get account pages:\n%s", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, tracerr.Errorf("cannot read getting pages response:\n%s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, tracerr.Errorf(
			"get account pages response status %d\npagesResponse:%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, tracerr.Errorf("cannot unmarshal getting pages body:\n%s", err)
	}

	for _, page := range data.Data {
		pages = append(pages, SocialNetworkPage{
			ID:   page.Id,
			Name: page.Name,
		})
	}

	return pages, nil
}

func (f *fbClient) CreatePost(credentials string, groupID string, post string) (string, error) {
	var (
		data fbCreatePostResponse
	)

	fbCredentials, err := f.stringToFBCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/feed", f.workApiUrl, groupID), nil)
	if err != nil {
		return "", tracerr.Errorf("cannot create createPost request:\n%s", err)
	}

	q := req.URL.Query()
	q.Add("access_token", fbCredentials.AccessToken)
	q.Add("message", post)
	req.URL.RawQuery = q.Encode()
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", tracerr.Errorf("cannot create post:\n%s", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", tracerr.Errorf("cannot read create post response:\n%s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", tracerr.Errorf(
			"create post response status is %d\nresponse:%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return "", tracerr.Errorf("cannot unmarshal access token body:\n%s", err)
	}

	return data.PostID, nil
}

func (f *fbClient) DeletePost() {
	//TODO implement me
	panic("implement me")
}

func (f *fbClient) UploadImage() {
	//TODO implement me
	panic("implement me")
}

func (f *fbClient) stringToFBCredentials(credentials string) (*FBCredentials, error) {
	fbCredentials := &FBCredentials{}
	err := json.Unmarshal([]byte(credentials), fbCredentials)
	if err != nil {
		return nil, tracerr.Errorf("cannot unmarshal fb credentials {%s}:\n%s", credentials, err)
	}
	return fbCredentials, nil
}

func NewFBClient() SocialNetworkClient {
	return &fbClient{
		httpClient:  &http.Client{},
		authApiUrl:  "https://www.facebook.com",
		workApiUrl:  "https://graph.facebook.com",
		redirectUrl: "http://localhost:8080/auth/get_token?socialNetwork=FB",
	}
}

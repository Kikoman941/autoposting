package social_network_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type fbCreatePostResponse struct {
	PostID string `json:"id"`
}

type FBCredentials struct {
	AppID       string `json:"app_id"`
	AccessToken string `json:"access_token"`
}

type fbClient struct {
	httpClient  *http.Client
	authApiUrl  string
	workApiUrl  string
	redirectUrl string
}

func (f *fbClient) GetAuthURL(credentials string) (string, error) {
	fbCredentials, err := f.stringToFBCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v15.0/dialog/oauth", f.authApiUrl), nil)
	if err != nil {
		return "", fmt.Errorf("cannot create auth url request")
	}

	q := url.Values{
		"client_id":     []string{fbCredentials.AppID},
		"redirect_uri":  []string{f.redirectUrl},
		"response_type": []string{"token"},
		"scope":         []string{"pages_show_list,pages_read_engagement,pages_manage_posts"},
		"display":       []string{"popup"},
	}
	req.URL.RawQuery = q.Encode()

	return req.URL.String(), nil
}

func (f *fbClient) GetAccessToken(credentials string, queryParams map[string][]string) (string, error) {
	var fbCredentials VKCredentials
	var data AccessTokenResponse

	err := json.Unmarshal([]byte(credentials), &fbCredentials)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal fb credentials {%s}:\n%s", credentials, err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", f.authApiUrl, fbCredentials.AppID), nil)
	if err != nil {
		return "", fmt.Errorf("cannot create access token request:\n%s", err)
	}
	q := url.Values{
		"client_id":     []string{fbCredentials.AppID},
		"client_secret": []string{fbCredentials.SecureKey},
		"redirect_uri":  []string{v.redirectUrl},
		"code":          []string{queryParams["code"][0]},
	}
	req.URL.RawQuery = q.Encode()
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot get access token:\n%s", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read access token response:\n%s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"access token response status is %d\ntokenResponse:%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal access token body:\n%s", err)
	}

	return data.AccessToken, nil
}

func (f *fbClient) CreatePost(credentials string, groupID string, post string) (string, error) {
	var data fbCreatePostResponse
	var fbCredentials FBCredentials

	err := json.Unmarshal([]byte(credentials), &fbCredentials)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal fb credentials {%s}:\n%s", credentials, err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/feed", f.workApiUrl, groupID), nil)
	if err != nil {
		return "", fmt.Errorf("cannot create createPost request:\n%s", err)
	}

	q := req.URL.Query()
	q.Add("access_token", fbCredentials.AccessToken)
	q.Add("message", post)
	req.URL.RawQuery = q.Encode()
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot create post:\n%s", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read create post response:\n%s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"create post response status is %d\nresponse:%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal access token body:\n%s", err)
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
		return nil, fmt.Errorf("cannot unmarshal fb credentials {%s}:\n%s", credentials, err)
	}
	return fbCredentials, nil
}

func NewFBClient() SocialNetworkClient {
	return &fbClient{
		httpClient:  &http.Client{},
		authApiUrl:  "https://www.facebook.com",
		workApiUrl:  "https://graph.facebook.com",
		redirectUrl: "http://localhost:8080/auth/",
	}
}

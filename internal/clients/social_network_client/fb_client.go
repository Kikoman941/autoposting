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
	AccessToken string `json:"access_token"`
}

type fbClient struct {
	httpClient *http.Client
	workApiUrl string
}

func (f *fbClient) GetAuthURL(credentials string) (string, error) {
	vkCredentials, err := f.stringToFBCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/authorize", f.authApiUrl), nil)
	if err != nil {
		return "", fmt.Errorf("cannot create auth url request")
	}

	q := url.Values{
		"client_id":     []string{vkCredentials.AppID},
		"redirect_uri":  []string{v.redirectUrl},
		"response_type": []string{"code"},
		"scope":         []string{"offline,groups,photos,video,pages,wall"},
	}
	req.URL.RawQuery = q.Encode()

	return req.URL.String(), nil
}

func (f *fbClient) CreatePost(credentials string, groupID string, post string) (string, error) {
	var data fbCreatePostResponse
	var fbCredentials FBCredentials

	err := json.Unmarshal([]byte(credentials), &fbCredentials)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal fb credentials {%s}: %s", credentials, err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/feed", f.workApiUrl, groupID), nil)
	if err != nil {
		return "", fmt.Errorf("cannot create createPost request: %s", err)
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
		return nil, fmt.Errorf("cannot unmarshal fb credentials {%s}: %s", credentials, err)
	}
	return fbCredentials, nil
}

func NewFBClient() SocialNetworkClient {
	return &fbClient{
		httpClient: &http.Client{},
		workApiUrl: "https://graph.facebook.com",
	}
}

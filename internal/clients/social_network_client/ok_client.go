package social_network_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type OKCredentials struct {
	AppID       string `json:"app_id"`
	PublicKey   string `json:"public_key"`
	AccessToken string `json:"access_token"`
}

type okClient struct {
	httpClient  *http.Client
	authApiUrl  string
	workApiUrl  string
	redirectUrl string
}

func (o *okClient) GetAuthURL(credentials string) (string, error) {
	okCredentials, err := o.stringToOKCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/oauth/authorize", o.authApiUrl), nil)
	if err != nil {
		return "", fmt.Errorf("cannot create auth url request")
	}

	q := url.Values{
		"client_id":     []string{okCredentials.AppID},
		"redirect_uri":  []string{o.redirectUrl},
		"response_type": []string{"token"},
		"scope":         []string{"VALUABLE_ACCESS;LONG_ACCESS_TOKEN;PHOTO_CONTENT;GROUP_CONTENT;VIDEO_CONTENT"},
	}
	req.URL.RawQuery = q.Encode()

	return req.URL.String(), nil
}

func (o *okClient) GetAccessToken(s string, m map[string][]string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (o *okClient) CreatePost(credentials string, groupID string, post string) (string, error) {
	var okCredentials OKCredentials

	err := json.Unmarshal([]byte(credentials), &okCredentials)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal ok credentials {%s}:\n%s", credentials, err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/mediatopic/post", o.workApiUrl), nil)
	if err != nil {
		return "", fmt.Errorf("cannot create createPost request:\n%s", err)
	}

	q := req.URL.Query()
	q.Add("application_key", okCredentials.PublicKey)
	q.Add("access_token", okCredentials.AccessToken)
	q.Add("type", "GROUP_THEME")
	q.Add("gid", groupID)
	q.Add("attachment", post)
	req.URL.RawQuery = q.Encode()
	resp, err := o.httpClient.Do(req)
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

	return string(respBody), nil
}

func (o *okClient) DeletePost() {
	//TODO implement me
	panic("implement me")
}

func (o *okClient) UploadImage() {
	//TODO implement me
	panic("implement me")
}

func (o *okClient) stringToOKCredentials(credentials string) (*OKCredentials, error) {
	okCredentials := &OKCredentials{}
	err := json.Unmarshal([]byte(credentials), okCredentials)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal ok credentials {%s}:\n%s", credentials, err)
	}
	return okCredentials, nil
}

func NewOKClient() SocialNetworkClient {
	client := okClient{
		httpClient:  &http.Client{},
		authApiUrl:  "https://connect.ok.ru",
		workApiUrl:  "https://api.ok.ru",
		redirectUrl: "http://localhost:8080/auth/",
	}
	return &client
}

package social_network_client

import (
	logging "amplifr/pkg"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type VKCredentials struct {
	AppID       string `json:"app_id"`
	SecureKey   string `json:"secure_key"`
	ServiceKey  string `json:"service_key"`
	AccessToken string `json:"access_token"`
	Code        string `json:"code"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserId      int    `json:"user_id"`
}

type vkCreatePostResponse struct {
	Response struct {
		PostID int `json:"post_id"`
	} `json:"response"`
}

type vkClient struct {
	logger      *logging.Logger
	httpClient  *http.Client
	authApiUrl  string
	workApiUrl  string
	redirectUrl string
	scope       string
}

func (v *vkClient) GetAuthURL(credentials string) (string, error) {
	vkCredentials, err := v.stringToVKCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/authorize", v.authApiUrl), nil)
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

func (v *vkClient) CreatePost(credentials string, groupID, post string) (string, error) {
	var data vkCreatePostResponse
	var vkCredentials VKCredentials

	err := json.Unmarshal([]byte(credentials), &vkCredentials)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal vk credentials {%s}: %s", credentials, err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/method/wall.post", v.workApiUrl), nil)
	if err != nil {
		return "", fmt.Errorf("cannot create createPost request: %s", err)
	}

	q := req.URL.Query()
	q.Add("owner_id", groupID)
	q.Add("access_token", vkCredentials.AccessToken)
	q.Add("from_group", "1")
	q.Add("message", post)
	q.Add("v", "5.131")
	req.URL.RawQuery = q.Encode()
	resp, err := v.httpClient.Do(req)
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

	return strconv.Itoa(data.Response.PostID), nil
}

func (v *vkClient) DeletePost() {
	//TODO implement me
	panic("implement me")
}

func (v *vkClient) UploadImage() {
	//TODO implement me
	panic("implement me")
}

func (v *vkClient) stringToVKCredentials(credentials string) (*VKCredentials, error) {
	vkCredentials := &VKCredentials{}
	err := json.Unmarshal([]byte(credentials), vkCredentials)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal vk credentials {%s}: %s", credentials, err)
	}
	return vkCredentials, nil
}

func (v *vkClient) getAccessToken(credentials *VKCredentials) (string, error) {
	var data AccessTokenResponse
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/access_token", v.authApiUrl), nil)
	if err != nil {
		return "", fmt.Errorf("cannot create access token request: %s", err)
	}
	q := req.URL.Query()
	q.Add("client_id", credentials.AppID)
	q.Add("client_secret", credentials.SecureKey)
	q.Add("redirect_uri", v.redirectUrl)
	q.Add("code", credentials.Code)
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

func NewVKClient() SocialNetworkClient {
	client := vkClient{
		httpClient:  &http.Client{},
		authApiUrl:  "https://oauth.vk.com",
		workApiUrl:  "https://api.vk.com",
		redirectUrl: "http://localhost:8080/auth/",
		scope:       "offline,groups,photos,video,pages,wall",
	}
	return &client
}

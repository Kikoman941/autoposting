package social_network_client

import (
	"encoding/json"
	"fmt"
	"github.com/ztrue/tracerr"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type VKCredentials struct {
	AppID       string `json:"app_id"`
	SecureKey   string `json:"secure_key"`
	ServiceKey  string `json:"service_key"`
	UserId      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type vkAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserId      int    `json:"user_id"`
}

type vkCreatePostResponse struct {
	Response struct {
		PostID int `json:"post_id"`
	} `json:"response"`
}

type vkGetAccountPagesResponse struct {
	Response struct {
		Count int `json:"count"`
		Items []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Image string `json:"photo_200"`
		} `json:"items"`
	}
}

type vkClient struct {
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
		return "", tracerr.New("cannot create auth url request")
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

func (v *vkClient) GetAccessToken(credentials string, queryParams map[string][]string) (string, error) {
	var data vkAccessTokenResponse

	vkCredentials, err := v.stringToVKCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/access_token", v.authApiUrl), nil)
	if err != nil {
		return "", tracerr.Errorf("cannot create access token request:\n%s", err)
	}
	q := url.Values{
		"client_id":     []string{vkCredentials.AppID},
		"client_secret": []string{vkCredentials.SecureKey},
		"redirect_uri":  []string{v.redirectUrl},
		"code":          []string{queryParams["code"][0]},
	}
	req.URL.RawQuery = q.Encode()
	resp, err := v.httpClient.Do(req)
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

func (v *vkClient) GetAccountPages(credentials, accessToken string) ([]SocialNetworkPage, error) {
	var (
		pages []SocialNetworkPage
		data  vkGetAccountPagesResponse
	)
	vkCredentials, err := v.stringToVKCredentials(credentials)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/method/groups.get", v.workApiUrl), nil)
	if err != nil {
		return nil, tracerr.Errorf("cannot create getting account pages request:\n%s", err)
	}
	q := url.Values{
		"access_token": []string{accessToken},
		"user_id":      []string{vkCredentials.UserId},
		"extended":     []string{"1"},
		"filter":       []string{"admin"},
		"fields":       []string{"id,name,photo_200"},
		"v":            []string{"5.131"},
	}
	req.URL.RawQuery = q.Encode()
	resp, err := v.httpClient.Do(req)
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

	for _, page := range data.Response.Items {
		pages = append(pages, SocialNetworkPage{
			ID:    strconv.Itoa(page.ID),
			Name:  page.Name,
			Image: page.Image,
		})
	}

	return pages, nil
}

func (v *vkClient) CreatePost(credentials string, groupID, post string) (string, error) {
	var data vkCreatePostResponse
	vkCredentials, err := v.stringToVKCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/method/wall.post", v.workApiUrl), nil)
	if err != nil {
		return "", fmt.Errorf("cannot create createPost request:\n%s", err)
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
		return nil, tracerr.Errorf("cannot unmarshal vk credentials {%s}:\n%s", credentials, err)
	}
	return vkCredentials, nil
}

func NewVKClient() SocialNetworkClient {
	client := vkClient{
		httpClient:  &http.Client{},
		authApiUrl:  "https://oauth.vk.com",
		workApiUrl:  "https://api.vk.com",
		redirectUrl: "http://localhost:8080/auth/get_token?socialNetwork=VK",
		scope:       "offline,groups,photos,video,pages,wall",
	}
	return &client
}

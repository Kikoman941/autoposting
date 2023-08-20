package social_network_client

import (
	"encoding/json"
	"fmt"
	"github.com/ztrue/tracerr"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/url"
)

type OKCredentials struct {
	AppID       string `json:"app_id"`
	PublicKey   string `json:"public_key"`
	SecretKey   string `json:"secret_key"`
	AccessToken string `json:"access_token"`
}

type okClient struct {
	httpClient  *http.Client
	authApiUrl  string
	workApiUrl  string
	redirectUrl string
}

type okAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

type okGetAccountPagesResponse struct {
	Groups []struct {
		GroupId string `json:"groupId"`
		Status  string `json:"status"`
	} `json:"groups"`
	Anchor string `json:"anchor"`
}

type okGetPagesInfoResponse struct {
	ID             string `json:"uid"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	PreviewImageId string `json:"photo_id"`
	_              bool   `json:"revenue_pp_enabled"`
	_              bool   `json:"pin_notifications_off"`
}

type okGetImageInfoResponse struct {
	Photo struct {
		Type     string `json:"type"`
		ImageUrl string `json:"pic128x128"`
		_        bool   `json:"text_detected"`
	}
}

func (o *okClient) GetAuthURL(credentials string) (string, error) {
	okCredentials, err := o.stringToOKCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/oauth/authorize", o.authApiUrl), nil)
	if err != nil {
		return "", tracerr.Errorf("cannot create auth url request")
	}

	q := url.Values{
		"client_id":     []string{okCredentials.AppID},
		"redirect_uri":  []string{o.redirectUrl},
		"response_type": []string{"code"},
		"scope":         []string{"VALUABLE_ACCESS;LONG_ACCESS_TOKEN;PHOTO_CONTENT;GROUP_CONTENT;VIDEO_CONTENT"},
	}
	req.URL.RawQuery = q.Encode()

	return req.URL.String(), nil
}

func (o *okClient) GetAccessToken(credentials string, queryParams map[string][]string) (string, error) {
	var (
		data okAccessTokenResponse
	)

	okCredentials, err := o.stringToOKCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", o.workApiUrl, "/oauth/token.do"), nil)
	if err != nil {
		return "", tracerr.Errorf("cannot create access token request:\n%s", err)
	}
	q := url.Values{
		"code":          []string{queryParams["code"][0]},
		"client_id":     []string{okCredentials.AppID},
		"client_secret": []string{okCredentials.SecretKey},
		"redirect_uri":  []string{o.redirectUrl},
		"grant_type":    []string{"authorization_code"},
	}
	req.URL.RawQuery = q.Encode()
	resp, err := o.httpClient.Do(req)
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
			resp.Request.URL,
		)
	}

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return "", tracerr.Errorf("cannot unmarshal access token body:\n%s", err)
	}

	return data.AccessToken, nil
}

func (o *okClient) GetAccountPages(credentials, accessToken string) ([]SocialNetworkPage, error) {
	var (
		pagesIds []string
		data     okGetAccountPagesResponse
	)

	okCredentials, err := o.stringToOKCredentials(credentials)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%s/api/group/getUserGroupsV2", o.workApiUrl), nil,
	)
	if err != nil {
		return nil, tracerr.Errorf("cannot create getting account pages request:\n%s", err)
	}
	q := url.Values{
		"application_key":    []string{okCredentials.PublicKey},
		"access_token":       []string{accessToken},
		"session_secret_key": []string{okCredentials.SecretKey},
		"format":             []string{"json"},
	}
	req.URL.RawQuery = q.Encode()
	resp, err := o.httpClient.Do(req)
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

	for _, page := range data.Groups {
		if page.Status == "ADMIN" {
			pagesIds = append(pagesIds, page.GroupId)
		}
	}

	return o.getPagesInfo(okCredentials, accessToken, pagesIds)
}

func (o *okClient) getPagesInfo(
	okCredentials *OKCredentials,
	accessToken string,
	pagesIds []string,
) ([]SocialNetworkPage, error) {
	var (
		pages []SocialNetworkPage
		data  []okGetPagesInfoResponse
	)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/group/getInfo", o.workApiUrl), nil)
	if err != nil {
		return nil, tracerr.Errorf("cannot create getting pages info request:\n%s", err)
	}
	q := url.Values{
		"application_key":    []string{okCredentials.PublicKey},
		"access_token":       []string{accessToken},
		"session_secret_key": []string{okCredentials.SecretKey},
		"format":             []string{"json"},
		"uids":               pagesIds,
		"fields":             []string{"name,description,photo_id,uid"},
	}
	req.URL.RawQuery = q.Encode()
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, tracerr.Errorf("cannot get pages info:\n%s", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, tracerr.Errorf("cannot read getting pages info response:\n%s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, tracerr.Errorf(
			"get pages info response status %d\npagesInfoResponse:%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, tracerr.Errorf("cannot unmarshal getting pages info body:\n%s", err)
	}

	for _, page := range data {
		imageUrl, err := o.getImageUrl(okCredentials, accessToken, page.PreviewImageId)
		if err != nil {
			slog.Warn(
				"failed to get image info",
				slog.String("imageId", page.PreviewImageId),
				slog.Any("err", err),
			)
		}
		pages = append(pages, SocialNetworkPage{
			ID:          page.ID,
			Name:        page.Name,
			Description: page.Description,
			Image:       imageUrl,
		})
	}

	return pages, nil
}

func (o *okClient) getImageUrl(okCredentials *OKCredentials, accessToken string, imageId string) (string, error) {
	var data okGetImageInfoResponse

	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%s/api/photos/getPhotoInfo", o.workApiUrl), nil,
	)
	if err != nil {
		return "error", tracerr.Errorf("cannot create getting image info request:\n%s", err)
	}
	q := url.Values{
		"photo_id":           []string{imageId},
		"fields":             []string{"photo.PIC128X128"},
		"application_key":    []string{okCredentials.PublicKey},
		"access_token":       []string{accessToken},
		"session_secret_key": []string{okCredentials.SecretKey},
		"format":             []string{"json"},
	}
	req.URL.RawQuery = q.Encode()
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return "error", tracerr.Errorf("cannot get image info:\n%s", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "error", tracerr.Errorf("cannot read getting image info response:\n%s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "error", tracerr.Errorf(
			"get image info response status %d\npagesInfoResponse:%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return "error", tracerr.Errorf("cannot unmarshal getting image info body:\n%s", err)
	}

	return data.Photo.ImageUrl, nil
}

func (o *okClient) CreatePost(credentials string, groupID string, post string) (string, error) {
	okCredentials, err := o.stringToOKCredentials(credentials)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/mediatopic/post", o.workApiUrl), nil)
	if err != nil {
		return "", tracerr.Errorf("cannot create createPost request:\n%s", err)
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
		return nil, tracerr.Errorf("cannot unmarshal ok credentials {%s}:\n%s", credentials, err)
	}
	return okCredentials, nil
}

func NewOKClient() SocialNetworkClient {
	client := okClient{
		httpClient:  &http.Client{},
		authApiUrl:  "https://connect.ok.ru",
		workApiUrl:  "https://api.ok.ru",
		redirectUrl: "http://localhost:8080/auth/get_token?socialNetwork=OK",
	}
	return &client
}

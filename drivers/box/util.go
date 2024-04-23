package box

import (
	"fmt"
	"net/http"

	"github.com/alist-org/alist/v3/drivers/base"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func (d *Box) refreshToken() error {
	var resp base.TokenResp
	var e TokenError
	res, err := base.RestyClient.R().SetResult(&resp).SetError(&e).
		SetFormData(map[string]string{
			"client_id":     d.ClientID,
			"client_secret": d.ClientSecret,
			"refresh_token": d.RefreshToken,
			"grant_type":    "refresh_token",
		}).Post("https://api.box.com/oauth2/token")

	if err != nil {
		return err
	}
	log.Info(res.String())
	if e.Error != "" {
		return fmt.Errorf("%s: %s", e.Error, e.ErrorDescription)
	}
	d.AccessToken = resp.AccessToken
	d.RefreshToken = resp.RefreshToken
	return nil
}

func (d *Box) request(url string, method string, callback base.ReqCallback, resp interface{}) ([]byte, error) {
	req := base.RestyClient.R()
	req.SetHeader("Authorization", "Bearer "+d.AccessToken)
	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}

	res, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 401 {
		err = d.refreshToken()
		if err != nil {
			return nil, err
		}
		return d.request(url, method, callback, resp)
	} else if res.StatusCode() > 399 {
		return nil, fmt.Errorf(res.String())
	}

	return res.Body(), nil
}

func (d *Box) getFiles(id string) ([]File, error) {
	nextMarker := "first"
	res := make([]File, 0)

	for nextMarker != "" {
		if nextMarker == "first" {
			nextMarker = ""
		}

		var resp Files
		_, err := d.request(fmt.Sprintf("https://api.box.com/2.0/folders/%s/items", id), http.MethodGet, func(req *resty.Request) {
			req.SetQueryParam("limit", "1000")
			req.SetQueryParam("usemarker", "true")
			req.SetQueryParam("fields", "type,id,name,size,created_at,modified_at,content_created_at,content_modified_at,item_status")
			if nextMarker != "" {
				req.SetQueryParam("marker", nextMarker)
			}
		}, &resp)

		if err != nil {
			return nil, err
		}

		nextMarker = resp.NextMarker
		res = append(res, resp.Entries...)
	}

	return res, nil
}

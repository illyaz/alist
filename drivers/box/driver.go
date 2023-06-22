package box

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/alist-org/alist/v3/drivers/base"
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/pkg/utils"
)

type Box struct {
	model.Storage
	Addition
	AccessToken string
}

func (d *Box) Config() driver.Config {
	return config
}

func (d *Box) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *Box) Init(ctx context.Context) error {
	return d.refreshToken()
}

func (d *Box) Drop(ctx context.Context) error {
	return nil
}

func (d *Box) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	files, err := d.getFiles(dir.GetID())
	if err != nil {
		return nil, err
	}

	return utils.SliceConvert(files, func(src File) (model.Obj, error) {
		return fileToObj(src, args), nil
	})
}

func (d *Box) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	link := model.Link{}

	if args.Type == "thumb" {
		link.URL = fmt.Sprintf("https://api.box.com/2.0/files/%s/thumbnail.jpg?max_width=160", file.GetID())
		link.Header = http.Header{
			"Authorization": []string{"Bearer " + d.AccessToken},
		}
	} else {
		req := base.NoRedirectClient.R()
		req.SetHeader("Authorization", "Bearer "+d.AccessToken)
		res, err := req.Get(fmt.Sprintf("https://api.box.com/2.0/files/%s/content", file.GetID()))

		if err != nil {
			return nil, err
		}

		link.Status = 302
		link.Header = http.Header{
			"Location": []string{res.Header().Get("location")},
		}
		link.Writer = func(w io.Writer) error { return nil }
	}

	return &link, nil
}

var _ driver.Driver = (*Box)(nil)

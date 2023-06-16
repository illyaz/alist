package box

import (
	stdpath "path"
	"time"

	"github.com/alist-org/alist/v3/internal/sign"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/alist-org/alist/v3/server/common"

	"github.com/alist-org/alist/v3/internal/model"
	log "github.com/sirupsen/logrus"
)

type TokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type Files struct {
	Entries    []File `json:"entries"`
	NextMarker string `json:"next_marker"`
}

type File struct {
	Type       string    `json:"type"`
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modified_at"`
	ItemStatus string    `json:"item_status"`
}

func fileToObj(f File, args model.ListArgs) *model.ObjThumb {
	log.Debugf("box file: %+v", f)

	var thumb = ""

	if f.Type == "file" {
		thumb = common.GetApiUrl(nil) + stdpath.Join("/d", args.ReqPath, f.Name)
		thumb = utils.EncodePath(thumb, true)
		thumb += "?type=thumb&sign=" + sign.Sign(stdpath.Join(args.ReqPath, f.Name))
	}

	obj := &model.ObjThumb{
		Object: model.Object{
			ID:       f.Id,
			Name:     f.Name,
			Size:     f.Size,
			Modified: f.ModifiedAt,
			IsFolder: f.Type == "folder",
		},
		Thumbnail: model.Thumbnail{
			Thumbnail: thumb,
		},
	}

	return obj
}

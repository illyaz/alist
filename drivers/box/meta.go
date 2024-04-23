package box

import (
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/op"
)

type Addition struct {
	driver.RootID
	RefreshToken string `json:"refresh_token" required:"true"`
	ClientID     string `json:"client_id" required:"true" default:"d0374ba6pgmaguie02ge15sv1mllndho"`
	ClientSecret string `json:"client_secret" required:"true" default:"grWXGU7zW6034GI54GuswaDQdE30QOfn"`
}

var config = driver.Config{
	Name:        "Box",
	OnlyProxy:   true,
	NoUpload:    true,
	DefaultRoot: "0",
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &Box{}
	})
}

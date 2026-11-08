package google

import "github.com/gokhanaltun/o2go/v2"

func New(cfg *o2go.Config, opts ...o2go.Option) (*Google, error) {
	google := &Google{}
	oauth2, err := o2go.New(google, cfg, opts...)
	if err != nil {
		return nil, err
	}

	google.oauth2 = oauth2

	return google, nil
}

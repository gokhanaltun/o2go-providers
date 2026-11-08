package github

import "github.com/gokhanaltun/o2go/v2"

func New(cfg *o2go.Config, opts ...o2go.Option) (*Github, error) {
	github := &Github{}
	opts = append(opts, withJSONTokenResponse())
	oauth2, err := o2go.New(github, cfg, opts...)
	if err != nil {
		return nil, err
	}

	github.oauth2 = oauth2

	return github, nil
}

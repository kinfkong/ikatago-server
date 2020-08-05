package storage

import "github.com/kinfkong/ikatago-server/model"

type Provider interface {
	SaveUserSSHInfo(userSSHInfo model.SSHLoginInfo) error
}

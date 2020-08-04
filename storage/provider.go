type Provider interface {
	SaveUserSSHInfo(userSSHInfo model.SSHLoginInfo) error
}
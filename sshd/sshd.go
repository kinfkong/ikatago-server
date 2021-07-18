package sshd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/kinfkong/ikatago-server/config"
	"github.com/kinfkong/ikatago-server/utils"
)

// UserInfo represents the user info
type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Users all the sshd users
var Users []UserInfo

// SSHCommandHandler the ssh command handler
type SSHCommandHandler func(ssh.Session, ...string) (*exec.Cmd, error)

// Handlers all the sshd handlers
var Handlers map[string]SSHCommandHandler = make(map[string]SSHCommandHandler)

// RegisterCommandHandler registers the command handler
func RegisterCommandHandler(commandName string, handler SSHCommandHandler) {
	Handlers[commandName] = handler
}

func readUsers(filename string) []UserInfo {
	if config.GetConfig().GetBool("clusterModeEnabled") {
		// ignore the local users for cluster mode
		log.Println("LOCAL USER is ignore because it is in cluster mode")
		return make([]UserInfo, 0)
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	result := make([]UserInfo, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		items := strings.Split(line, ":")
		if len(items) != 2 {
			log.Printf("WARN: cannot read user line: %s\n", line)
			continue
		}
		match, _ := regexp.MatchString("^[0-9a-zA-Z_\\-]+$", items[0])
		if !match {
			log.Printf("WARN: invalid user name (only digits and letters only): %s\n", items[0])
			continue
		}
		result = append(result, UserInfo{
			Username: items[0],
			Password: items[1],
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result

}

// RunAsync runs the sshd
func RunAsync() error {
	Users = readUsers(config.GetConfig().GetString("users.file"))
	ssh.Handle(func(s ssh.Session) {
		defer s.Exit(0)
		cmds := s.Command()
		if len(cmds) == 0 {
			io.WriteString(s, "No command found.\n")
			return
		}
		handler, ok := Handlers[cmds[0]]
		if !ok {
			io.WriteString(s, fmt.Sprintf("command [%s] is not supported.\n", cmds[0]))
			return
		}
		log.Printf("INFO user[%s] executing cmd: %+v\n", s.User(), cmds)
		cmd, err := handler(s, cmds[1:]...)
		if err != nil {
			io.WriteString(s, fmt.Sprintf("Something error happens...\nerr:%+v\n", err))
			log.Printf("INFO user [%s] session done\n", s.User())
			return
		}
		if cmd == nil {
			// nothing to do
			log.Printf("INFO user [%s] session done\n", s.User())
			return
		}
		username := s.User()
		err = utils.GetCmdManager().RunCommand(&username, "katago", cmd)
		if err != nil {
			log.Println(err)
		}
		log.Printf("INFO user [%s] session done\n", username)
	})

	passwordAuthOption := ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		for _, userInfo := range Users {
			if userInfo.Username == ctx.User() && userInfo.Password == password {
				return true
			}
		}
		return false
	})
	publicKeyAuthOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		publicKeyInConfig := config.GetConfig().GetString("auth.publicKey")
		if publicKeyInConfig == "" {
			return false
		}
		pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKeyInConfig))
		if err != nil {
			log.Printf("Failed to read key: %+v", err)
			return false
		}
		return ssh.KeysEqual(pubKey, key)

	})
	go func() {
		err := ssh.ListenAndServe(config.GetConfig().GetString("server.listen"), nil, passwordAuthOption, publicKeyAuthOption)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

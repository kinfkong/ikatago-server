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
)

// UserInfo represents the user info
type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func readUsers(filename string) []UserInfo {
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
		log.Printf("USER: %s\n", items[0])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result

}
func runBaiduKatago(args ...string) (*exec.Cmd, error) {
	// exec.Command()
	// 	_, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
	// decrypte data
	decryptePassword := "abcde12345"
	decrypteCommandTemplate := "openssl enc -in %s -d -aes-256-cbc -pass pass:%s > %s"

	output, err := exec.Command("/bin/sh", "-c", "rm -rf /tmp/l && mkdir -p /tmp/l").CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}
	output, err = exec.Command("/bin/sh", "-c", "rm -rf /tmp/b && mkdir -p /tmp/b").CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}
	output, err = exec.Command("/bin/sh", "-c", fmt.Sprintf(decrypteCommandTemplate, "./data/k", decryptePassword, "/tmp/b/k")).CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}
	output, err = exec.Command("/bin/sh", "-c", fmt.Sprintf(decrypteCommandTemplate, "./data/lc", decryptePassword, "/tmp/l/lc")).CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}
	output, err = exec.Command("/bin/sh", "-c", fmt.Sprintf(decrypteCommandTemplate, "./data/lc", decryptePassword, "/tmp/l/lz")).CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}

	output, err = exec.Command("/bin/sh", "-c", "chmod +x /tmp/b/k").CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}
	output, err = exec.Command("/bin/sh", "-c", "rm -f /tmp/l/libzip.so.4 && ln -s /tmp/l/lz /tmp/l/libzip.so.4").CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}
	output, err = exec.Command("/bin/sh", "-c", "rm -f /tmp/l/libstdc++.so.6 && ln -s /tmp/l/lc /tmp/l/libstdc++.so.6").CombinedOutput()
	if err != nil {
		log.Printf("DEBUG error output: %s\n", string(output))
		return nil, err
	}
	return exec.Command("/bin/sh", "-c", "export LD_LIBRARY_PATH=/tmp/l:$LD_LIBRARY_PATH; /tmp/b/k gtp -config ./gtp_example.cfg -model ./weight.bin.gz"), nil
}

// Run runs the sshd
func Run() {
	users := readUsers("./userlist.txt")
	ssh.Handle(func(s ssh.Session) {
		cmds := s.Command()
		if len(cmds) == 0 {
			io.WriteString(s, "No command found.\n")
			return
		}
		if cmds[0] != "run-katago" {
			io.WriteString(s, "Only run-katago command is supported.\n")
			return
		}
		log.Printf("DEBUG executing cmd: %+v\n", cmds)
		cmd, err := runBaiduKatago(cmds[1:]...)
		if err != nil {
			io.WriteString(s, fmt.Sprintf("Something error happens...err:%+v\n", err))
			return
		}
		// cmd := exec.Command(cmds[0], cmds[1:]...)
		cmd.Env = s.Environ()
		// log.Printf("DEBUG executing cmd with env: %+v\n", cmd.Env)
		cmd.Stdin = s
		cmd.Stdout = s
		cmd.Stderr = s.Stderr()
		if err := cmd.Run(); err != nil {
			log.Println(err)
			return
		}
	})

	passwordAuthOption := ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		for _, userInfo := range users {
			if userInfo.Username == ctx.User() && userInfo.Password == password {
				return true
			}
		}
		return false
	})

	err := ssh.ListenAndServe(":2222", nil, passwordAuthOption)
	log.Printf("sshd exit. %+v\n", err)
}

# build stage
FROM flexconstructor/docker-centos-golang

ENV GO111MODULE=on
RUN go get golang.org/dl/go1.14.2
RUN go1.14.2 download

# GOOS=linux GOARCH=amd64 go1.14.2 build -buildmode=c-archive -o token.a verify_token.go
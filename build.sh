CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.env=prod' -X 'main.version=`cat VERSION`' -X 'main.gitCommit=`git show -s --format=%H`' -X 'main.buildDate=`date`' -X 'main.goVersion=`go version`'"

#CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.env=prod' -X 'main.version=`cat VERSION`' -X 'main.gitCommit=`git show -s --format=%H`' -X 'main.buildDate=`date`' -X 'main.goVersion=`go version`'"
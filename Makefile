
BUILD_DATE=`date +%Y%m%d%H%M%S`
BUILD_COMMIT=`git rev-parse --short HEAD`

all:
	go build -ldflags "-X github.com/newtonproject/newchain-blind/cli.buildCommit=${BUILD_COMMIT}\
	    -X github.com/newtonproject/newchain-blind/cli.buildDate=${BUILD_DATE}"
install:
	go install -ldflags "-X github.com/newtonproject/newchain-blind/cli.buildCommit=${BUILD_COMMIT}\
	    -X github.com/newtonproject/newchain-blind/cli.buildDate=${BUILD_DATE}"

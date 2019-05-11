
BUILD_DATE=`date +%Y%m%d%H%M%S`
BUILD_COMMIT=`git rev-parse --short HEAD`

all:
	go build -ldflags "-X gitlab.newtonproject.org/yangchenzhong/NewChainBlind/cli.buildCommit=${BUILD_COMMIT}\
	    -X gitlab.newtonproject.org/yangchenzhong/NewChainBlind/cli.buildDate=${BUILD_DATE}"


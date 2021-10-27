BUILD_TIME=$(date +"%Y%m%d.%H%M%S")
CommitHash=$(git describe --tags | cut -d- -f3)
GoVersion=$(go version | cut -c 14- | cut -d' ' -f1)
GitTag=$(git describe --tags | cut -d- -f1)
FLAG="-X $TRG_PKG.BuildTime=$BUILD_TIME"
FLAG="$FLAG -X $TRG_PKG.CommitHash=$CommitHash"
FLAG="$FLAG -X $TRG_PKG.GoVersion=$GoVersion"
FLAG="$FLAG -X $TRG_PKG.GitTag=$GitTag"
go build -v -ldflags "$FLAG"
go install
pushd cmd/gctanalyze
go build -v -ldflags "$FLAG"
go install
popd

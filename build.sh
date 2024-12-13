export LDFLAGS='-s -w '

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kscan_linux_amd64 kscan.go
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="$LDFLAGS" -trimpath -o kscan_windows_386.exe  kscan.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kscan_windows_amd64.exe  kscan.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="$LDFLAGS" -trimpath -o kscan_windows_arm64.exe  kscan.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kscan_darwin_amd64 kscan.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$LDFLAGS" -trimpath -o kscan_darwin_arm64 kscan.go
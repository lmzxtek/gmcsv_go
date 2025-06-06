# build.sh（Linux/macOS）
#!/bin/bash

platforms=(
    "windows/amd64"
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

for platform in "${platforms[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    output_name="bin/gmcsv-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    go build -o $output_name main.go
done
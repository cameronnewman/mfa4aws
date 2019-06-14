
#!/usr/bin/env bash

BINARY=$1
if [[ -z "$BINARY" ]]; then
  echo "usage: $0 <BINARY> <VERSION>"
  exit 1
fi

VERSION=$2
if [[ -z "$VERSION" ]]; then
  echo "usage: $0 <BINARY> <VERSION>"
  exit 1
fi

platforms=("linux/amd64" "darwin/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    OUTPUT_NAME=$BINARY'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        OUTPUT_NAME+='.exe'
    fi  
    echo "Building for $GOOS $GOARCH"
    env GOOS=$GOOS GOARCH=$GOARCH go build -x -ldflags "-X main.version=$VERSION" -o bin/$OUTPUT_NAME cmd/$BINARY/main.go;
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
    chmod +x bin/$OUTPUT_NAME
done

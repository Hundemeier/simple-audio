#!/bin/bash
#make the assetFS filesystem that includes everything in the webgui/build folder into the executable
go-bindata-assetfs webgui/dist/...

GREEN='\033[0;32m'
NC='\033[0m' # No Color
#TICK='\033[0;32m \xE2\x9C\x93 \033[0m'
TICK=$GREEN'\xE2\x9C\x93'$NC

# build webgui
#cd ./webgui
#npm run build
#cd ..

# build for following architectures:
printf 'building for...\n'
printf 'Linux AMD64'
if GOOS=linux GOARCH=amd64 go build -o out/simpleAudio_linux_amd64; then
  printf ' '$TICK'\n' $TICK
fi

printf 'Linux 386'
if GOOS=linux GOARCH=386 go build -o out/simpleAudio_linux_386; then
  printf ' '$TICK'\n' $TICK
fi

printf 'Linux ARM'
if GOOS=linux GOARCH=arm go build -o out/simpleAudio_linux_arm; then
  printf ' '$TICK'\n' $TICK
fi

printf 'Windows AMD64'
if GOOS=windows GOARCH=amd64 go build -o out/simpleAudio_win_amd64.exe; then
  printf ' '$TICK'\n' $TICK
fi

printf 'Windows 386'
if GOOS=windows GOARCH=386 go build -o out/simpleAudio_win_386.exe; then
  printf ' '$TICK'\n' $TICK
fi
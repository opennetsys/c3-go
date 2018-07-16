#!/bin/bash

#nc="\\033[0m"
red="\\033[31m"
green="\\033[32m"
yellow="\\033[33m"
#blue="\\033[34m"
#purple="\\033[35m"
cyan="\\033[36m"
white="\\033[37m"
bold="$(tput bold)"
normal="$(tput sgr0)"

# function to install `jq`
install_jq() {
  if [[ "$OSTYPE" == "linux-gnu" ]]; then
      sudo apt-get install jq
  elif [[ "$OSTYPE" == "darwin"* ]]; then
      brew install jq
  fi
}

# check if `gorep` is installed, otherwise install it.
command -v gorep >/dev/null
exit_code=$?
if [ "$exit_code" -ne 0 ]; then
  printf "$yellow%s$normal\\n" "gorep is not installed. Installing..."
  go get github.com/novalagung/gorep && printf "$green%s$normal\\n" "gorep installed"
fi

# check if `jq` is installed, otherwise install it.
command -v jq >/dev/null
exit_code=$?
if [ "$exit_code" -ne 0 ]; then
  printf "$yellow%s$normal\\n" "jq is not installed. Installing..."
  install_jq && printf "$green%s$normal\\n" "jq installed"
fi

# the path to search passed an the only argument
dirpath="$1"

if [[ "$dirpath" == "" ]]; then
  printf "$red$bold%s$normal\\n" "Error: path is required"
  exit 1
fi

printf "$yellow%s$normal\\n" "finding imports recursively under $dirpath"

filter="sed -r 's|/[^/]+$||'"
if [[ "$OSTYPE" == "darwin"* ]]; then
  filter="sed -E 's|/[^/]+$||'"
fi

# read all Go imports that contain "gx/ipfs/"
for dir in $(find "$dirpath" -maxdepth 100 -type f -name '*.go' | eval "$filter" | sort -u)
do
  for line in $(go list -json "./$dir" | jq '.Imports' | grep 'gx/ipfs/' | sed -e 's/gx\///g' | sed -e 's/"//g' | sed -e 's/,//g' | sed -e 's/ //g')
  do
    # fetch the gx package.json and read the github url
    root="$(echo "$line" | tr "/" "\\n" | awk 'FNR <= 3 {print}' | tr '\n' '/' | sed -e s'/.$//g')"
    pkg="$(echo "$line" | tr "/" "\\n" | awk 'FNR > 3 {print}' | tr '\n' '/' | sed -e s'/.$//g')"
    jsonurl="https://gateway.ipfs.io/$root/package.json"
    printf "$white%s$normal\\n" "fetching $jsonurl"
    new="$(curl -s "$jsonurl" | jq '.gx.dvcsimport' | sed -e 's/"//g')"
    if [ "$pkg" != "" ]; then
      new="$new/$pkg"
    fi
    old="gx/$line"

    printf "$cyan%s => %s$normal\\n" "$old" "$new"

    # replace the imports
    gorep -path="$dir" \
        -from="$old" \
        -to="$new"

  done
done

printf "$green%s$normal\\n" "complete"


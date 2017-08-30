#!/bin/bash

set +e

version=$(git describe --abbrev=0 --tags)

which github-release || echo 'please install the tool github-releases'

github-release info -u remijouannet -r terraform-provider-osc -t $version || echo "the release doesn't exist"

cd pkg/

rm -f *.zip
 
ls | while read binary
do
    echo "zipping $binary"
    zip $binary.zip $binary
    echo "upload $binary"
    github-release upload \
        -u remijouannet \
        --name "$binary.zip" \
        -r terraform-provider-osc \
        -f "$binary.zip" \
        --replace \
        -t $version || echo "failed to upload $binary"
done

#!/bin/sh
VERSION=`cat VERSION`
cat << EOF > release.json
{
  "version": "$VERSION",
  "go_version": "1.19.0"
}
EOF

cat << EOF > src/version.go
package src

// Version : Version, the Build Number is parsed in the main func
const Version = "$VERSION"
EOF

sed -Ei "s/ARG XTEVE_VERSION.*/ARG XTEVE_VERSION=$VERSION/" Dockerfile

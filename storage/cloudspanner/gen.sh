#! /bin/sh
cat > spanner.sdl.go <<EOF
// Code generated by storage/cloudspanner/gen.sh DO NOT EDIT.

package cloudspanner

const base64DDL = \`
EOF
if [[ "$OSTYPE" == "darwin"* ]]; then
base64 -b 76 < spanner.sdl >> spanner.sdl.go
else
base64 -w 76 < spanner.sdl >> spanner.sdl.go
fi
echo '`' >> spanner.sdl.go


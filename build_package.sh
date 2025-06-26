#!/bin/bash

APP_NAME=nexa-mqtt

ARCHS="amd64 arm"
LDFLAGS="-s -w"

GITCOMMIT=$(git rev-parse HEAD)
GITVERSION=$(git describe --tags --long)
# replace v1.2.3-4-gxxxxx with 1.2.3.4 or v1.2-3-gxx with 1.2.3
VERSION=$(echo $GITVERSION | sed -E 's/v([0-9]+\.[0-9]+\.?[0-9]*)-([0-9]+)-g.*/\1.\2/')

BUILD_DIR=$(pwd)/build
DEB_DIR=$BUILD_DIR/deb

for arch in $ARCHS; do
		mkdir -p $DEB_DIR/usr/bin;

		echo "Building for $arch...";
		GOOS=linux GOARCH=$arch go build -o $DEB_DIR/usr/bin/${APP_NAME} -ldflags "$LDFLAGS -X main.version=$GITVERSION -X main.commit=$GITCOMMIT" cmd/nexa-mqtt/main.go;

		if [ "$arch" = "arm" ]; then
			deb_arch="armhf";
		else
			deb_arch="$arch";
		fi;
		echo "Creating DEB package for $arch (DEB arch: $deb_arch)...";
		mkdir -p $DEB_DIR/DEBIAN;
		cp -r package/* $DEB_DIR/;
		echo "Version: $VERSION" >> $DEB_DIR/DEBIAN/control;
		echo "Architecture: $deb_arch" >> $DEB_DIR/DEBIAN/control;
		chmod 755 $DEB_DIR/DEBIAN/config;
		chmod 755 $DEB_DIR/DEBIAN/postinst;
		chmod 755 $DEB_DIR/DEBIAN/prerm;
		chmod 755 $DEB_DIR/DEBIAN/postrm;
		echo "Creating $BUILD_DIR/${APP_NAME}_${VERSION}_$arch.deb";
		fakeroot dpkg-deb --build $DEB_DIR $BUILD_DIR/${APP_NAME}_${VERSION}_$arch.deb;
		rm -rf $DEB_DIR
done

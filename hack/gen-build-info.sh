# This script is VERY important. This generates build info and places them into *.go files
# into the source code for every build.
#
# When Syringe is running on a particular release, this is placed within SYRINGE_VERSION
# and made available via the "buildVersion" key in the map shown below. This key is relied
# upon by code within Syringe to determine the corresponding version of infrastructure Docker
# images to call upon for things like configuration and jupyter notebooks.
#
# Without a way to couple the version of Syringe with the version of these important platform-related
# docker image, these images might not work properly, resulting in the whole platform not working properly.
#
# Edit this script, and anything that calls it, with caution.


SYRINGE_SHA=$(git rev-parse HEAD)
SYRINGE_VERSION=$(git describe --exact-match --tags $(git log -n1 --pretty='%h') 2> /dev/null)

if [ -z $SYRINGE_VERSION ];
then
	SYRINGE_VERSION="dev-$SYRINGE_SHA"
fi

cat <<EOT >> ./cmd/syringed/buildinfo.go
package main

/*
This file is automatically generated by Syringe build scripts.
Please do not edit.
*/

var (
	buildInfo = map[string]string{
		"buildSha": "$SYRINGE_SHA",
		"buildVersion": "$SYRINGE_VERSION",
	}
)

EOT

cat <<EOT >> ./cmd/syrctl/buildinfo.go
package main

/*
This file is automatically generated by Syringe build scripts.
Please do not edit.
*/

var (
	buildInfo = map[string]string{
		"buildSha": "$SYRINGE_SHA",
		"buildVersion": "$SYRINGE_VERSION",
	}
)

EOT

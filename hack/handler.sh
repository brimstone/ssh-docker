#!/bin/sh -eu
#set -x

cleanup() {
	echo "Cleaning up"
	docker cp "$CID:/$REPO" . || exit 0
	docker rm -vf "$CID"
	(
	cd "$REPO"
	git status
	)
	rm -rf "$REPO"
}
trap cleanup EXIT

if [ -e /tmp/auth.log ]; then
	cat /tmp/auth.log || true
fi

export PATH=$PATH:/usr/local/bin

USERNAME=${USERNAME:-}

if [ -z "$USERNAME" ]; then
	echo
	echo "$USER is now setup, login with project_name@sprinkle.cloud!"
	echo
	exit 0
fi

echo "USERNAME: $USERNAME"
ORG="${USER%%/*}"
if [ -z "$ORG" ]; then
	ORG="$USERNAME"
fi
REPO="${USER##*/}"

cd /workdir
mkdir -p "$USERNAME"
cd "$USERNAME"
# TODO detect already cloned repo
git clone "https://github.com/$ORG/$REPO" </dev/null || exit 1
# TODO stuff

set -x
CID="$(docker create -t busybox)"
docker cp "$REPO" "$CID":/
#export CID
#/bin/sh

docker start -ai "$CID"

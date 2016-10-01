#!/bin/sh -eu
#set -x

if [ -e /tmp/auth.log ]; then
	cat /tmp/auth.log || true
fi

export PATH=$PATH:/usr/local/bin


if [ -z "${USERNAME:-}" ]; then

	if [ -z "${SETUP:-}" ]; then
		echo "User github.com/$USER doesn't exist or doesn't have ssh keys registered."
		exit 1
	fi

	echo
	echo "$USER is now setup, login with repo@sprinkle.cloud!"
	echo
	exit 0
fi

deleterepo() {
echo "== Deleting fork"
curl -XDELETE \
-u "brimbot:$GITHUB_TOKEN" \
"https://api.github.com/repos/brimbot/$REPO" \
-s >/dev/null

}

commit() {

# Fork repo
echo "== Forking repo"
curl -XPOST \
-u "brimbot:$GITHUB_TOKEN" \
"https://api.github.com/repos/$USERNAME/$REPO/forks" \
-s >/dev/null

echo "== Committing changes"
# save git things
git checkout -b ssh-docker
git add .
git commit -m "Work from Sprinkle Cloud"

git remote add brimbot "git@github.com:brimbot/$REPO"

echo "== Pushing changes back to github"
git push -u brimbot ssh-docker

echo "== Submitting PR"
pr_url="$(curl -XPOST \
-d "{
  \"title\": \"Work from Sprinkle Cloud\",
  \"body\": \"Thanks for using Sprinkle Cloud!\n\nIf you want to save your work, merge this PR.\",
  \"head\": \"brimbot:ssh-docker\",
  \"base\": \"master\"
}" \
-u "brimbot:$GITHUB_TOKEN" \
"https://api.github.com/repos/$USERNAME/$REPO/pulls" \
-s \
| grep -A 1 html | awk -F '"' '/href/ {print $4}')"

echo "URL: $pr_url"

deleterepo

}

cleanup() {
	echo "Cleaning up"
	docker cp "$CID:/$REPO" . || exit 0
	docker rm -vf "$CID" >/dev/null
	(
	cd "$REPO"
	git status
	git diff-files --quiet || commit
	)
	rm -rf "$REPO"
}
trap cleanup EXIT

if [ "${USER}" != "${USER%%/*}" ]; then
	echo "Only repos under https://github.com/$USERNAME can be closed."
	exit 1
fi
REPO="$USER"

mkdir -p /workdir
cd /workdir
mkdir -p "$USERNAME"
cd "$USERNAME"
# TODO detect already cloned repo
git clone "https://github.com/$USERNAME/$REPO" </dev/null || exit 1
# TODO stuff

CID="$(docker create -it -m 64m --cpu-shares 100 busybox)"
docker cp "$REPO" "$CID":/
#export CID
#/bin/sh

docker start -ai "$CID"

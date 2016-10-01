#!/bin/sh -eu

mkdir -p users
cd users

#set -x
#exec 2>/tmp/auth.log

USERNAME="$(grep -lr "$2" -- * || true)"
if [ -n "$USERNAME" ]; then
	echo "USERNAME=$USERNAME"
	exit 0
fi

if curl "https://api.github.com/users/$1/keys" -s \
| awk -F '"' '$2 == "key" {print $4}' \
| grep -q "$2"; then
	echo "SETUP=true"
	echo "$2" >> "$1"
	exit 0
fi

exit 0

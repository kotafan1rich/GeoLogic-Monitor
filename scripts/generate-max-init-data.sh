#!/usr/bin/env bash

set -euo pipefail

readonly max_int64="9223372036854775807"

urlencode() {
	local value=$1
	local encoded=""
	local char
	local index

	for ((index = 0; index < ${#value}; index++)); do
		char=${value:index:1}
		case "$char" in
			[a-zA-Z0-9.~_-]) encoded+=$char ;;
			*) printf -v encoded '%s%%%02X' "$encoded" "'$char" ;;
		esac
	done

	printf '%s' "$encoded"
}

if ! command -v openssl >/dev/null 2>&1; then
	printf 'openssl is required\n' >&2
	exit 1
fi

IFS= read -rsp 'MAX bot token: ' bot_token
printf '\n'
if [[ -z $bot_token ]]; then
	printf 'MAX bot token must not be empty\n' >&2
	exit 1
fi

IFS= read -rp 'MAX user ID: ' raw_user_id
if [[ ! $raw_user_id =~ ^[0-9]+$ ]]; then
	printf 'MAX user ID must be an integer\n' >&2
	exit 1
fi

user_id=$raw_user_id
while [[ $user_id == 0* && ${#user_id} -gt 1 ]]; do
	user_id=${user_id#0}
done
if [[ $user_id == 0 ]] ||
	(( ${#user_id} > ${#max_int64} )) ||
	[[ ${#user_id} -eq ${#max_int64} && $user_id > $max_int64 ]]; then
	printf 'MAX user ID must be a positive int64\n' >&2
	exit 1
fi

auth_date=$(date +%s)
user_json=$(printf '{"id":%s}' "$user_id")
launch_params=$(printf 'auth_date=%s\nuser=%s' "$auth_date" "$user_json")

secret_key=$(
	printf '%s' "$bot_token" |
		openssl dgst -sha256 -hmac 'WebAppData' -hex |
		sed 's/^.*= //'
)
hash=$(
	printf '%s' "$launch_params" |
		openssl dgst -sha256 -mac HMAC -macopt "hexkey:$secret_key" -hex |
		sed 's/^.*= //'
)

printf 'auth_date=%s&user=%s&hash=%s\n' \
	"$(urlencode "$auth_date")" \
	"$(urlencode "$user_json")" \
	"$hash"

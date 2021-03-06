#!/bin/sh

HOST='127.0.0.1:8080'
INITIAL_TOKEN_FILE="./initial_token.txt"

if [[ ! -f "${INITIAL_TOKEN_FILE}" ]]; then
  echo "[E] ${INITIAL_TOKEN_FILE} not found"
  exit 1
fi
INITIAL_TOKEN="$(cat ${INITIAL_TOKEN_FILE})"

set -e

function run_curl {
  METHOD="${1}"
  API_PATH="${2}"
  DATA="${3}"

  TOKEN_STRING=""
  if [[ "${METHOD}" == "POST" ]] && [[ ${API_PATH} == "/v1/owner" ]]; then
    TOKEN_STRING="X-Admin-Token:${INITIAL_TOKEN}"
  else
    TOKEN_STRING="X-User-Token:${INITIAL_TOKEN}"
  fi

  if [[ ! -z "${DATA}" ]]; then
    echo "send: ${DATA}"
    curl -X${METHOD} -s -w "${METHOD} ${API_PATH} -> %{http_code}\n-------\n" \
      -d "${DATA}" \
      -H 'Content-Type:application/json' \
      -H "${TOKEN_STRING}" \
      http://${HOST}${API_PATH}
  else
    curl -X${METHOD} -s -w "${METHOD} ${API_PATH} -> %{http_code}\n-------\n" \
      -H "${TOKEN_STRING}" \
      http://${HOST}${API_PATH}
  fi
}

function test_owner {
  echo "======> [$(date) Starting test of /v1/owner"
  U='/v1/owner'
  run_curl GET ${U}
  run_curl GET ${U}/test
  run_curl POST ${U} '{"username":"test","fullname":"test","email":"abcd@abcd.xy"}'
  run_curl GET ${U}
  run_curl GET ${U}/test
  run_curl PUT ${U} '{"username":"test","fullname":"upd","email":"meh@meh.meh"}'
  run_curl GET ${U}/test
  run_curl DELETE ${U}/test
  run_curl GET ${U}/test
  echo "======> [$(date) Finished test of /v1/owner"
}

function test_asnum {
  echo "======> [$(date) Starting test of /v1/asnum"
  U='/v1/asnum'
  run_curl GET ${U}
  run_curl GET ${U}/64512
  run_curl POST /v1/owner '{"username":"test","fullname":"test","email":"abcd@abcd.xy"}'
  run_curl POST ${U} '{"asnum":64512,"description":"test as","username":"test"}'
  run_curl GET ${U}
  run_curl GET ${U}/64512
  run_curl PUT ${U} '{"asnum":64512,"description":"updated as","username":"test"}'
  run_curl GET ${U}/64512
  run_curl DELETE ${U}/64512
  run_curl GET ${U}/64512
  run_curl DELETE /v1/owner/test
  echo "======> [$(date) Finished test of /v1/asnum"
}

function test_prefix {
  echo "======> [$(date) Starting test of /v1/prefix"
  U='/v1/prefix'
  run_curl GET ${U}
  run_curl GET ${U}/10.0.0.0/8
  run_curl POST /v1/owner '{"username":"test","fullname":"test","email":"abcd@abcd.xy"}'
  run_curl POST ${U} '{"network":"0.0.0.0/0","description":"internet","username":"test"}'
  run_curl POST ${U} '{"network":"10.0.0.0/8","description":"rfc1918 10/8","username":"test"}'
  run_curl POST ${U} '{"network":"10.42.0.0/19","description":"AS65342-NET","username":"test"}'
  run_curl GET ${U}
  run_curl GET ${U}/10.0.0.0/8
  run_curl PUT ${U} '{"network":"10.0.0.0/8","description":"ipv6 rules","username":"test"}'
  run_curl GET ${U}/10.0.0.0/8
  run_curl DELETE ${U}/10.42.0.0/19
  run_curl DELETE ${U}/10.0.0.0/8
  run_curl DELETE ${U}/0.0.0.0/0
  run_curl GET ${U}/10.0.0.0/8
  run_curl DELETE /v1/owner/test
  echo "======> [$(date) Finished test of /v1/prefix"
}

case "${1}" in
  'owner')
    test_owner ${INITIAL_TOKEN}
    ;;
  'asnum')
    test_asnum ${INITIAL_TOKEN}
    ;;
  'prefix')
    test_prefix ${INITIAL_TOKEN}
    ;;
  *)
    test_owner
    test_asnum
    test_prefix
    ;;
esac

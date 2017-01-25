#!/bin/sh

JSON_DATA='{"username":"test","fullname":"test","email":"abcd@abcd.xy"}'

set -e

function test_owner {
  echo "======> [$(date) Starting test of /v1/owner"
  set -x
  curl -kv http://127.0.0.1:8080/v1/owner
  curl -kv http://127.0.0.1:8080/v1/owner/test
  curl -d '{"username":"test","fullname":"test","email":"abcd@abcd.xy"}' \
    -H 'Content-Type:application/json' -kv http://127.0.0.1:8080/v1/owner
  curl -kv http://127.0.0.1:8080/v1/owner
  curl -kv http://127.0.0.1:8080/v1/owner/test
  curl -XPUT -d '{"username":"test","fullname":"upd","email":"meh@meh.meh"}' \
    -H 'Content-Type:application/json' -kv http://127.0.0.1:8080/v1/owner
  curl -kv http://127.0.0.1:8080/v1/owner/test
  curl -XDELETE -kv http://127.0.0.1:8080/v1/owner/test
  curl -kv http://127.0.0.1:8080/v1/owner/test
  set +x
  echo "======> [$(date) Finished test of /v1/owner"
}

function test_asnum {
  echo "======> [$(date) Starting test of /v1/asnum"
  set -x
  curl -kv http://127.0.0.1:8080/v1/asnum
  curl -kv http://127.0.0.1:8080/v1/asnum/64512
  curl -d '{"username":"test","fullname":"test","email":"abcd@abcd.xy"}' \
    -H 'Content-Type:application/json' -kv http://127.0.0.1:8080/v1/owner
  curl -d '{"asnum":64512,"description":"test as","username":"test"}' \
    -H 'Content-Type:application/json' -kv http://127.0.0.1:8080/v1/asnum
  curl -kv http://127.0.0.1:8080/v1/asnum
  curl -kv http://127.0.0.1:8080/v1/asnum/64512
  curl -XPUT -d '{"asnum":64512,"description":"updated as","username":"test"}' \
    -H 'Content-Type:application/json' -kv http://127.0.0.1:8080/v1/asnum
  curl -kv http://127.0.0.1:8080/v1/asnum/64512
  curl -XDELETE -kv http://127.0.0.1:8080/v1/asnum/64512
  curl -kv http://127.0.0.1:8080/v1/asnum/64512
  curl -XDELETE -kv http://127.0.0.1:8080/v1/owner/test
  set +x
  echo "======> [$(date) Finished test of /v1/asnum"
}

case "${1}" in
  'owner')
    test_owner
    ;;
  'asnum':)
    test_asnum
    ;;
  *)
    test_owner
    test_asnum
    ;;
esac

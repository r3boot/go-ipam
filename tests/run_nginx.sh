#!/bin/sh

echo "[+] Starting nginx"
nginx -p ${PWD} -c ./tests/nginx/nginx.conf

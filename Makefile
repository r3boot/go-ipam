SERVER = nic-server
WEBAPP = nic-app

BUILD_DIR = ./build
TEST_DIR = ./tests

NGINX_DIR = ${TEST_DIR}/nginx
NGINX_TMP_DIR = ${NGINX_DIR}/tmp

INITIAL_TOKEN = ./initial_token.txt

all: validate generate ${SERVER} nginx webapp

initial_token:
	uuidgen -r > ${INITIAL_TOKEN}

validate:
	swagger validate go-ipam.yaml

generate: initial_token
	swagger generate server -f go-ipam.yaml

${SERVER}:
	[[ -d "${BUILD_DIR}" ]] || mkdir -p "${BUILD_DIR}"
	go build -v -o "${BUILD_DIR}"/${SERVER} cmd/${SERVER}/main.go

nginx:
	mkdir -p ${NGINX_TMP_DIR}

webapp:
	cd ${WEBAPP}; ember build --environment=production

clean:
	grep restapi .gitignore | xargs rm -rf
	rm -rf ${BUILD_DIR} ${NGINX_TMP_DIR} client cmd models || true
	rm -rf ${WEBAPP}/dist || true
	rm -f ${INITIAL_TOKEN} || true

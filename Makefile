SERVER = nic-server
CLIENT = nic-client

BUILD_DIR = ./build
INITIAL_TOKEN = ./initial_token.txt

all: validate generate ${SERVER}

initial_token:
	uuidgen -r > ${INITIAL_TOKEN}

validate:
	swagger validate go-ipam.yaml

generate: initial_token
	swagger generate server -f go-ipam.yaml
	#swagger generate client -f go-ipam.yaml

${SERVER}:
	[[ -d "${BUILD_DIR}" ]] || mkdir -p "${BUILD_DIR}"
	go build -v -o "${BUILD_DIR}"/${SERVER} cmd/${SERVER}/main.go

${CLIENT}:
	[[ -d "${BUILD_DIR}" ]] || mkdir -p "${BUILD_DIR}"
	go build -v -o "${BUILD_DIR}"/${CLIENT} client/nic_client.go

clean:
	grep restapi .gitignore | xargs rm -rf
	rm -rf build build client cmd models || true
	rm -f ${INITIAL_TOKEN} || true

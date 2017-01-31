package restapi

import (
	"crypto/tls"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	interpose "github.com/carbocation/interpose/middleware"
	recover "github.com/dre1080/recover"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/r3boot/go-ipam/restapi/operations"

	"github.com/r3boot/go-ipam/models"
	"github.com/r3boot/go-ipam/storage"
	"github.com/r3boot/go-ipam/storage/email"
	"github.com/r3boot/go-ipam/storage/postgres"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name  --spec ../go-ipam.yaml

func configureFlags(api *operations.NicAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.NicAPI) http.Handler {
	var (
		signup_enabled bool
		backend        *storage.Storage
		err            error
		activationQ    chan email.ActivationQItem
	)
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Initialize storage backend
	backend = storage.Setup(postgres.Config{
		Host:     "postgresql.service.local:5432",
		User:     "nic_user",
		Pass:     "gLJxBsf7wyZL2",
		Database: "nic",
	}, email.Config{
		Smarthost:   "smtp.as65342.net:25",
		Sender:      "do-not-reply@as65342.net",
		SenderName:  "IPAM network management",
		NetworkName: "Test Network",
	})

	signup_enabled = true

	activationQ = make(chan email.ActivationQItem)
	go email.SignupEmailRoutine(activationQ)

	fs, err := os.Stat("initial_token.txt")
	if err != nil {
		panic("stat(): initial_token.txt not found")
	}

	fd, err := os.Open("initial_token.txt")
	if err != nil {
		panic("open(): failed to open initial_token.txt: " + err.Error())
	}

	data := make([]byte, fs.Size())
	_, err = fd.Read(data)
	if err != nil {
		panic("read(): failed to read initial_token.txt: " + err.Error())
	}

	initial_token := strings.Trim(string(data), "\n")

	// Applies when the "X-Admin-Token" header is set
	api.AdminSecurityAuth = func(token string) (interface{}, error) {
		// Check if the token matches a known Owner token
		owner := backend.GetOwnerByApiToken(token)
		if owner.Username != nil {
			log.Print("Admin request authenticated via token from " + *owner.Username)
			return owner, nil
		}

		// Check if token matches initial_token
		if token == initial_token {
			log.Print("Admin request authenticated via initial_token")
			return true, nil
		}
		return nil, errors.Unauthenticated("admin level clearance")
	}

	// Applies when the "X-User-Token" header is set
	api.UserSecurityAuth = func(token string) (interface{}, error) {
		// Check if the token matches a known Owner token
		owner := backend.GetOwnerByApiToken(token)
		if owner.Username != nil {
			log.Print("User request authenticated via token from " + *owner.Username)
			return owner, nil
		}

		// Check if token matches initial_token
		if token == initial_token {
			log.Print("User request authenticated via initial_token")
			return true, nil
		}
		return nil, errors.Unauthenticated("user level clearance")
	}

	/*
	 * Handlers for /v1/signup
	 */
	api.PostSignupHandler = operations.PostSignupHandlerFunc(func(params operations.PostSignupParams) middleware.Responder {
		var (
			salt   string
			hash   []byte
			hash_s string
			ts     string
		)

		if !signup_enabled {
			return operations.NewPostSignupNotFound()
		}

		if params.Owner.Password == "" {
			log.Print("Signup: No Password specified")
			return operations.NewPostSignupBadRequest()
		}

		ts = time.Now().Format(time.RFC3339)

		salt, err = backend.GenerateSalt(32)
		if err != nil {
			log.Print("Signup: " + err.Error())
			return operations.NewPostSignupBadRequest()
		}

		hash, err = backend.GenerateHash(params.Owner.Password, salt)
		if err != nil {
			log.Print("Signup: " + err.Error())
			return operations.NewPostSignupBadRequest()
		}
		hash_s = base64.URLEncoding.EncodeToString(hash)

		owner := models.Owner{
			Username:       params.Owner.Username,
			Password:       hash_s,
			Salt:           salt,
			IsActive:       false,
			IsAdmin:        false,
			Fullname:       params.Owner.Fullname,
			Email:          params.Owner.Email,
			SignupTime:     ts,
			ActivationTime: "",
			LastLogin:      "",
			LastLoginHost:  "",
		}

		if err = backend.RunSignup(activationQ, owner); err != nil {
			return operations.NewPostSignupBadRequest()
		}

		return operations.NewPostSignupNoContent()
	})

	/*
	 * Handlers for /v1/owner
	 */
	api.DeleteOwnerUsernameHandler = operations.DeleteOwnerUsernameHandlerFunc(func(params operations.DeleteOwnerUsernameParams, principal interface{}) middleware.Responder {
		if !backend.HasOwner(params.Username) {
			return operations.NewDeleteOwnerUsernameNotFound()
		}

		if err = backend.DeleteOwner(params.Username); err != nil {
			return operations.NewDeleteOwnerUsernameInternalServerError()
		}

		return operations.NewDeleteOwnerUsernameNoContent()
	})

	api.GetOwnerHandler = operations.GetOwnerHandlerFunc(func(params operations.GetOwnerParams, principal interface{}) middleware.Responder {
		owners := backend.GetOwners()
		if len(owners) == 0 {
			return operations.NewGetOwnerNotFound()
		}

		return operations.NewGetOwnerOK().WithPayload(owners)
	})

	api.GetOwnerUsernameHandler = operations.GetOwnerUsernameHandlerFunc(func(params operations.GetOwnerUsernameParams, principal interface{}) middleware.Responder {
		owner := backend.GetOwner(params.Username)
		if owner.Username == nil {
			return operations.NewGetOwnerUsernameNotFound()
		}

		return operations.NewGetOwnerUsernameOK().WithPayload(&owner)
	})

	api.PostOwnerHandler = operations.PostOwnerHandlerFunc(func(params operations.PostOwnerParams, principal interface{}) middleware.Responder {
		owner := models.Owner{
			Username: params.Owner.Username,
			Fullname: params.Owner.Fullname,
			Email:    params.Owner.Email,
		}

		if err = backend.AddOwner(owner); err != nil {
			return operations.NewPostOwnerBadRequest()
		}

		return operations.NewPostOwnerNoContent()
	})

	api.PutOwnerHandler = operations.PutOwnerHandlerFunc(func(params operations.PutOwnerParams, principal interface{}) middleware.Responder {
		owner := models.Owner{
			Username: params.Owner.Username,
			Fullname: params.Owner.Fullname,
			Email:    params.Owner.Email,
		}

		if err = backend.UpdateOwner(owner); err != nil {
			return operations.NewPutOwnerBadRequest()
		}

		return operations.NewPutOwnerNoContent()
	})

	/*
	 * Handlers for /v1/asnum
	 */
	api.DeleteAsnumAsnumHandler = operations.DeleteAsnumAsnumHandlerFunc(func(params operations.DeleteAsnumAsnumParams, principal interface{}) middleware.Responder {
		if !backend.HasAsnum(params.Asnum) {
			return operations.NewDeleteAsnumAsnumNotFound()
		}

		if err = backend.DeleteAsnum(params.Asnum); err != nil {
			return operations.NewDeleteAsnumAsnumInternalServerError()
		}

		return operations.NewDeleteAsnumAsnumNoContent()
	})

	api.GetAsnumHandler = operations.GetAsnumHandlerFunc(func(params operations.GetAsnumParams, principal interface{}) middleware.Responder {
		asnums := backend.GetAsnums()
		if len(asnums) == 0 {
			return operations.NewGetAsnumNotFound()
		}

		return operations.NewGetAsnumOK().WithPayload(asnums)
	})

	api.GetAsnumAsnumHandler = operations.GetAsnumAsnumHandlerFunc(func(params operations.GetAsnumAsnumParams, principal interface{}) middleware.Responder {
		asnum := backend.GetAsnum(params.Asnum)
		if asnum.Asnum == nil {
			return operations.NewGetAsnumAsnumNotFound()
		}

		return operations.NewGetAsnumAsnumOK().WithPayload(&asnum)
	})

	api.PostAsnumHandler = operations.PostAsnumHandlerFunc(func(params operations.PostAsnumParams, principal interface{}) middleware.Responder {
		asnum := models.Asnum{
			Asnum:       params.Asnum.Asnum,
			Description: params.Asnum.Description,
			Username:    params.Asnum.Username,
		}

		if err = backend.AddAsnum(asnum); err != nil {
			log.Print("error: " + err.Error())
			return operations.NewPostAsnumBadRequest()
		}

		return operations.NewPostAsnumNoContent()
	})

	api.PutAsnumHandler = operations.PutAsnumHandlerFunc(func(params operations.PutAsnumParams, principal interface{}) middleware.Responder {
		asnum := models.Asnum{
			Asnum:       params.Asnum.Asnum,
			Description: params.Asnum.Description,
			Username:    params.Asnum.Username,
		}

		if err = backend.UpdateAsnum(asnum); err != nil {
			return operations.NewPutAsnumBadRequest()
		}

		return operations.NewPutAsnumNoContent()
	})

	/*
	 * Handlers for /v1/prefix
	 */
	api.DeletePrefixSubnetPrefixlenHandler = operations.DeletePrefixSubnetPrefixlenHandlerFunc(func(params operations.DeletePrefixSubnetPrefixlenParams, principal interface{}) middleware.Responder {
		prefix := params.Subnet + "/" + strconv.Itoa(int(params.Prefixlen))
		if !backend.HasPrefix(prefix) {
			return operations.NewDeletePrefixSubnetPrefixlenNotFound()
		}

		if err = backend.DeletePrefix(prefix); err != nil {
			return operations.NewDeletePrefixSubnetPrefixlenInternalServerError()
		}

		return operations.NewDeletePrefixSubnetPrefixlenNoContent()
	})

	api.GetPrefixHandler = operations.GetPrefixHandlerFunc(func(params operations.GetPrefixParams, principal interface{}) middleware.Responder {
		prefixes := backend.GetPrefixes()
		if len(prefixes) == 0 {
			return operations.NewGetPrefixNotFound()
		}

		return operations.NewGetPrefixOK().WithPayload(prefixes)
	})

	api.GetPrefixSubnetPrefixlenHandler = operations.GetPrefixSubnetPrefixlenHandlerFunc(func(params operations.GetPrefixSubnetPrefixlenParams, principal interface{}) middleware.Responder {
		network := params.Subnet + "/" + strconv.Itoa(int(params.Prefixlen))
		prefix := backend.GetPrefix(network)
		if prefix.Network == nil {
			return operations.NewGetPrefixSubnetPrefixlenNotFound()
		}

		return operations.NewGetPrefixSubnetPrefixlenOK().WithPayload(&prefix)
	})

	api.PostPrefixHandler = operations.PostPrefixHandlerFunc(func(params operations.PostPrefixParams, principal interface{}) middleware.Responder {

		prefix := models.Prefix{
			Network:     params.Prefix.Network,
			Description: params.Prefix.Description,
			Username:    params.Prefix.Username,
		}

		if err := backend.AddPrefix(prefix); err != nil {
			log.Print("error: " + err.Error())
			return operations.NewPostPrefixBadRequest()
		}

		return operations.NewPostPrefixNoContent()
	})

	api.PutPrefixHandler = operations.PutPrefixHandlerFunc(func(params operations.PutPrefixParams, principal interface{}) middleware.Responder {

		prefix := models.Prefix{
			Network:     params.Prefix.Network,
			Description: params.Prefix.Description,
			Username:    params.Prefix.Username,
		}

		if err := backend.UpdatePrefix(prefix); err != nil {
			return operations.NewPutPrefixBadRequest()
		}

		return operations.NewPutPrefixNoContent()
	})

	/*
	 * Handlers end here
	 */

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handlePanic := recover.New(&recover.Options{
		Log: log.Print,
	})

	logViaLogrus := interpose.NegroniLogrus()

	return handlePanic(logViaLogrus(handler))
}

package restapi

import (
	"crypto/tls"
	"log"
	"net/http"

	interpose "github.com/carbocation/interpose/middleware"
	recover "github.com/dre1080/recover"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/r3boot/go-ipam/restapi/operations"

	"github.com/r3boot/go-ipam/models"
	"github.com/r3boot/go-ipam/storage"
	"github.com/r3boot/go-ipam/storage/postgres"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name  --spec ../go-ipam.yaml

func configureFlags(api *operations.NicAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.NicAPI) http.Handler {
	var (
		backend *storage.Storage
		err     error
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
	})

	/*
	 * Handlers for /v1/owner
	 */
	api.DeleteOwnerUsernameHandler = operations.DeleteOwnerUsernameHandlerFunc(func(params operations.DeleteOwnerUsernameParams) middleware.Responder {
		if !backend.HasOwner(params.Username) {
			return operations.NewDeleteOwnerUsernameNotFound()
		}

		if err = backend.DeleteOwner(params.Username); err != nil {
			return operations.NewDeleteOwnerUsernameInternalServerError()
		}

		return operations.NewDeleteOwnerUsernameNoContent()
	})

	api.GetOwnerHandler = operations.GetOwnerHandlerFunc(func(params operations.GetOwnerParams) middleware.Responder {
		owners := backend.GetOwners()
		if len(owners) == 0 {
			return operations.NewGetOwnerNotFound()
		}

		return operations.NewGetOwnerOK().WithPayload(owners)
	})

	api.GetOwnerUsernameHandler = operations.GetOwnerUsernameHandlerFunc(func(params operations.GetOwnerUsernameParams) middleware.Responder {
		owner := backend.GetOwner(params.Username)
		if owner.Username == nil {
			return operations.NewGetOwnerUsernameNotFound()
		}

		return operations.NewGetOwnerUsernameOK().WithPayload(&owner)
	})

	api.PostOwnerHandler = operations.PostOwnerHandlerFunc(func(params operations.PostOwnerParams) middleware.Responder {
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

	api.PutOwnerHandler = operations.PutOwnerHandlerFunc(func(params operations.PutOwnerParams) middleware.Responder {
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
	api.DeleteAsnumAsnumHandler = operations.DeleteAsnumAsnumHandlerFunc(func(params operations.DeleteAsnumAsnumParams) middleware.Responder {
		if !backend.HasAsnum(params.Asnum) {
			return operations.NewDeleteAsnumAsnumNotFound()
		}

		if err = backend.DeleteAsnum(params.Asnum); err != nil {
			return operations.NewDeleteAsnumAsnumInternalServerError()
		}

		return operations.NewDeleteAsnumAsnumNoContent()
	})

	api.GetAsnumHandler = operations.GetAsnumHandlerFunc(func(params operations.GetAsnumParams) middleware.Responder {
		asnums := backend.GetAsnums()
		if len(asnums) == 0 {
			return operations.NewGetAsnumNotFound()
		}

		return operations.NewGetAsnumOK().WithPayload(asnums)
	})

	api.GetAsnumAsnumHandler = operations.GetAsnumAsnumHandlerFunc(func(params operations.GetAsnumAsnumParams) middleware.Responder {
		asnum := backend.GetAsnum(params.Asnum)
		if asnum.Asnum == nil {
			return operations.NewGetAsnumAsnumNotFound()
		}

		return operations.NewGetAsnumAsnumOK().WithPayload(&asnum)
	})

	api.PostAsnumHandler = operations.PostAsnumHandlerFunc(func(params operations.PostAsnumParams) middleware.Responder {
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

	api.PutAsnumHandler = operations.PutAsnumHandlerFunc(func(params operations.PutAsnumParams) middleware.Responder {
		owner := models.Asnum{
			Asnum:       params.Asnum.Asnum,
			Description: params.Asnum.Description,
			Username:    params.Asnum.Username,
		}

		if err = backend.UpdateAsnum(owner); err != nil {
			return operations.NewPutAsnumBadRequest()
		}

		return operations.NewPutAsnumNoContent()
	})

	/*
	 * Handlers for /v1/prefix
	 */
	api.DeletePrefixNetworkHandler = operations.DeletePrefixNetworkHandlerFunc(func(params operations.DeletePrefixNetworkParams) middleware.Responder {
		if !backend.HasPrefix(params.Network) {
			return operations.NewDeletePrefixNetworkNotFound()
		}

		if err = backend.DeletePrefix(params.Network); err != nil {
			return operations.NewDeletePrefixNetworkInternalServerError()
		}

		return operations.NewDeletePrefixNetworkNoContent()
	})

	api.GetPrefixHandler = operations.GetPrefixHandlerFunc(func(params operations.GetPrefixParams) middleware.Responder {
		prefixes := backend.GetPrefixes()
		if len(prefixes) == 0 {
			return operations.NewGetPrefixNotFound()
		}

		return operations.NewGetPrefixOK().WithPayload(prefixes)
	})

	api.GetPrefixNetworkHandler = operations.GetPrefixNetworkHandlerFunc(func(params operations.GetPrefixNetworkParams) middleware.Responder {
		prefix := backend.GetPrefix(params.Network)
		if prefix.Network == nil {
			return operations.NewGetPrefixNetworkNotFound()
		}

		return operations.NewGetPrefixNetworkOK().WithPayload(&prefix)
	})

	api.PostPrefixHandler = operations.PostPrefixHandlerFunc(func(params operations.PostPrefixParams) middleware.Responder {
		prefix := models.Prefix{
			Network:     params.Prefix.Network,
			Description: params.Prefix.Description,
			Username:    params.Prefix.Username,
		}

		if err = backend.AddPrefix(prefix); err != nil {
			log.Print("error: " + err.Error())
			return operations.NewPostPrefixBadRequest()
		}

		return operations.NewPostPrefixNoContent()
	})

	api.PutPrefixHandler = operations.PutPrefixHandlerFunc(func(params operations.PutPrefixParams) middleware.Responder {
		owner := models.Prefix{
			Network:     params.Prefix.Network,
			Description: params.Prefix.Description,
			Username:    params.Prefix.Username,
		}

		if err = backend.UpdatePrefix(owner); err != nil {
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

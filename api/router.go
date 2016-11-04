package api

import (
	"net/http"

	"github.com/anthonynsimon/parrot/datastore"
	"github.com/anthonynsimon/parrot/paths"
	"github.com/pressly/chi"
)

var store datastore.Store
var signingKey []byte

func NewRouter(ds datastore.Store, sk []byte) http.Handler {
	store = ds
	signingKey = sk

	router := chi.NewRouter()
	router.Use(cors)
	registerRoutes(router)

	return router
}

func registerRoutes(router *chi.Mux) {
	router.Get(paths.PingPath, apiHandlerFunc(ping).ServeHTTP)
	router.Post(paths.AuthenticatePath, apiHandlerFunc(authenticate).ServeHTTP)
	router.Post(paths.UsersPath, apiHandlerFunc(createUser).ServeHTTP)

	router.Route(paths.ProjectsPath, func(pr chi.Router) {
		// Past this point, all routes require a valid token
		pr.Use(tokenGate)
		pr.Get("/", apiHandlerFunc(getUserProjects).ServeHTTP)
		pr.Post("/", apiHandlerFunc(createProject).ServeHTTP)
		pr.Get("/:projectID", apiHandlerFunc(showProject).ServeHTTP)
		pr.Put("/:projectID", apiHandlerFunc(updateProject).ServeHTTP)
		pr.Delete("/:projectID", apiHandlerFunc(deleteProject).ServeHTTP)

		pr.Route("/:projectID"+paths.UsersPath, func(dr chi.Router) {
			dr.Get("/", apiHandlerFunc(getProjectUsers).ServeHTTP)
			dr.Post("/", apiHandlerFunc(assignProjectUser).ServeHTTP)
		})

		pr.Route("/:projectID"+paths.LocalesPath, func(dr chi.Router) {
			dr.Post("/", apiHandlerFunc(createLocale).ServeHTTP)
			dr.Get("/", apiHandlerFunc(findLocales).ServeHTTP)
			dr.Get("/:localeID", apiHandlerFunc(showLocale).ServeHTTP)
			dr.Put("/:localeID", apiHandlerFunc(updateLocale).ServeHTTP)
			dr.Delete("/:localeID", apiHandlerFunc(deleteLocale).ServeHTTP)
		})
	})
}

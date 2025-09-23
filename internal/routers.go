package internal

import (
	"database/sql"

	"github.com/Peter-Tabarani/PiconexBackend/internal/routes"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	routes.RegisterPersonRoutes(router, db)
	routes.RegisterStudentRoutes(router, db)
	routes.RegisterAdminRoutes(router, db)
	routes.RegisterActivityRoutes(router, db)
	routes.RegisterDocumentationRoutes(router, db)
	routes.RegisterPersonalDocumentationRoutes(router, db)
	routes.RegisterSpecificDocumentationRoutes(router, db)
	routes.RegisterPointOfContactRoutes(router, db)
	routes.RegisterDisabilityRoutes(router, db)
	routes.RegisterAccommodationRoutes(router, db)
	routes.RegisterRelationshipRoutes(router, db)
	routes.RegisterAuthRoutes(router, db)

	return router
}

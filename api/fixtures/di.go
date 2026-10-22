package fixtures

import (
	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/di"
	"github.com/Rabi-IT/rabi-food-core/libs/http"

	"github.com/samber/do"
)

// Built by Api.Start once the containers are up: the GoTrue, Stripe and HTTP adapters copy config
// values on creation (the database adapter keeps a pointer and only reads it on Connect).
var (
	testDB         *database.PgxAdapter
	testHTTPServer http.HTTPServer
)

func buildInjector() {
	injector := di.Injector()
	testDB = do.MustInvoke[*database.PgxAdapter](injector)
	testHTTPServer = do.MustInvoke[http.HTTPServer](injector)
}

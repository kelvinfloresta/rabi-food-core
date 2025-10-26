package fixtures

import (
	"rabi-food-core/libs/database/gorm_adapter"
	"rabi-food-core/libs/di"
	"rabi-food-core/libs/http"

	"github.com/samber/do"
)

var (
	testInjector   = di.NewTest()
	testDB         = do.MustInvoke[*gorm_adapter.GormAdapter](testInjector)
	testHTTPServer = do.MustInvoke[http.HTTPServer](testInjector)
)

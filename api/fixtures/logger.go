package fixtures

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/stretchr/testify/require"
)

type loggerFixture struct{}

var Logger = loggerFixture{}

// WideEvent logs we with the production logger and returns the fields of the emitted line, so tests
// check the log as it is shipped (field names included).
func (*loggerFixture) WideEvent(t *testing.T, we *logger.WideEvent) map[string]any {
	t.Helper()

	var buf bytes.Buffer
	l := logger.NewWithOutput(&buf, true)
	l.Info().EmbedObject(we).Send()

	fields := map[string]any{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &fields))

	return fields
}

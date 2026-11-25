package fixtures

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/stretchr/testify/require"
)

type loggerFixture struct{}

var Logger = loggerFixture{}

// capturedLogs keeps every log line the test API emits (see captureLogs), so tests check the log as
// production writes it.
var capturedLogs = &logBuffer{}

type logBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *logBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.Write(p)
}

func (b *logBuffer) lines() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()

	return bytes.Clone(b.buf.Bytes())
}

// captureLogs keeps the API logs for Logger.Request while still printing them.
func captureLogs() {
	logger.SetOutput(io.MultiWriter(os.Stdout, capturedLogs))
}

// Request returns the fields of the log line the API wrote for the request with this ID (the
// X-Request-ID response header).
func (*loggerFixture) Request(t *testing.T, requestID string) map[string]any {
	t.Helper()

	scanner := bufio.NewScanner(bytes.NewReader(capturedLogs.lines()))
	for scanner.Scan() {
		fields := map[string]any{}
		if json.Unmarshal(scanner.Bytes(), &fields) == nil && fields["request_id"] == requestID {
			return fields
		}
	}

	require.Failf(t, "request not logged", "no log line with request_id %q", requestID)

	return nil
}

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

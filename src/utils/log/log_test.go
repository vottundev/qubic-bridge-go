package log_test

import (
	"testing"

	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

func TestDebug(t *testing.T) {
	kk := "All is fine"
	log.Debugf("Information message: %s", kk)
	log.Debugln("Debug ln")
}

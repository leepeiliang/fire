package configmap

import (
	"fire/pkg/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNeg(t *testing.T) {
	var devices map[string]*common.BaseDevice

	assert.NotNil(t, NewParse("/Users/lipeiliang/go/src/edgeplanet/fire/config/deviceProfile_test.json", devices))
}

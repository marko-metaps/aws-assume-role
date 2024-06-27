package testutils

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/sts"
)

type MockSTSAPI struct {
	Resp sts.GetSessionTokenOutput
	Err  error
}

func (m *MockSTSAPI) GetSessionToken(input *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
	return &m.Resp, m.Err
}

type MockCmdRunner struct {
	MockOutput map[string]string
	MockError  error
}

func (m *MockCmdRunner) RunCommand(name string, args ...string) ([]byte, error) {
	cmdKey := name + " " + strings.Join(args, " ")
	output, exists := m.MockOutput[cmdKey]
	if !exists {
		return nil, m.MockError
	}
	return []byte(output), m.MockError
}

/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package message_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"magma/dp/cloud/go/services/dp/active_mode_controller/message_generator/message"
	"magma/dp/cloud/go/services/dp/active_mode_controller/protos/active_mode"
)

const id = 123

func TestDeleteMessageString(t *testing.T) {
	m := message.NewDeleteMessage(id)
	expected := fmt.Sprintf("delete: %d", id)
	assert.Equal(t, expected, m.String())
}

func TestDeleteMessageSend(t *testing.T) {
	client := &stubDeleteClient{}

	m := message.NewDeleteMessage(id)
	require.NoError(t, m.Send(context.Background(), client))

	expected := &active_mode.DeleteCbsdRequest{Id: id}
	assert.Equal(t, expected, client.req)
}

type stubDeleteClient struct {
	active_mode.ActiveModeControllerClient
	req *active_mode.DeleteCbsdRequest
}

func (s *stubDeleteClient) DeleteCbsd(_ context.Context, in *active_mode.DeleteCbsdRequest, _ ...grpc.CallOption) (*empty.Empty, error) {
	s.req = in
	return &empty.Empty{}, nil
}

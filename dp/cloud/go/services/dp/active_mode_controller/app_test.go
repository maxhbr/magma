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

package active_mode_controller_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"magma/dp/cloud/go/services/dp/active_mode_controller"
	"magma/dp/cloud/go/services/dp/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/services/dp/active_mode_controller/test_utils/builders"
)

const (
	bufferSize       = 16
	timeout          = time.Millisecond * 50
	heartbeatTimeout = time.Second * 10
	pollingTimeout   = time.Second * 20
)

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, &AppTestSuite{})
}

type AppTestSuite struct {
	suite.Suite
	clock                *stubClock
	activeModeController *stubActiveModeControllerService
	appDone              chan error
	cancel               context.CancelFunc
	dialer               active_mode_controller.Dialer
	grpcServerDone       chan error
	grpcServer           *grpc.Server
}

func (s *AppTestSuite) SetupTest() {
	s.clock = &stubClock{ticker: make(chan time.Time, bufferSize)}
	s.activeModeController = &stubActiveModeControllerService{
		requests: make(chan *active_mode.RequestPayload, bufferSize),
		states:   make(chan *active_mode.State, bufferSize),
		err:      make(chan error, bufferSize),
	}
	s.givenGrpcServer()
	s.givenAppRunning()
}

func (s *AppTestSuite) TearDownTest() {
	s.whenAppWasShutdown()
	s.thenAppWasShutdown()
	s.thenNoOtherRequestWasReceived()

	s.whenGrpcServerWasShutdown()
	s.thenGrpcServerWasShutdown()
}

func (s *AppTestSuite) TestGetStateAndSendRequests() {
	s.givenState(buildSomeState("some"))
	s.whenTickerFired()
	s.thenRequestsWereEventuallyReceived(getExpectedRequests("some"))
}

// TODO cleanup this
func (s *AppTestSuite) TestCalculateHeartbeatDeadline() {
	const interval = 50 * time.Second
	const delta = heartbeatTimeout + pollingTimeout
	now := s.clock.Now()
	base := now.Add(delta - interval)
	timestamps := []time.Time{base.Add(time.Second), base}
	s.givenState(buildStateWithAuthorizedGrants("some", interval, timestamps...))
	s.whenTickerFired()
	s.thenRequestsWereEventuallyReceived(getExpectedHeartbeatRequests("some", "1"))
}

func (s *AppTestSuite) TestAppWorkInALoop() {
	s.givenState(buildSomeState("some"))
	s.whenTickerFired()
	s.thenRequestsWereEventuallyReceived(getExpectedRequests("some"))

	s.givenState(buildSomeState("other"))
	s.whenTickerFired()
	s.thenRequestsWereEventuallyReceived(getExpectedRequests("other"))
}

func (s *AppTestSuite) TestContinueWhenFailedToGetState() {
	s.givenStateError(errors.New("some error"))
	s.whenTickerFired()

	s.givenState(buildSomeState("some"))
	s.whenTickerFired()
	s.thenRequestsWereEventuallyReceived(getExpectedRequests("some"))
}

func (s *AppTestSuite) givenAppRunning() {
	s.appDone = make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	a := active_mode_controller.NewApp(
		active_mode_controller.WithDialer(s.dialer),
		active_mode_controller.WithClock(s.clock),
		active_mode_controller.WithRNG(stubRNG{}),
		active_mode_controller.WithDialTimeout(timeout),
		active_mode_controller.WithHeartbeatSendTimeout(heartbeatTimeout, timeout),
		active_mode_controller.WithRequestTimeout(timeout),
		active_mode_controller.WithPollingInterval(pollingTimeout),
		active_mode_controller.WithCbsdInactivityTimeout(timeout),
	)
	go func() {
		s.appDone <- a.Run(ctx)
	}()
}

func (s *AppTestSuite) givenGrpcServer() {
	listener := bufconn.Listen(bufferSize)
	s.grpcServer = grpc.NewServer()
	active_mode.RegisterActiveModeControllerServer(s.grpcServer, s.activeModeController)
	s.grpcServerDone = make(chan error)
	go func() {
		s.grpcServerDone <- s.grpcServer.Serve(listener)
	}()
	s.dialer = func(_ context.Context, _ string) (net.Conn, error) {
		return listener.Dial()
	}
}

func (s *AppTestSuite) givenState(state *active_mode.State) {
	s.activeModeController.states <- state
	s.activeModeController.err <- nil
}

func (s *AppTestSuite) givenStateError(err error) {
	s.activeModeController.states <- nil
	s.activeModeController.err <- err
}

func (s *AppTestSuite) whenAppWasShutdown() {
	s.cancel()
}

func (s *AppTestSuite) whenGrpcServerWasShutdown() {
	s.grpcServer.Stop()
}

func (s *AppTestSuite) whenTickerFired() {
	s.clock.ticker <- time.Time{}
}

func (s *AppTestSuite) thenAppWasShutdown() {
	select {
	case err := <-s.appDone:
		s.EqualError(err, context.Canceled.Error())
		return
	case <-time.After(timeout):
		s.Fail("failed to stop app")
	}
}

func (s *AppTestSuite) thenGrpcServerWasShutdown() {
	select {
	case err := <-s.grpcServerDone:
		s.NoError(err)
		return
	case <-time.After(timeout):
		s.Fail("failed to stop grpc server")
	}
}

func (s *AppTestSuite) thenRequestsWereEventuallyReceived(expectedRequests []*active_mode.RequestPayload) {
	timer := time.After(timeout)
	for _, expected := range expectedRequests {
		select {
		case actual := <-s.activeModeController.requests:
			s.JSONEq(expected.Payload, actual.Payload)
		case <-timer:
			s.Fail("Waiting for requests timed out")
		}
	}
}

func (s *AppTestSuite) thenNoOtherRequestWasReceived() {
	select {
	case actual := <-s.activeModeController.requests:
		s.Failf("Expected no more requests, got: %s", actual.Payload)
	default:
	}
}

func buildSomeState(names ...string) *active_mode.State {
	cbsds := make([]*active_mode.Cbsd, len(names))
	for i, name := range names {
		cbsds[i] = builders.NewCbsdBuilder().
			WithState(active_mode.CbsdState_Unregistered).
			WithName(name).
			Build()
	}
	return &active_mode.State{Cbsds: cbsds}
}

func buildStateWithAuthorizedGrants(name string, interval time.Duration, timestamps ...time.Time) *active_mode.State {
	b := builders.NewCbsdBuilder().
		WithName(name).
		WithChannel(builders.SomeChannel).
		WithAvailableFrequencies(builders.NoAvailableFrequencies).
		WithCarrierAggregation()
	for i, timestamp := range timestamps {
		b.WithGrant(&active_mode.Grant{
			Id:                     fmt.Sprintf("%d", i),
			State:                  active_mode.GrantState_Authorized,
			HeartbeatIntervalSec:   int64(interval / time.Second),
			LastHeartbeatTimestamp: timestamp.Unix(),
			LowFrequencyHz:         int64(3550+10*i) * 1e6,
			HighFrequencyHz:        int64(3550+10*(i+1)) * 1e6,
		})
	}
	return &active_mode.State{Cbsds: []*active_mode.Cbsd{b.Build()}}
}

func getExpectedRequests(name string) []*active_mode.RequestPayload {
	const template = `{"registrationRequest":[%s]}`
	request := fmt.Sprintf(template, getExpectedSingleRequest(name))
	return []*active_mode.RequestPayload{{Payload: request}}
}

func getExpectedSingleRequest(name string) string {
	const template = `{"userId":"%[1]s","fccId":"%[1]s","cbsdSerialNumber":"%[1]s"}`
	return fmt.Sprintf(template, name)
}

func getExpectedHeartbeatRequests(id string, grantIds ...string) []*active_mode.RequestPayload {
	if len(grantIds) == 0 {
		return nil
	}
	reqs := make([]string, len(grantIds))
	for i, grantId := range grantIds {
		reqs[i] = getExpectedHeartbeatRequest(id, grantId)
	}
	const template = `{"heartbeatRequest":[%s]}`
	payload := fmt.Sprintf(template, strings.Join(reqs, ","))
	return []*active_mode.RequestPayload{{Payload: payload}}
}

func getExpectedHeartbeatRequest(id string, grantId string) string {
	const template = `{"cbsdId":"%s","grantId":"%s","operationState":"AUTHORIZED"}`
	return fmt.Sprintf(template, id, grantId)
}

type stubRNG struct{}

func (stubRNG) Int() int {
	return 0
}

type stubClock struct {
	ticker chan time.Time
}

func (s *stubClock) Now() time.Time {
	return time.Unix(builders.Now, 0)
}

func (s *stubClock) Tick(_ time.Duration) *time.Ticker {
	return &time.Ticker{C: s.ticker}
}

type stubActiveModeControllerService struct {
	active_mode.UnimplementedActiveModeControllerServer
	requests chan *active_mode.RequestPayload
	states   chan *active_mode.State
	err      chan error
}

func (s *stubActiveModeControllerService) GetState(_ context.Context, _ *active_mode.GetStateRequest) (*active_mode.State, error) {
	return <-s.states, <-s.err
}

func (s *stubActiveModeControllerService) UploadRequests(_ context.Context, in *active_mode.RequestPayload) (*empty.Empty, error) {
	s.requests <- in
	return &empty.Empty{}, nil
}

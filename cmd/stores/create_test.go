package stores

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock_client "github.com/openfga/cli/mocks"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

var errMockCreate = errors.New("mock error")

func TestCreateError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFgaClient := mock_client.NewMockSdkClient(mockCtrl)

	mockExecute := mock_client.NewMockSdkClientCreateStoreRequestInterface(mockCtrl)

	var expectedResponse client.ClientCreateStoreResponse

	mockExecute.EXPECT().Execute().Return(&expectedResponse, errMockCreate)

	mockBody := mock_client.NewMockSdkClientCreateStoreRequestInterface(mockCtrl)

	body := client.ClientCreateStoreRequest{
		Name: "foo",
	}
	mockBody.EXPECT().Body(body).Return(mockExecute)

	mockFgaClient.EXPECT().CreateStore(context.Background()).Return(mockBody)

	_, err := create(mockFgaClient, "foo")
	if err == nil {
		t.Error("Expect error but there is none")
	}
}

func TestCreateSuccess(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFgaClient := mock_client.NewMockSdkClient(mockCtrl)

	mockExecute := mock_client.NewMockSdkClientCreateStoreRequestInterface(mockCtrl)
	expectedTime := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	expectedResponse := client.ClientCreateStoreResponse{
		Id:        openfga.PtrString("12345"),
		Name:      openfga.PtrString("foo"),
		CreatedAt: &expectedTime,
		UpdatedAt: &expectedTime,
	}

	mockExecute.EXPECT().Execute().Return(&expectedResponse, nil)

	mockBody := mock_client.NewMockSdkClientCreateStoreRequestInterface(mockCtrl)

	body := client.ClientCreateStoreRequest{
		Name: "foo",
	}
	mockBody.EXPECT().Body(body).Return(mockExecute)

	mockFgaClient.EXPECT().CreateStore(context.Background()).Return(mockBody)

	output, err := create(mockFgaClient, "foo")
	if err != nil {
		t.Error(err)
	}

	expectedOutput := "{\"created_at\":\"2009-11-10T23:00:00Z\",\"id\":\"12345\",\"name\":\"foo\",\"updated_at\":\"2009-11-10T23:00:00Z\"}" //nolint:lll
	if output != expectedOutput {
		t.Errorf("Expected output %v actual %v", expectedOutput, output)
	}
}

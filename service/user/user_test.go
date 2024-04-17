package user

import (
	"context"
	"reflect"
	"testing"
)

func TestGetUserById(t *testing.T) {
	mockUserProvider := NewMockUserRepo()
	service := New(&mockUserProvider)

	testCases := GetValidTestGetUserByIdTestCases()
	testCases = append(testCases, GetInvalidTestGetUserByIdTestCase())

	for _, tc := range testCases {
		actualStatusCode, actualResponse := service.GetUserById(context.Background(), tc.UserId)

		if actualStatusCode != tc.ExpectedStatusCode {
			t.Fatalf("test case userId %d, expected exected status code %d, got %d", tc.UserId, tc.ExpectedStatusCode, actualStatusCode)
		}

		if !reflect.DeepEqual(actualResponse.Errors, tc.Response.Errors) {
			t.Fatalf("test case userId %d, expected response errors %+v, got %+v", tc.UserId, tc.Response.Errors, actualResponse.Errors)
		}

		if actualResponse.User != nil && tc.Response.User != nil && (actualResponse.User.Id != tc.Response.User.Id || actualResponse.User.Email != tc.Response.User.Email ||
			actualResponse.User.FirstName != tc.Response.User.FirstName || actualResponse.User.LastName != tc.Response.User.LastName ||
			actualResponse.User.Password != tc.Response.User.Password) {
			t.Fatalf("test case userId %d, expected response user %+v, got %+v", tc.UserId, tc.Response.User, actualResponse.User)
		}
	}
}

package gomock79587160

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/mock"
)

// EventData interface defines a method MySub.
// Both MyEventData and MockedEventData implement this interface.
// This way, you can swap out the real implementation with a mock in tests.
// This is useful if you want to test something that _TAKES_ an EventData
// interface as a parameter, and you want to control the behavior of that
// interface in your tests.
type EventData interface {
	MySub(data int) bool
}

// MyEventData implements EventData interface
// and provides the MySub method.
type MyEventData struct {
	ID   int
	Name string
}

// MySub does something with the data.
// It also implements the MySub method from the EventData interface,
// so this is a real implementation.
func (event *MyEventData) MySub(data int) bool {
	slog.Info("Code MySub method called", "data", fmt.Sprintf("%v", data))
	if data < 0 {
		return false
	}
	return true
}

func (event *MyEventData) MyData() error {
	slog.Info("Code MyData method called")
	ev := MyEventData{
		ID:   1,
		Name: "data",
	}

	status := ev.MySub(ev.ID)
	if !status {
		return errors.New("data is negative")
	}
	slog.Info(" method called", "value", fmt.Sprintf("%#v", ev))
	return nil
}

// MockedEventData is a mock implementation of the EventData interface.
// It uses the testify/mock package to create a mock object.
// This allows you to set expectations and return values for the methods
// of the mock object.
// This is useful for unit testing, where you want to isolate the code
// being tested from its dependencies.
type MockedEventData struct {
	mock.Mock
	t *testing.T
}

// MySub satisfies the EventData interface.
// Note that we don't need to implement the MyData method here,
// because we are only interested in the EventData interface for this test.
func (m *MockedEventData) MySub(data int) bool {
	args := m.Called(data)
	rv := args.Bool(0)
	m.t.Logf("MockedEventData MySub method called with data %d returning %t", data, rv)
	return args.Bool(0)
}

// Let's assume you want to test EventdataProcessors.
// Now, instead of passing a real EventData, you can pass a mock.
// This allows you to control the behavior of the EventData
// without calling the real MySub method, which may have side effects.
// For simplicity reasons, EventDataProcessor is just a function
// that takes an EventData as a parameter and returns an error - or not.
// Note that it is the EventDataProcessor which is under test here - not
// the EventData interface and it's implementation(s).
type EventDataProccessor func(event EventData) error

// This is just a simple example.
// Now, with our mock in place...
var eventDataSubcallWithOne EventDataProccessor = func(event EventData) error {
	slog.Info("Code EventDataProccessor method called")
	if event == nil {
		return errors.New("event is nil")
	}
	// ... we can react to the same call of the MySub function
	// with a _different_ behavior.
	if !event.MySub(1) {
		return fmt.Errorf("Something terrible happened")
	}
	return nil
}

func TestEvenDataProcessorWithMock(t *testing.T) {

	// Create a new instance of the mock object
	mockedReturningTrue := &MockedEventData{
		t: t,
	}

	// Set up the expectation for the MySub method
	mockedReturningTrue.On("MySub", 1).Return(true)

	// Call the EventDataProccessor with the mock object
	err := eventDataSubcallWithOne(mockedReturningTrue)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create a new instance of the mock object
	mockedReturningFalse := &MockedEventData{
		t: t,
	}
	// Set up the expectation for the MySub method
	mockedReturningFalse.On("MySub", 1).Return(false)
	// Call the EventDataProccessor with the mock object
	// Note that the SAME call will go to the mock object
	// but it will show a different behavior.
	err = eventDataSubcallWithOne(mockedReturningFalse)
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

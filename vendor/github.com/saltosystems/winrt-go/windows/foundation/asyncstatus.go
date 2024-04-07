// Code generated by winrt-go-gen. DO NOT EDIT.

//go:build windows

//nolint:all
package foundation

type AsyncStatus int32

const SignatureAsyncStatus string = "enum(Windows.Foundation.AsyncStatus;i4)"

const (
	AsyncStatusCanceled  AsyncStatus = 2
	AsyncStatusCompleted AsyncStatus = 1
	AsyncStatusError     AsyncStatus = 3
	AsyncStatusStarted   AsyncStatus = 0
)
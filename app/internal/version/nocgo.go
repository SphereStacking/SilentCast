//go:build !cgo
// +build !cgo

package version

// _cgoEnabled returns false when CGO is disabled
func _cgoEnabled() bool {
	return false
}

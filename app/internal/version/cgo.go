//go:build cgo
// +build cgo

package version

// _cgoEnabled returns true when CGO is enabled
func _cgoEnabled() bool {
	return true
}
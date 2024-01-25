//go:build !(go1.20 && goexperiment.arenas)
// +build !go1.20 !goexperiment.arenas

package features

// GoArenaAvailndicates whether the experimental arena
// feature is enabled in the build.
const GoArenaAvail = false

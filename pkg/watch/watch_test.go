package watch_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/porter614/gobones/pkg/watch"
)

func TestWatchFile(t *testing.T) {
	// TC: Try to watch a file that doesn't exist - should return error
	_, _, err := watch.WatchFile("foobar")
	assert.Error(t, err, "Expected error watching non-existent file")

	// Create a temporary file
	f, err := ioutil.TempFile("", "foobar")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	// TC: Watch the file, change it - should put event on channel
	ch, _, err := watch.WatchFile(f.Name())
	assert.NoErrorf(t, err, "Error watching file: %v", err)
	if _, err = f.Write([]byte("foobar")); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	event := <-ch
	assert.Equal(t, []byte("foobar"), event)
	f.Close()
}

// -----------------------------------------------------------------------------
// (c) balarabe@protonmail.com                                      License: MIT
// :v: 2019-04-04 17:29:22 518331                 zr-fs/[dir_watcher_windows.go]
// -----------------------------------------------------------------------------

package fs

import (
	"os"

	"github.com/balacode/zr"
)

// A DirWatcher holds a channel that delivers a folder's
// path when any file in the folder changes.
type DirWatcher struct {
	Chan <-chan string // the channel on which path changes are sent
	dir  string
} //                                                                  DirWatcher

// NewDirWatcher returns a new DirWatcher that contains
// a channel that be sent the name of a file every time
// a file in the folder or one of its subfolders changes.
func NewDirWatcher(dir string) *DirWatcher {
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		zr.Error("Folder^", dir, "does not exist")
		return nil
	}
	c := make(chan string, 1)
	ret := &DirWatcher{
		Chan: c,
		dir:  dir,
	}
	go waitForDirChange(c, dir)
	return ret
} //                                                               NewDirWatcher

//end

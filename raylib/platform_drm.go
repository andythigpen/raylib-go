//go:build linux && drm && !android
// +build linux,drm,!android

package rl

/*
#include "raylib.h"
#include <stdlib.h>
#include <stdint.h>

extern int GetDrmConnectorPropertyValue(uint32_t property_id, uint64_t *property_value);
extern int SetDrmConnectorProperty(uint32_t property_id, uint64_t property_value);
*/
import "C"

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"unsafe"
)

// SetMain - Sets callback function
func SetMain(func()) {
	return
}

// ShowCursor - Shows cursor
func ShowCursor() {
	C.ShowCursor()
}

// HideCursor - Hides cursor
func HideCursor() {
	C.HideCursor()
}

// IsCursorHidden - Returns true if cursor is not visible
func IsCursorHidden() bool {
	ret := C.IsCursorHidden()
	v := bool(ret)
	return v
}

// IsCursorOnScreen - Check if cursor is on the current screen.
func IsCursorOnScreen() bool {
	ret := C.IsCursorOnScreen()
	v := bool(ret)
	return v
}

// EnableCursor - Enables cursor
func EnableCursor() {
	C.EnableCursor()
}

// DisableCursor - Disables cursor
func DisableCursor() {
	C.DisableCursor()
}

// IsFileDropped - Check if a file have been dropped into window
func IsFileDropped() bool {
	return false
}

// LoadDroppedFiles - Load dropped filepaths
func LoadDroppedFiles() (files []string) {
	return
}

// UnloadDroppedFiles - Unload dropped filepaths
func UnloadDroppedFiles() {
	return
}

// Open implements fs.FS interface - opens the named file for reading
func (a *Asset) Open(name string) (fs.File, error) {
	return openAssetFile(a.root, name)
}

// ReadFile implements fs.ReadFileFS interface - reads the entire file
func (a *Asset) ReadFile(name string) ([]byte, error) {
	return readAssetFile(a.root, name)
}

// ReadDir implements fs.ReadDirFS interface - reads the directory
func (a *Asset) ReadDir(name string) ([]fs.DirEntry, error) {
	return readAssetDir(a.root, name)
}

// desktopAsset wraps os.File to implement AssetFile interface
type desktopAsset struct {
	*os.File
}

func (d *desktopAsset) Stat() (fs.FileInfo, error) {
	return d.File.Stat()
}

func (d *desktopAsset) Seek(offset int64, whence int) (int64, error) {
	return d.File.Seek(offset, whence)
}

func openAssetFile(root, name string) (fs.File, error) {
	fullPath := name
	if root != "" {
		fullPath = filepath.Join(root, name)
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	return &desktopAsset{f}, nil
}

func readAssetFile(root, name string) ([]byte, error) {
	fullPath := name
	if root != "" {
		fullPath = filepath.Join(root, name)
	}

	return os.ReadFile(fullPath)
}

func readAssetDir(root, name string) ([]fs.DirEntry, error) {
	fullPath := name
	if root != "" {
		fullPath = filepath.Join(root, name)
	}

	return os.ReadDir(fullPath)
}

func SetDrmConnectorProperty(id uint32, val uint64) error {
	ret := C.SetDrmConnectorProperty(C.uint(id), C.ulong(val))
	if ret < 0 {
		return fmt.Errorf("failed to set property value %u", id)
	}
	return nil
}

func GetDrmConnectorPropertyValue(id uint32) (uint64, error) {
	val := uint64(0)
	ret := C.GetDrmConnectorPropertyValue(C.uint(id), (*C.ulong)(unsafe.Pointer(&val)))
	if ret < 0 {
		return val, fmt.Errorf("failed to get property value %u", id)
	}
	return val, nil
}

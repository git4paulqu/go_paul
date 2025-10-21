// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package driver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// FileItem represents a node in the directory tree
type FileItem struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"isDir"`
	Children []FileItem `json:"children,omitempty"`
}

type FileBrowserResponse struct {
	CurrentPath string     `json:"currentPath"`
	Items       []FileItem `json:"items"`
	Error       string     `json:"error,omitempty"`
}

// Cache for directory probes
var hasPprofCache = struct {
	mu sync.Mutex
	m  map[string]bool
}{m: map[string]bool{}}

func hasPprofBelowLimitedCached(p string) bool {
	hasPprofCache.mu.Lock()
	v, ok := hasPprofCache.m[p]
	hasPprofCache.mu.Unlock()
	if ok {
		return v
	}
	// depth and entry limits chosen to balance speed and accuracy
	v = hasPprofBelowLimited(p, 3, 1000)
	hasPprofCache.mu.Lock()
	hasPprofCache.m[p] = v
	hasPprofCache.mu.Unlock()
	return v
}

// fileBrowser serves the file browser page
func (ui *webInterface) fileBrowser(w http.ResponseWriter, req *http.Request) {
	html := &bytes.Buffer{}
	data := webArgs{Title: "File Browser"}
	if err := renderHTML(html, "filebrowser", nil, nil, nil, data); err != nil {
		http.Error(w, "internal template error", http.StatusInternalServerError)
		ui.options.UI.PrintErr(err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(html.Bytes())
}

// apiFiles returns immediate children (for Windows root: drives) or a recursive
// filtered tree containing all .pprof files and their directories
func (ui *webInterface) apiFiles(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get("path")
	if path == "/" {
		// Use the configured file browser path as default
		if ui.fileBrowserPath != "" {
			path = ui.fileBrowserPath
		} else {
			path = "/"
		}
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join("/", path)
	}

	// Check if we need to rescan due to file changes
	forceRescan := req.URL.Query().Get("force") == "true"
	shouldRescan := forceRescan || ui.shouldRescanDirectory(path)

	var (
		items []FileItem
		err   error
	)
	if path == "/" && runtime.GOOS == "windows" {
		items, err = listWindowsDrives()
	} else {
		items, err = scanDirFull(path)
	}

	// Update scan time if successful and we actually rescanned
	if err == nil && shouldRescan {
		ui.lastScanTime[path] = time.Now().Unix()
	}

	resp := FileBrowserResponse{CurrentPath: path, Items: items}
	if err != nil {
		resp.Error = err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// shouldRescanDirectory checks if a directory should be rescanned based on file changes
func (ui *webInterface) shouldRescanDirectory(path string) bool {
	lastScan, exists := ui.lastScanTime[path]
	if !exists {
		return true // Never scanned before
	}

	// Check if directory modification time is newer than last scan
	info, err := os.Stat(path)
	if err != nil {
		return true // Error reading directory, rescan
	}

	modTime := info.ModTime().Unix()
	return modTime > lastScan
}

func scanDir(root string) ([]FileItem, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var result []FileItem
	for _, e := range entries {
		p := filepath.Join(root, e.Name())
		if e.IsDir() {
			if hasPprofBelowLimitedCached(p) {
				result = append(result, FileItem{Name: e.Name(), Path: p, IsDir: true})
			}
			continue
		}
		nameLower := strings.ToLower(e.Name())
		if strings.HasSuffix(nameLower, ".pprof") {
			result = append(result, FileItem{Name: e.Name(), Path: p, IsDir: false})
		}
	}
	return result, nil
}

// scanDirFull returns a recursive tree of directories that contain at least
// one .pprof descendant and includes all .pprof files under them.
func scanDirFull(root string) ([]FileItem, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var result []FileItem
	for _, e := range entries {
		p := filepath.Join(root, e.Name())
		if e.IsDir() {
			children, err := scanDirFull(p)
			if err != nil {
				continue
			}
			if len(children) > 0 {
				result = append(result, FileItem{Name: e.Name(), Path: p, IsDir: true, Children: children})
			}
			continue
		}
		nameLower := strings.ToLower(e.Name())
		if strings.HasSuffix(nameLower, ".pprof") {
			result = append(result, FileItem{Name: e.Name(), Path: p, IsDir: false})
		}
	}
	return result, nil
}

// hasPprofBelowLimited quickly checks if a directory likely contains a .pprof file
// by scanning up to maxDepth levels and up to maxEntries entries in total.
func hasPprofBelowLimited(root string, maxDepth int, maxEntries int) bool {
	type node struct {
		path  string
		depth int
	}
	q := []node{{root, 0}}
	scanned := 0
	for len(q) > 0 {
		if scanned >= maxEntries {
			return false
		}
		cur := q[0]
		q = q[1:]
		entries, err := os.ReadDir(cur.path)
		if err != nil {
			continue
		}
		for _, e := range entries {
			scanned++
			if !e.IsDir() {
				nameLower := strings.ToLower(e.Name())
				if strings.HasSuffix(nameLower, ".pprof") {
					return true
				}
				continue
			}
			if cur.depth < maxDepth {
				q = append(q, node{filepath.Join(cur.path, e.Name()), cur.depth + 1})
			}
			if scanned >= maxEntries {
				break
			}
		}
	}
	return false
}

// listWindowsDrives lists existing drives and returns them as directories
func listWindowsDrives() ([]FileItem, error) {
	var items []FileItem
	for letter := 'A'; letter <= 'Z'; letter++ {
		root := string([]rune{letter}) + ":\\"
		if _, err := os.Stat(root); err == nil {
			items = append(items, FileItem{Name: root, Path: root, IsDir: true})
		}
	}
	return items, nil
}

func (ui *webInterface) apiSetBase(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Query().Get("path")
	if p == "" {
		http.Error(w, "path required", http.StatusBadRequest)
		return
	}
	ui.baseProfile = p
}

func (ui *webInterface) apiSetSource(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Query().Get("path")
	if p == "" {
		http.Error(w, "path required", http.StatusBadRequest)
		return
	}
	ui.sourceProfile = p
}

func (ui *webInterface) apiReset(w http.ResponseWriter, req *http.Request) {
	ui.baseProfile, ui.sourceProfile = "", ""
}

func checkFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (ui *webInterface) apiWork(w http.ResponseWriter, req *http.Request) {
	if ui.baseProfile == "" || !checkFileExists(ui.baseProfile) {
		http.Error(w, "base profile required", http.StatusBadRequest)
		return
	}

	// Build a temporary source from selected files and refresh the UI profile
	src := &source{
		Sources:  []string{ui.baseProfile},
		DiffBase: false,
	}

	if ui.sourceProfile != "" && checkFileExists(ui.sourceProfile) {
		src.Base = []string{ui.baseProfile}
		src.Sources = []string{ui.sourceProfile}
		src.DiffBase = true
	}

	p, err := fetchProfiles(src, ui.options)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ui.prof = p
	ui.copier = makeProfileCopier(p)

	// Redirect to the graph view so user can see the result immediately
	http.Redirect(w, req, "../", http.StatusFound)
}

// apiSetDirectory sets the selected directory
func (ui *webInterface) apiSetDirectory(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path required", http.StatusBadRequest)
		return
	}
	ui.selectedDir = path
	w.WriteHeader(http.StatusOK)
}

// apiSetExpanded sets the expanded state of a directory
func (ui *webInterface) apiSetExpanded(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get("path")
	expanded := req.URL.Query().Get("expanded") == "true"
	if path == "" {
		http.Error(w, "path required", http.StatusBadRequest)
		return
	}
	ui.expandedDirs[path] = expanded
	w.WriteHeader(http.StatusOK)
}

// apiGetState returns the current state (selected directory, expanded directories, base and source profiles)
func (ui *webInterface) apiGetState(w http.ResponseWriter, req *http.Request) {
	state := map[string]interface{}{
		"selectedDir":   ui.selectedDir,
		"expandedDirs":  ui.expandedDirs,
		"baseProfile":   ui.baseProfile,
		"sourceProfile": ui.sourceProfile,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

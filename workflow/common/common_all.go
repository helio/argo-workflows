// +build !windows

package common

import "path/filepath"

func GenerateMountPath(mountPath string) string {
	return filepath.Join(ExecutorMainFilesystemDir, mountPath)
}

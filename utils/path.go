package utils

import "go/build"


// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func ImportPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}


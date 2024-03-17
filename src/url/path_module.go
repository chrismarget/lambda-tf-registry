package url

import "net/http"

func parseModulePath(parts []string) (Path, error) {
	return nil, newPathError(http.StatusInternalServerError, "module path not implemented")
}

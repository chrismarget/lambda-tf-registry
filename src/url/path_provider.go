package url

import (
	"net/http"

	"github.com/orsinium-labs/enum"
)

const providerPathMinParts = 5

type ProviderPathType enum.Member[string]

var (
	ProviderPathTypeDownload = ProviderPathType{Value: "download"}
	ProviderPathTypeVersions = ProviderPathType{Value: "versions"}
	ProviderPathTypes        = enum.New(
		ProviderPathTypeDownload,
		ProviderPathTypeVersions,
	)
)

func parseProviderPath(parts []string) (Path, error) {
	var ppt *ProviderPathType
	switch len(parts) {
	case 5:
		ppt = ProviderPathTypes.Parse(parts[4])
		if ppt == nil {
			return nil, newPathError(http.StatusUnprocessableEntity, "failed to parse provider path %q", parts)
		}
	case 8:
		ppt = ProviderPathTypes.Parse(parts[5])
		if ppt == nil {
			return nil, newPathError(http.StatusUnprocessableEntity, "failed to parse provider path %q", parts)
		}
	default:
		return nil, newPathError(http.StatusUnprocessableEntity, "failed to parse provider path %q", parts)
	}

	switch *ppt {
	case ProviderPathTypeVersions:
		return new(ProviderVersionsPath).loadParts(parts)
	case ProviderPathTypeDownload:
		return new(ProviderDownloadPath).loadParts(parts)
	default:
		return nil, newPathError(http.StatusInternalServerError, "unhandled provider path type %q", ppt.Value)
	}
}

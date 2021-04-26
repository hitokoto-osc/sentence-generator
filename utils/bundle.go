package utils

import (
	"encoding/json"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"os"
	"path/filepath"
)

// GetBundleVersion Get bundle version by reading version lockfile
func GetBundleVersion() (string, error) {
	buffer, err := os.ReadFile(filepath.Join(config.Core.Workdir, "./version.json"))
	if err != nil {
		return "", err
	}
	var data VersionData
	if err = json.Unmarshal(buffer, &data); err != nil {
		return "", err
	}
	return data.BundleVersion, nil
}

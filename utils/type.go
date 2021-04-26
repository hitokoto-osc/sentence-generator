package utils

import (
	"fmt"
	"github.com/blang/semver/v4"
)

// SentencesUnit  is a unit of VersionData SentencesUnitCollection
type SentencesUnit struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	Path      string `json:"path"`
	Timestamp int64  `json:"timestamp"`
}

// VersionData is the structure of version lockfile
type VersionData struct {
	ProtocolVersion string `json:"protocol_version"`
	BundleVersion   string `json:"bundle_version"`
	UpdatedAt       int64  `json:"updated_at"`
	Categories      struct {
		Path      string `json:"path"`
		Timestamp int64  `json:"timestamp"`
	} `json:"categories"`
	Sentences SentencesUnitCollection `json:"sentences"`
}

// StepVersion will increase patch version
func (v *VersionData) StepVersion() error {
	version := semver.MustParse(v.BundleVersion)
	if err := version.IncrementPatch(); err != nil {
		return err
	}
	v.BundleVersion = version.String()
	return nil
}

// UpdateCategoriesRecord will update timestamp in categories field and updated_at field
func (v *VersionData) UpdateCategoriesRecord() {
	ts := GetMillionSecondTimestamp()
	v.Categories.Timestamp = ts
	v.UpdatedAt = ts
}

// UpdateSentenceRecord will update timestamp in specific category of sentences field and updated_at field
func (v *VersionData) UpdateSentenceRecord(categoryKey string) error {
	c, ok := v.Sentences.Find(categoryKey)
	if !ok {
		return fmt.Errorf("The key %s is not exist", categoryKey)
	}
	ts := GetMillionSecondTimestamp()
	c.Timestamp = ts
	v.UpdatedAt = ts
	return nil
}

// SentencesUnitCollection is a collection of SentencesUnit
type SentencesUnitCollection []SentencesUnit

// Update will update the timestamp of specific SentencesUnit
func (c *SentencesUnitCollection) Update(key string) error {
	s, ok := c.Find(key)
	if !ok {
		return fmt.Errorf("The key %s is not exist", key)
	}
	s.Timestamp = GetMillionSecondTimestamp()
	return nil
}

// Find will return specific SentencesUnit Pointer
func (c *SentencesUnitCollection) Find(key string) (*SentencesUnit, bool) {
	for _, v := range *c {
		if v.Key == key {
			return &v, true
		}
	}
	return nil, false
}

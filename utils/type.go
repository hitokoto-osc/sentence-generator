// Package utils is intended to provide some useful functions
package utils

import (
	"fmt"
	"github.com/Masterminds/semver/v3"

	"github.com/hitokoto-osc/sentence-generator/database"
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

// DeepCopyVersionData deep copy from source
func DeepCopyVersionData(data VersionData) VersionData {
	return VersionData{
		ProtocolVersion: data.ProtocolVersion,
		BundleVersion:   data.BundleVersion,
		UpdatedAt:       data.UpdatedAt,
		Categories: struct {
			Path      string `json:"path"`
			Timestamp int64  `json:"timestamp"`
		}{
			Path:      data.Categories.Path,
			Timestamp: data.Categories.Timestamp,
		},
		Sentences: DeepCopySentencesUnitCollection(data.Sentences),
	}
}

// StepVersion will increase patch version
func (v *VersionData) StepVersion() error {
	version := semver.MustParse(v.BundleVersion).IncPatch()
	v.BundleVersion = version.String()
	return nil
}

// UpdateCategoriesRecord will update timestamp in categories field and updated_at field
func (v *VersionData) UpdateCategoriesRecord(categoriesData CategoryUnitCollection) {
	ts := GetMillionSecondTimestamp()
	tmp := SentencesUnitCollection{}
	for _, c := range categoriesData {
		if unit, ok := v.Sentences.Find(c.Key); !ok || unit.Name != c.Name || unit.Path != c.Path {
			tmp = append(tmp, SentencesUnit{
				Name:      c.Name,
				Key:       c.Key,
				Path:      c.Path,
				Timestamp: ts,
			})
		} else {
			tmp = append(tmp, *unit)
		}
	}
	v.Sentences = tmp
	v.Categories.Timestamp = ts
	v.UpdatedAt = ts
}

// UpdateSentenceRecord will update timestamp in specific category of sentences field and updated_at field
func (v *VersionData) UpdateSentenceRecord(categoryKey string) error {
	c, ok := v.Sentences.Find(categoryKey)
	if !ok {
		return fmt.Errorf("the key %s is not exist", categoryKey)
	}
	ts := GetMillionSecondTimestamp()
	c.Timestamp = ts
	v.UpdatedAt = ts
	return nil
}

// SentencesUnitCollection is a collection of SentencesUnit
type SentencesUnitCollection []SentencesUnit

// DeepCopySentencesUnitCollection deep copy from source
func DeepCopySentencesUnitCollection(s SentencesUnitCollection) SentencesUnitCollection {
	tmp := SentencesUnitCollection{}
	for _, v := range s {
		tmp = append(tmp, SentencesUnit{
			Name:      v.Name,
			Key:       v.Key,
			Path:      v.Path,
			Timestamp: v.Timestamp,
		})
	}
	return tmp
}

// Update will update the timestamp of specific SentencesUnit
func (c *SentencesUnitCollection) Update(key string) error {
	s, ok := c.Find(key)
	if !ok {
		return fmt.Errorf("the key %s is not exist", key)
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

// CategoryUnit is unit of category collection, stored in categories.json file
type CategoryUnit struct {
	database.Category
	Path string `json:"path"`
}

// CategoryUnitCollection is collection of category unit
type CategoryUnitCollection []CategoryUnit

// DeepCopyCategoryUnitCollection deep copy from source
func DeepCopyCategoryUnitCollection(collection CategoryUnitCollection) CategoryUnitCollection {
	tmp := CategoryUnitCollection{}
	for _, v := range collection {
		tmp = append(tmp, CategoryUnit{
			Category: database.Category{
				ID:        v.ID,
				Name:      v.Name,
				Desc:      v.Desc,
				Key:       v.Key,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			},
			Path: v.Path,
		})
	}
	return tmp
}

// ImportFrom can convert []category to []categoryUnit
func (p *CategoryUnitCollection) ImportFrom(c []database.Category) {
	for _, v := range c {
		*p = append(*p, CategoryUnit{
			Category: v,
			Path:     fmt.Sprintf("./sentences/%s.json", v.Key),
		})
	}
}

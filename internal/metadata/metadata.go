package metadata

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

/*
match everything inbetween the '---' pairs
e.g.
---
this is matched content
---
*/
var metadataRegex = regexp.MustCompile(`(?s)^---\n(.*?)\n---`)

type Date time.Time

func (d *Date) UnmarshalYAML(v *yaml.Node) error {
	parsed, err := time.Parse("2006-01-02", v.Value)
	if err != nil {
		return err
	}
	*d = Date(parsed)
	return nil
}

func (d Date) Equal(t Date) bool {
	a := time.Time(d)
	b := time.Time(t)
	ayr, amo, aday := a.Date()
	byr, bmo, bday := b.Date()
	return ayr == byr && amo == bmo && aday == bday
}

func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}

type Metadata struct {
	ID          int
	Filepath    string
	Title       string   `yaml:"title"`
	Author      string   `yaml:"author"`
	Created     Date     `yaml:"created"`
	LastUpdated Date     `yaml:"last_updated"`
	Tags        []string `yaml:"tags"`
	IsHome      bool     `yaml:"home"`
}

func (m *Metadata) String() string {
	s := `---
title: %s
author: %s
created: %s
last_updated: %s
tags: %s
home: %v
---`
	return fmt.Sprintf(s, m.Title, m.Author, m.Created, m.LastUpdated, m.Tags, m.IsHome)
}

func MetadataTemplate() string {
	m := Metadata{
		Title: "<title>",
		Author: "<author>",
		Created: Date(time.Now()),
		LastUpdated: Date(time.Now()),
		Tags: []string{"list", "of", "tags"},
		IsHome: false,
	}
	return m.String()
}

func (m *Metadata) ContainsTag(tagName string) bool {
	for _, tag := range m.Tags {
		if tagName == tag {
			return true
		}
	}
	return false
}

func parseMetadata(content []byte) (Metadata, error) {
	matches := metadataRegex.FindSubmatch(content)
	if len(matches) < 2 {
		return Metadata{}, fmt.Errorf("no metadata found. use this template:\n%s", MetadataTemplate())
	}
	var metadata Metadata
	if err := yaml.Unmarshal(matches[1], &metadata); err != nil {
		return Metadata{}, fmt.Errorf("unable to parse metadata: %w", err)
	}
	return metadata, nil
}

func ExtractFromFile(filepath string) (Metadata, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return Metadata{}, err
	}
	m, err := parseMetadata(content)
	if err != nil {
		return Metadata{}, fmt.Errorf("unable to extract metadata from file %s: %w", filepath, err)
	}
	m.Filepath = filepath
	return m, nil
}

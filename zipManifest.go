// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    zipManifest, err := UnmarshalZipManifest(bytes)
//    bytes, err = zipManifest.Marshal()

package main

import "encoding/json"

func UnmarshalZipManifest(data []byte) (ZipManifest, error) {
	var r ZipManifest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ZipManifest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ZipManifest struct {
	Minecraft       Minecraft `json:"minecraft"`
	ManifestType    string    `json:"manifestType"`
	Overrides       string    `json:"overrides"`
	ManifestVersion int64     `json:"manifestVersion"`
	Version         string    `json:"version"`
	Author          string    `json:"author"`
	Name            string    `json:"name"`
	Files           []File    `json:"files"`
}

type File struct {
	ProjectID int64 `json:"projectID"`
	FileID    int64 `json:"fileID"`
	Required  bool  `json:"required"`
}

type Minecraft struct {
	Version    string      `json:"version"`
	ModLoaders []ModLoader `json:"modLoaders"`
}

type ModLoader struct {
	ID      string `json:"id"`
	Primary bool   `json:"primary"`
}

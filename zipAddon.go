// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    zipAddon, err := UnmarshalZipAddon(bytes)
//    bytes, err = zipAddon.Marshal()

package main

import "encoding/json"

func UnmarshalZipAddon(data []byte) (ZipAddon, error) {
	var r ZipAddon
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ZipAddon) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ZipAddon struct {
	ID                      int64         `json:"id"`
	DisplayName             string        `json:"displayName"`
	FileName                string        `json:"fileName"`
	FileDate                string        `json:"fileDate"`
	FileLength              int64         `json:"fileLength"`
	ReleaseType             int64         `json:"releaseType"`
	FileStatus              int64         `json:"fileStatus"`
	DownloadURL             string        `json:"downloadUrl"`
	IsAlternate             bool          `json:"isAlternate"`
	AlternateFileID         int64         `json:"alternateFileId"`
	Dependencies            []interface{} `json:"dependencies"`
	IsAvailable             bool          `json:"isAvailable"`
	Modules                 []Module      `json:"modules"`
	PackageFingerprint      int64         `json:"packageFingerprint"`
	GameVersion             []string      `json:"gameVersion"`
	InstallMetadata         interface{}   `json:"installMetadata"`
	ServerPackFileID        interface{}   `json:"serverPackFileId"`
	HasInstallScript        bool          `json:"hasInstallScript"`
	GameVersionDateReleased string        `json:"gameVersionDateReleased"`
	GameVersionFlavor       interface{}   `json:"gameVersionFlavor"`
}

type Module struct {
	Foldername  string `json:"foldername"`
	Fingerprint int64  `json:"fingerprint"`
}

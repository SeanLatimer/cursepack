// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    cPVersion, err := UnmarshalCPVersion(bytes)
//    bytes, err = cPVersion.Marshal()

package main

import "encoding/json"

func UnmarshalCPVersion(data []byte) (CPVersion, error) {
	var r CPVersion
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CPVersion) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CPVersion struct {
	Version string `json:"version"`
}

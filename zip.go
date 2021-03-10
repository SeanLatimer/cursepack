package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"golang.org/x/sync/errgroup"
)

const FORGESVC_URL = "https://addons-ecs.forgesvc.net/api/v2"

func handleZipPack(opts PackInstallOptions) error {
	jww.INFO.Println("Installing ZIP pack")
	tempPath := filepath.Join(os.TempDir(), TEMPDIR_NAME)
	os.MkdirAll(tempPath, 0700)

	packPath, err := filepath.Abs(opts.Path)
	if err != nil {
		errors.Wrapf(err, "Failed to resolve path: %s", packPath)
	}
	zipPath := filepath.Join(packPath, opts.Pack)

	if !fileExists(zipPath) {
		return fmt.Errorf("Provided zip does not exist at %s", zipPath)
	}

	bytes, err := extractZipManifest(zipPath)
	if err != nil {
		return err
	}
	manifest, err := UnmarshalZipManifest(bytes)
	if err != nil {
		return err
	}

	// TODO: Possibly change this later to not blow away the whole mods folder
	modsDest := filepath.Join(packPath, "mods")
	err = os.RemoveAll(modsDest)
	if err != nil {
		return err
	}

	err = downloadPackMods(manifest, modsDest)
	if err != nil {
		return err
	}

	err = extractZipOverrides(zipPath, manifest.Overrides, packPath)
	if err != nil {
		return err
	}

	if opts.Server {
		mcVersion := manifest.Minecraft.Version
		modLoader := manifest.Minecraft.ModLoaders[0].ID
		forgeVersion := mcVersion + "-" + modLoader[6:]
		err := installForgeServer(forgeVersion, packPath)
		if err != nil {
			return err
		}
	}

	return err
}

// https://addons-ecs.forgesvc.net/api/v2/addon/231275/file/3222705
func downloadPackMods(manifest ZipManifest, dest string) error {
	jww.INFO.Println("Downloading pack files")
	err := os.MkdirAll(dest, 0700)
	if err != nil {
		return err
	}

	var g errgroup.Group
	sem := make(chan struct{}, 5)
	for _, file := range manifest.Files {
		file := file // create locals for closure below
		sem <- struct{}{}
		g.Go(func() error {
			defer func() {
				<-sem
			}()
			zipAddon, err := downloadZipAddon(file.ProjectID, file.FileID)
			if err != nil {
				return err
			}

			jww.INFO.Printf("Downloading %s", zipAddon.FileName)
			modPath := filepath.Join(dest, path.Base(zipAddon.DownloadURL))
			modFile, err := os.Create(modPath)
			if err != nil {
				return err
			}
			defer modFile.Close()

			modResp, err := http.Get(zipAddon.DownloadURL)
			if err != nil {
				return err
			}
			defer modResp.Body.Close()

			_, err = io.Copy(modFile, modResp.Body)
			if err != nil {
				return err
			}

			return nil
		})
	}
	return g.Wait()
}

func downloadZipAddon(projectID int64, fileID int64) (ZipAddon, error) {
	resp, err := http.Get(fmt.Sprintf("%s/addon/%d/file/%d", FORGESVC_URL, projectID, fileID))
	var zipAddon ZipAddon
	if err != nil {
		return zipAddon, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return zipAddon, err
	}

	zipAddon, err = UnmarshalZipAddon(data)
	if err != nil {
		return zipAddon, err
	}

	return zipAddon, nil
}

func extractZipManifest(zipPath string) ([]byte, error) {
	jww.INFO.Println("Extracting Manifest")
	zReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zReader.Close()
	for _, src := range zReader.File {
		if src.Name == "manifest.json" {
			r, err := src.Open()
			if err != nil {
				return nil, err
			}
			defer r.Close()
			return ioutil.ReadAll(r)
		}
	}
	return nil, errors.Errorf("Failed to find manifest in pack")
}

func extractZipOverrides(zipPath string, overrides string, dest string) error {
	jww.INFO.Println("Extracting Overrides")
	zReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zReader.Close()

	var g errgroup.Group
	sem := make(chan struct{}, 5)
	for _, src := range zReader.File {
		src := src // create local for closure
		if filepath.HasPrefix(src.Name, overrides) {
			sem <- struct{}{}
			g.Go(func() error {
				defer func() {
					<-sem
				}()
				path := filepath.Join(dest, filepath.Clean(src.Name[len(overrides):]))
				dir := filepath.Dir(path)
				err := os.MkdirAll(dir, 0700)
				if err != nil {
					return err
				}

				if !strings.HasSuffix(src.Name, "/") {
					jww.DEBUG.Printf("%s -> %s", src.Name, path)
					file, err := os.Create(path)
					if err != nil {
						return err
					}
					defer file.Close()

					r, err := src.Open()
					if err != nil {
						return err
					}
					defer r.Close()

					_, err = io.Copy(file, r)
					if err != nil {
						return err
					}

				}
				return nil
			})
		}
	}

	return g.Wait()
}

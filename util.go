package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magiconair/properties"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
)

const FORGE_DL_URL = "https://files.minecraftforge.net/maven/net/minecraftforge/forge/"

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func folderExists(folder string) bool {
	info, err := os.Stat(folder)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func installForgeServer(forgeVersion string, path string) error {
	jww.INFO.Println("Downloading Forge Server")
	installerJar := "forge-" + forgeVersion + "-installer.jar"
	installerJarPath := filepath.Join(path, "installer.jar")
	if fileExists(installerJarPath) {
		err := os.Remove(installerJarPath)
		if err != nil {
			return errors.Wrap(err, "Failed removing existing Forge installer")
		}
	}

	forgeFile, err := os.Create(installerJarPath)
	if err != nil {
		return errors.Wrap(err, "Failed creating Forge installer")
	}
	defer os.Remove(forgeFile.Name())
	defer forgeFile.Close()

	forgeResp, err := http.Get(FORGE_DL_URL + forgeVersion + "/" + installerJar)
	if err != nil {
		return errors.Wrap(err, "Failed downloading Forge")
	}
	defer forgeResp.Body.Close()

	_, err = io.Copy(forgeFile, forgeResp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed writing Forge to disk")
	}

	javaPath, err := exec.LookPath("java")
	if err != nil {
		return errors.Wrap(err, "Failed finding Java")
	}

	cmdInstall := &exec.Cmd{
		Path:   javaPath,
		Args:   []string{javaPath, "-jar", installerJarPath, "--installServer", path},
		Stdout: jww.DEBUG.Writer(),
		Stderr: jww.ERROR.Writer(),
	}

	jww.INFO.Println("Installing Forge Server")
	err = cmdInstall.Run()
	if err != nil {
		return errors.Wrap(err, "Failed installing Forge")
	}

	src := filepath.Join(path, "forge-"+forgeVersion+".jar")
	dest := filepath.Join(path, "server.jar")

	return errors.Wrap(os.Rename(src, dest), "Failed renaming server jar")
}

func precreateUserDir(path string) error {
	err := os.MkdirAll(filepath.Join(path, "mods"), 0700)
	if err != nil {
		return errors.Wrap(err, "Failed creating mods directory")
	}
	err = os.MkdirAll(filepath.Join(path, "config"), 0700)
	if err != nil {
		return errors.Wrap(err, "Failed creating config directory")
	}
	file, err := os.Create(filepath.Join(path, "README.txt"))
	if err != nil {
		return errors.Wrap(err, "Failed creating README")
	}
	defer file.Close()

	_, err = file.WriteString("The contents of this folder should mirror the structure of the pack\n")
	if err != nil {
		return errors.Wrap(err, "Failed writting to README")
	}

	_, err = file.WriteString("Any user provided files will be copied from here into the pack\n")
	if err != nil {
		return errors.Wrap(err, "Failed writting to README")
	}

	return errors.Wrap(file.Sync(), "Failed flushing README to disk")
}

func writeVersionFile(version string, path string) error {
	ver := CPVersion{
		Version: version,
	}
	bytes, err := ver.Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal version file")
	}

	return errors.Wrap(ioutil.WriteFile(path, bytes, 0700), "Failed to write version file")
}

// compareVersion will return true if versions match, false if otherwise
func compareVersion(version string, path string) (bool, error) {
	if !fileExists(path) {
		return false, nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return false, errors.Wrap(err, "Failed to open version file")
	}
	ver, err := UnmarshalCPVersion(bytes)
	if err != nil {
		return false, errors.Wrap(err, "Failed to parse version file")
	}
	return ver.Version == version, nil
}

func updateServerPropsVersion(version string, path string) error {
	if fileExists(path) {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.Wrap(err, "Failed to read server properties")
		}

		props, err := properties.Load(bytes, properties.UTF8)
		if err != nil {
			return errors.Wrap(err, "Failed to parse server properties")
		}
		props.SetValue("motd", fmt.Sprintf("Version %s", version))
		if err != nil {
			return errors.Wrap(err, "Failed to set MOTD")
		}

		file, err := os.Create(path)
		if err != nil {
			return errors.Wrap(err, "Failed to open server properties")
		}
		defer file.Close()

		_, err = props.Write(file, properties.UTF8)
		if err != nil {
			return errors.Wrap(err, "Failed to write server properties")
		}
	}
	return nil
}

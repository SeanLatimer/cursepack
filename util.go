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
			return err
		}
	}

	forgeFile, err := os.Create(installerJarPath)
	if err != nil {
		return err
	}
	defer os.Remove(forgeFile.Name())
	defer forgeFile.Close()

	forgeResp, err := http.Get(FORGE_DL_URL + forgeVersion + "/" + installerJar)
	if err != nil {
		return err
	}
	defer forgeResp.Body.Close()

	_, err = io.Copy(forgeFile, forgeResp.Body)
	if err != nil {
		return err
	}

	javaPath, err := exec.LookPath("java")
	if err != nil {
		return err
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
		return err
	}

	src := filepath.Join(path, "forge-"+forgeVersion+".jar")
	dest := filepath.Join(path, "server.jar")

	return os.Rename(src, dest)
}

func precreateUserDir(path string) error {
	err := os.MkdirAll(filepath.Join(path, "mods"), 0700)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(path, "config"), 0700)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(path, "README.txt"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("The contents of this folder should mirror the structure of the pack\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString("Any user provided files will be copied from here into the pack\n")
	if err != nil {
		return err
	}

	return file.Sync()
}

func writeVersionFile(version string, path string) error {
	ver := CPVersion{
		Version: version,
	}
	bytes, err := ver.Marshal()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0700)
}

// compareVersion will return true if versions match, false if otherwise
func compareVersion(version string, path string) (bool, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	ver, err := UnmarshalCPVersion(bytes)
	if err != nil {
		return false, err
	}
	return ver.Version == version, nil
}

func updateServerPropsVersion(version string, path string) error {
	if fileExists(path) {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		props, err := properties.Load(bytes, properties.UTF8)
		props.SetValue("motd", fmt.Sprintf("Version %s", version))

		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = props.Write(file, properties.UTF8)
		if err != nil {
			return err
		}
	}
	return nil
}

package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

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

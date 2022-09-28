package up2date

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/kardianos/osext"
)

// DoUpdate : Update binary
func DoUpdate(fileType, filenameBIN string) (err error) {

	var url string
	switch fileType {
	case "bin":
		url = Updater.Response.UpdateBIN
	case "zip":
		url = Updater.Response.UpdateZIP
	}

	switch runtime.GOOS {
	case "windows":
		filenameBIN = filenameBIN + ".exe"
	}

	if len(url) > 0 {
		log.Println("["+strings.ToUpper(fileType)+"]", "New version ("+Updater.Name+"):", Updater.Response.Version)

		// Download new binary
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		log.Println("["+strings.ToUpper(fileType)+"]", "Download new version...")

		if resp.StatusCode != http.StatusOK {
			log.Println("["+strings.ToUpper(fileType)+"]", "Download new version...OK")
			return fmt.Errorf("bad status: %s", resp.Status)
		}

		// Change binary filename to .filename
		binary, err := osext.Executable()
		var filename = getFilenameFromPath(binary)
		var path = getPlatformPath(binary)
		var oldBinary = path + "_old_" + filename
		var newBinary = binary

		// ZIP
		var tmpFolder = path + "tmp"
		var tmpFile = tmpFolder + string(os.PathSeparator) + filenameBIN

		//fmt.Println(binary, path+"."+filename)
		os.Rename(newBinary, oldBinary)

		// Save the new binary with the old file name
		out, err := os.Create(binary)
		if err != nil {
			restorOldBinary(oldBinary, newBinary)
			return err
		}
		defer out.Close()

		// Write the body to file

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			restorOldBinary(oldBinary, newBinary)
			return err
		}

		// Update as a ZIP file
		if fileType == "zip" {

			log.Println("["+strings.ToUpper(fileType)+"]", "Update file:", filenameBIN)
			log.Println("["+strings.ToUpper(fileType)+"]", "Unzip ZIP file...")
			err = extractZIP(binary, tmpFolder)

			binary = newBinary

			if err != nil {

				log.Println("["+strings.ToUpper(fileType)+"]", "Unzip ZIP file...ERROR")

				restorOldBinary(oldBinary, newBinary)

				return err
			} else {

				log.Println("["+strings.ToUpper(fileType)+"]", "Unzip ZIP file...OK")
				log.Println("["+strings.ToUpper(fileType)+"]", "Copy binary file...")

				err = copyFile(tmpFile, binary)
				if err == nil {
					log.Println("["+strings.ToUpper(fileType)+"]", "Copy binary file...OK")
				} else {

					log.Println("["+strings.ToUpper(fileType)+"]", "Copy binary file...ERROR")
					restorOldBinary(oldBinary, newBinary)

					return err
				}

				os.RemoveAll(tmpFolder)
			}

		}

		// Set the permission
		err = os.Chmod(binary, 0755)

		// Close the new file !Windows
		out.Close()

		log.Println("["+strings.ToUpper(fileType)+"]", "Update Successful")

		// Restart binary (Windows)
		if runtime.GOOS == "windows" {

			bin, err := os.Executable()

			if err != nil {
				restorOldBinary(oldBinary, newBinary)
				return err
			}

			var pid = os.Getpid()
			var process, _ = os.FindProcess(pid)

			if proc, err := start(bin); err == nil {

				os.RemoveAll(oldBinary)
				process.Kill()
				proc.Wait()

			} else {
				restorOldBinary(oldBinary, newBinary)
			}

		} else {

			// Restart binary (Linux and UNIX)
			file, _ := osext.Executable()
			os.RemoveAll(oldBinary)
			err = syscall.Exec(file, os.Args, os.Environ())
			if err != nil {
				restorOldBinary(oldBinary, newBinary)
				log.Fatal(err)
				return err
			}

		}

	}

	return
}

func start(args ...string) (p *os.Process, err error) {

	if args[0], err = exec.LookPath(args[0]); err == nil {
		//fmt.Println(args[0])
		var procAttr os.ProcAttr
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
		p, err := os.StartProcess(args[0], args, &procAttr)

		if err == nil {
			return p, nil
		}

	}

	return nil, err
}

func restorOldBinary(oldBinary, newBinary string) {
	os.RemoveAll(newBinary)
	os.Rename(oldBinary, newBinary)
}

func getFilenameFromPath(path string) string {

	file := filepath.Base(path)

	return file
}

func getPlatformPath(path string) string {

	var newPath = filepath.Dir(path) + string(os.PathSeparator)

	return newPath
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func extractZIP(archive, target string) (err error) {

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {

		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}

	}

	return
}

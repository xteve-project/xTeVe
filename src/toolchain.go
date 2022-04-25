package src

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/samber/lo"
)

// --- System Tools ---

// Checks whether the Folder exists, if not, the Folder is created
func checkFolder(path string) (err error) {

	var debug string
	_, err = os.Stat(filepath.Dir(path))

	if os.IsNotExist(err) {
		// Folder does not exist, will now be created

		err = os.MkdirAll(getPlatformPath(path), 0755)
		if err == nil {

			debug = fmt.Sprintf("Create Folder:%s", path)
			showDebug(debug, 1)

		} else {
			return err
		}

		return nil
	}

	return nil
}

// Checks whether the File exists in the Filesystem
func checkFile(filename string) (err error) {

	var file = getPlatformFile(filename)

	if _, err = os.Stat(file); os.IsNotExist(err) {
		return err
	}

	fi, err := os.Stat(file)
	if err != nil {
		return err
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		err = fmt.Errorf("%s: %s", file, getErrMsg(1072))
	case mode.IsRegular():
		break
	}

	return
}

func allFilesExist(list ...string) bool {
	for _, f := range list {
		if err := checkFile(f); err != nil {
			return false
		}
	}
	return true
}

// GetUserHomeDirectory : User Home Directory
func GetUserHomeDirectory() (userHomeDirectory string) {

	usr, err := user.Current()

	if err != nil {

		for _, name := range []string{"HOME", "USERPROFILE"} {

			if dir := os.Getenv(name); dir != "" {
				userHomeDirectory = dir
				break
			}

		}

	} else {
		userHomeDirectory = usr.HomeDir
	}

	return
}

// Checks File Permissions
func checkFilePermission(dir string) (err error) {

	var filename = dir + "permission.test"

	err = ioutil.WriteFile(filename, []byte(""), 0644)
	if err == nil {
		err = os.RemoveAll(filename)
	}

	return
}

// Generate folder path for the running OS
func getPlatformPath(path string) string {
	return filepath.Dir(path) + string(os.PathSeparator)
}

// getDefaultTempDir returns default temporary folder path + application name, e.g.: "/tmp/xteve/" or %Tmp%\xteve.
//
// Function assumes default OS temporary folder exists and writable. 
func getDefaultTempDir() string {
	return os.TempDir() + string(os.PathSeparator) + System.AppName + string(os.PathSeparator)
}

// getValidTempDir returns standartized temporary folder <path> with trailing path separator:
//
// Slashes will be replaced with OS specific ones and duplicated slashes removed.
//
// On Windows, "/tmp" will be replaced with expanded system environment variable %Tmp%.
func getValidTempDir(path string) string {
	if runtime.GOOS == "windows" {
		if strings.HasPrefix(path, "/tmp") {
			path = strings.Replace(path, "/tmp", os.TempDir(), 1)
		}
	}
	path = filepath.Clean(path)
	path = path + string(os.PathSeparator)

	err := checkFolder(path)
	if err == nil {
		err = checkFilePermission(path)
	}

	if err != nil {
		ShowError(err, 1015)
		path = getDefaultTempDir()
	}

	return path
}

// Generate File Path for the running OS
func getPlatformFile(filename string) (osFilePath string) {

	path, file := filepath.Split(filename)
	var newPath = filepath.Dir(path)
	osFilePath = newPath + string(os.PathSeparator) + file

	return
}

// Output Filenames from the File Path
func getFilenameFromPath(path string) (file string) {
	return filepath.Base(path)
}

// Searches for a File in the OS
func searchFileInOS(file string) (path string) {

	switch runtime.GOOS {

	case "linux", "darwin", "freebsd":
		var args = file
		var cmd = exec.Command("which", strings.Split(args, " ")...)

		out, err := cmd.CombinedOutput()
		if err == nil {

			var slice = strings.Split(strings.Replace(string(out), "\r\n", "\n", -1), "\n")

			if len(slice) > 0 {
				path = strings.Trim(slice[0], "\r\n")
			}

		}

	default:
		return

	}

	return
}

//
func removeChildItems(dir string) error {

	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}

	for _, file := range files {

		err = os.RemoveAll(file)
		if err != nil {
			return err
		}

	}

	return nil
}

// JSON
func mapToJSON(tmpMap interface{}) string {

	jsonString, err := json.MarshalIndent(tmpMap, "", "  ")
	if err != nil {
		return "{}"
	}

	return string(jsonString)
}

func jsonToMap(content string) map[string]interface{} {

	var tmpMap = make(map[string]interface{})
	json.Unmarshal([]byte(content), &tmpMap)

	return (tmpMap)
}

func jsonToInterface(content string) (tmpMap interface{}, err error) {

	err = json.Unmarshal([]byte(content), &tmpMap)
	return

}

func saveMapToJSONFile(file string, tmpMap interface{}) error {

	var filename = getPlatformFile(file)
	jsonString, err := json.MarshalIndent(tmpMap, "", "  ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, []byte(jsonString), 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadJSONFileToMap(file string) (tmpMap map[string]interface{}, err error) {

	f, err := os.Open(getPlatformFile(file))
	defer f.Close()

	content, err := ioutil.ReadAll(f)

	if err == nil {
		err = json.Unmarshal([]byte(content), &tmpMap)
	}

	f.Close()

	return
}

// Binary
func readByteFromFile(file string) (content []byte, err error) {

	f, err := os.Open(getPlatformFile(file))
	defer f.Close()

	content, err = ioutil.ReadAll(f)
	f.Close()

	return
}

func writeByteToFile(file string, data []byte) (err error) {

	var filename = getPlatformFile(file)
	err = ioutil.WriteFile(filename, data, 0644)

	return
}

func readStringFromFile(file string) (str string, err error) {

	var content []byte
	var filename = getPlatformFile(file)

	err = checkFile(filename)
	if err != nil {
		return
	}

	content, err = ioutil.ReadFile(filename)
	if err != nil {
		ShowError(err, 0)
		return
	}

	str = string(content)

	return
}

// Network
func resolveHostIP() (err error) {

	netInterfaceAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, netInterfaceAddress := range netInterfaceAddresses {

		networkIP, ok := netInterfaceAddress.(*net.IPNet)
		System.IPAddressesList = append(System.IPAddressesList, networkIP.IP.String())

		if ok {

			var ip = networkIP.IP.String()

			if networkIP.IP.To4() != nil {

				System.IPAddressesV4 = append(System.IPAddressesV4, ip)
				System.IPAddressesV4Raw = append(System.IPAddressesV4Raw, networkIP.IP)

				if !networkIP.IP.IsLoopback() && ip[0:7] != "169.254" {
					System.IPAddressesV4Host = append(System.IPAddressesV4Host, ip)
				}

			} else {
				System.IPAddressesV6 = append(System.IPAddressesV6, ip)
			}

		}

	}

	// If IP previously set in settings (including the default, empty) is not available anymore
	if lo.Contains(System.IPAddressesV4Host, Settings.HostIP) == false {
		Settings.HostIP = System.IPAddressesV4Host[0]
	}

	if len(Settings.HostIP) == 0 {

		switch len(System.IPAddressesV4) {

		case 0:
			if len(System.IPAddressesV6) > 0 {
				Settings.HostIP = System.IPAddressesV6[0]
			}

		default:
			Settings.HostIP = System.IPAddressesV4[0]

		}

	}

	System.Hostname, err = os.Hostname()
	if err != nil {
		return
	}

	return
}

// Miscellaneous
func randomString(n int) string {

	const alphanum = "AB1CD2EF3GH4IJ5KL6MN7OP8QR9ST0UVWXYZ"

	var bytes = make([]byte, n)

	rand.Read(bytes)

	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}

	return string(bytes)
}

func parseTemplate(content string, tmpMap map[string]interface{}) (result string) {

	t := template.Must(template.New("template").Parse(content))

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, tmpMap); err != nil {
		ShowError(err, 0)
	}
	result = tpl.String()

	return
}

func indexOfString(element string, data []string) int {

	for k, v := range data {
		if element == v {
			return k
		}
	}

	return -1
}

func indexOfFloat64(element float64, data []float64) int {

	for k, v := range data {
		if element == v {
			return (k)
		}
	}

	return -1
}

func getMD5(str string) string {

	md5Hasher := md5.New()
	md5Hasher.Write([]byte(str))

	return hex.EncodeToString(md5Hasher.Sum(nil))
}

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
)

// --- System Tools ---

// Prüft ob der Ordner existiert, falls nicht, wir der Ordner erstellt
func checkFolder(path string) (err error) {

	var debug string
	_, err = os.Stat(filepath.Dir(path))

	if os.IsNotExist(err) {
		// Ordner existiert nicht, wird jetzt erstellt

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

// Prüft ob die Datei im Dateisystem existiert
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

// GetUserHomeDirectory : Benutzer Homer Verzeichnis
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

// Prüft Dateiberechtigung
func checkFilePermission(dir string) (err error) {

	var filename = dir + "permission.test"

	err = ioutil.WriteFile(filename, []byte(""), 0644)
	if err == nil {
		err = os.RemoveAll(filename)
	}

	return
}

// Ordnerpfad für das laufende OS generieren
func getPlatformPath(path string) string {
	return filepath.Dir(path) + string(os.PathSeparator)
}

// Dateipfad für das laufende OS generieren
func getPlatformFile(filename string) (osFilePath string) {

	path, file := filepath.Split(filename)
	var newPath = filepath.Dir(path)
	osFilePath = newPath + string(os.PathSeparator) + file

	return
}

// Dateinamen aus dem Dateipfad ausgeben
func getFilenameFromPath(path string) (file string) {
	return filepath.Base(path)
}

// Nicht mehr verwendete Systemdaten löschen
func removeOldSystemData() {
	// Temporären Ordner löschen
	os.RemoveAll(System.Folder.Temp)
}

// Sucht eine Datei im OS
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

func jsonToMapInt64(content string) map[int64]interface{} {

	var tmpMap = make(map[int64]interface{})
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

// Netzwerk
func resolveHostIP() (err error) {

	netInterfaceAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return
	}
	defer conn.Close()
	System.IPAddress = conn.LocalAddr().(*net.UDPAddr).IP.String()

	for _, netInterfaceAddress := range netInterfaceAddresses {

		networkIP, ok := netInterfaceAddress.(*net.IPNet)
		System.IPAddressesList = append(System.IPAddressesList, networkIP.IP.String())

		if ok {

			var ip = networkIP.IP.String()

			if networkIP.IP.To4() != nil {
				System.IPAddressesV4 = append(System.IPAddressesV4, ip)
			} else {
				System.IPAddressesV6 = append(System.IPAddressesV6, ip)
			}

		}

	}

	if len(System.IPAddress) == 0 {

		switch len(System.IPAddressesV4) {

		case 0:
			if len(System.IPAddressesV6) > 0 {
				System.IPAddress = System.IPAddressesV6[0]
			}

		default:
			System.IPAddress = System.IPAddressesV4[0]

		}

	}

	System.Hostname, err = os.Hostname()
	if err != nil {
		return
	}

	return
}

// Sonstiges
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

func indexOfInt(element int, data []int) int {

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

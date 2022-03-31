package main

import (
        "flag"
        "fmt"
        "github.com/go-resty/resty/v2"
        "github.com/satori/go.uuid"
        "io"
        "io/ioutil"
        "log"
        "net"
        "net/http"
        "os"
        "os/exec"
        "runtime"
        "strings"
        // "testing"
        // "github.com/stretchr/testify/assert"
        // "net/http"
        "crypto/sha256"
        "encoding/hex"
        "encoding/json"
        "time"
        // "syscall"
        // "path/filepath"
)

// https://www.ontestautomation.com/an-introduction-to-rest-api-testing-in-go-with-resty/

var INTERNET_IP_API_URL = "https://api.ipify.org"
var operatingSystem string = "Undetermined"
var RunningProcessPID int
var CURRENT_VERSION string = "v0.1.0"
var scriptFileSuffix string = ""
var shell string = ""
var intranetIPv4 = ""
var internetIPv4 = ""
var hostname = ""
var DNS_DOMAIN = "example.com"
var PROXY_TOKEN = ""
var dns string = ""

// var uuidFlag string = ""
// var forceFlag bool = false

type ipapicomIPaddr struct {
        Query string
}

// url := "https://api64.ipify.org"

type Message struct {
        ipversion string
        ipaddr    string
}

type OpenResponse struct {
        IPv4 string `json:"ip"`
}

type LocationResponse struct {
        CountryAbbreviation string                   `json:"country abbreviation"`
        Country             string                   `json:"country"`
        Places              []LocationResponsePlaces `json:"places"`
        // Places              []byte `json:"places"`
}

type LocationResponsePlaces struct {
        PlaceName string `json:"place name"`
}

func getOperatingSystem() string {
        // os := runtime.GOOS
        operatingSystem = runtime.GOOS
        // fmt.Printf("s %s", operatingSystem)
        switch operatingSystem {
        case "windows":
                return "windows"
        case "darwin":
                return "macos"
        case "linux":
                return "linux"
        default:
                return operatingSystem
                // return "aa"
        }
}

func setScriptFileSuffix() {
        if operatingSystem == "windows" {
                scriptFileSuffix = "ps1"
                shell = "powershell"
        } else if operatingSystem == "linux" {
                scriptFileSuffix = "sh"
                shell = "/bin/bash"
        } else {
                scriptFileSuffix = ""
        }
}

func getIntranetIPv4() net.IP {
        conn, err := net.Dial("udp", "8.8.8.8:80")
        if err != nil {
                log.Fatal(err)
        }
        defer conn.Close()

        localAddr := conn.LocalAddr().(*net.UDPAddr)

        // fmt.Println(localAddr.IP)
        intranetIPv4 = localAddr.IP.String()
        return localAddr.IP
}

func getHostname() string {
        name, err := os.Hostname()
        if err != nil {
                panic(err)
        }
        hostname = name
        dns = hostname + "." + DNS_DOMAIN
        return hostname
}

func arrayContains(s []string, e string) bool {
        for _, a := range s {
                if a == e {
                        return true
                }
        }
        return false
}

func getHostIpaddrPublic() string {
        req, err := http.Get("http://ip-api.com/json/")
        // req, err := http.Get("https://api64.ipify.org?format=json")
        if err != nil {
                return err.Error()
        }
        defer req.Body.Close()

        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
                return err.Error()
        }

        var ip ipapicomIPaddr
        json.Unmarshal(body, &ip)

        return ip.Query
}

func usage() {
        // fmt.Println(len(os.Args), os.Args)
        // msg := fmt.Sprintf("Usage: %s uuid", os.Args[0])
        // usage_msg := "Usage: iutil <agentUUID>"
        // fmt.Println(usage_msg)
        usage_msg := "Usage: iutil -h"
        fmt.Println(usage_msg)
        // example_msg := "Example: iutil 13c023b2-ed03-11eb-b237-00163ebb406c"
        // fmt.Println(example_msg)
}

func isValidUUID(agentUUID string) bool {
        _, err := uuid.FromString(agentUUID)
        if err != nil {
                // fmt.Printf("Something went wrong: %s", err)
                return false
        }
        // fmt.Printf("Successfully parsed: %s", u2)
        return true
        // NULL value
}

func getInternetIPv4() string {
        IPv4Addr := getURLBodyAsString(INTERNET_IP_API_URL)
        if isValidIPAddr(IPv4Addr) {
                internetIPv4 = IPv4Addr
                return IPv4Addr
        } else {
                return ""
        }
}

func isValidIPAddr(ip string) bool {
        if net.ParseIP(ip) == nil {
                // fmt.Printf("IP Address: %s - Invalid\n", ip)
                return false
        } else {
                // fmt.Printf("IP Address: %s - Valid\n", ip)
                return true
        }
}

func getScript1(tmpFile string, currentFile string) string {
        script := fmt.Sprintf(""+
                "while( $true ){"+
                "Get-Process x -ErrorAction SilentlyContinue;"+
                "if($pid){break}; write-output 'hello'; start-sleep -s 5"+
                "};"+
                "write-output 'start';"+
                "cp %s %s; rm %s;"+
                "write-output 'end'",
                tmpFile, currentFile, tmpFile)
        return script
}

func getCurrentFilePath() string {
        //dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
        // exPath := filepath.Dir(ex) + "\\" + ex
        ex, err := os.Executable()
        if err != nil {
                panic(err)
        }
        return ex
}

func getFileHash(file string) string {
        // f, err := os.Open("file.txt")
        f, err := os.Open(file)
        if err != nil {
                log.Fatal(err)
        }
        defer f.Close()

        h := sha256.New()
        if _, err := io.Copy(h, f); err != nil {
                log.Fatal(err)
                fmt.Printf("Error")
        }

        sha := hex.EncodeToString(h.Sum(nil))
        return sha
}

func getTmpFilePath(suffix string) string {
        tmpDir := os.TempDir()
        // tmpFile := uuid.Must(uuid.NewV4())
        tmpFile := uuid.NewV4()
        tmpFilePath := fmt.Sprintf("%s/%s.insights.tmp.%s", tmpDir, tmpFile, suffix)
        fmt.Printf("Tempfile %s", tmpFilePath)
        return tmpFilePath
}

func test() {

        client := resty.New()

        resp, _ := client.R().Get("http://api.zippopotam.us/us/90210")
        // {"post code": "90210", "country": "United States", "country abbreviation": "US", "places": [{"place name": "Beverly Hills", "longitude": "-118.4065", "state": "California", "state abbreviation": "CA", "latitude": "34.0901"}]}

        myResponse := LocationResponse{}

        err := json.Unmarshal(resp.Body(), &myResponse)

        if err != nil {
                fmt.Println(err)
                return
        }

        fmt.Printf("United States %s\n", myResponse.Country)
        fmt.Printf("Country abbreviation %s\n", myResponse.CountryAbbreviation)
        fmt.Printf("My places 0 %s\n", myResponse.Places[0])
        // place := fmt.Sprintf("Foo Says: %s", myResponse.Places[0])
        // fmt.Printf("Places sprintf: %s",  place)
        fmt.Printf("%T\n", myResponse.Places[0])
        e, err := json.Marshal(myResponse.Places[0])
        fmt.Printf("Marshalled String %s\n", e)
        fmt.Printf("Places %b\n", myResponse.Places)
        fmt.Printf("%+v\n", myResponse)
}

func test3() {
        // Create a resty client
        client := resty.New()

        // POST JSON string
        // No need to set content type, if you have client level setting
        resp, err := client.R().
                //SetHeader("Content-Type", "application/json").
                // SetBody(`{"username":"testuser", "password":"testpass"}`).
                // SetResult(&OpenResponse{}). // or SetResult(&AuthSuccess{}).
                Get("https://api.ipify.org?format=json")

        fmt.Println(resp, err)
        // fmt.Printf("ip %s\n", OpenResponse.ip)
        myResponse := OpenResponse{}

        err = json.Unmarshal(resp.Body(), &myResponse)
        fmt.Printf("ip %s\n", myResponse.IPv4)
}

// func DownloadFile(filepath string, url string) error {
func DownloadFile(url string, filepath string) error {

        // Get the data
        resp, err := http.Get(url)
        if err != nil {
                return err
        }
        defer resp.Body.Close()

        // Create the file
        out, err := os.Create(filepath)
        if err != nil {
                return err
        }
        defer out.Close()

        // Write the body to file
        _, err = io.Copy(out, resp.Body)
        return err
}

func getFile() {
        // fileUrl := "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"
        fileUrl := "https://www.gimp.org/images/frontpage/wilber-big.png"
        fileLocal := "logo.png"
        err := DownloadFile(fileUrl, fileLocal)
        if err != nil {
                panic(err)
        }
        fileHash := getFileHash(fileLocal)
        fmt.Println("Downloaded: " + fileUrl)
        fmt.Println("Hash: " + fileHash)
        filePath := getCurrentFilePath()
        // filePath :=  getCurrentFilePath
        fmt.Printf("Path: %s", filePath)
}

func getURLBodyAsString(url string) string {

        resp, err := http.Get(url)
        if err != nil {
                log.Fatal(err)
        }
        defer resp.Body.Close()
        rspBody, err := ioutil.ReadAll(resp.Body)
        sbody := strings.TrimSpace(string(rspBody))
        return sbody
}

func getFilePathRunning() string {
        ex, err := os.Executable()
        if err != nil {
                panic(err)
        }
        return ex
}

func replaceFiles(currentFile string, replacementFile string) {
        src, dst := currentFile, replacementFile
        data, err := ioutil.ReadFile(src)
        if err != nil {
                panic(err)
        }
        if err := ioutil.WriteFile(dst, data, 0744); err != nil {
                panic(err)
        }

        if err := os.Rename("./exa_", "./exa"); err != nil {
                panic(err)
        }
}

func getPid() {
        RunningProcessPID = os.Getpid()
}

func downloadScriptFromUrlAndExecute(url string) {
        tmpFile := getTmpFilePath(scriptFileSuffix)
        err := DownloadFile(url, tmpFile)
        if err != nil {
                panic(err)
        }
        fmt.Printf("%s to\n", tmpFile)
        runScriptFile(tmpFile)
        defer os.Remove(tmpFile)
}

func installViaScripts() {
        downloadScriptFromUrlAndExecute("https://raw.githubusercontent.com/jeremybusk/sandboxfiles/main/example.ps1")
}

func updateApp() {
        // fmt.Printf("os %s\n", operatingSystem)
        // var operatingSystemgetOperatignSystem()
        fmt.Printf("os %s\n", operatingSystem)
        getPid()
        fmt.Printf("pid: %d\n", RunningProcessPID)
        replacementFileUrl := "https://github.com/jeremybusk/sandboxfiles/blob/main/example.exe?raw=true"
        replacementFileUrlSHA256 := "https://raw.githubusercontent.com/jeremybusk/sandboxfiles/main/example.exe.sha256"
        replacementFileHash := getURLBodyAsString(replacementFileUrlSHA256)
        tmpFile := getTmpFilePath("")
        // defer os.Remove(tmpFile)

        fmt.Println("Created File: " + tmpFile)

        currentFilePath := getFilePathRunning()
        // fmt.Printf("s %s", script1)
        // os.Exit(3)
        fmt.Printf("Replacement file URI path %s\n", replacementFileUrl)
        fmt.Printf("Replacement file hash %s\n", replacementFileHash)
        fmt.Printf("Current file path %s\n", currentFilePath)
        currentFileHash := getFileHash(currentFilePath)
        fmt.Printf("Current version: %s hash: %s\n", CURRENT_VERSION, currentFileHash)
        if currentFileHash != replacementFileHash {
                fmt.Printf("New version available. Updating version %s to latest now.\n", CURRENT_VERSION)
                err := DownloadFile(replacementFileUrl, tmpFile)
                if err != nil {
                        panic(err)
                }
                fmt.Printf("%s to %s\n", tmpFile, currentFilePath)
                // replaceCurrentFile := fmt.Sprintf("Start-Sleep -s 5; write-output 'start'; cp %s %s; rm %s; write-output 'end'", tmpFile, currentFilePath, tmpFile)
                // replaceCurrentFile := fmt.Sprintf("while($true){ Get-Process x -ErrorAction SilentlyContinue; if($pid){break}; write-output 'hello'; start-sleep -s 5}; write-output 'start'; cp %s %s; rm %s; write-output 'end'", tmpFile, currentFilePath, tmpFile)
                script1 := getScript1(tmpFile, currentFilePath)
                // replaceCurrentFileScript := getScript1()
                runBackgroundCmd(script1)
        }
        return
}

func waitUntilPIDNotExist(pid int) {
        for i := 0; i < 4; i++ {
                fmt.Println("Check for process")
                time.Sleep(10 * time.Second)
                process, err := os.FindProcess(int(pid))
                if err != nil {
                        fmt.Printf("Failed to find process: %s\n", err)
                } else {
                        // err := process.Signal(syscall.Signal(0))
                        fmt.Printf("process.%s Signal on pid %d returned: %v\n", process, pid, err)
                }

        }
}

func validShellCmd(shellName string) {
        s := []string{"powershell", "bash"}
        if !arrayContains(s, shellName) {
                fmt.Printf("E: Invalid shell %s! \n", shellName)
                os.Exit(1)
        }
}

func runCmd(cmdargs string) {
        cmd := exec.Command("powershell", "-c", cmdargs)
        out, err := cmd.Output()
        if err != nil {
                fmt.Println(err)
                fmt.Println(out)
        }
        fmt.Println(string(out))
}

func shellcmd(name string, args []string) {
        cmd := exec.Command(name, args...)
        out, err := cmd.Output()
        if err != nil {
                fmt.Println(err)
        }
        fmt.Println(string(out))
}

func runScriptFile(file string) {
        cmd := exec.Command(shell, file)
        out, err := cmd.Output()
        if err != nil {
                fmt.Println(err)
                fmt.Println(out)
        }
        fmt.Println(string(out))
}

func runBackgroundCmd(cmdargs string) {
        cmd := exec.Command("powershell", "-c", cmdargs)
        cmd.Stdout = os.Stdout
        err := cmd.Start()
        if err != nil {
                log.Fatal(err)
        }
        // log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
}

func getAttributes() {
        getOperatingSystem()
        setScriptFileSuffix()
        getIntranetIPv4()
        getInternetIPv4()
        getHostname()
}
func displayAttributes() {
        fmt.Printf("Intranet IPv4: %s\n", intranetIPv4)
        fmt.Printf("Internet IPv4: %s\n", internetIPv4)
        fmt.Println("Hostname: " + hostname)
        fmt.Println("dns: " + dns)
        fmt.Println("OS: " + operatingSystem)
        // fmt.Println("uuid:", *uuidFlag)
        // fmt.Println("force:", *forceFlag)
        // flag.Parse()
}

func main() {
        // test3()
        // test()
        // os.Exit(1)
        //installViaScripts()

        var uuidFlag = flag.String("uuid", "", "V4 uuid like xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx")
        shellcmdFlag := flag.String("shellcmd", "", "command")
        updateFlag := flag.Bool("update", false, "Update client")
        forceFlag := flag.Bool("force", false, "Danger! This will overwrite existing contents if they exist.")
        displayFlag := flag.Bool("display", false, "Danger! This will overwrite existing contents if they exist.")
        flag.Parse()
        if *uuidFlag == "" {
                usage()
                os.Exit(1)
        }
        if !isValidUUID(*uuidFlag) {
                fmt.Printf("E: -uuid %q is invalid! \n", *uuidFlag)
                os.Exit(1)
        }
        fmt.Printf("d: %s", *displayFlag)
        // getAttributes()
        // if *displayFlag == true {
        //      displayAttributes()
        //      fmt.Println("poo")
        //      os.Exit(1)
        //}
        if *updateFlag == true {
                updateApp()
        }
        _ = forceFlag

        // fmt.Println("tail:", flag.Args())
        // internetIPv4 := getHostIpaddrPublic()
        // hostname := getHostname()
        // fqdn := strings.ToLower(hostname) + "." + DNS_DOMAIN
        if *shellcmdFlag != "" {
                shellCmdArray := strings.Fields(*shellcmdFlag)
                shellName := shellCmdArray[0]
                args := shellCmdArray[1:]
                validShellCmd(shellName)
                shellcmd(shellName, args)
        }
        displayAttributes()

}

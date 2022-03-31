package main

import (
        "fmt"
        "github.com/go-resty/resty/v2"
        "github.com/satori/go.uuid"
        "io"
        "io/ioutil"
        "log"
        "net/http"
        "os"
        "os/exec"
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
type LocationResponse struct {
        CountryAbbreviation string                   `json:"country abbreviation"`
        Country             string                   `json:"country"`
        Places              []LocationResponsePlaces `json:"places"`
        // Places              []byte `json:"places"`
}

type LocationResponsePlaces struct {
        PlaceName string `json:"place name"`
}

func getScript1(tmpFile string, currentFile string)string{
  script := fmt.Sprintf("" +
  "while( $true ){"+
  "Get-Process x -ErrorAction SilentlyContinue;" +
  "if($pid){break}; write-output 'hello'; start-sleep -s 5" +
  "};" +
  "write-output 'start';" +
  "cp %s %s; rm %s;" +
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

func getTmpFilePath() string {
        tmpDir := os.TempDir()
        // tmpFile := uuid.Must(uuid.NewV4())
        tmpFile := uuid.NewV4()
        tmpFilePath := fmt.Sprintf("%s/%s.insights.tmp", tmpDir, tmpFile)
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

func updateFile() {
        RunningProcessPID := os.Getpid()
        fmt.Printf("pid: %d\n", RunningProcessPID)
        currentVersion := "v0.1.0"
        replacementFileUrl := "https://github.com/jeremybusk/sandboxfiles/blob/main/example.exe?raw=true"
        replacementFileUrlSHA256 := "https://raw.githubusercontent.com/jeremybusk/sandboxfiles/main/example.exe.sha256"
        replacementFileHash := getURLBodyAsString(replacementFileUrlSHA256)
        tmpFile := `"C:\tmp\a"`
        tmpFile = getTmpFilePath()
        // defer os.Remove(tmpFile)

        fmt.Println("Created File: " + tmpFile)

        currentFilePath := getFilePathRunning()
        // fmt.Printf("s %s", script1)
        // os.Exit(3)
        fmt.Printf("Replacement file URI path %s\n", replacementFileUrl)
        fmt.Printf("Replacement file hash %s\n", replacementFileHash)
        fmt.Printf("Current file path %s\n", currentFilePath)
        currentFileHash := getFileHash(currentFilePath)
        fmt.Printf("Current version: %s hash: %s\n", currentVersion, currentFileHash)
        if currentFileHash != replacementFileHash {
                fmt.Printf("New version available. Updating version %s to latest now.\n", currentVersion)
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

func runCmd(cmdargs string) {
        cmd := exec.Command("powershell", "-c", cmdargs)
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

func main() {
        updateFile()
}

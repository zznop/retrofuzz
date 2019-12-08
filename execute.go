package main;

import (
    "strings"
    "syscall"
    "fmt"
    "os"
    "sync"
    "path/filepath"
)

// LogCrashFile creates the crash directory if it doesn't exist, and copies the file to it
func LogCrashFile(outputFile string) {
    outputFileDir := filepath.Dir(outputFile)
    crashFileDir := filepath.Join(outputFileDir, "crashes")
    if _, err := os.Stat(crashFileDir); os.IsNotExist(err) {
        os.Mkdir(crashFileDir, os.ModeDir)
    }

    os.Rename(outputFile, filepath.Join(crashFileDir, filepath.Base(outputFile)))
}

// ExecuteCommandLine spawns a target process against the fuzzed input
func ExecuteCommandLine(outputFile string, commandLine string, timeout uint32, wg *sync.WaitGroup) {
    var startupInfo syscall.StartupInfo
    var processInfo syscall.ProcessInformation
    var status uint32
    var err error

    commandLineStr := strings.Replace(commandLine, "@@", "\"" + outputFile + "\"", 1)
    argv := syscall.StringToUTF16Ptr(commandLineStr)

    err = syscall.CreateProcess(nil, argv, nil, nil, true, 0, nil, nil, &startupInfo, &processInfo)
    if err != nil {
        fmt.Println("CreateProcess()\n")
        goto DONE
    }

    status, err = syscall.WaitForSingleObject(processInfo.Process, timeout)
    if err != nil {
        fmt.Println("WaitForSingleObject()\n");
        goto DONE
    }

    if status != syscall.WAIT_TIMEOUT {
        LogCrashFile(outputFile)
    } else {
        syscall.TerminateProcess(processInfo.Process, 0)
        os.Remove(outputFile) 
    }

DONE:
    wg.Done()
}

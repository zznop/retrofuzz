package main

import (
    "fmt"
    "flag"
    "os"
    "io/ioutil"
    "time"
)

// FixupCallback is a function pointer type used by fixup routines specific to a file format
type FixupCallback func(context *RFContext)

// RFContext is used as the main context and stores user input and information on the current state
type RFContext struct {
    startTime time.Time
    inputFile string
    inputData []byte
    outputDir string
    commandLine string
    numWorkers int
    timeout int
    mutatorInfo MutatorInfo
    fixupCallback FixupCallback
}

// ReadInputFile reads data from the user-specified input file into a buffer.
// It returns the content of the file in a byte array.
func ReadInputFile(inputFile *string) []byte {
    data, err := ioutil.ReadFile(*inputFile)
    if err != nil {
        fmt.Println("Error reading input file", err)
        return nil
    }

    return data
}

func main() {
    var context RFContext;

    inputFile   := flag.String("in", "", "Input template ROM")
    outputDir   := flag.String("out", "", "Output directory for fuzzed ROMs")
    commandLine := flag.String("cmd", "", "Command line")
    numWorkers  := flag.Int("workers", 4, "Max workers (threads) to run in parallel")
    timeout     := flag.Int("timeout", 1000, "Timeout in milliseconds")

    flag.Parse()

    if *inputFile == "" {
        fmt.Println("Missing required flag: --in")
        os.Exit(1)
    }

    if *outputDir == "" {
        fmt.Println("Missing required flag: --outputDir")
        os.Exit(1)
    }

    if *commandLine == "" {
        fmt.Println("Missing required flag: --cmd")
        os.Exit(1)
    }

    // Create the output directory if it doesn't exist already.
    os.MkdirAll(*outputDir, os.ModePerm)

    // Read the template data from the input file.
    context.inputData = ReadInputFile(inputFile)
    if context.inputData == nil {
        os.Exit(1);
    }

    context.inputFile     = *inputFile
    context.outputDir     = *outputDir
    context.commandLine   = *commandLine
    context.numWorkers    = *numWorkers
    context.timeout       = *timeout
    context.fixupCallback = GenesisFixupChecksum
    RunApplication(&context)
}

package main;

import(
    "io/ioutil"
    "fmt"
    "path/filepath"
    "strings"
    "errors"
    "time"
)

// GenerateRandomFilename creates a unique filename from the original input file name.
func GenerateRandomFilename(context *RFContext) string {
    fileName := strings.TrimSuffix(filepath.Base(context.inputFile), filepath.Ext(context.inputFile))
    return fmt.Sprintf(
        "%s_pos%d_bit%d_%v%s",
        fileName,
        context.mutatorInfo.index,
        context.mutatorInfo.bitPosition,
        int64(time.Now().UTC().UnixNano()),
        filepath.Ext(context.inputFile))
}

// WriteOutputFile writes out fuzzed inputs to a file on disk.
func WriteOutputFile(context *RFContext) (string, error) {
    var outputFile string

    filename := GenerateRandomFilename(context)
    fullPath := filepath.Join(context.outputDir, filename)
    err := ioutil.WriteFile(fullPath, context.inputData, 0644)
    if err != nil {
        return "", errors.New("Failed to write fuzzed output file")
    }

    outputFile, err = filepath.Abs(fullPath)
    if err != nil {
        return "", errors.New("Failed to resolve absolute path")
    }

    return outputFile, nil
}

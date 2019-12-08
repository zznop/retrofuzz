package main;

import(
	"sync"
	"os"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"strconv"
	ui "github.com/gizak/termui/v3"
    "github.com/gizak/termui/v3/widgets"
    "time"
)

// UIComponents is a context structure for managing the UI
type UIComponents struct {
	header *widgets.Paragraph
	stageCounter *widgets.Paragraph
	fuzzerElapsed *widgets.Paragraph
	currentStage *widgets.Paragraph
	stageTests *widgets.Paragraph
	progressBar *widgets.Gauge
	crashCounter *widgets.Paragraph
}

// GetCrashFileCount counts the number of files in the crash directory
func GetCrashFileCount(crashFileDir string) int {
	// Check if the crash directory exists, yet
	if _, err := os.Stat(crashFileDir); os.IsNotExist(err) {
		return 0
	}

	// Return file count
	files, _ := ioutil.ReadDir(crashFileDir)
	return len(files)
}

// RunMutationStage runs the bit flip stage and reports on progress
func RunMutationStage(context *RFContext, uic *UIComponents, keepRunning *bool) {
	const bitsInByte = 8
    var wg sync.WaitGroup

    context.mutatorInfo.bitPosition = 0
    context.mutatorInfo.index = 0
    context.mutatorInfo.currTestCase = 1
    crashFileDir := filepath.Join(context.outputDir, "crashes")

    for (context.mutatorInfo.currTestCase <= context.mutatorInfo.maxTestCases) && *keepRunning {
        context.mutatorInfo.callback(context, false)

        // Fixup if necessary
        if context.fixupCallback != nil {
        	context.fixupCallback(context)
    	}

        // Write out file to disk
        outputFile, err := WriteOutputFile(context)
        if err != nil {
            break
        }

        // Spawn worker
        wg.Add(1)
        go ExecuteCommandLine(outputFile, context.commandLine, (uint32)(context.timeout), &wg)

        // Max workers reached, wait on the group
        if context.mutatorInfo.currTestCase % context.numWorkers == 0 {
            wg.Wait()
        }

        context.mutatorInfo.callback(context, true)
        context.mutatorInfo.currTestCase++

        uic.fuzzerElapsed.Text = time.Since(context.startTime).String()
        uic.stageTests.Text = strconv.Itoa(context.mutatorInfo.currTestCase) +
        	" / " + strconv.Itoa(context.mutatorInfo.maxTestCases)
        uic.progressBar.Percent =
        	(int)((float64(context.mutatorInfo.currTestCase) / float64(context.mutatorInfo.maxTestCases)) * 100.0)
        uic.crashCounter.Text = strconv.Itoa(GetCrashFileCount(crashFileDir))
		ui.Render(uic.header, uic.stageCounter, uic.fuzzerElapsed,
			uic.currentStage, uic.stageTests, uic.progressBar, uic.crashCounter)
    }

    wg.Wait()
}

// RunFuzzer runs the fuzzer
func RunFuzzer(context *RFContext, uic *UIComponents, keepRunning *bool, wg *sync.WaitGroup) {
	context.startTime = time.Now()

	// Generate 68k instructions
	uic.currentStage.Text            = "Instr Gen"
	uic.stageCounter.Text            = "1 / 4"
	context.mutatorInfo.callback     = InstructionGeneration
	context.mutatorInfo.maxTestCases = 0xFFFF
	RunMutationStage(context, uic, keepRunning)

	// Flip a byte at a time
	uic.currentStage.Text            = "Byte flip"
	uic.stageCounter.Text            = "2 / 4"
    context.mutatorInfo.callback     = ByteFlip
    context.mutatorInfo.maxTestCases = len(context.inputData)
    RunMutationStage(context, uic, keepRunning)

    // Flip a nibble at a time
	uic.currentStage.Text            = "Nibble flip"
	uic.stageCounter.Text            = "3 / 4"
    context.mutatorInfo.callback     = NibbleFlip
    context.mutatorInfo.maxTestCases = len(context.inputData) * 2
    RunMutationStage(context, uic, keepRunning)

    // Flip a bit at a time
	uic.currentStage.Text            = "Bit flip"
	uic.stageCounter.Text            = "4 / 4"
    context.mutatorInfo.callback     = BitFlip
    context.mutatorInfo.maxTestCases = len(context.inputData) * 8
    RunMutationStage(context, uic, keepRunning)

    wg.Done()
}

// RunApplication initializes the UI and runs the app
func RunApplication(context *RFContext) {
	var uic UIComponents
	var wg sync.WaitGroup

	if err := ui.Init(); err != nil {
        fmt.Println("Failed to initialize UI\n")
    }
    defer ui.Close()

    // Create header block
    header := widgets.NewParagraph()
    header.Text = " ____  ____  ____  ____  _____  ____  __  __  ____  ____ \n" +
        "(  _ \\( ___)(_  _)(  _ \\(  _  )( ___)(  )(  )(_   )(_   )\n"         +
        " )   / )__)   )(   )   / )(_)(  )__)  )(__)(  / /_  / /_ \n"           +
        "(_)\\_)(____) (__) (_)\\_)(_____)(__)  (______)(____)(____)\n\n"       +
        ""                                                                      +
        "                   (c) zznop 2019\n\n"
    header.Border = false
    header.SetRect(1, 0, 75, 8)

    // Stage count
    stageCounter := widgets.NewParagraph()
    stageCounter.Title = "Stages"
    stageCounter.TitleStyle.Fg = ui.ColorGreen
    stageCounter.SetRect(1, 8, 30, 11)

    // Time fuzzer started
    fuzzerElapsed := widgets.NewParagraph()
    fuzzerElapsed.Title = "Fuzzer Elapsed"
    fuzzerElapsed.TitleStyle.Fg = ui.ColorGreen
    fuzzerElapsed.SetRect(31, 8, 62, 11)


    // Current Stage
    currentStage := widgets.NewParagraph()
    currentStage.Title = "Current Stage"
    currentStage.TitleStyle.Fg = ui.ColorGreen
    currentStage.SetRect(1, 11, 30, 14)

    stageTests := widgets.NewParagraph()
    stageTests.Title = "Stage Tests"
    stageTests.TitleStyle.Fg = ui.ColorGreen
    stageTests.SetRect(31, 11, 62, 14)

    // Create progress bar
    progressBar := widgets.NewGauge()
    progressBar.Title     = "Stage Progress"
    progressBar.Percent       = 0
    progressBar.TitleStyle.Fg = ui.ColorGreen
    progressBar.SetRect(1, 14, 50, 17)

    // Create crash count box
    crashCounter := widgets.NewParagraph()
    crashCounter.Title         = "Crashes"
    crashCounter.Text          = "0"
    crashCounter.TitleStyle.Fg = ui.ColorRed
    crashCounter.SetRect(51, 14, 62, 17)

    // Render UI
	uic.header        = header
	uic.stageCounter  = stageCounter
	uic.fuzzerElapsed = fuzzerElapsed
	uic.currentStage  = currentStage
	uic.stageTests  = stageTests
	uic.progressBar   = progressBar
	uic.crashCounter  = crashCounter
	ui.Render(uic.header, uic.stageCounter, uic.fuzzerElapsed,
		uic.currentStage, uic.stageTests, uic.progressBar, uic.crashCounter)

	// Run the fuzzer
	keepRunning := true
	wg.Add(1)
	go RunFuzzer(context, &uic, &keepRunning, &wg)

	// Handle user input
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			keepRunning = false
			wg.Wait()
			return
		}
	}
}

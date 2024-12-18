package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	// "sync"
	"time"

	"github.com/briandowns/spinner"
)

var (
	quiet          bool
	spin           *spinner.Spinner
	currentMessage string
	isRunning      bool
	// mutex          sync.Mutex
	originalStdout *os.File
	pipeReader     *os.File
	pipeWriter     *os.File
)

func init() {
	spin = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	redirectStdout()
}

func SpinEnable(isEnabled bool) {
	quiet = !isEnabled
}

func SpinStart(msg string) {
	// mutex.Lock()
	// defer mutex.Unlock()

	isRunning = true
	currentMessage = msg
	if quiet {
		fmt.Println(msg)
		return
	}
	_ = spin.Color("red")
	spin.Stop()
	spin.Suffix = " " + msg + "\n"
	spin.Start()
}

func SpinPause() bool {
	// mutex.Lock()
	// defer mutex.Unlock()
	wasRunning := isRunning
	if wasRunning {
		spin.Stop()
		isRunning = false
	}
	return wasRunning
}

func SpinUnpause() {
	// mutex.Lock()
	// defer mutex.Unlock()

	if isRunning {
		return
	}
	isRunning = true
	if quiet {
		return
	}
	spin.Stop()
	spin.Suffix = " " + currentMessage + "\n"
	spin.Start()
}

func SpinStop() {
	isRunning = false
	if quiet {
		return
	}
	spin.Stop()
}

func redirectStdout() {
	originalStdout = os.Stdout
	pipeReader, pipeWriter, _ = os.Pipe()
	os.Stdout = pipeWriter
	go monitorStdout()
}

func monitorStdout() {
	scanner := bufio.NewScanner(pipeReader)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "Enter OTP:") {
			SpinPause()
			// Restore stdout for OTP prompt
			os.Stdout = originalStdout
			fmt.Println(line) // Show OTP prompt
			reader := bufio.NewReader(os.Stdin)
			otp, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(originalStdout, "Failed to read OTP:", err)
			}
			otp = strings.TrimSpace(otp)
		} else {
			fmt.Fprintln(originalStdout, line)
		}
	}
	spin.Suffix = " " + currentMessage + "\n"
	spin.Start()
}

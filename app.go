package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func execute() bool {
	running, list := checkWslInstance()
	if running {
		fmt.Println("Those WSL are are running:")
		for i := 0; i < len(list)-1; i++ {
			fmt.Println("\t-", list[i])
		}
		shouldContinueWsl, _ := continueOrNot("We need to shutdown all of them to reset the firewall", 'y')
		if !shouldContinueWsl {
			return false
		}
		if !shutdownWsl() {
			fmt.Println("Error while shutting down WSL")
			return false
		}
	} else {
		fmt.Println("No WSL are running")
	}

	shouldContinueFirewall, _ := continueOrNot("We need to reset the firewall and restart Windows", 'y')
	if !shouldContinueFirewall {
		return false
	}
	return resetFirewall()

}

func checkWslInstance() (bool, []string) {
	out, err := exec.Command("cmd", "/C", "wsl --list --running").Output()
	if err != nil {
		if !strings.Contains(string(err.Error()), "0xffffffff") {
			fmt.Printf("%s\n", err)
		}
		return false, []string{}
	}
	output := strings.Split(string(out[:]), "\n")
	return len(output) > 2, output[1:]

}

func shutdownWsl() bool {
	fmt.Println("Shutting down all WSL")
	err := exec.Command("cmd", "/C", "wsl --shutdown").Run()
	if err != nil {
		fmt.Printf("%s\n", err)
		return false
	}
	return true
}

func resetFirewall() bool {
	fmt.Println("Resetting firewall")
	// Execute the command:
	// 	netsh winsock reset
	// netsh int ip reset all
	// netsh winhttp reset proxy
	// ipconfig /flushdns
	err := exec.Command("cmd", "/C", "netsh winsock reset").Run()
	if err != nil {
		fmt.Printf("Error while resetting winsock: %s\n", err)
		return false
	}

	err = exec.Command("cmd", "/C", "netsh int ip reset all").Run()
	// An error always occurs here, but it's not a problem
	/*if err != nil {
		fmt.Printf("Error while resetting ip: %s\n", err)
		return false
	}*/

	err = exec.Command("cmd", "/C", "netsh winhttp reset proxy").Run()
	if err != nil {
		fmt.Printf("Error while resetting prox: %s\n", err)
		return false
	}

	err = exec.Command("cmd", "/C", "ipconfig /flushdns").Run()
	if err != nil {
		fmt.Printf("Error while flushing dns: %s\n", err)
		return false
	}

	return true
}

func continueOrNot(msg string, acceptedInput rune) (bool, rune) {
	fmt.Print(msg + ", do you want to continue?(" + string(acceptedInput) + "/any): ")
	reader := bufio.NewReader(os.Stdin)
	output, _, _ := reader.ReadRune()
	return output == acceptedInput, output
}

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("Can Only execute this on a windows machine")
	} else {
		fmt.Println("Resetting firewall to fix wsl networking issues")
		if execute() {
			fmt.Println("Shutting down Windows in 1 minute")
			exec.Command("cmd", "/C", "shutdown -r -t 60").Run()
		} else {
			fmt.Println("Aborting")
		}
	}
	continueOrNot("Press any key to exit", 0)
}

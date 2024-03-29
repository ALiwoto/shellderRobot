package tests

import (
	"fmt"
	"log"
	"os/exec"
	"testing"

	"github.com/AnimeKaizoku/ssg/ssg"
)

const (
	PSCode01 = `
	$theWrite = Write-Output -InputObject "test"
	
	Write-Output $theWrite.GetType()
	`
)

func TestPS01(t *testing.T) {
	cmd := exec.Command("powershell", "-nologo", "-noprofile")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(stdin, PSCode01)
	stdin.Close()

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%s\n", out)
}

func TestPS02(t *testing.T) {
	result := ssg.RunPowerShell("$PSVersionTable.PSVersion")
	if !result.IsDone() {
		t.Error("powershell command hasn't finished yet.")
	}
}

func TestPWSH01(t *testing.T) {
	cmd := exec.Command("powershell", "-nologo", "-noprofile")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(stdin, PSCode01)
	stdin.Close()

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%s\n", out)
}

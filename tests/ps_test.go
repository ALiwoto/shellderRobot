package tests

import (
	"fmt"
	"log"
	"os/exec"
	"testing"
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

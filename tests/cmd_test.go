package tests

import (
	"testing"

	"github.com/AnimeKaizoku/ssg/ssg"
)

func TestCmd01(t *testing.T) {
	result := ssg.RunCommand("echo ok")
	if !result.IsDone() {
		t.Error("command isn't done yet")
		return
	}
}

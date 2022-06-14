package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/wwqdrh/logger"
)

//Cmd is exec on os ,no return
func Cmd(name string, arg ...string) error {
	logger.DefaultLogger.Info(fmt.Sprintf("[os]exec cmd is : %s %v", name, arg))
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("[os]os call error. %w", err)
	}
	return nil
}

//String is exec on os , return result
func String(name string, arg ...string) string {
	logger.DefaultLogger.Info(fmt.Sprintf("[os]exec cmd is : %s %v", name, arg))
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	if err := cmd.Run(); err != nil {
		logger.DefaultLogger.Error("[os]os call error. " + err.Error())
		return ""
	}
	return b.String()
}

package sshutil

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/wwqdrh/logger"
)

//Cmd is in host exec cmd
func (ss *SSH) Cmd(host string, cmd string) []byte {
	logger.DefaultLogger.Info(fmt.Sprintf("[ssh][%s] %s", host, cmd))
	session, err := ss.Connect(host)
	defer func() {
		if r := recover(); r != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("[ssh][%s]Error create ssh session failed,%s", host, err))
		}
	}()
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
	defer session.Close()
	b, err := session.CombinedOutput(cmd)
	logger.DefaultLogger.Debug(fmt.Sprintf("[ssh][%s]command result is: %s", host, string(b)))
	defer func() {
		if r := recover(); r != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("[ssh][%s]Error exec command failed: %s", host, err))
		}
	}()
	if err != nil {
		panic(1)
	}
	return b
}

func readPipe(host string, pipe io.Reader, isErr bool) {
	r := bufio.NewReader(pipe)
	for {
		line, _, err := r.ReadLine()
		if line == nil {
			return
		} else if err != nil {
			logger.DefaultLogger.Info(fmt.Sprintf("[%s] %s", host, line))
			logger.DefaultLogger.Error(fmt.Sprintf("[ssh] [%s] %s", host, err))
			return
		} else {
			if isErr {
				logger.DefaultLogger.Error(fmt.Sprintf("[%s] %s", host, line))
			} else {
				logger.DefaultLogger.Info(fmt.Sprintf("[%s] %s", host, line))
			}
		}
	}
}

func (ss *SSH) CmdAsync(host string, cmd string) error {
	logger.DefaultLogger.Debug(fmt.Sprintf("[%s] %s", host, cmd))
	session, err := ss.Connect(host)
	if err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("[ssh][%s]Error create ssh session failed,%s", host, err))
		return err
	}
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("[ssh][%s]Unable to request StdoutPipe(): %s", host, err))
		return err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("[ssh][%s]Unable to request StderrPipe(): %s", host, err))
		return err
	}
	if err := session.Start(cmd); err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("[ssh][%s]Unable to execute command: %s", host, err))
		return err
	}
	doneout := make(chan bool, 1)
	doneerr := make(chan bool, 1)
	go func() {
		readPipe(host, stderr, true)
		doneerr <- true
	}()
	go func() {
		readPipe(host, stdout, false)
		doneout <- true
	}()
	<-doneerr
	<-doneout
	return session.Wait()
}

//CmdToString is in host exec cmd and replace to spilt str
func (ss *SSH) CmdToString(host, cmd, spilt string) string {
	if data := ss.Cmd(host, cmd); data != nil {
		str := string(data)
		str = strings.ReplaceAll(str, "\r\n", spilt)
		return str
	}
	return ""
}

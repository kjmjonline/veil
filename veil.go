// File: veil.go
// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (c) 2024 Justin Hanekom
// -*- mode: Go -*-

/*
  This file is part of veil - minor enhancements to Go libraries.

  veil is free software: you can redistribute it and/or modify it
  under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  veil is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with go-veil. If not, see <https://www.gnu.org/licenses/>.
*/

package veil

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// CaptureOutput captures and returns the output of function `f`.
//
// The standard input and standard output are both returned,
// and are not written to `stdin` or `stdout`.
//
// When this function returns, `stdin` and `stdout` are restored to the
// streams that they originally referred to.
func CaptureOutput(f func()) (string, error) {
	reader, writer, err := os.Pipe()
	if err != nil {
		return "", err
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdout = writer
	os.Stderr = writer
	wg := new(sync.WaitGroup)
	wg.Add(1)
	out := make(chan string)
	go func() {
		var buff bytes.Buffer
		wg.Done()
		// do nothing if an error occurs
		// because there is nothing we can do
		io.Copy(&buff, reader) // nolint:errcheck
		out <- buff.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out, nil
} // CaptureOutput

// FilePathInCwd returns the full path of the file named
// `fileName` in the current working directory.
func FilePathInCwd(fileName string) (filePath string, err error) {
	var cwd string
	if cwd, err = os.Getwd(); err == nil {
		filePath = filepath.Join(cwd, fileName)
	}
	return filePath, err
} // FilePathInCwd

// IgnoreUnused hides the fact that any of its given
// values are not used. The values can be of any type.
//
// This function can be used to remove Go errors caused by constants,
// functions, or variables not being used anywhere in your code.
func IgnoreUnused(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

// SetGlobalZerologToFile sets up the global log with
// the given logging `level` to a file named `logName`.
//
// The file used for logging is appended to if it already exists,
// otherwise the file named `logName` is created with reading and writing
// permissions for the current user, and reading permissions for the group
// or other users.
//
// Logging is set up to create log entries with the current time timestamp,
// and file and line number where the log entry was created. The timestamps
// use the RFC 3339 Nano time format, which is a format that has
// millisecond accuracy.
//
// If you want to log stack traces then you should follow this pattern:
//
//	```go
//	withStack := errors.WithStack(err)
//	log.Error().Stack().Err(withStack).Msg("an error occurred")
//
// i.e., you need to wrap the error using github.com/pkg/errors.
func SetGlobalZerologToFile(logName string, level zerolog.Level) (err error) {
	var f *os.File
	f, err = os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        f,
		TimeFormat: "Mon 02 Jan 2006, 15:04:05.000",
	}).
		With().Timestamp().Caller().Logger()
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return err
} // SetGlobalZerologToFile

// vim: set ft=go sw=4 sts=4 ts=4 ai ar si sta

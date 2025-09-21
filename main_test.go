package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/maaslalani/slides/internal/model"
	"github.com/maaslalani/slides/styles"
)

func TestRunStdout(t *testing.T) {
	presentation := model.Model{
		Slides: []string{"first slide", "second slide"},
		Theme:  styles.SelectTheme("notty"),
	}

	var buf bytes.Buffer
	if err := runStdout(presentation, &buf); err != nil {
		t.Fatalf("runStdout error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "first slide") {
		t.Fatalf("stdout output missing first slide: %q", out)
	}
	if !strings.Contains(out, "second slide") {
		t.Fatalf("stdout output missing second slide: %q", out)
	}
	if !strings.Contains(out, "\n---\n") {
		t.Fatalf("stdout output missing slide delimiter: %q", out)
	}
}

func TestRunStdoutWideTable(t *testing.T) {
	tableSlide := `|   Status Code    | Error Code                                     | Description                                                               |
| :--------------: | ---------------------------------------------- | ------------------------------------------------------------------------- |
|     **共通**     |                                                |
|       400        | invalid_parameter                              | リクエストパラメータが不正                                                |
|       400        | out_of_range                                   | 指定された値が範囲外                                                      |
|       400        | insufficient_conditions                        | 条件が満たされていない                                                    |`

	presentation := model.Model{
		Slides: []string{tableSlide},
		Theme:  styles.SelectTheme("notty"),
	}

	var buf bytes.Buffer
	if err := runStdout(presentation, &buf); err != nil {
		t.Fatalf("runStdout error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "STATUS CODE |       ERROR CODE        |        DESCRIPTION") {
		t.Fatalf("table header was rendered unexpectedly: %q", out)
	}
	if !strings.Contains(out, "400     | invalid_parameter       | リクエストパラメータが不正") {
		t.Fatalf("table row for invalid_parameter was wrapped: %q", out)
	}
}

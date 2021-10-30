package highlight_test

import (
	"polarite/platform/highlight"
	"testing"
)

func TestHighlight(t *testing.T) {
	var code string = `using System;

	class Main {
		public static void Main() {
			Console.WriteLine("sup world");
		};
	};`

	_, err := highlight.Highlight(code, "c#", "nord", "true")
	if err != nil {
		t.Error("an error was thrown:", err)
	}
}

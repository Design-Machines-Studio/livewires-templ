package livewires

import "testing"

func TestInitials(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"two words", "John Doe", "JD"},
		{"single word", "Alice", "A"},
		{"three words", "Mary Jane Watson", "MJW"},
		{"empty", "", "?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Initials(tt.input); got != tt.expect {
				t.Errorf("Initials(%q) = %q, want %q", tt.input, got, tt.expect)
			}
		})
	}
}

func TestTitle(t *testing.T) {
	if got := Title("hello world"); got != "Hello World" {
		t.Errorf("Title(\"hello world\") = %q, want \"Hello World\"", got)
	}
}

func TestDateShort(t *testing.T) {
	if got := DateShort("2024-03-15"); got != "Mar 15, 2024" {
		t.Errorf("DateShort(\"2024-03-15\") = %q", got)
	}
	if got := DateShort(""); got != "" {
		t.Errorf("DateShort(\"\") = %q, want empty", got)
	}
}

func TestYear(t *testing.T) {
	if got := Year("2024-03-15"); got != "2024" {
		t.Errorf("Year(\"2024-03-15\") = %q", got)
	}
}

func TestItoa(t *testing.T) {
	if got := Itoa(42); got != "42" {
		t.Errorf("Itoa(42) = %q", got)
	}
}

func TestPtrVal(t *testing.T) {
	s := "hello"
	if got := PtrVal(&s); got != "hello" {
		t.Errorf("PtrVal(&\"hello\") = %q", got)
	}
	if got := PtrVal(nil); got != "" {
		t.Errorf("PtrVal(nil) = %q, want empty", got)
	}
}

func TestSplitCSV(t *testing.T) {
	got := SplitCSV("a, b, c")
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("SplitCSV(\"a, b, c\") = %v", got)
	}
	if got := SplitCSV(""); len(got) != 0 {
		t.Errorf("SplitCSV(\"\") = %v, want empty", got)
	}
}

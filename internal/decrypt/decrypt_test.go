package decrypt

import (
	"strconv"
	"strings"
	"testing"
)

func encodeSegment(val int, digits string, base int) string {
	if val == 0 {
		return string(digits[0])
	}
	var out []byte
	for val > 0 {
		out = append([]byte{digits[val%base]}, out...)
		val /= base
	}
	return string(out)
}

func encodeForTest(plain, n string, tNum, eNum int) string {
	digits := alphabet[:eNum]
	segs := make([]string, len(plain))
	for i := 0; i < len(plain); i++ {
		seg := encodeSegment(int(plain[i])+tNum, digits, eNum)
		var obfuscated strings.Builder
		for _, c := range seg {
			obfuscated.WriteByte(n[c-'0'])
		}
		segs[i] = obfuscated.String()
	}
	return strings.Join(segs, string(n[eNum]))
}

func TestDecodeSnapAppRoundTrip(t *testing.T) {
	n := "qwertyuiop"
	tNum, eNum := 3, 5
	plain := "Hello, world!"
	args := []string{encodeForTest(plain, n, tNum, eNum), "_", n, strconv.Itoa(tNum), strconv.Itoa(eNum), "_"}
	got, err := decodeSnapApp(args)
	if err != nil {
		t.Fatalf("decodeSnapApp error: %v", err)
	}
	if got != plain {
		t.Fatalf("decodeSnapApp() = %q, want %q", got, plain)
	}
}

func TestExtractArgs(t *testing.T) {
	data := `foo;decodeURIComponent(escape(r))}("abc","u","n","3","5","r"));bar`
	args, err := extractArgs(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"abc", "u", "n", "3", "5", "r"}
	if len(args) != len(want) {
		t.Fatalf("got %v, want %v", args, want)
	}
	for i := range want {
		if args[i] != want[i] {
			t.Fatalf("args[%d] = %q, want %q", i, args[i], want[i])
		}
	}
}

func TestDecodedSnapSaveError(t *testing.T) {
	data := `document.querySelector("#alert").innerHTML = "Invalid link";`
	_, err := decodedSnapSave(data)
	if err == nil || err.Error() != "Invalid link" {
		t.Fatalf("expected alert error, got %v", err)
	}
}

func TestDecodedSnapSaveBody(t *testing.T) {
	data := `getElementById("download-section").innerHTML = "<div>ok<\/div>"; document.getElementById("inputData").remove(); `
	got, err := decodedSnapSave(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "<div>ok</div>" {
		t.Fatalf("decodedSnapSave() = %q, want %q", got, "<div>ok</div>")
	}
}

func TestDecodedSnaptikBody(t *testing.T) {
	data := `$("#download").innerHTML = "<a>dl<\/a>"; document.getElementById("inputData").remove(); `
	got, err := decodedSnaptik(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "<a>dl</a>" {
		t.Fatalf("decodedSnaptik() = %q, want %q", got, "<a>dl</a>")
	}
}

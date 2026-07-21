package decrypt

import (
	"errors"
	"strconv"
	"strings"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/"

func ipow(base, exp int) int {
	r := 1
	for range exp {
		r *= base
	}
	return r
}

func baseDecode(d string, e, f int) string {
	hArr := alphabet[:e]
	j := 0
	for c, i := 0, len(d)-1; i >= 0; c, i = c+1, i-1 {
		idx := strings.IndexByte(hArr, d[i])
		if idx == -1 {
			continue
		}
		j += idx * ipow(e, c)
	}
	if j <= 0 {
		return "0"
	}
	var k []byte
	for j > 0 {
		k = append([]byte{alphabet[j%f]}, k...)
		j /= f
	}
	return string(k)
}

func decodeSnapApp(args []string) (string, error) {
	if len(args) < 6 {
		return "", errors.New("decrypt: expected 6 encoded arguments")
	}
	h, n, t, e := args[0], args[2], args[3], args[4]
	tNum, err := strconv.Atoi(t)
	if err != nil {
		return "", errors.New("decrypt: invalid t argument")
	}
	eNum, err := strconv.Atoi(e)
	if err != nil || eNum >= len(n) {
		return "", errors.New("decrypt: invalid e argument")
	}
	delim := n[eNum]

	var out []byte
	for i := 0; i < len(h); {
		start := i
		for i < len(h) && h[i] != delim {
			i++
		}
		seg := h[start:i]
		i++

		for j := range len(n) {
			seg = strings.ReplaceAll(seg, string(n[j]), strconv.Itoa(j))
		}
		val, err := strconv.Atoi(baseDecode(seg, eNum, 10))
		if err != nil {
			return "", errors.New("decrypt: malformed segment")
		}
		out = append(out, byte(val-tNum)) //nolint:gosec
	}
	return string(out), nil
}

func extractArgs(data string) ([]string, error) {
	const marker = "decodeURIComponent(escape(r))}("
	_, after, found := strings.Cut(data, marker)
	if !found {
		return nil, errors.New("decrypt: encoded payload not found")
	}
	inner, _, _ := strings.Cut(after, "))")
	rawArgs := strings.Split(inner, ",")
	args := make([]string, len(rawArgs))
	for i, v := range rawArgs {
		args[i] = strings.TrimSpace(strings.ReplaceAll(v, `"`, ""))
	}
	return args, nil
}

func decodedSnapSave(data string) (string, error) {
	if _, after, found := strings.Cut(data, `document.querySelector("#alert").innerHTML = "`); found {
		if msg, _, _ := strings.Cut(after, `";`); strings.TrimSpace(msg) != "" {
			return "", errors.New(strings.TrimSpace(msg))
		}
	}
	const marker = `getElementById("download-section").innerHTML = "`
	const end = `"; document.getElementById("inputData").remove(); `
	_, after, found := strings.Cut(data, marker)
	if !found {
		return "", errors.New("decrypt: download-section not found")
	}
	body, _, _ := strings.Cut(after, end)
	return strings.ReplaceAll(body, `\`, ""), nil
}

func DecryptSnapSave(data string) (string, error) {
	args, err := extractArgs(data)
	if err != nil {
		return "", err
	}
	decoded, err := decodeSnapApp(args)
	if err != nil {
		return "", err
	}
	return decodedSnapSave(decoded)
}

package test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	goahocorasick "github.com/anknown/ahocorasick"
)

func ReadRunes(filename string) ([][]rune, error) {
	dict := [][]rune{}

	f, err := os.OpenFile(filename, os.O_RDONLY, 0660)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(f)
	for {
		l, err := r.ReadBytes('\n')
		if err != nil || err == io.EOF {
			break
		}
		l = bytes.TrimSpace(l)
		dict = append(dict, bytes.Runes(l))
	}

	return dict, nil
}

// Aho-Corasick use Double Array Trie instead of common Linked List Trie
func TestAcdat(t *testing.T) {
	dict, err := ReadRunes("../engine/test_data/dictionary.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	content := []rune("xxxx《我的团长我的团》作者：兰晓龙")

	m := new(goahocorasick.Machine)
	if err := m.Build(dict); err != nil {
		fmt.Println(err)
		return
	}

	terms := m.MultiPatternSearch(content, false)
	for _, t := range terms {
		fmt.Printf("found %d %s\n", t.Pos, string(t.Word))
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/TobiasYin/go-lsp/logs"
	"github.com/TobiasYin/go-lsp/lsp"
	"github.com/TobiasYin/go-lsp/lsp/defines"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type Snippets map[string]Snippet

type Snippet struct {
	Prefix       Prefix `json:"prefix"`
	Body         Body   `json:"body"`
	Descripttion string `json:"description"`
}

type Prefix struct {
	Value interface{}
}

func (p *Prefix) UnmarshalJSON(data []byte) error {
	var singleWord string
	if err := json5.Unmarshal(data, &singleWord); err == nil {
		p.Value = singleWord
		return nil
	}

	var multipleWords []string
	if err := json5.Unmarshal(data, &multipleWords); err == nil {
		p.Value = multipleWords
		return nil
	}

	return fmt.Errorf("Cannot unmarshal snippet body: %s", data)
}

func (p Prefix) ToStringSlice() []string {
	switch v := p.Value.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	default:
		return nil
	}
}

type Body struct {
	Value interface{}
}

func (b *Body) UnmarshalJSON(data []byte) error {
	var singleLine string
	if err := json5.Unmarshal(data, &singleLine); err == nil {
		b.Value = singleLine
		return nil
	}

	var multipleLines []string
	if err := json5.Unmarshal(data, &multipleLines); err == nil {
		b.Value = multipleLines
		return nil
	}

	return fmt.Errorf("Cannot unmarshal snippet body: %s", data)
}

func (b Body) String() string {
	switch v := b.Value.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, "\n")
	default:
		return ""
	}
}

func main() {
	lang := flag.String("lang", "go", "A language that is used in a vscode-snippet path")
	flag.Parse()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home path: %v", err)
	}

	configPath := path.Join(home, "Library/Application Support/Code/User/snippets", *lang+".json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var snippets Snippets
	if err := json5.Unmarshal(data, &snippets); err != nil {
		log.Fatalf("Error unmarshalling config data: %v", err)
	}

	logger := log.New(os.Stdout, "snippets-ls: ", log.LstdFlags)
	logs.Init(logger)

	server := lsp.NewServer(&lsp.Options{CompletionProvider: &defines.CompletionOptions{
		TriggerCharacters: &[]string{"."},
	}})

	items := make([]defines.CompletionItem, 0)
	k := defines.CompletionItemKindSnippet
	for _, snippet := range snippets {
		for _, prefix := range snippet.Prefix.ToStringSlice() {
			item := defines.CompletionItem{
				Kind:       &k,
				Label:      prefix,
				InsertText: strPtr(fmt.Sprintf("%s", snippet.Body)),
			}
			items = append(items, item)
		}
	}

	server.OnCompletion(func(ctx context.Context, req *defines.CompletionParams) (result *[]defines.CompletionItem, err error) {
		logs.Println(req)
		return &items, nil
	})

	server.Run()
}

func strPtr(s string) *string {
	return &s
}

package components

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/walteh/protobuf-language-server/proto/view"
	"github.com/walteh/retab/v2/pkg/format"
	"github.com/walteh/retab/v2/pkg/format/editorconfig"
	"github.com/walteh/retab/v2/pkg/format/protofmt"

	"github.com/walteh/protobuf-language-server/go-lsp/lsp/defines"
)

func FormatWithRetab(ctx context.Context, req *defines.DocumentFormattingParams) (result *[]defines.TextEdit, err error) {
	formatter := protofmt.NewFormatter()
	proto_file, err := view.ViewManager.GetFile(req.TextDocument.Uri)
	if err != nil {
		return nil, err
	}
	text, _, _ := proto_file.Read(ctx)
	reader := strings.NewReader(string(text))

	cfgProvider, err := editorconfig.NewDynamicConfigurationProvider(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("creating configuration provider: %w", err)
	}

	r, err := format.Format(ctx, formatter, cfgProvider, string(req.TextDocument.Uri), reader)
	if err != nil {
		return nil, fmt.Errorf("formatting content: %w", err)
	}

	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading formatted content: %w", err)
	}

	return &[]defines.TextEdit{
		{
			Range: defines.Range{
				Start: defines.Position{Line: 0, Character: 0},
				End:   defines.Position{Line: math.MaxInt32, Character: math.MaxInt32},
			},
			NewText: string(out),
		},
	}, nil
}

func Format(ctx context.Context, req *defines.DocumentFormattingParams) (result *[]defines.TextEdit, err error) {
	if !view.IsProtoFile(req.TextDocument.Uri) {
		return nil, nil
	}
	format := exec.Command("clang-format", fmt.Sprintf("--assume-filename=%v", filepath.Base(string(req.TextDocument.Uri))))
	in, err := format.StdinPipe()
	if err != nil {
		return nil, err
	}

	proto_file, err := view.ViewManager.GetFile(req.TextDocument.Uri)
	if err != nil {
		return nil, err
	}
	data, _, _ := proto_file.Read(ctx)
	go func() {
		io.WriteString(in, string(data))
		defer in.Close()
	}()

	res, err := format.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return &[]defines.TextEdit{
		{
			Range: defines.Range{
				Start: defines.Position{Line: 0, Character: 0},
				End:   defines.Position{Line: math.MaxInt32, Character: math.MaxInt32},
			},
			NewText: string(res),
		},
	}, nil
}

func FormatRange(ctx context.Context, req *defines.DocumentRangeFormattingParams) (result *[]defines.TextEdit, err error) {
	if !view.IsProtoFile(req.TextDocument.Uri) {
		return nil, nil
	}
	format := exec.Command("clang-format", fmt.Sprintf("--assume-filename=%v", filepath.Base(string(req.TextDocument.Uri))), fmt.Sprintf("--lines=%v:%v", req.Range.Start.Line+1, req.Range.End.Line+1))
	in, err := format.StdinPipe()
	if err != nil {
		return nil, err
	}
	proto_file, err := view.ViewManager.GetFile(req.TextDocument.Uri)
	if err != nil {
		return nil, err
	}
	data, _, _ := proto_file.Read(ctx)
	go func() {
		io.WriteString(in, string(data))
		defer in.Close()
	}()

	res, err := format.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return &[]defines.TextEdit{
		{
			Range: defines.Range{
				Start: defines.Position{Line: 0, Character: 0},
				End:   defines.Position{Line: math.MaxInt32, Character: math.MaxInt32},
			},
			NewText: string(res),
		},
	}, nil
}

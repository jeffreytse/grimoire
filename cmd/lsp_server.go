package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/jeffreytse/grimoire/internal/compliance"
)

// ─── JSON-RPC 2.0 wire types ─────────────────────────────────────────────────

type rpcMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"` // number | string | null; absent on notifications
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ─── LSP capability types ─────────────────────────────────────────────────────

type initializeParams struct {
	RootURI  string `json:"rootUri"`
	RootPath string `json:"rootPath"` // deprecated but still sent by many clients
}

type initializeResult struct {
	Capabilities serverCapabilities `json:"capabilities"`
	ServerInfo   serverInfo         `json:"serverInfo"`
}

type serverCapabilities struct {
	TextDocumentSync textDocumentSyncOptions `json:"textDocumentSync"`
}

type textDocumentSyncOptions struct {
	OpenClose bool    `json:"openClose"`
	Save      saveOpt `json:"save"`
	Change    int     `json:"change"` // 0 = None — we don't need incremental diffs
}

type saveOpt struct {
	IncludeText bool `json:"includeText"`
}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ─── LSP diagnostic types ─────────────────────────────────────────────────────

type publishDiagnosticsParams struct {
	URI         string          `json:"uri"`
	Diagnostics []lspDiagnostic `json:"diagnostics"`
}

type lspDiagnostic struct {
	Range    lspRange `json:"range"`
	Severity int      `json:"severity"`
	Code     string   `json:"code,omitempty"`
	Source   string   `json:"source,omitempty"`
	Message  string   `json:"message"`
}

type lspRange struct {
	Start lspPosition `json:"start"`
	End   lspPosition `json:"end"`
}

type lspPosition struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// lspExitError is returned by serveLSP when the client sends the LSP "exit"
// notification. Code is 0 if "shutdown" was received first, 1 otherwise.
type lspExitError struct{ code int }

func (e *lspExitError) Error() string { return fmt.Sprintf("lsp: exit %d", e.code) }

// ─── Server ──────────────────────────────────────────────────────────────────

type lspServer struct {
	projectDir       string
	prevURIs         map[string]bool // URIs that had diagnostics in the previous check
	shutdownReceived bool
	opMu             sync.Mutex // serializes the full check+publish operation
}

// serveLSP runs the JSON-RPC 2.0 / LSP message loop over in/out.
func serveLSP(in io.Reader, out io.Writer) error {
	r := bufio.NewReader(in)
	srv := &lspServer{prevURIs: make(map[string]bool)}

	var writeMu sync.Mutex
	write := func(v any) {
		writeMu.Lock()
		defer writeMu.Unlock()
		_ = writeRPCMessage(out, v)
	}

	for {
		raw, err := readRPCMessage(r)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		var msg rpcMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		if msg.Method == "exit" {
			code := 1
			if srv.shutdownReceived {
				code = 0
			}
			return &lspExitError{code}
		}

		// Requests have a non-null ID; notifications do not.
		isRequest := len(msg.ID) > 0 && string(msg.ID) != "null"
		if isRequest {
			result, rpcErr := srv.handleRequest(msg.Method, msg.Params)
			resp := rpcMessage{JSONRPC: "2.0", ID: msg.ID}
			if rpcErr != nil {
				resp.Error = rpcErr
			} else {
				b, _ := json.Marshal(result)
				resp.Result = json.RawMessage(b)
			}
			write(resp)
		} else {
			go func(method string, params json.RawMessage) {
				srv.handleNotification(method, params, write)
			}(msg.Method, msg.Params)
		}
	}
}

// handleRequest responds to LSP request methods (those with an ID).
func (s *lspServer) handleRequest(method string, params json.RawMessage) (any, *rpcError) {
	switch method {
	case "initialize":
		var p initializeParams
		_ = json.Unmarshal(params, &p)
		switch {
		case p.RootURI != "":
			s.projectDir = uriToPath(p.RootURI)
		case p.RootPath != "":
			s.projectDir = p.RootPath
		default:
			if wd, err := os.Getwd(); err == nil {
				s.projectDir = wd
			}
		}
		return initializeResult{
			Capabilities: serverCapabilities{
				TextDocumentSync: textDocumentSyncOptions{
					OpenClose: false,
					Save:      saveOpt{IncludeText: false},
					Change:    0,
				},
			},
			ServerInfo: serverInfo{Name: "grimoire", Version: cliVersion},
		}, nil

	case "shutdown":
		s.shutdownReceived = true
		return nil, nil

	default:
		return nil, &rpcError{Code: -32601, Message: "method not found: " + method}
	}
}

// handleNotification processes LSP notification methods (no ID, no response).
func (s *lspServer) handleNotification(method string, _ json.RawMessage, write func(any)) {
	switch method {
	case "initialized",
		"textDocument/didSave":
		// these trigger a compliance check
	default:
		return
	}

	s.opMu.Lock()
	defer s.opMu.Unlock()
	report, err := s.runCheck()
	if err != nil {
		fmt.Fprintf(os.Stderr, "grimoire lsp: check failed: %v\n", err)
		write(lspShowMessage(1, "grimoire check failed: "+err.Error()))
		return
	}
	s.publishAll(report, write)
}

// runCheck execs `grimoire check --json` and parses the compliance report.
// Caller must hold s.opMu.
func (s *lspServer) runCheck() (*compliance.Report, error) {
	exe, err := os.Executable()
	if err != nil {
		exe = os.Args[0] // fallback
	}
	var stdout bytes.Buffer
	cmd := exec.Command(exe, "check", "--json", "--no-color") //nolint:gosec // exe comes from os.Executable(), not user input
	cmd.Stdout = &stdout
	cmd.Dir = s.projectDir
	_ = cmd.Run() // non-zero exit (compliance failures) is expected — parse output regardless

	if stdout.Len() == 0 {
		return nil, fmt.Errorf("check produced no output")
	}
	var report compliance.Report
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		return nil, fmt.Errorf("parsing check output: %w", err)
	}
	return &report, nil
}

// publishAll converts a compliance report to LSP publishDiagnostics notifications
// and sends one notification per affected file. Files that were dirty in the
// previous run but are now clean receive an empty-array notification to clear them.
func (s *lspServer) publishAll(report *compliance.Report, write func(any)) {
	byURI := make(map[string][]lspDiagnostic)

	for i := range report.Diagnostics {
		d := &report.Diagnostics[i]
		if d.Severity == 4 || d.Status == "pass" {
			continue // skip hints / passing items
		}
		uri := d.URI
		if uri == "" {
			uri = "grimoire.toml" // project-level finding — pin to manifest
		}
		absPath := filepath.Join(s.projectDir, filepath.FromSlash(uri))
		fileURI := pathToURI(absPath)
		byURI[fileURI] = append(byURI[fileURI], lspDiagnostic{
			Range: lspRange{
				Start: lspPosition{Line: d.Range.Start.Line, Character: d.Range.Start.Character},
				End:   lspPosition{Line: d.Range.End.Line, Character: d.Range.End.Character},
			},
			Severity: d.Severity,
			Code:     d.Code,
			Source:   "grimoire",
			Message:  d.Message,
		})
	}

	// clear URIs that were dirty last run but are now clean
	for uri := range s.prevURIs {
		if _, stillDirty := byURI[uri]; !stillDirty {
			write(lspNotification("textDocument/publishDiagnostics", publishDiagnosticsParams{
				URI:         uri,
				Diagnostics: []lspDiagnostic{},
			}))
		}
	}

	// publish current diagnostics
	newPrev := make(map[string]bool, len(byURI))
	for uri, diags := range byURI {
		write(lspNotification("textDocument/publishDiagnostics", publishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diags,
		}))
		newPrev[uri] = true
	}
	s.prevURIs = newPrev
}

// ─── JSON-RPC 2.0 framing ────────────────────────────────────────────────────

// readRPCMessage reads one Content-Length-framed JSON message from r.
func readRPCMessage(r *bufio.Reader) (json.RawMessage, error) {
	var contentLength int
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break // blank line ends headers
		}
		if strings.HasPrefix(line, "Content-Length: ") {
			if _, err := fmt.Sscanf(line[16:], "%d", &contentLength); err != nil {
				return nil, fmt.Errorf("bad Content-Length header: %w", err)
			}
		}
	}
	if contentLength <= 0 {
		return nil, fmt.Errorf("Content-Length missing or zero")
	}
	buf := make([]byte, contentLength)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return json.RawMessage(buf), nil
}

// writeRPCMessage serialises v as JSON and writes it with a Content-Length header.
func writeRPCMessage(w io.Writer, v any) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}
	_, err = w.Write(body)
	return err
}

// lspShowMessage builds a window/showMessage notification (msgType: 1=Error 2=Warning 3=Info 4=Log).
func lspShowMessage(msgType int, text string) rpcMessage {
	b, _ := json.Marshal(map[string]any{"type": msgType, "message": text})
	return rpcMessage{JSONRPC: "2.0", Method: "window/showMessage", Params: json.RawMessage(b)}
}

// lspNotification builds a JSON-RPC notification message (no ID).
func lspNotification(method string, params any) rpcMessage {
	b, _ := json.Marshal(params)
	return rpcMessage{
		JSONRPC: "2.0",
		Method:  method,
		Params:  json.RawMessage(b),
	}
}

// ─── URI helpers ─────────────────────────────────────────────────────────────

// uriToPath converts a file:// URI to an OS filesystem path.
func uriToPath(uri string) string {
	path := strings.TrimPrefix(uri, "file://")
	// Windows: file:///C:/path arrives as /C:/path — strip the leading slash.
	if runtime.GOOS == "windows" && len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		path = path[1:]
	}
	return filepath.FromSlash(path)
}

// pathToURI converts an absolute OS filesystem path to a file:// URI.
func pathToURI(path string) string {
	p := filepath.ToSlash(path)
	if !strings.HasPrefix(p, "/") {
		p = "/" + p // Windows absolute paths need a leading slash in the URI
	}
	return "file://" + p
}

package file

import (
	"context"
	"strings"
	"testing"

	"github.com/opencharly/sdk/kit"
	"github.com/opencharly/spec/spec"
)

// statSymlink is the stat probe's answer for a symlink; readlinkTo is what `readlink -f`
// resolves it to.
func linkCC(statType, readlinkOut string, readlinkExit int) *fakeCC {
	return &fakeCC{exec: &fakeExec{responses: []fakeResponse{
		{matchPrefix: "if [ -e", stdout: "exists=1|" + statType + "|777|root|root\n"},
		{matchPrefix: "readlink -f", stdout: readlinkOut, exit: readlinkExit},
	}}}
}

func runFile(cc kit.CheckContext, in map[string]any) kit.Result {
	return verb{}.RunVerb(context.Background(), cc, &spec.Op{PluginInput: in})
}

func TestLinkTarget_ResolvedTargetMatches(t *testing.T) {
	cc := linkCC("symbolic link", "/usr/share/zoneinfo/UTC\n", 0)
	res := runFile(cc, map[string]any{
		"file": "/etc/localtime", "link_target": "/usr/share/zoneinfo/UTC"})
	if res.Status != kit.StatusPass {
		t.Errorf("expected pass for a matching target, got %+v", res)
	}
}

func TestLinkTarget_WrongTargetFailsAndShowsIt(t *testing.T) {
	cc := linkCC("symbolic link", "/usr/share/zoneinfo/Europe/Berlin\n", 0)
	res := runFile(cc, map[string]any{
		"file": "/etc/localtime", "link_target": "/usr/share/zoneinfo/UTC"})
	if res.Status != kit.StatusFail {
		t.Fatalf("expected fail for a mismatched target, got %+v", res)
	}
	// The message must carry the ACTUAL target — the whole point of the assertion is
	// which one it is, and on a disposable bed the guest is gone by the time it is read.
	if !strings.Contains(res.Message, "Europe/Berlin") {
		t.Errorf("the failure must name the actual target; got %q", res.Message)
	}
}

// The matcher list is the full vocabulary, not just equality.
func TestLinkTarget_AcceptsOperatorMatchers(t *testing.T) {
	cc := linkCC("symbolic link", "/usr/share/zoneinfo/UTC\n", 0)
	res := runFile(cc, map[string]any{
		"file": "/etc/localtime", "link_target": []any{map[string]any{"contains": "zoneinfo"}}})
	if res.Status != kit.StatusPass {
		t.Errorf("a contains: matcher must work on the target, got %+v", res)
	}
}

// THE trap this guards. `readlink -f /etc/hosts` prints /etc/hosts for a REGULAR file, so
// an `equals` matcher naming the path would pass and assert precisely nothing.
func TestLinkTarget_NonSymlinkFailsRatherThanSelfMatching(t *testing.T) {
	cc := linkCC("regular file", "/etc/hosts\n", 0)
	res := runFile(cc, map[string]any{"file": "/etc/hosts", "link_target": "/etc/hosts"})
	if res.Status != kit.StatusFail {
		t.Fatalf("link_target on a regular file must FAIL, not self-match; got %+v", res)
	}
	if !strings.Contains(res.Message, "not a symlink") {
		t.Errorf("the failure must say the node is not a symlink; got %q", res.Message)
	}
}

// A dangling link resolves to nothing. An empty subject satisfies every not_contains
// matcher, so returning it instead of failing would be the check failing open.
func TestLinkTarget_DanglingLinkFails(t *testing.T) {
	cc := linkCC("symbolic link", "", 1)
	res := runFile(cc, map[string]any{
		"file": "/etc/localtime", "link_target": "/usr/share/zoneinfo/UTC"})
	if res.Status != kit.StatusFail {
		t.Fatalf("a dangling symlink must fail, got %+v", res)
	}
	if !strings.Contains(res.Message, "dangling") {
		t.Errorf("the failure must identify the link as dangling; got %q", res.Message)
	}
}

// NEGATIVE CONTROL: with link_target unset, readlink must not run at all — no fake
// response is registered for it, so a stray probe would return exit 127 and fail.
func TestLinkTarget_NotProbedWhenUnset(t *testing.T) {
	cc := &fakeCC{exec: &fakeExec{responses: []fakeResponse{
		{matchPrefix: "if [ -e", stdout: "exists=1|symbolic link|777|root|root\n"},
	}}}
	res := runFile(cc, map[string]any{"file": "/etc/localtime", "filetype": "symlink"})
	if res.Status != kit.StatusPass {
		t.Errorf("a step without link_target must not probe readlink; got %+v", res)
	}
}

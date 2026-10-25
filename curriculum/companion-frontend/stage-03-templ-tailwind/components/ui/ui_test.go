package ui_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/rasel9t6/the-go-engineer/curriculum/companion-frontend/stage-03-templ-tailwind/components/ui"
)

func renderToString(ctx context.Context, c templ.Component) (string, error) {
	var buf bytes.Buffer
	err := c.Render(ctx, &buf)
	return buf.String(), err
}

func textComponent(text string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, text)
		return err
	})
}

func TestButton_Variants(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		props       ui.ButtonProps
		childText   string
		expectedTag string
		mustContain string
	}{
		{
			name:        "Default button",
			props:       ui.ButtonProps{},
			childText:   "Click Me",
			expectedTag: "<button",
			mustContain: "bg-indigo-600",
		},
		{
			name:        "Destructive button",
			props:       ui.ButtonProps{Variant: "destructive", Size: "sm"},
			childText:   "Delete",
			expectedTag: "<button",
			mustContain: "bg-rose-600",
		},
		{
			name:        "Outline link button",
			props:       ui.ButtonProps{Variant: "outline", Href: "/dashboard"},
			childText:   "Go to Dashboard",
			expectedTag: "<a",
			mustContain: "href=\"/dashboard\"",
		},
		{
			name:        "Disabled submit button",
			props:       ui.ButtonProps{Type: "submit", Disabled: true},
			childText:   "Saving...",
			expectedTag: "<button",
			mustContain: "disabled",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			child := textComponent(tc.childText)
			comp := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
				return ui.Button(tc.props).Render(templ.WithChildren(c, child), w)
			})

			out, err := renderToString(ctx, comp)
			if err != nil {
				t.Fatalf("unexpected render error: %v", err)
			}
			if !strings.Contains(out, tc.expectedTag) {
				t.Errorf("expected tag %q in output, got: %s", tc.expectedTag, out)
			}
			if !strings.Contains(out, tc.mustContain) {
				t.Errorf("expected output to contain %q, got: %s", tc.mustContain, out)
			}
			if !strings.Contains(out, tc.childText) {
				t.Errorf("expected child text %q in output, got: %s", tc.childText, out)
			}
		})
	}
}

func TestCard_Structure(t *testing.T) {
	ctx := context.Background()

	cardComp := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
		header := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
			titleChild := textComponent("Tenant Configuration")
			descChild := textComponent("Settings for Acme Corp")
			if err := ui.CardTitle(ui.CardProps{}).Render(templ.WithChildren(c, titleChild), w); err != nil {
				return err
			}
			return ui.CardDescription(ui.CardProps{}).Render(templ.WithChildren(c, descChild), w)
		})

		content := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
			return textComponent("<p>Active isolation: strict</p>").Render(c, w)
		})

		cardChildren := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
			if err := ui.CardHeader(ui.CardProps{}).Render(templ.WithChildren(c, header), w); err != nil {
				return err
			}
			return ui.CardContent(ui.CardProps{}).Render(templ.WithChildren(c, content), w)
		})

		return ui.Card(ui.CardProps{ID: "tenant-card"}).Render(templ.WithChildren(c, cardChildren), w)
	})

	out, err := renderToString(ctx, cardComp)
	if err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}

	for _, expected := range []string{
		"id=\"tenant-card\"",
		"Tenant Configuration",
		"Settings for Acme Corp",
		"Active isolation: strict",
		"rounded-xl border border-slate-800 bg-slate-900",
	} {
		if !strings.Contains(out, expected) {
			t.Errorf("card missing expected substring: %q\nOutput: %s", expected, out)
		}
	}
}

func TestBadge_Variants(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		variant     string
		mustContain string
	}{
		{"default", "bg-indigo-600"},
		{"success", "text-emerald-400"},
		{"warning", "text-amber-400"},
		{"destructive", "bg-rose-600"},
		{"secondary", "bg-slate-800"},
	}

	for _, tc := range tests {
		t.Run(tc.variant, func(t *testing.T) {
			badgeChild := textComponent(tc.variant)
			badge := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
				return ui.Badge(ui.BadgeProps{Variant: tc.variant}).Render(templ.WithChildren(c, badgeChild), w)
			})

			out, err := renderToString(ctx, badge)
			if err != nil {
				t.Fatalf("unexpected render error: %v", err)
			}
			if !strings.Contains(out, tc.mustContain) {
				t.Errorf("badge %q missing expected class %q, got: %s", tc.variant, tc.mustContain, out)
			}
		})
	}
}

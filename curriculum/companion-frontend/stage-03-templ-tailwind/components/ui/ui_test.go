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

func TestButton_VariantsAndProps(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		props       ui.ButtonProps
		childText   string
		expectedTag string
		mustContain []string
	}{
		{
			name:        "Default button",
			props:       ui.ButtonProps{},
			childText:   "Click Me",
			expectedTag: "<button",
			mustContain: []string{"bg-indigo-600", "h-9 px-4 py-2", "Click Me", "type=\"button\""},
		},
		{
			name:        "Destructive button lg",
			props:       ui.ButtonProps{Variant: ui.VariantDestructive, Size: ui.SizeLg},
			childText:   "Delete Tenant",
			expectedTag: "<button",
			mustContain: []string{"bg-rose-600", "h-10 rounded-md px-8", "Delete Tenant"},
		},
		{
			name:        "Outline button sm",
			props:       ui.ButtonProps{Variant: ui.VariantOutline, Size: ui.SizeSm},
			childText:   "Cancel",
			expectedTag: "<button",
			mustContain: []string{"border border-slate-700", "h-8 rounded-md px-3", "Cancel"},
		},
		{
			name:        "Secondary button",
			props:       ui.ButtonProps{Variant: ui.VariantSecondary},
			childText:   "Secondary",
			expectedTag: "<button",
			mustContain: []string{"bg-slate-800", "Secondary"},
		},
		{
			name:        "Ghost button",
			props:       ui.ButtonProps{Variant: ui.VariantGhost},
			childText:   "Ghost",
			expectedTag: "<button",
			mustContain: []string{"hover:bg-slate-800", "Ghost"},
		},
		{
			name:        "Link button",
			props:       ui.ButtonProps{Variant: ui.VariantLink},
			childText:   "Documentation",
			expectedTag: "<button",
			mustContain: []string{"text-indigo-400", "underline-offset-4", "Documentation"},
		},
		{
			name:        "Disabled submit button",
			props:       ui.ButtonProps{Type: ui.TypeSubmit, Disabled: true},
			childText:   "Saving...",
			expectedTag: "<button",
			mustContain: []string{"type=\"submit\"", "disabled", "disabled:opacity-50"},
		},
		{
			name:        "Anchor link button with target",
			props:       ui.ButtonProps{Href: "/dashboard", Target: "_blank"},
			childText:   "Open Dashboard",
			expectedTag: "<a",
			mustContain: []string{"href=\"/dashboard\"", "target=\"_blank\"", "Open Dashboard"},
		},
		{
			name: "Universal props ID, Class, and Attributes",
			props: ui.ButtonProps{
				ID:         "action-btn",
				Class:      "extra-btn-class",
				Attributes: templ.Attributes{"data-testid": "submit-btn", "aria-label": "Submit Form"},
			},
			childText:   "Submit",
			expectedTag: "<button",
			mustContain: []string{"id=\"action-btn\"", "extra-btn-class", "data-testid=\"submit-btn\"", "aria-label=\"Submit Form\""},
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
			for _, s := range tc.mustContain {
				if !strings.Contains(out, s) {
					t.Errorf("expected output to contain %q, got: %s", s, out)
				}
			}
		})
	}
}

func TestCard_StructureAndPrimitives(t *testing.T) {
	ctx := context.Background()

	t.Run("Card with full composition", func(t *testing.T) {
		cardComp := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
			header := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
				titleChild := textComponent("Tenant Configuration")
				descChild := textComponent("Settings for Acme Corp")
				actionChild := textComponent("Action Badge")
				if err := ui.CardTitle(ui.CardTitleProps{}).Render(templ.WithChildren(c, titleChild), w); err != nil {
					return err
				}
				if err := ui.CardDescription(ui.CardDescriptionProps{}).Render(templ.WithChildren(c, descChild), w); err != nil {
					return err
				}
				return ui.CardAction(ui.CardActionProps{}).Render(templ.WithChildren(c, actionChild), w)
			})

			content := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
				return textComponent("<p>Active isolation: strict</p>").Render(c, w)
			})

			footer := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
				return textComponent("<span>Status: Healthy</span>").Render(c, w)
			})

			cardChildren := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
				if err := ui.CardHeader(ui.CardHeaderProps{Class: "border-b"}).Render(templ.WithChildren(c, header), w); err != nil {
					return err
				}
				if err := ui.CardContent(ui.CardContentProps{}).Render(templ.WithChildren(c, content), w); err != nil {
					return err
				}
				return ui.CardFooter(ui.CardFooterProps{}).Render(templ.WithChildren(c, footer), w)
			})

			return ui.Card(ui.CardProps{ID: "tenant-card", Class: "shadow-2xl"}).Render(templ.WithChildren(c, cardChildren), w)
		})

		out, err := renderToString(ctx, cardComp)
		if err != nil {
			t.Fatalf("unexpected render error: %v", err)
		}

		for _, expected := range []string{
			"id=\"tenant-card\"",
			"Tenant Configuration",
			"Settings for Acme Corp",
			"Action Badge",
			"Active isolation: strict",
			"Status: Healthy",
			"rounded-xl border border-slate-800 bg-slate-900",
			"border-t border-slate-800/60",
		} {
			if !strings.Contains(out, expected) {
				t.Errorf("card missing expected substring: %q\nOutput: %s", expected, out)
			}
		}
	})

	t.Run("Card with small size and attributes", func(t *testing.T) {
		comp := ui.Card(ui.CardProps{
			Size:       ui.CardSizeSm,
			Class:      "compact-card",
			Attributes: templ.Attributes{"data-state": "active"},
		})
		out, err := renderToString(ctx, comp)
		if err != nil {
			t.Fatalf("unexpected render error: %v", err)
		}
		if !strings.Contains(out, "text-sm") {
			t.Errorf("expected small card to contain 'text-sm', got %s", out)
		}
		if !strings.Contains(out, "data-state=\"active\"") {
			t.Errorf("expected attribute data-state, got %s", out)
		}
	})
}

func TestBadge_VariantsAndLinks(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		variant     string
		mustContain string
	}{
		{ui.BadgeVariantDefault, "bg-indigo-600"},
		{ui.BadgeVariantSecondary, "bg-slate-800"},
		{ui.BadgeVariantDestructive, "bg-rose-600"},
		{ui.BadgeVariantOutline, "border-slate-700"},
		{ui.BadgeVariantGhost, "hover:bg-slate-800"},
		{ui.BadgeVariantLink, "underline-offset-4"},
		{ui.BadgeVariantSuccess, "text-emerald-400"},
		{ui.BadgeVariantWarning, "text-amber-400"},
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

	t.Run("Badge as link with attributes", func(t *testing.T) {
		linkChild := textComponent("Documentation")
		badge := templ.ComponentFunc(func(c context.Context, w io.Writer) error {
			return ui.Badge(ui.BadgeProps{
				Href:       "/docs",
				Target:     "_blank",
				Attributes: templ.Attributes{"data-badge": "link-mode"},
			}).Render(templ.WithChildren(c, linkChild), w)
		})

		out, err := renderToString(ctx, badge)
		if err != nil {
			t.Fatalf("unexpected render error: %v", err)
		}
		if !strings.HasPrefix(strings.TrimSpace(out), "<a") {
			t.Errorf("expected link badge to render <a> tag, got: %s", out)
		}
		if !strings.Contains(out, "href=\"/docs\"") {
			t.Errorf("expected href=\"/docs\", got: %s", out)
		}
		if !strings.Contains(out, "target=\"_blank\"") {
			t.Errorf("expected target=\"_blank\", got: %s", out)
		}
		if !strings.Contains(out, "data-badge=\"link-mode\"") {
			t.Errorf("expected data-badge attribute, got: %s", out)
		}
	})
}

func TestInput_Rendering(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		props       ui.InputProps
		mustContain []string
	}{
		{
			name: "Email input",
			props: ui.InputProps{
				Type:        "email",
				ID:          "email-field",
				Name:        "email",
				Value:       "alice@acme.corp",
				Placeholder: "name@company.com",
				Required:    true,
			},
			mustContain: []string{
				"<input",
				"type=\"email\"",
				"id=\"email-field\"",
				"name=\"email\"",
				"value=\"alice@acme.corp\"",
				"placeholder=\"name@company.com\"",
				"required",
			},
		},
		{
			name: "Disabled password input with aria-invalid",
			props: ui.InputProps{
				Type:       "password",
				ID:         "pwd-field",
				Name:       "password",
				Disabled:   true,
				Attributes: templ.Attributes{"aria-invalid": "true"},
			},
			mustContain: []string{
				"type=\"password\"",
				"disabled",
				"aria-invalid=\"true\"",
				"disabled:cursor-not-allowed",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := renderToString(ctx, ui.Input(tc.props))
			if err != nil {
				t.Fatalf("unexpected render error: %v", err)
			}
			for _, s := range tc.mustContain {
				if !strings.Contains(out, s) {
					t.Errorf("input missing expected substring: %q\nOutput: %s", s, out)
				}
			}
		})
	}
}

func TestSeparator_Rendering(t *testing.T) {
	ctx := context.Background()

	t.Run("Horizontal separator", func(t *testing.T) {
		out, err := renderToString(ctx, ui.Separator(ui.SeparatorProps{Orientation: ui.OrientationHorizontal}))
		if err != nil {
			t.Fatalf("unexpected render error: %v", err)
		}
		if !strings.Contains(out, "role=\"separator\"") {
			t.Errorf("expected role=\"separator\", got: %s", out)
		}
		if !strings.Contains(out, "aria-orientation=\"horizontal\"") {
			t.Errorf("expected aria-orientation=\"horizontal\", got: %s", out)
		}
		if !strings.Contains(out, "h-[1px] w-full") {
			t.Errorf("expected horizontal classes, got: %s", out)
		}
	})

	t.Run("Vertical separator", func(t *testing.T) {
		out, err := renderToString(ctx, ui.Separator(ui.SeparatorProps{Orientation: ui.OrientationVertical}))
		if err != nil {
			t.Fatalf("unexpected render error: %v", err)
		}
		if !strings.Contains(out, "aria-orientation=\"vertical\"") {
			t.Errorf("expected aria-orientation=\"vertical\", got: %s", out)
		}
		if !strings.Contains(out, "h-full w-[1px]") {
			t.Errorf("expected vertical classes, got: %s", out)
		}
	})

	t.Run("Decorative separator", func(t *testing.T) {
		out, err := renderToString(ctx, ui.Separator(ui.SeparatorProps{Decorative: true}))
		if err != nil {
			t.Fatalf("unexpected render error: %v", err)
		}
		if !strings.Contains(out, "role=\"none\"") {
			t.Errorf("expected role=\"none\", got: %s", out)
		}
	})
}

package footer

import (
	"fmt"
	"path"
	"strings"

	bbHelp "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/dlvhdr/gh-dash/v4/config"
	"github.com/dlvhdr/gh-dash/v4/git"
	"github.com/dlvhdr/gh-dash/v4/ui/constants"
	"github.com/dlvhdr/gh-dash/v4/ui/context"
	"github.com/dlvhdr/gh-dash/v4/ui/keys"
	"github.com/dlvhdr/gh-dash/v4/utils"
)

type Model struct {
	ctx             *context.ProgramContext
	leftSection     *string
	rightSection    *string
	help            bbHelp.Model
	ShowAll         bool
	ShowConfirmQuit bool
}

func NewModel(ctx *context.ProgramContext) Model {
	help := bbHelp.New()
	help.ShowAll = true
	help.Styles = ctx.Styles.Help.BubbleStyles
	l := ""
	r := ""
	return Model{
		ctx:          ctx,
		help:         help,
		leftSection:  &l,
		rightSection: &r,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Quit):
			if m.ShowConfirmQuit {
				return m, tea.Quit
			} else {
				m.ShowConfirmQuit = true
			}
		case m.ShowConfirmQuit && !key.Matches(msg, keys.Keys.Quit):
			m.ShowConfirmQuit = false
		case key.Matches(msg, keys.Keys.Help):
			m.ShowAll = !m.ShowAll
		}
	}

	return m, nil
}

func (m Model) View() string {
	var footer string

	if m.ShowConfirmQuit {
		footer = lipgloss.NewStyle().Render("Really quit? (Press q/esc again to quit)")
	} else {
		helpIndicator := lipgloss.NewStyle().
			Background(m.ctx.Theme.FaintText).
			Foreground(m.ctx.Theme.SelectedBackground).
			Padding(0, 1).
			Render("? help")
		donationIndicator := zone.Mark("donate", lipgloss.NewStyle().
			Background(m.ctx.Theme.SelectedBackground).
			Foreground(m.ctx.Theme.WarningText).
			Padding(0, 1).
			Underline(true).
			Render(fmt.Sprintf("%s donate", constants.DonateIcon)))
		viewSwitcher := m.renderViewSwitcher(m.ctx)
		leftSection := ""
		if m.leftSection != nil {
			leftSection = *m.leftSection
		}
		rightSection := ""
		if m.rightSection != nil {
			rightSection = *m.rightSection
		}
		spacing := lipgloss.NewStyle().
			Background(m.ctx.Theme.SelectedBackground).
			Render(
				strings.Repeat(
					" ",
					utils.Max(0,
						m.ctx.ScreenWidth-lipgloss.Width(
							viewSwitcher,
						)-lipgloss.Width(leftSection)-
							lipgloss.Width(rightSection)-
							lipgloss.Width(
								helpIndicator,
							)-lipgloss.Width(donationIndicator),
					)))

		footer = m.ctx.Styles.Common.FooterStyle.
			Render(lipgloss.JoinHorizontal(lipgloss.Top, viewSwitcher, leftSection, spacing, rightSection, donationIndicator, helpIndicator))
	}

	if m.ShowAll {
		keymap := keys.CreateKeyMapForView(m.ctx.View)
		fullHelp := m.help.View(keymap)
		return lipgloss.JoinVertical(lipgloss.Top, footer, fullHelp)
	}

	return footer
}

func (m *Model) SetWidth(width int) {
	m.help.Width = width
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	m.help.Styles = ctx.Styles.Help.BubbleStyles
}

func (m *Model) renderViewButton(view config.ViewType) string {
	v := " PRs"
	if view == config.IssuesView {
		v = " Issues"
	}

	if m.ctx.View == view {
		return m.ctx.Styles.ViewSwitcher.ActiveView.Render(v)

	}
	return m.ctx.Styles.ViewSwitcher.InactiveView.Render(v)
}

func (m *Model) renderViewSwitcher(ctx *context.ProgramContext) string {
	var repo string
	if m.ctx.RepoPath != "" {
		name := path.Base(m.ctx.RepoPath)
		if m.ctx.RepoUrl != "" {
			name = git.GetRepoShortName(m.ctx.RepoUrl)
		}
		repo = ctx.Styles.Common.FooterStyle.Render(fmt.Sprintf(" %s", name))
	}

	var user string
	if ctx.User != "" {
		user = ctx.Styles.Common.FooterStyle.Render("@" + ctx.User)
	}

	view := lipgloss.JoinHorizontal(
		lipgloss.Top,
		ctx.Styles.ViewSwitcher.ViewsSeparator.PaddingLeft(1).Render(m.renderViewButton(config.PRsView)),
		ctx.Styles.ViewSwitcher.ViewsSeparator.Render(" │ "),
		m.renderViewButton(config.IssuesView),
		lipgloss.NewStyle().Background(ctx.Styles.Common.FooterStyle.GetBackground()).Foreground(ctx.Styles.ViewSwitcher.ViewsSeparator.GetBackground()).Render(" "),
		repo,
		ctx.Styles.Common.FooterStyle.Foreground(m.ctx.Theme.FaintText).Render(" • "),
		user,
		ctx.Styles.Common.FooterStyle.Foreground(m.ctx.Theme.FaintBorder).Render(" │"),
	)

	return ctx.Styles.ViewSwitcher.Root.Render(view)
}

func (m *Model) SetLeftSection(leftSection string) {
	*m.leftSection = leftSection
}

func (m *Model) SetRightSection(rightSection string) {
	*m.rightSection = rightSection
}

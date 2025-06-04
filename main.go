package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/urfave/cli/v2"
)

var (
	version = "1.0.1"
	commit  = "dev"
	author  = "Swan Htet Aung Phyo"
	title   = "Computer Science Student at AGH, Poland - Backend Developer"
)

func SafeHashShort(hash string, n int) string {
	if len(hash) < n {
		return hash
	}
	return hash[:n]
}

type CommitInfo struct {
	Hash      string
	Message   string
	Author    string
	Email     string
	Timestamp time.Time
	Parents   []string
	Changes   int
	Files     []string
	IsMerge   bool
}

type Visualizer struct {
	repo      *git.Repository
	commits   []CommitInfo
	branchMap map[string][]string
	commitMap map[string]*object.Commit
	authors   map[string]int
	stats     struct {
		totalFilesChanged int
		totalAdditions    int
		totalDeletions    int
	}
}

func NewVisualizer(repoPath string) (*Visualizer, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}
	return &Visualizer{
		repo:      repo,
		branchMap: make(map[string][]string),
		commitMap: make(map[string]*object.Commit),
		authors:   make(map[string]int),
	}, nil
}

func (v *Visualizer) LoadCommits(limit int) error {
	ref, err := v.repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	iter, err := v.repo.Log(&git.LogOptions{
		From:  ref.Hash(),
		All:   true,
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return fmt.Errorf("failed to get commit log: %w", err)
	}

	count := 0
	err = iter.ForEach(func(c *object.Commit) error {
		if limit > 0 && count >= limit {
			return nil
		}
		count++

		parents := make([]string, len(c.ParentHashes))
		for i, h := range c.ParentHashes {
			parents[i] = SafeHashShort(h.String(), 7)
		}

		changes, files, additions, deletions, err := countChanges(c, v.repo)
		if err != nil {
			return fmt.Errorf("failed to count changes: %w", err)
		}

		v.stats.totalFilesChanged += changes
		v.stats.totalAdditions += additions
		v.stats.totalDeletions += deletions

		commitHash := SafeHashShort(c.Hash.String(), 7)
		v.commits = append(v.commits, CommitInfo{
			Hash:      commitHash,
			Message:   strings.Split(c.Message, "\n")[0],
			Author:    c.Author.Name,
			Email:     c.Author.Email,
			Timestamp: c.Author.When,
			Parents:   parents,
			Changes:   changes,
			Files:     files,
			IsMerge:   len(parents) > 1,
		})

		v.commitMap[commitHash] = c
		v.authors[c.Author.Email]++
		return nil
	})

	return err
}

func countChanges(c *object.Commit, repo *git.Repository) (int, []string, int, int, error) {
	if len(c.ParentHashes) == 0 {
		return 0, nil, 0, 0, nil
	}

	parent, err := c.Parents().Next()
	if err != nil {
		return 0, nil, 0, 0, err
	}

	tree, err := c.Tree()
	if err != nil {
		return 0, nil, 0, 0, err
	}

	parentTree, err := parent.Tree()
	if err != nil {
		return 0, nil, 0, 0, err
	}

	changes, err := parentTree.Diff(tree)
	if err != nil {
		return 0, nil, 0, 0, err
	}

	count := 0
	var files []string
	totalAdditions := 0
	totalDeletions := 0

	for _, ch := range changes {
		patch, err := ch.Patch()
		if err != nil {
			continue
		}

		stats := patch.Stats()
		for _, stat := range stats {
			totalAdditions += stat.Addition
			totalDeletions += stat.Deletion
		}

		_, to, err := ch.Files()
		if err != nil {
			continue
		}

		if to != nil {
			count++
			files = append(files, to.Name)
		}
	}

	return count, files, totalAdditions, totalDeletions, nil
}

func (v *Visualizer) LoadBranches() error {
	refs, err := v.repo.References()
	if err != nil {
		return fmt.Errorf("failed to get references: %w", err)
	}

	return refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			hash := SafeHashShort(ref.Hash().String(), 7)
			name := ref.Name().Short()
			v.branchMap[hash] = append(v.branchMap[hash], name)
		}
		return nil
	})
}

func (v *Visualizer) DisplayGraph(compact bool) {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#5A56E0")).
		Padding(0, 1).
		Bold(true).
		Render(" Git Repository Visualizer ")

	fmt.Println(style.Render(header))
	fmt.Println()

	for i, commit := range v.commits {
		v.printCommit(commit, compact)

		if i < len(v.commits)-1 {
			if !compact {
				fmt.Println(color.CyanString("│"))
				fmt.Println(color.CyanString("▼"))
			} else {
				fmt.Print(" → ")
			}
		}
	}

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5A56E0")).
		Render(fmt.Sprintf(" Showing %d commits ", len(v.commits)))

	fmt.Println()
	fmt.Println(footer)
}

func (v *Visualizer) printCommit(commit CommitInfo, compact bool) {
	hashColor := color.New(color.FgHiRed, color.Bold)
	authorColor := color.New(color.FgHiGreen)
	messageColor := color.New(color.FgHiWhite)
	branchColor := color.New(color.FgHiYellow)
	fileColor := color.New(color.FgHiBlue)
	metaColor := color.New(color.FgHiCyan)
	mergeColor := color.New(color.FgHiMagenta, color.Bold)

	if compact {
		// Compact view
		hashColor.Printf("%s", commit.Hash)
		if commit.IsMerge {
			mergeColor.Print(" (merge)")
		}
		fmt.Print(" ")
		messageColor.Printf("%.50s", commit.Message)
		if len(commit.Message) > 50 {
			fmt.Print("...")
		}
		fmt.Print(" ")
		authorColor.Printf("(%s)", commit.Author)
		if branches := v.branchMap[commit.Hash]; len(branches) > 0 {
			branchColor.Printf(" [%s]", strings.Join(branches, ", "))
		}
	} else {
		// Detailed view
		hashColor.Printf("Commit: %s\n", commit.Hash)
		if commit.IsMerge {
			mergeColor.Println("Merge commit")
		}
		authorColor.Printf("Author: %s <%s>\n", commit.Author, commit.Email)
		metaColor.Printf("Date:   %s\n", commit.Timestamp.Format("Mon Jan 2 15:04:05 2006 -0700"))
		fmt.Println()
		messageColor.Printf("    %s\n", commit.Message)
		fmt.Println()

		if branches := v.branchMap[commit.Hash]; len(branches) > 0 {
			branchColor.Printf("Branches: %s\n", strings.Join(branches, ", "))
		}

		fileColor.Printf("Files changed: %d\n", commit.Changes)
		for _, file := range commit.Files {
			fmt.Printf("    %s\n", file)
		}

		if len(commit.Parents) > 0 {
			metaColor.Printf("Parents: %s\n", strings.Join(commit.Parents, ", "))
		}
		fmt.Println()
	}
}

func (v *Visualizer) DisplayStats() {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#5A56E0")).
		Padding(0, 1).
		Bold(true).
		Render(" Repository Statistics ")

	content := lipgloss.NewStyle().Padding(0, 1)

	stats := []string{
		fmt.Sprintf("Total Commits:    %d", len(v.commits)),
		fmt.Sprintf("Files Changed:    %d", v.stats.totalFilesChanged),
		fmt.Sprintf("Lines Added:      %d", v.stats.totalAdditions),
		fmt.Sprintf("Lines Deleted:    %d", v.stats.totalDeletions),
		"",
		"Top Authors:",
	}

	authors := make([]struct {
		Email string
		Count int
	}, 0, len(v.authors))
	for email, count := range v.authors {
		authors = append(authors, struct {
			Email string
			Count int
		}{email, count})
	}
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].Count > authors[j].Count
	})

	for i, author := range authors {
		if i >= 5 { // Limit to top 5 authors
			break
		}
		stats = append(stats, fmt.Sprintf("- %s: %d commits", author.Email, author.Count))
	}

	fmt.Println(style.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			content.Render(strings.Join(stats, "\n")),
		)))
}

func (v *Visualizer) DisplayTimeline() {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#5A56E0")).
		Padding(0, 1).
		Bold(true).
		Render(" Commit Timeline ")

	// Group commits by week
	weeklyCommits := make(map[string]int)
	for _, commit := range v.commits {
		year, week := commit.Timestamp.ISOWeek()
		key := fmt.Sprintf("%d-W%02d", year, week)
		weeklyCommits[key]++
	}

	// Convert to slice for sorting
	var weeks []string
	for week := range weeklyCommits {
		weeks = append(weeks, week)
	}
	sort.Strings(weeks)

	// Prepare timeline content
	var timeline []string
	for _, week := range weeks {
		count := weeklyCommits[week]
		bar := strings.Repeat("■", count)
		if len(bar) > 50 {
			bar = bar[:50] + "..."
		}
		timeline = append(timeline, fmt.Sprintf("%s: %s %d commits", week, bar, count))
	}

	fmt.Println(style.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			lipgloss.NewStyle().Padding(0, 1).Render(strings.Join(timeline, "\n")),
		)))
}

func main() {
	authorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF7F50")).
		Bold(true).
		Underline(true)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6495ED")).
		Italic(true)

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32CD32")).
		Bold(true)

	buildStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Bold(true)

	app := &cli.App{
		Name:  "gitviz",
		Usage: "Visualize Git repositories with style",
		Version: fmt.Sprintf("%s\n%s\n%s",
			versionStyle.Render(fmt.Sprintf("Version: %s", version)),
			buildStyle.Render(fmt.Sprintf("Build: %s", SafeHashShort(commit, 7))),
			lipgloss.JoinVertical(lipgloss.Left,
				authorStyle.Render(fmt.Sprintf("Developer: %s", author)),
				titleStyle.Render(title),
			),
		),
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"n"},
				Usage:   "Limit number of commits to display",
				Value:   50,
			},
			&cli.BoolFlag{
				Name:    "compact",
				Aliases: []string{"c"},
				Usage:   "Compact display mode",
			},
			&cli.BoolFlag{
				Name:    "stats",
				Aliases: []string{"s"},
				Usage:   "Show only statistics",
			},
			&cli.BoolFlag{
				Name:    "timeline",
				Aliases: []string{"t"},
				Usage:   "Show commit timeline",
			},
		},
		Action: func(c *cli.Context) error {
			vis, err := NewVisualizer(".")
			if err != nil {
				return cli.Exit(fmt.Sprintf("Error opening repository: %v", err), 1)
			}

			if err := vis.LoadCommits(c.Int("limit")); err != nil {
				return cli.Exit(fmt.Sprintf("Error loading commits: %v", err), 1)
			}

			if err := vis.LoadBranches(); err != nil {
				return cli.Exit(fmt.Sprintf("Error loading branches: %v", err), 1)
			}

			if c.Bool("stats") {
				vis.DisplayStats()
				return nil
			}

			if c.Bool("timeline") {
				vis.DisplayTimeline()
				return nil
			}

			vis.DisplayGraph(c.Bool("compact"))
			fmt.Println()
			vis.DisplayStats()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

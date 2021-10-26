package graph

import (
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestRenderCommitGraph(t *testing.T) {
	tests := []struct {
		name           string
		commits        []*models.Commit
		expectedOutput string
	}{
		{
			name: "with some merges",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"3"}},
				{Sha: "3", Parents: []string{"4"}},
				{Sha: "4", Parents: []string{"5", "7"}},
				{Sha: "7", Parents: []string{"5"}},
				{Sha: "5", Parents: []string{"8"}},
				{Sha: "8", Parents: []string{"9"}},
				{Sha: "9", Parents: []string{"A", "B"}},
				{Sha: "B", Parents: []string{"D"}},
				{Sha: "D", Parents: []string{"D"}},
				{Sha: "A", Parents: []string{"E"}},
				{Sha: "E", Parents: []string{"F"}},
				{Sha: "F", Parents: []string{"D"}},
				{Sha: "D", Parents: []string{"G"}},
			},
			expectedOutput: `
			1 ⎔
			2 ⎔
			3 ⎔
			4 ⏣─┐
			7 │ ⎔
			5 ⎔─┘
			8 ⎔
			9 ⏣─┐
			B │ ⎔
			D │ ⎔
			A ⎔ │
			E ⎔ │
			F ⎔ │
			D ⎔─┘`,
		},
		{
			name: "with a path that has room to move to the left",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"3", "4"}},
				{Sha: "4", Parents: []string{"3", "5"}},
				{Sha: "3", Parents: []string{"5"}},
				{Sha: "5", Parents: []string{"6"}},
				{Sha: "6", Parents: []string{"7"}},
			},
			expectedOutput: `
			1 ⎔
			2 ⏣─┐
			4 │ ⏣─┐
			3 ⎔─┘ │
			5 ⎔───┘
			6 ⎔`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"3", "4"}},
				{Sha: "3", Parents: []string{"5", "4"}},
				{Sha: "5", Parents: []string{"7", "8"}},
				{Sha: "4", Parents: []string{"7"}},
				{Sha: "7", Parents: []string{"11"}},
			},
			expectedOutput: `
			1 ⎔
			2 ⏣─┐
			3 ⏣─│─┐
			5 ⏣─│─│─┐
			4 │ ⎔─┘ │
			7 ⎔─┘ ┌─┘`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"3", "4"}},
				{Sha: "3", Parents: []string{"5", "4"}},
				{Sha: "5", Parents: []string{"7", "8"}},
				{Sha: "7", Parents: []string{"4", "A"}},
				{Sha: "4", Parents: []string{"B"}},
				{Sha: "B", Parents: []string{"C"}},
			},
			expectedOutput: `
			1 ⎔
			2 ⏣─┐
			3 ⏣─│─┐
			5 ⏣─│─│─┐
			7 ⏣─│─│─│─┐
			4 ⎔─┴─┘ │ │
			B ⎔ ┌───┘ │`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2", "3"}},
				{Sha: "3", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"4", "5"}},
				{Sha: "4", Parents: []string{"6", "7"}},
				{Sha: "6", Parents: []string{"8"}},
			},
			expectedOutput: `
			1 ⏣─┐
			3 │ ⎔
			2 ⏣─│
			4 ⏣─│─┐
			6 ⎔ │ │`,
		},
		{
			name: "new merge path fills gap before continuing path on right",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2", "3", "4", "5"}},
				{Sha: "4", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"A"}},
				{Sha: "A", Parents: []string{"6", "B"}},
				{Sha: "B", Parents: []string{"C"}},
			},
			expectedOutput: `
			1 ⏣─┬─┬─┐
			4 │ │ ⎔ │
			2 ⎔─│─┘ │
			A ⏣─│─┐ │
			B │ │ ⎔ │`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			getStyle := func(c *models.Commit) style.TextStyle { return style.FgDefault }
			_, lines, _, _ := RenderCommitGraph(test.commits, &models.Commit{Sha: "blah"}, getStyle)
			output := ""
			for i, line := range lines {
				description := test.commits[i].Sha
				output += strings.TrimSpace(description+" "+utils.Decolorise(line)) + "\n"
			}
			t.Log("\n" + output)

			trimmedExpectedOutput := ""
			for _, line := range strings.Split(strings.TrimPrefix(test.expectedOutput, "\n"), "\n") {
				trimmedExpectedOutput += strings.TrimSpace(line) + "\n"
			}

			assert.Equal(t,
				trimmedExpectedOutput,
				output)
		})
	}
}

// // func TestGetCellsFromPipeSet(t *testing.T) {
// // 	tests := []struct {
// // 		pipeSet       PipeSet
// // 		expectedCells []*Cell
// // 	}{
// // {
// // 	pipeSet: PipeSet{
// // 		pipes: []Pipe{
// // 			{
// // 				fromPos:         0,
// // 				toPos:           0,
// // 				kind:            STARTS,
// // 				style:           style.FgDefault,
// // 				sourceCommitSha: "a",
// // 			},
// // 			{
// // 				fromPos:         0,
// // 				toPos:           0,
// // 				kind:            TERMINATES,
// // 				style:           style.FgDefault,
// // 				sourceCommitSha: "b",
// // 			},
// // 		},
// // 		isMerge: false,
// // 	},
// // 	expectedCells: []*Cell{
// // 		{
// // 			up:       true,
// // 			down:     true,
// // 			cellType: COMMIT,
// // 			style:    style.FgDefault,
// // 		},
// // 	},
// // },
// // {
// // 	pipeSet: PipeSet{
// // 		pipes: []Pipe{
// // 			{
// // 				fromPos:         0,
// // 				toPos:           0,
// // 				kind:            CONTINUES,
// // 				style:           style.FgDefault,
// // 				sourceCommitSha: "a",
// // 			},
// // 			{
// // 				fromPos:         1,
// // 				toPos:           1,
// // 				kind:            TERMINATES,
// // 				style:           style.FgDefault,
// // 				sourceCommitSha: "a",
// // 			},
// // 			{
// // 				fromPos:         1,
// // 				toPos:           1,
// // 				kind:            STARTS,
// // 				style:           style.FgDefault,
// // 				sourceCommitSha: "b",
// // 			},
// // 		},
// // 		isMerge: false,
// // 	},
// // 	expectedCells: []*Cell{
// // 		{
// // 			up:       true,
// // 			down:     true,
// // 			cellType: CONNECTION,
// // 			style:    style.FgDefault,
// // 		},
// // 		{
// // 			up:       true,
// // 			down:     true,
// // 			cellType: COMMIT,
// // 			style:    style.FgDefault,
// // 		},
// // 	},
// // },
// // 		{
// // 			pipeSet: PipeSet{
// // 				pipes: []Pipe{
// // 					{
// // 						fromPos:         0,
// // 						toPos:           0,
// // 						kind:            TERMINATES,
// // 						style:           style.FgDefault,
// // 						sourceCommitSha: "a",
// // 					},
// // 					{
// // 						fromPos:         0,
// // 						toPos:           0,
// // 						kind:            STARTS,
// // 						style:           style.FgDefault,
// // 						sourceCommitSha: "b",
// // 					},
// // 					{
// // 						fromPos:         0,
// // 						toPos:           2,
// // 						kind:            STARTS,
// // 						style:           style.FgDefault,
// // 						sourceCommitSha: "b",
// // 					},
// // 					{
// // 						fromPos:         1,
// // 						toPos:           1,
// // 						kind:            CONTINUES,
// // 						style:           style.FgDefault,
// // 						sourceCommitSha: "c",
// // 					},
// // 				},
// // 				isMerge: true,
// // 			},
// // 			expectedCells: []*Cell{
// // 				{
// // 					up:         true,
// // 					down:       true,
// // 					right:      true,
// // 					cellType:   MERGE,
// // 					style:      style.FgDefault,
// // 					rightStyle: &style.FgDefault,
// // 				},
// // 				{
// // 					up:         true,
// // 					down:       true,
// // 					left:       true,
// // 					right:      true,
// // 					cellType:   CONNECTION,
// // 					style:      style.FgDefault,
// // 					rightStyle: &style.FgDefault,
// // 				},
// // 				{
// // 					down:     true,
// // 					left:     true,
// // 					cellType: CONNECTION,
// // 					style:    style.FgDefault,
// // 				},
// // 			},
// // 		},
// // 		{
// // 			pipeSet: PipeSet{
// // 				pipes: []Pipe{
// // 					{
// // 						fromPos:         0,
// // 						toPos:           0,
// // 						kind:            TERMINATES,
// // 						style:           style.FgDefault,
// // 						sourceCommitSha: "a",
// // 					},
// // 					{
// // 						fromPos:         0,
// // 						toPos:           0,
// // 						kind:            STARTS,
// // 						style:           style.FgCyan,
// // 						sourceCommitSha: "selected",
// // 					},
// // 					{
// // 						fromPos:         0,
// // 						toPos:           2,
// // 						kind:            STARTS,
// // 						style:           style.FgCyan,
// // 						sourceCommitSha: "selected",
// // 					},
// // 					{
// // 						fromPos:         1,
// // 						toPos:           1,
// // 						kind:            CONTINUES,
// // 						style:           style.FgDefault,
// // 						sourceCommitSha: "c",
// // 					},
// // 				},
// // 				isMerge: true,
// // 			},
// // 			expectedCells: []*Cell{
// // 				{
// // 					up:         false,
// // 					down:       true,
// // 					right:      true,
// // 					cellType:   MERGE,
// // 					style:      style.FgCyan,
// // 					rightStyle: &style.FgCyan,
// // 				},
// // 				{
// // 					up:         false,
// // 					down:       false,
// // 					left:       true,
// // 					right:      true,
// // 					cellType:   CONNECTION,
// // 					style:      style.FgCyan,
// // 					rightStyle: &style.FgCyan,
// // 				},
// // 				{
// // 					down:     true,
// // 					left:     true,
// // 					cellType: CONNECTION,
// // 					style:    style.FgCyan,
// // 				},
// // 			},
// // 		},
// // 		{
// // 			pipeSet: PipeSet{
// // 				pipes: []Pipe{
// // 					{
// // 						fromPos:         0,
// // 						toPos:           0,
// // 						kind:            TERMINATES,
// // 						style:           style.FgGreen,
// // 						sourceCommitSha: "a",
// // 					},
// // 					{
// // 						fromPos:         0,
// // 						toPos:           0,
// // 						kind:            STARTS,
// // 						style:           style.FgYellow,
// // 						sourceCommitSha: "b",
// // 					},
// // 					{
// // 						fromPos:         0,
// // 						toPos:           2,
// // 						kind:            STARTS,
// // 						style:           style.FgYellow,
// // 						sourceCommitSha: "b",
// // 					},
// // 					{
// // 						fromPos:         1,
// // 						toPos:           1,
// // 						kind:            CONTINUES,
// // 						style:           style.FgDefault,
// // 						sourceCommitSha: "a",
// // 					},
// // 				},
// // 				isMerge: true,
// // 			},
// // 			expectedCells: []*Cell{
// // 				{
// // 					up:         true,
// // 					down:       true,
// // 					right:      true,
// // 					cellType:   MERGE,
// // 					style:      style.FgYellow,
// // 					rightStyle: &style.FgYellow,
// // 				},
// // 				{
// // 					up:         true,
// // 					down:       true,
// // 					left:       true,
// // 					right:      true,
// // 					cellType:   CONNECTION,
// // 					style:      style.FgGreen,
// // 					rightStyle: &style.FgYellow,
// // 				},
// // 				{
// // 					down:     true,
// // 					left:     true,
// // 					cellType: CONNECTION,
// // 					style:    style.FgYellow,
// // 				},
// // 			},
// // 		},
// // 	}

// // 	for _, test := range tests {
// // 		cells := getCellsFromPipeSet(test.pipeSet, "selected")
// // 		if len(cells) != len(test.expectedCells) {
// // 			t.Errorf("expected cells to be %s, got %s", spew.Sdump(test.expectedCells), spew.Sdump(cells))
// // 			continue
// // 		}
// // 		t.Log(spew.Sdump(cells))
// // 		for i, cell := range cells {
// // 			assert.EqualValues(t, test.expectedCells[i], cell)
// // 		}
// // 	}
// // }

// // func TestCellRender(t *testing.T) {
// // 	tests := []struct {
// // 		cell           *Cell
// // 		expectedString string
// // 	}{
// // 		{
// // 			cell: &Cell{
// // 				up:       true,
// // 				down:     true,
// // 				cellType: CONNECTION,
// // 				style:    style.FgDefault,
// // 			},
// // 			expectedString: "\x1b[39m│\x1b[0m\x1b[39m \x1b[0m",
// // 		},
// // 		{
// // 			cell: &Cell{
// // 				up:       true,
// // 				down:     true,
// // 				cellType: COMMIT,
// // 				style:    style.FgDefault,
// // 			},
// // 			expectedString: "\x1b[39m⎔\x1b[0m\x1b[39m \x1b[0m",
// // 		},
// // 	}

// // 	for _, test := range tests {
// // 		assert.EqualValues(t, test.expectedString, test.cell.render())
// // 	}
// // }

func TestGetNextPipes(t *testing.T) {
	tests := []struct {
		prevPipes []Pipe
		commit    *models.Commit
		expected  []Pipe
	}{
		{
			prevPipes: []Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a", toSha: "b", kind: STARTS, style: style.FgDefault},
			},
			commit: &models.Commit{
				Sha:     "b",
				Parents: []string{"c"},
			},
			expected: []Pipe{
				{fromPos: 0, toPos: 0, fromSha: "b", toSha: "c", kind: STARTS, style: style.FgDefault},
				{fromPos: 0, toPos: 0, fromSha: "a", toSha: "b", kind: TERMINATES, style: style.FgDefault},
			},
		},
		{
			prevPipes: []Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a", toSha: "b", kind: TERMINATES, style: style.FgDefault},
				{fromPos: 0, toPos: 0, fromSha: "b", toSha: "c", kind: STARTS, style: style.FgDefault},
				{fromPos: 0, toPos: 1, fromSha: "b", toSha: "d", kind: STARTS, style: style.FgDefault},
			},
			commit: &models.Commit{
				Sha:     "d",
				Parents: []string{"e"},
			},
			expected: []Pipe{
				{fromPos: 0, toPos: 0, fromSha: "b", toSha: "c", kind: CONTINUES, style: style.FgDefault},
				{fromPos: 1, toPos: 1, fromSha: "d", toSha: "e", kind: STARTS, style: style.FgDefault},
				{fromPos: 1, toPos: 1, fromSha: "b", toSha: "d", kind: TERMINATES, style: style.FgDefault},
			},
		},
	}

	for _, test := range tests {
		getStyle := func(c *models.Commit) style.TextStyle { return style.FgDefault }
		pipes := getNextPipes(test.prevPipes, test.commit, getStyle)
		assert.EqualValues(t, test.expected, pipes)
	}
}
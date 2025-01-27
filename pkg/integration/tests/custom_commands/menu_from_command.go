package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// NOTE: we're getting a weird offset in the popup prompt for some reason. Not sure what's behind that.

var MenuFromCommand = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using menuFromCommand prompt type",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("foo").
			EmptyCommit("bar").
			EmptyCommit("baz").
			NewBranch("feature/foo")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "localBranches",
				Command: `echo "{{index .PromptResponses 0}} {{index .PromptResponses 1}} {{ .SelectedLocalBranch.Name }}" > output.txt`,
				Prompts: []config.CustomCommandPrompt{
					{
						Type:        "menuFromCommand",
						Title:       "Choose commit message",
						Command:     `git log --oneline --pretty=%B`,
						Filter:      `(?P<commit_message>.*)`,
						ValueFormat: `{{ .commit_message }}`,
						LabelFormat: `{{ .commit_message | yellow }}`,
					},
					{
						Type:         "input",
						Title:        "Description",
						InitialValue: `{{ if .SelectedLocalBranch.Name }}Branch: #{{ .SelectedLocalBranch.Name }}{{end}}`,
					},
				},
			},
		}
	},
	Run: func(
		shell *Shell,
		input *Input,
		assert *Assert,
		keys config.KeybindingConfig,
	) {
		assert.WorkingTreeFileCount(0)
		input.SwitchToBranchesView()

		input.Press("a")

		input.Menu(Equals("Choose commit message"), Contains("bar"))

		input.Prompt(Equals("Description"), " my branch")

		input.SwitchToFilesView()

		assert.WorkingTreeFileCount(1)
		assert.CurrentView().SelectedLine(Contains("output.txt"))
		assert.MainView().Content(Contains("bar Branch: #feature/foo my branch feature/foo"))
	},
})

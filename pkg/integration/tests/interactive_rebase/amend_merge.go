package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var (
	postMergeFileContent = "post merge file content"
	postMergeFilename    = "post-merge-file"
)

var AmendMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Amends a staged file to a merge commit.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("development-branch").
			CreateFileAndAdd("initial-file", "initial file content").
			Commit("initial commit").
			NewBranch("feature-branch"). // it's also checked out automatically
			CreateFileAndAdd("new-feature-file", "new content").
			Commit("new feature commit").
			Checkout("development-branch").
			Merge("feature-branch").
			CreateFileAndAdd(postMergeFilename, postMergeFileContent)
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.CommitCount(3)

		input.SwitchToCommitsView()

		mergeCommitMessage := "Merge branch 'feature-branch' into development-branch"
		assert.HeadCommitMessage(Contains(mergeCommitMessage))

		input.Press(keys.Commits.AmendToCommit)
		input.AcceptConfirmation(Equals("Amend Commit"), Contains("Are you sure you want to amend this commit with your staged files?"))

		// assuring we haven't added a brand new commit
		assert.CommitCount(3)
		assert.HeadCommitMessage(Contains(mergeCommitMessage))

		// assuring the post-merge file shows up in the merge commit.
		assert.MainView().
			Content(Contains(postMergeFilename)).
			Content(Contains("++" + postMergeFileContent))
	},
})

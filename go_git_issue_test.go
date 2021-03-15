package go_git_issue

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/require"
)

// This test passes.
func Test_push_to_local_remote_with_1_commits_cloned_with_depth_of_1(t *testing.T) {
	remoteRepo, dir := initRepo(t, false)
	commit(t, remoteRepo)

	repo := cloneRepo(t, dir, 1)
	commit(t, repo)
	push(t, repo)
}

// This test passes.
func Test_push_to_local_remote_with_1_commits_cloned_deep(t *testing.T) {
	remoteRepo, dir := initRepo(t, false)
	commit(t, remoteRepo)

	repo := cloneRepo(t, dir, 0)
	commit(t, repo)
	push(t, repo)
}

// This test fails
func Test_push_to_local_remote_with_2_commits_cloned_with_depth_of_1(t *testing.T) {
	remoteRepo, dir := initRepo(t, false)
	commit(t, remoteRepo)

	// This commit is required to cause the failure.
	commit(t, remoteRepo)

	repo := cloneRepo(t, dir, 1)
	commit(t, repo)

	// Fails with: `Received unexpected error: object not found`
	push(t, repo)
}

// This test passes.
func Test_push_to_local_remote_with_2_commits_cloned_deep(t *testing.T) {
	remoteRepo, dir := initRepo(t, false)
	commit(t, remoteRepo)
	commit(t, remoteRepo)

	repo := cloneRepo(t, dir, 0)
	commit(t, repo)
	push(t, repo)
}

func initRepo(t *testing.T, bare bool) (*git.Repository, string){
	remoteRepoDir, err := ioutil.TempDir("", t.Name())
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(remoteRepoDir))
	})

	repo, err := git.PlainInit(remoteRepoDir, bare)
	require.NoError(t, err)

	return repo, remoteRepoDir
}

func cloneRepo(t *testing.T, repoToClone string, depth int) *git.Repository {
	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:           repoToClone,
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName("master"),
		SingleBranch:  false,
		Depth: depth,
	})
	require.NoError(t, err)
	return repo
}

func commit(t *testing.T, repo *git.Repository) {
	wt, err := repo.Worktree()
	require.NoError(t, err)

	_, err = wt.Commit("Commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "",
			Email: "",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)
}

func push(t *testing.T, repo *git.Repository) {
	err := repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"refs/heads/master:refs/heads/master"},
	})
	require.NoError(t, err)
}

title = nodes git
date = 2019-06-07
desc = notes on git

---

#### git

git is hard

##### repositories

have 2 parts: `.git/`: where git stuff is stored,
and the 'working tree': a visible, working copy of everything, the files you work on

##### state

- history: committed changes
- staging / index: tracked changes
- working directory: untracked changes

##### git clone

`--recurse-submodules`: include any submodules

##### git init

`--bare`: no working tree,
store the git stuff directly in the directory,
conventionally the dir has the `.git` extenstion: `repo.git`

`--template`, or the config `init.templateDir`:
copy the non dotfile contents of a directory

##### git reset

`git reset --soft $commit`:
uncommit, uncommitted things staged, working directory untouched

(default) `git reset --mixed $commit`:
uncommit, uncommitted things unstaged, working directory untouched

`git reset --hard $commit`:
uncommit, uncommitted things discarded, working directory changes discarded

##### git bisect

binary search through history to find first commit that introduced a change

- `git bisect start`
- `git bisect (good/old) / (bad/new) / skip`
- `git bisect log / visualize / view`
- `git bisect reset`
- `git bisect replay`: replay the output of `git bisect log` and continue where it left off
- `git bisect run $script`: automate with the use of script

##### git branch

branch listing and management

##### git checkout

`git checkout $branch`: switch to branch, keep uncommitted changes

`git checkout -b $branch`: create a new branch

`git checkout $branch -- $paths...`: change `$paths` to match from `$branch`

##### git merge

merge changes from \$branch onto current branch in a new commit

`--squash`: create a new commit with the same effect as a merge

`--allow-unrelated-histories`: merge different porjects into one

##### git pull

get and apply remote changes

`-r`: rebase local onto remote changes

`--autostash`: auto stash and unstash changes

##### git push

`--tags`: push tags

`-u`: set upstream

##### git cherry-pick

on a clean working tree,
apply changes from commits

##### git clean

remove untracked files

`-x`: also ignore .gitignored files

`-X`: only remove ignored files

`-d`: remove directories, also requires `-f`

##### git revert

ass a new commit changing state to \$commit

##### git stash

`git stash (push)`: store changes, reset to HEAD

`git stash apply`: replay changes

##### git blame

show who last touched each line

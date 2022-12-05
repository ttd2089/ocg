# ocg

Obsessive compulsive Git -- A tool to help make sure your local repos are synced and merged

## Why?

With many repos and many projects on the go it's easy to lose track of what's in flight. OCG aims to combat the problem by making it easy to get a complete summary of every repo in your `src` directory -- you do keep them all together right? -- and which ones have unfinished work.

Currently there is just one command, `ocg list`, which prints a YAML summary of all branches in all repos. The branch info includes the name, the SHA, and the tracked remote branch name and SHA if applicable. The next iterations will capture and present the following details as a quickly recognizable status:

- Local branches with no tracked remote
- Local branches that are ahead of tracked branches
- Local branches that are behind tracked branches
- Branches that are not merged to the default branch of the remote

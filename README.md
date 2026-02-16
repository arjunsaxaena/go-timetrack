# go-timetrack

A personalized time tracker CLI built with Cobra, inspired by Toggl Track.

## Install `tt` globally

From the project root:

```bash
go install .
```

Make sure your Go bin is in your `PATH` (usually `$HOME/go/bin`):

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

Verify:

```bash
tt --help
```

## Enable autocomplete

### Bash

```bash
tt completion bash > ~/.tt-completion.bash
echo 'source ~/.tt-completion.bash' >> ~/.bashrc
source ~/.bashrc
```

### Zsh

```bash
tt completion zsh > "${fpath[1]}/_tt"
```

### Fish

```bash
tt completion fish > ~/.config/fish/completions/tt.fish
```

## Collaboration and issues

If you find a bug, want a feature, or want to collaborate, open an issue (or PR) in this repository.
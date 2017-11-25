# kc
a wrapper for a few kubectl commands

## Background

kc is a wrapper around a kubernetes utility called 'kubectl'

It simplifies a few operations which I tend to use a lot

 * log/tail
 * exec
 * shell
 * get versions of pods
 * apply/replace
 * 'bounce' (scale down & up)

It's not meant to be a full replacement or in the least bit comprehensive. If you want it do be different, I recommend you fork it and bend it to your own will. PRs would also be great, but I'd like to see other people twist it all around for different use cases.

### Pod selectors

When you need to specify a pod, it can optionally look up by a kubernetes 'selector'
kc will look up a pod by 'selector' whenever there's an '=' in the pod's 'name'

NOTE: for now, it just takes the first result. In the future I'd like to offer either an interactive 'dropdown', or (for tailing), multiplex logs

### Context

kc commands all take a -c<context>

## Subcommands

### Exec

The main point of `kc exec` / `kc sh` is to add the usual flags (e.g. -it, `--`), (bash), and also to allow selector-based querying

#### Examples
 kc l -t10 -f podname
 kc t name=mysvc

 kubectl t name
 kc l -t10 -f name=mysvc

### Logging

The main point of `kc log` / `kc tail` is to use selector rather than pod name

#### Examples
 kc l -t10 -f podname
 kc t name=mysvc

 kubectl t name
 kc l -t10 -f name=mysvc


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

kc commands all take a -c<context>. See ~/.kube/config for your context definitions. I have one context defined for each namespace within every cluster I work with.

### Help

    kc help
    kc help x

### Shell Completion

Put the following into your environment (typically in a file such as ~/.bashrc, ~/.zshrc or ~/.profile)

For bash:

    eval "$(kc --completion-script-bash)"

Or for ZSH:

    eval "$(kc --completion-script-zsh)"

## Subcommands

### Execing into a container

The main point of `kc exec` / `kc sh` is to add the usual flags (e.g. -it, `--`), (bash), and also to allow selector-based querying

#### Examples

    kc l -t10 -f podname
    kc t app=mysvc
    kc l -t10 -f app=mysvc

### Logging

The idea of `kc log` / `kc tail` is to use selector rather than pod name

#### Examples

    kc l -t10 -f podname
    kc l -t10 -f app=mysvc

    kc t app=mysvc
    kc t -cdev mysvc

### Bouncing a deployment

Bounce a deployment's pod by scaling down to zero and back to 1

#### Examples

    kc b mydeployment

### Applying/replacing

Apply or replace the resources specified in a file.

Replace uses --cascade --force. This ensures that resources get restarted

#### Examples

    kc a file.yaml
    kc r file.yaml


### Get pod versions

List versions of pods in the current namespace.

This can be useful for verifying state differences between k8s clusters

#### Examples

    kc v -cdev


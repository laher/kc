# kc

wrappers for a few kubectl commands

## Background

kc* are wrappers around a kubernetes utility called 'kubectl'

It simplifies a few operations which I tend to use a lot. Additionally all commands can be run on multiple 'contexts' at once.

 * kct/kcl: logging, potentially multiplexed
 * kcx/kcsh: exec, shell
 * kcv: versions of pods
 * kca,kcr: apply, replace resources based on config files
 * kcb: 'bounce' (scale down & up)
 * kc: basic kubectl wrapper providing multi-context support

It's not meant to be a full replacement or in the least bit comprehensive. If you want it do be different, I recommend you fork it and bend it to your own will. PRs would also be great, but I'd like to see other people twist it all around for different use cases.

# Installation

 * Prerequisite: install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) and configure one or more contexts

`go get github.com/laher/kc/...`

... or just go-get the ones you want individually.

## Usage

[kc*] [namespace] [options] [args]

e.g. Given a namespace 'dev', and a pod called toolbox, tail its logs:

    kct dev toolbox

e.g. Given namespace 'dev' and 'test', and several pods labelled name=toolbox, run `ps aux` on each:

    kcsh dev,test -l name=toolbox -- ps aux

### Contexts

kc commands all take an optional context name as a first arg. 
See ~/.kube/config for your context definitions. 
For convenience I have one context defined for each namespace within every cluster I work with.


# kc

wrappers for a few kubectl commands. Totally pre-alpha, YMMV.

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

kc[*] [contexts] [options] [args] [--] [options-directly-for-kubectl]

e.g. Given a context 'dev', and a pod called toolbox, find all the services

    kc dev get service

e.g. Given contexts 'dev' and 'test', and several pods labelled name=toolbox, run `ps aux` on each:

    kcx dev,test -l name=toolbox -- ps aux

e.g. Given contexts 'dev' and 'test', and several pods labelled name=toolbox, tail them all at once…

    kct dev,test -l name=toolbox

### Contexts

kc commands all take an optional comma-delimited list of context names, as a first arg. So, you can easily run kubectl commands across data centres.

See ~/.kube/config for your context definitions. 

Tip: I have one context defined for each namespace within every cluster I work with. I don't bother with namespaces directly because it seems like less cognitive load to just have more contexts.

### kubectl flags

Note that it should be possible to pass the right args to kubectl itself, particularly with kc. Use `--` as necessary to indicate to kc[*] that the subsequent flags are for kubectl itself

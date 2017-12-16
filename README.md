# kc

wrappers for a few kubectl commands

## Background

kc* are wrappers around a kubernetes utility called 'kubectl'

It simplifies a few operations which I tend to use a lot

 * kct: tail
 * kcl: log
 * kcx: exec
 * kcsh: shell
 * kcv: get versions of pods
 * kca: apply
 * kcr: replace
 * kcb: 'bounce' (scale down & up)

It's not meant to be a full replacement or in the least bit comprehensive. If you want it do be different, I recommend you fork it and bend it to your own will. PRs would also be great, but I'd like to see other people twist it all around for different use cases.

### Pod selectors

When you need to specify a pod, it can optionally look up by a kubernetes 'selector'
kc will look up a pod by 'selector' whenever there's an '=' in the pod's 'name'

NOTE: for now, it just takes the first result. In the future I'd like to offer either an interactive 'dropdown', or (for tailing), multiplex logs

### Context

kc commands all take an optional context name as a first arg. 
See ~/.kube/config for your context definitions. 
I have one context defined for each namespace within every cluster I work with.


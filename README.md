# scgeme

Because I'm feeling some retroactive FOMO for the kids who actually paid attention in 61A: behold, a non-standard and partial Scheme interpretation written in Go.

For the current project, I'm considering going as far as parking this interpreter on top of Go's concurrency runtime, which would be an interesting exercise in mashing what probably begs to be a single-threaded interpreter on top of a really decent concurrency model.

But the real goal is just to climb the minor mountain of implementing an actual language. All of this amounts to preparation for the medium-term goal of coming up with a language I might come close to actually wanting to use in production, something like an optionally-typed JS look-alike, with first-class functions but without all the WAT-inspiring weirdness, and old-fashioned userland threads on a runtime that can handle lots of context switches. And a batteries-included standard library that's worth a damn. While we're at it, let's also throw in a cheap blueprint for cold fusion and a workable plan to achieve communist utopia. I guess I'm saying I want Go but without whatever it is that makes me feel like Go is not fun to use, which maybe means all the boilerplate and the monastic disdain for abstraction-bulding.

# microscheme

Because I'm feeling some retroactive FOMO for the kids who actually paid attention in 61A.

This is a minimal Scheme interpretation written in Ruby, and probably will migrate to Go if the lack of types starts to make me cranky. The idea that I am preoccupied with after a few hours of development is that one can actually get away with very few special forms, and this is especially convenient since I am already feeling extremely lazy about writing special forms. I've got `if` and `let` and `lambda`, which get me pretty far, and it seems with enough twisting you can get `let` out of `lambda` anyway. I will probably never get around to writing `define` or `cond`.

For the current project, I'm considering going as far as parking this interpreter on top of Go's concurrency runtime, which would be an interesting exercise in mashing what probably begs to be a single-threaded interpreter on top of a really decent concurrency model.

But the real goal is just to climb the minor mountain of implementing an actual language. All of this amounts to preparation for the medium-term goal of coming up with a language I might come close to actually wanting to use in production, something like an optionally-typed JS look-alike, with first-class functions but without all the WAT-inspiring weirdness, and old-fashioned userland threads on a runtime that can handle lots of context switches. And a batteries-included standard library that's worth a damn. While we're at it, let's also throw in a cheap blueprint for cold fusion and a workable plan to achieve communist utopia. I guess I'm saying I want Go but without whatever it is that makes me feel like Go is not fun to use, which maybe means all the boilerplate and the monastic disdain for abstraction-bulding.

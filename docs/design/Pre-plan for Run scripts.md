# 'Run scripts' Pre-plan document

Author: Sam Phillips

**TL/DR:** Design a system to allow scripts (and other command-line tools) to be run and data to be piped in and out of scripts (and other command-line tools).

## Considerations

- Can we create a grammar and API similar to how we read and write from databases?
- Would calling a file with command line look different to piping data in and out?
- What does Go have in its standard library that we can leverage?
- What can we learn from other languages, especially some of the more obscure languages that might have elegant solutions to this?
- What are the security implications? Being pragmatic, and staying simple, is the work done in fileSystemSecurity enough?
- How would spawning child processes work?
- For pragmatism's sake, can we be synchronous (as opposed to asynchronous)?
- What pitfalls should we look out for?

## Parsley Core Design Philosophy

When adding a feature we need to adhere to Parsley's design aesthetic: simplicity, minimalism, completeness, and composability. Are there languages that manage this language feature more elegantly than the standard approach? How would we create and use this feature without requiring a multitude of built-in functions or methods?

As Parsley is a concatenative language at heart, composability and `type-completeness` are always a goal (though not alway achievable). However, pragmatism and an easy-to-understand grammar and PIA are just as important.

## Linked TODO Items:

- run scripts: Execute external commands/scripts/tools with command line inputs/options, receive exit status + optional output as result
- pipe scripts: Execute external commands/scripts/tools with command line inputs/options, receive exit status + optional output as result

## Relevant Documents

Please reference the following documents when considering a design:

- docs/design/Database Design.md
- docs/design/plan-fileIoApi.prompt.md
- docs/design/plan-fileSystemSecurity.prompt.md
- docs/design/Design Philosophy.md
- docs/reference.md

# 'Run scripts' Pre-plan document

Author: Sam Phillips

**TL/DR:** Design a system to allow 

## Parsley Core Design Philosophy

When adding a feature we need to adhere to Parsley's design aesthetic: simplicity, minimalism, completeness, and composability. Are there languages that manage this language feature more elegantly than the standard approach? How would we create and use this featutr without requiring a multitude of built-in functions or methods?

As Parsley is a concatenative language at heart, composibility and `type-completeness` are always a goal (though not alway acheivable). However, pragmatism and an easy-to-understand grammar and PIA are just as important.

## Linked TODO Items:

- run scripts: Execute external commands/scripts/tools with command line imputs/options, receive exit status + optional output as result
- pipe scripts: Execute external commands/scripts/tools with command line imputs/options, receive exit status + optional output as result

## Relevent Documents

- docs/design/Database Design.md
- docs/design/plan-fileIoApi.prompt.md
- docs/design/plan-fileSystemSecurity.prompt.md
- docs/design/Design Philosophy.md
- docs/reference.md

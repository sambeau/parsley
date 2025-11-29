# 'Chroots' Pre-plan document

Author: Sam Phillips

**TL/DR:** create a security feature to limit Parsley's access to the file system (similar in purpose to a Unix 'chroot'). 

## Linked TODO Items:

- 'chroots' for write: limit writes to one or more directories and their children
- 'chroots' for read: limit reads to one or more directories and their children
- 'chroots' for execute: limit execution of external scripts/tools to one or more directories and their children, e.g. ./bin

## Proposed API

 - Command prompts

## Note about 'chroots' for execute

- Can be left as phase 2 for when we add a feature needing it
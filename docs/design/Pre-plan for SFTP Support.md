# Pre-plan for SFTP Support

**Author**: Sam Phillips  
**Status**: ✅ Completed  
**Result**: See `plan-sftpSupport.md`  
**Date**: 2025-11-29

**TL/DR:** Add a fetch-like SFTP support (Based on Go's official library) to format objects for SFTP - in a Parsley manner.

## Considerations

- Fetch isn't a format, it's a wrapper around a network protocol (HTTP)
- SFTP isn't a format, it's a network protocol for secure file transfer that operates over the Secure Shell (SSH) protocol
- We don't explicity call a fetch object 'FETCH' we use format objects + URL
- Can we 'simply add another protocol/URL type to Fetch support?
- For static file upload, write is of a priotity than read
- reading diectories is a conundrum: do we create a format object DIR(), do we infer a directory by path and return a directory object or do we use an API like we have for filesystems?
- What pitfalls should we look out for?
- What else should we consider supporting?
- How would we do more complicated operations (methods on the Format Object? … what would that mean for JSON()? or CSV()??)

## Linked TODO Items:

- SFTP Support: Read/Write files from FTP server - useful for static site generation

## Relevant Documents

Please reference the following documents when considering a design:

- docs/design/plan-fileIoApi.prompt.md
- docs/design/Database Design.md
- docs/design/Design Philosophy.md
- docs/reference.md

Online docs for Go's client

- https://pkg.go.dev/github.com/pkg/sftp


# Parsley VS Code Extension

To install this extension locally:

## Quick Install (macOS/Linux)

```bash
# From the parsley repository root
mkdir -p ~/.vscode/extensions/parsley-language-0.1.0
cp -r .vscode-extension/* ~/.vscode/extensions/parsley-language-0.1.0/

# Reload VS Code
# Press Cmd+Shift+P (macOS) or Ctrl+Shift+P (Windows/Linux)
# Type "Developer: Reload Window" and press Enter
```

## Quick Install (Windows)

```powershell
# From the parsley repository root
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.vscode\extensions\parsley-language-0.1.0"
Copy-Item -Recurse .vscode-extension\* "$env:USERPROFILE\.vscode\extensions\parsley-language-0.1.0\"

# Reload VS Code
# Press Ctrl+Shift+P
# Type "Developer: Reload Window" and press Enter
```

## Package as VSIX (Optional)

To create a distributable `.vsix` file:

```bash
# Install vsce
npm install -g @vscode/vsce

# Package
cd .vscode-extension
vsce package

# Install the generated .vsix
code --install-extension parsley-language-0.1.0.vsix
```

## Testing

1. Open any `.pars` file in VS Code
2. Syntax highlighting should be applied automatically
3. Test features:
   - Comment toggling: `Cmd/Ctrl + /`
   - Auto-closing brackets and quotes
   - Code folding
   - Bracket matching

## Updating

If you make changes to the grammar:

1. Update `syntaxes/parsley.tmLanguage.json`
2. Reinstall using the quick install commands above
3. Reload VS Code

## Features Included

- ✅ Syntax highlighting for all Parsley language constructs
- ✅ Comment support
- ✅ Auto-closing pairs
- ✅ Bracket matching
- ✅ Code folding
- ✅ `.pars` file association

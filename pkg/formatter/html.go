package formatter

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

// FormatHTML pretty-prints HTML with proper indentation.
// Returns the original string if it doesn't appear to be HTML or if parsing fails.
func FormatHTML(input string) string {
	// Trim whitespace for detection
	trimmed := strings.TrimSpace(input)

	// Auto-detect: must start with < and contain >
	if !strings.HasPrefix(trimmed, "<") || !strings.Contains(trimmed, ">") {
		return input
	}

	// Try to parse as HTML
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		// Not valid HTML, return original
		return input
	}

	// Pretty-print the document
	var buf bytes.Buffer
	renderNode(&buf, doc, 0)

	result := buf.String()
	// Remove leading/trailing newlines
	return strings.TrimSpace(result)
}

// renderNode recursively renders an HTML node with indentation
func renderNode(buf *bytes.Buffer, n *html.Node, level int) {
	indent := strings.Repeat("  ", level)

	switch n.Type {
	case html.DocumentNode:
		// Document node - just render children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			renderNode(buf, c, level)
		}

	case html.ElementNode:
		// Determine if this is an inline element that should not add newlines
		isInline := isInlineElement(n.Data)
		isVoid := isVoidElement(n.Data)

		// Only add indent and newline for block elements at level > 0
		if !isInline && level >= 0 {
			buf.WriteString(indent)
		}

		// Opening tag
		buf.WriteString("<")
		buf.WriteString(n.Data)

		// Attributes
		for _, attr := range n.Attr {
			buf.WriteString(" ")
			buf.WriteString(attr.Key)
			buf.WriteString("=\"")
			buf.WriteString(html.EscapeString(attr.Val))
			buf.WriteString("\"")
		}

		// Self-closing tag
		if isVoid {
			buf.WriteString(" />")
			if !isInline && level >= 0 {
				buf.WriteString("\n")
			}
			return
		}

		buf.WriteString(">")

		// Check if element has only text content
		hasOnlyText := hasOnlyTextChildren(n)

		// For inline elements or elements with only text, keep content on same line
		if isInline || hasOnlyText {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				renderNode(buf, c, level+1)
			}
		} else {
			// Block element with nested content - add newlines
			if n.FirstChild != nil && !isInline && level >= 0 {
				buf.WriteString("\n")
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				renderNode(buf, c, level+1)
			}
			// Add indent before closing tag (unless it's style/script which handles its own indentation)
			if n.FirstChild != nil && !isInline && level >= 0 && n.Data != "style" && n.Data != "script" {
				buf.WriteString(indent)
			}
		}

		// Closing tag
		buf.WriteString("</")
		buf.WriteString(n.Data)
		buf.WriteString(">")

		// Add newline after block elements
		if !isInline && level >= 0 {
			buf.WriteString("\n")
		}

	case html.TextNode:
		// Text content
		text := n.Data
		// Only output non-empty text
		if strings.TrimSpace(text) != "" {
			// Normalize whitespace for style and script tags
			if n.Parent != nil && (n.Parent.Data == "style" || n.Parent.Data == "script") {
				// Trim and re-indent the content
				lines := strings.Split(strings.TrimSpace(text), "\n")
				if len(lines) > 0 {
					buf.WriteString("\n")
					// Use parent's indent level for the content
					parentIndent := strings.Repeat("  ", level-1)
					for _, line := range lines {
						trimmedLine := strings.TrimSpace(line)
						if trimmedLine != "" {
							buf.WriteString(parentIndent)
							buf.WriteString("  ") // Extra indent for content inside tag
							buf.WriteString(trimmedLine)
							buf.WriteString("\n")
						}
					}
					// Position for closing tag at parent's level
					buf.WriteString(parentIndent)
				}
			} else {
				// Preserve text content as-is for other elements
				buf.WriteString(text)
			}
		}

	case html.CommentNode:
		// Comments
		if level >= 0 {
			buf.WriteString(indent)
		}
		buf.WriteString("<!--")
		buf.WriteString(n.Data)
		buf.WriteString("-->")
		if level >= 0 {
			buf.WriteString("\n")
		}
	}
}

// hasOnlyTextChildren returns true if the node has only text children
func hasOnlyTextChildren(n *html.Node) bool {
	if n.FirstChild == nil {
		return false
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.TextNode {
			return false
		}
	}
	return true
}

// isInlineElement returns true if the element should be rendered inline
func isInlineElement(tag string) bool {
	inlineElements := map[string]bool{
		"a": true, "abbr": true, "acronym": true, "b": true, "bdi": true,
		"bdo": true, "big": true, "br": true, "button": true, "cite": true,
		"code": true, "dfn": true, "em": true, "i": true, "img": true,
		"input": true, "kbd": true, "label": true, "map": true, "mark": true,
		"object": true, "output": true, "q": true, "samp": true,
		"select": true, "small": true, "span": true, "strong": true, "sub": true,
		"sup": true, "textarea": true, "time": true, "tt": true, "var": true,
	}
	return inlineElements[tag]
}

// isVoidElement returns true if the element is self-closing (void)
func isVoidElement(tag string) bool {
	voidElements := map[string]bool{
		"area": true, "base": true, "br": true, "col": true, "embed": true,
		"hr": true, "img": true, "input": true, "link": true, "meta": true,
		"param": true, "source": true, "track": true, "wbr": true,
	}
	return voidElements[tag]
}

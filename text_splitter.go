package text_splitter

// Text splitter interface
type TextSplitter interface {
	SplitText() ([]string, error)
	SplitMultipleTexts() ([]*SplitResult, error)
}

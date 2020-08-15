	"strings"
const exampleHunk = `@@ -1,5 +1,5 @@
 apple
-grape
+orange
...
...
...
`


func TestLineNumberOfLine(t *testing.T) {
	type scenario struct {
		testName string
		hunk     *PatchHunk
		idx      int
		expected int
	}

	scenarios := []scenario{
		{
			testName: "nothing selected",
			hunk:     newHunk(strings.SplitAfter(exampleHunk, "\n"), 10),
			idx:      15,
			expected: 3,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			result := s.hunk.LineNumberOfLine(s.idx)
			if !assert.Equal(t, s.expected, result) {
				fmt.Println(result)
			}
		})
	}
}
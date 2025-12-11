package parser

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	utils "github.com/davewil/aoc-utils"
)

// Domain models used by the solver
type Row struct {
	Lights   []int
	Buttons  [][]int
	Joltages []int
}

type Input struct {
	Rows []Row
}

// AST models for Participle
type parsedRow struct {
	LightsRaw string         `parser:"@Lights"`
	Buttons   []parsedButton `parser:"@@+"`
	Joltages  parsedJoltages `parser:"@@"`
}

type parsedButton struct {
	Values []int `parser:"'(' @Int (',' @Int)* ')'"`
}

type parsedJoltages struct {
	Values []int `parser:"'{' @Int (',' @Int)* '}'"`
}

// Custom lexer to handle the specific format
var day10Lexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Lights", `\[[.#]+\]`},
	{"Int", `\d+`},
	{"Punct", `[(),{}]`},
	{"Whitespace", `\s+`},
})

var parser = participle.MustBuild[parsedRow](
	participle.Lexer(day10Lexer),
	participle.Elide("Whitespace"),
)

func ParseInput(raw string) (Input, error) {
	lines := utils.ReadLines(raw)
	input := Input{}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parsed, err := parser.ParseString("", line)
		if err != nil {
			return Input{}, err
		}

		row := Row{}

		// Convert Lights: [.##.] -> indices of #
		// Strip [ and ]
		lightsContent := strings.Trim(parsed.LightsRaw, "[]")
		for i, c := range lightsContent {
			if c == '#' {
				row.Lights = append(row.Lights, i)
			}
		}

		// Convert Buttons
		for _, btn := range parsed.Buttons {
			row.Buttons = append(row.Buttons, btn.Values)
		}

		// Convert Joltages
		row.Joltages = parsed.Joltages.Values

		input.Rows = append(input.Rows, row)
	}
	return input, nil
}

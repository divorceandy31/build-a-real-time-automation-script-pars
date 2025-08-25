package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/antlr/antlr4/runtime/Go/antlr/parser"
)

// Define the script parser structure
type ScriptParser struct {
	*parser.BaseParser
 SCRIPT
    rules      map[string][]string
    variables map[string]string
    timers    map[string]time.Time
}

// NewScriptParser returns a new instance of ScriptParser
func NewScriptParser(input antlr.CharStream) *ScriptParser {
	is := antlr.NewInputStream(input)
	lexer := NewScriptLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenStreamDefaultTokenChannel)
	p := NewScriptParser(stream)
	p.rules = make(map[string][]string)
	p.variables = make(map[string]string)
	p.timers = make(map[string]time.Time)
	return p
}

func (p *ScriptParser) Parse() error {
	p.SetBuildParseTrees(true)
	tree, err := p.ParseScripts()
	if err != nil {
		return err
	}
	p.Visit(tree)
	return nil
}

func (p *ScriptParser) Visit(node antlr.Tree) interface{} {
	switch t := node.(type) {
	case *ScriptParserRuleContext:
		for _, rule := range t.AllRuleContext() {
			p.visitRule(rule)
		}
	default:
		log.Printf("Unknown node type: %T", node)
	}
	return nil
}

func (p *ScriptParser) visitRule(rule antlr.RuleContext) {
	switch t := rule.(type) {
	case *RuleRuleContext:
		p.parseRule(t)
	default:
		log.Printf("Unknown rule type: %T", rule)
	}
}

func (p *ScriptParser) parseRule(rule *RuleContext) {
	switch rule.GetRuleIndex() {
	case RULE_SCRIPT:
		p.parseScript(rule)
	case RULE_TIMER:
		p.parseTimer(rule)
	case RULE_VARIABLE:
		p.parseVariable(rule)
	default:
		log.Printf("Unknown rule index: %d", rule.GetRuleIndex())
	}
}

func (p *ScriptParser) parseScript(rule *RuleContext) {
	script := rule.GetText()
	p.rules["default"] = strings.Split(script, "\n")
}

func (p *ScriptParser) parseTimer(rule *RuleContext) {
	timerName := rule.GetChild(1).GetText()
	timerValue := rule.GetChild(3).GetText()
	p.timers[timerName] = time.Parse(time.RFC3339, timerValue)
}

func (p *ScriptParser) parseVariable(rule *RuleContext) {
	varName := rule.GetChild(1).GetText()
	varValue := rule.GetChild(3).GetText()
	p.variables[varName] = varValue
}

func main() {
	file, err := os.Open("script.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

SCRIPT, err := antlr.NewFileStream("script.txt")
	if err != nil {
		log.Fatal(err)
	}

	parser := NewScriptParser(SCRIPT)
	err = parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Parsed rules:")
	for _, rules := range parser.rules {
		fmt.Println(rules)
	}

	fmt.Println("Parsed variables:")
	for _, variable := range parser.variables {
		fmt.Println(variable)
	}

	fmt.Println("Parsed timers:")
	for _, timer := range parser.timers {
		fmt.Println(timer)
	}
}
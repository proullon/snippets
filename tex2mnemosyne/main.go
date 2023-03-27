package main

import (
	"fmt"
	"io"
	"os"
	//"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

type Card struct {
	Tag   string
	Front string
	Back  string
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "file",
			Value: "main.tex",
			Usage: "File to read",
		},
		&cli.StringFlag{
			Name:  "import-format",
			Value: "tex",
			Usage: "Import format [delim, tex]",
		},
		&cli.StringFlag{
			Name:  "export-format",
			Value: "csv",
			Usage: "Export format [csv, txt]",
		},
		&cli.StringFlag{
			Name:  "tag",
			Value: "",
			Usage: "global tag",
		},
	}
	app.Action = start
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}

func start(c *cli.Context) error {

	var cards []Card
	var err error

	filename := c.String("file")

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	importFormat := c.String("import-format")
	exportFormat := c.String("export-format")
	globalTag := c.String("tag")
	switch importFormat {
	case "delim":
		cards, err = fromDelim(string(data))
	case "tex":
		cards, err = fromTex(string(data))
	default:
		return fmt.Errorf("unknown format '%s'", importFormat)
	}

	if err != nil {
		return err
	}

	switch exportFormat {
	case "csv":
		toCSV(cards, globalTag)
	case "txt":
		toTXT(cards)
	default:
		return fmt.Errorf("unknown format: '%s'", exportFormat)
	}

	return nil
}

func fromDelim(data string) ([]Card, error) {
	return nil, nil
}

func fromTex(data string) ([]Card, error) {
	var cursor int
	var cards []Card
	var section, subsection, subsubsection, count int

	section = 0
	subsection = 0
	subsubsection = 0
	count = 0

	for {
		cmd, opt, arg, begin, end, err := getAnyCommand(data, cursor)
		if err != nil && err == io.EOF {
			return cards, nil
		}
		if err != nil {
			return nil, fmt.Errorf("%d:%d: %w", begin, end, err)
		}
		//fmt.Printf("Command: '%s' opt '%s' arg '%s'\n", cmd, opt, arg)
		cursor = end
		switch cmd {
		case "section":
			//fmt.Printf("SECTION %s\n", opt)
			//		section, _ = strconv.Atoi(opt)
			section++
			subsection = 1
			subsubsection = 1
		case "subsection":
			//fmt.Printf("SUBSECTION %s\n", opt)
			//			subsection, _ = strconv.Atoi(opt)
			subsection++
			subsubsection = 1
		case "subsubsection":
			//fmt.Printf("SUBSUBSECTION %s\n", opt)
			//			subsubsection, _ = strconv.Atoi(opt)
			subsubsection++
		case "begin":
			switch opt {
			case "definition":
				c, end, err := getDefinition(data, cursor, arg, section, subsection, subsubsection)
				if err != nil {
					return nil, fmt.Errorf("Card %d: %d:%d: cannot get definition: %w", count, begin, end, err)
				}
				cards = append(cards, c)
				cursor = end
			case "exercise":
				//fmt.Printf("Exercice %s\n", opt)
				c, end, err := getExercise(data, cursor, section, subsection, subsubsection)
				if err != nil {
					return nil, fmt.Errorf("Card %d: %d:%d: cannot get exercise: %w", count, begin, end, err)
				}
				cards = append(cards, c)
				cursor = end
			case "theorem", "theo":
				c, end, err := getTheorem(data, cursor, arg, opt, section, subsection, subsubsection)
				if err != nil {
					return nil, fmt.Errorf("Card %d: %d:%d: cannot get theorem: %w", count, begin, end, err)
				}
				cards = append(cards, c)
				cursor = end
			case "lemma", "lem":
				//fmt.Printf("Lemma %s\n", opt)
				c, end, err := getLemma(data, cursor, arg, opt, section, subsection, subsubsection)
				if err != nil {
					return nil, fmt.Errorf("Card %d: %d:%d: cannot get lemma: %w", count, begin, end, err)
				}
				cards = append(cards, c)
				cursor = end
			case "proof":
				c, end, err := getProof(data, cursor, arg, section, subsection, subsubsection)
				if err != nil {
					return nil, fmt.Errorf("Card %d: %d:%d: cannot get proof: %w", count, begin, end, err)
				}
				cards = append(cards, c)
				cursor = end
			}
			count++
		}
	}

	return nil, nil
}

func getLemma(data string, cursor int, arg string, keyw string, section, subsection, subsubsection int) (Card, int, error) {
	var c Card

	begin, end, err := getCommand("end", keyw, data, cursor)
	if err != nil {
		return c, 0, err
	}

	c = Card{
		Tag:   fmt.Sprintf("%d::%d::%d", section, subsection, subsubsection),
		Front: fmt.Sprintf("<latex>Lemme: %s</latex>", arg),
		Back:  fmt.Sprintf("<latex>%s</latex>", data[cursor:begin]),
	}

	return c, end, nil
}

func getProof(data string, cursor int, arg string, section, subsection, subsubsection int) (Card, int, error) {
	var c Card

	begin, end, err := getCommand("end", "proof", data, cursor)
	if err != nil {
		return c, 0, err
	}

	c = Card{
		Tag:   fmt.Sprintf("%d::%d::%d", section, subsection, subsubsection),
		Front: fmt.Sprintf("<latex>%s</latex>", arg),
		Back:  fmt.Sprintf("<latex>%s</latex>", data[cursor:begin]),
	}

	return c, end, nil
}

func getDefinition(data string, cursor int, arg string, section, subsection, subsubsection int) (Card, int, error) {
	var c Card

	beginEndDefinition, end, err := getCommand("end", "definition", data, cursor)
	if err != nil {
		return c, 0, err
	}

	c = Card{
		Tag:   fmt.Sprintf("%d::%d::%d", section, subsection, subsubsection),
		Front: fmt.Sprintf("<latex>Définition: %s</latex>", arg),
		Back:  fmt.Sprintf("<latex>%s</latex>", data[cursor:beginEndDefinition]),
	}

	return c, end, nil
}

func getTheorem(data string, cursor int, arg string, keyw string, section, subsection, subsubsection int) (Card, int, error) {
	var c Card

	beginEndTheorem, end, err := getCommand("end", keyw, data, cursor)
	if err != nil {
		return c, 0, err
	}

	//fmt.Printf("We are at:\n%s\n", data[cursor:])

	//fmt.Printf("Theorem end at %d/%d\n", end, len(data))

	//fmt.Printf("Theo is %d to %d over %d\n", cursor, beginEndTheorem, len(data))

	c = Card{
		Tag:   fmt.Sprintf("%d::%d::%d", section, subsection, subsubsection),
		Front: fmt.Sprintf("<latex>Théorème: %s</latex>", arg),
		Back:  fmt.Sprintf("<latex>%s</latex>", data[cursor:beginEndTheorem]),
	}

	return c, end, nil
}

func getExercise(data string, cursor int, section, subsection, subsubsection int) (Card, int, error) {
	var c Card

	_, end, err := getCommand("end", "exercise", data, cursor)
	if err != nil {
		return c, 0, err
	}

	beginBeginQuestion, endBeginQuestion, err := getCommand("begin", "question", data, cursor)
	if err != nil {
		return c, 0, err
	}
	_ = beginBeginQuestion
	beginEndQuestion, endEndQuestion, err := getCommand("end", "question", data, cursor)
	if err != nil {
		return c, 0, err
	}
	_ = endEndQuestion
	_ = beginEndQuestion

	beginBeginSolution, endBeginSolution, err := getCommand("begin", "solution", data, cursor)
	if err != nil {
		return c, 0, err
	}
	_ = beginBeginSolution
	_ = endBeginSolution
	beginEndSolution, endEndSolution, err := getCommand("end", "solution", data, cursor)
	if err != nil {
		return c, 0, err
	}
	_ = endEndSolution

	if endBeginQuestion > beginEndQuestion {
		return c, 0, fmt.Errorf("Malformed exercice %s question", fmt.Sprintf("%d::%d::%d", section, subsection, subsubsection))
	}
	if endBeginSolution > beginEndSolution {
		return c, 0, fmt.Errorf("Malformed exercice %s solution", fmt.Sprintf("%d::%d::%d", section, subsection, subsubsection))
	}

	c = Card{
		Tag:   fmt.Sprintf("%d::%d::%d", section, subsection, subsubsection),
		Front: fmt.Sprintf("<latex>Exercice: %s</latex>", data[endBeginQuestion:beginEndQuestion]),
		Back:  fmt.Sprintf("<latex>%s</latex>", data[endBeginSolution:beginEndSolution]),
	}

	return c, end, nil
}

func getAnyCommand(data string, i int) (command string, option string, arg string, begin, end int, err error) {

	for i < len(data) {
		if data[i] == '$' {
			i++
			for i < len(data) {
				if data[i] == '$' {
					break
				}
				i++
			}
		}

		if i >= len(data) {
			err = io.EOF
			return
		}
		if data[i] == '\\' {
			if data[i+1] == '\\' {
				i += 2
				continue
			}
			for x := i; x < len(data); x++ {
				if data[x] == '\n' {
					begin = i
					command = data[i+1 : x]
					end = x + 1
					return
				}
				if data[x] == '{' {
					command = data[i+1 : x]
					begin = i
					i = x + 1
					break
				}
			}
			for x := i; x < len(data); x++ {
				if data[x] == '}' {
					option = data[i:x]
					end = x + 1
					if a, e, err := getArgument(data, x+1); err == nil {
						arg = a
						end = e
					}
					return
				}
			}
		}
		i++
	}

	err = io.EOF
	return
}

func getCommand(wantedCmd string, wantedOpt string, data string, i int) (begin, end int, err error) {
	cursor := i
	for {
		cmd, opt, _, b, end, err := getAnyCommand(data, cursor)
		if err != nil {
			return 0, 0, err
		}
		begin = b
		if cmd == wantedCmd && opt == wantedOpt {
			return begin, end, nil
		}
		cursor = end
	}

	return
}

func getArgument(data string, i int) (arg string, end int, err error) {
	if data[i] == '[' {
		begin := i + 1
		for i < len(data) {
			if data[i] == ']' {
				end = i
				arg = data[begin:end]
				end += 1
				return
			}
			i++
		}
	}

	return "", 0, io.EOF
}

func toCSV(cards []Card, globalTag string) error {
	//	fmt.Printf("%s\t%s\t%s\n", "Front", "Back", "Tag")
	for _, c := range cards {
		c.Front = strings.Replace(c.Front, "\n", " ", -1)
		c.Front = strings.Replace(c.Front, "\t", " ", -1)
		c.Back = strings.Replace(c.Back, "\n", " ", -1)
		c.Back = strings.Replace(c.Back, "\t", " ", -1)
		if globalTag != "" {
			c.Tag = globalTag + "::" + c.Tag
		}
		if c.Back == "<latex> </latex>" {
			c.Back = "<latex> FIXME </latex>"
		}
		fmt.Printf("%s\t%s\t%s\n", c.Front, c.Back, c.Tag)
	}
	return nil
}

func toTXT(cards []Card) error {
	for _, c := range cards {
		fmt.Printf("------------------\n%s\n", c.Tag)
		fmt.Printf("\n%s\n", c.Front)
		fmt.Printf("\n%s\n", c.Back)
	}
	return nil
}

func tex2csv(data string) error {

	cards := strings.Split(string(data), "% --- CARD ---")
	for _, c := range cards {
		if c == "" {
			continue
		}
		elements := strings.Split(string(c), "% ------------")
		if len(elements) != 3 {
			return fmt.Errorf("expected 3 elements (tags, front, back), got %d in %s", len(elements), c)
		}
		fmt.Printf("%s\t%s\n", elements[1], elements[2])
	}
	return nil
}

func tex2txt(data string) error {

	cards := strings.Split(string(data), "% --- CARD ---")
	for _, c := range cards {
		if c == "" {
			continue
		}
		elements := strings.Split(string(c), "% ------------")
		if len(elements) != 3 {
			return fmt.Errorf("expected 3 elements (tags, front, back), got %d in %s", len(elements), c)
		}
		fmt.Printf("------------------\n\n%s\n\n", elements[0])
		fmt.Printf("\n\n<latex>\n%s\n</latex>\n\n", elements[1])
		fmt.Printf("\n\n<latex>\n%s\n</latex>\n\n", elements[2])
	}
	return nil
}

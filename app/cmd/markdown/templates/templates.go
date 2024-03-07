package templates

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Template interface {
	Render() string
	GetName() string
	GetUnrenderedContent() string
}

func NewStaticTemplate(name string, content string) Template {
	return staticTemplate{name: name, content: content}
}

type staticTemplate struct {
	name    string
	content string
}

func (st staticTemplate) Render() string {
	return st.content
}

func (st staticTemplate) GetName() string {
	return st.name
}

func (st staticTemplate) GetUnrenderedContent() string {
	return st.content
}

func NewIdTemplate(name string, content string) Template {
	return idTemplate{staticTemplate{name: name, content: content}}
}

type idTemplate struct {
	staticTemplate
}

func (it idTemplate) Render() string {
	id := uuid.New().String()
	content := it.staticTemplate.Render()
	content = fmt.Sprintf(strings.ReplaceAll(content, `~~~`, "```")+"\n", id)
	return content
}

func (it idTemplate) GetName() string {
	return it.name
}

func (it idTemplate) GetUnrenderedContent() string {
	return it.content
}

func NewAttributeTemplate(name string, content string, minimal bool, withExplanation bool, withRubric bool, withHints int) Template {
	return attributeTemplate{
		minimal:         minimal,
		withExplanation: withExplanation,
		withRubric:      withRubric,
		withHints:       withHints,
		idTemplate: idTemplate{
			staticTemplate: staticTemplate{
				name:    name,
				content: content,
			},
		},
	}
}

type attributeTemplate struct {
	minimal         bool
	withExplanation bool
	withRubric      bool
	withHints       int
	idTemplate
}

func (at attributeTemplate) Render() string {
	comments := []string{"", "", "", ""}
	attrs := []string{"", "", ""}

	if !at.minimal && (!at.withExplanation || !at.withRubric || at.withHints == 0) {
		comments = append(comments, blockHeader)
	}

	if !at.minimal && at.withHints == 0 {
		comments = append(comments, hintTemplateSilent)
	}

	if !at.minimal && !at.withRubric {
		comments = append(comments, rubricTemplateSilent)
	}

	if !at.minimal && !at.withExplanation {
		comments = append(comments, explanationTemplateSilent)
	}

	if at.withHints > 0 {
		t := hintTemplate
		if at.minimal {
			t = hintTemplateMin
		}
		repeats := make([]string, at.withHints)
		for i := 0; i < at.withHints; i += 1 {
			repeats[i] = t
		}
		attrs = append(attrs, strings.Join(repeats, "\n\n"))
	}

	if at.withRubric {
		t := rubricTemplate
		if at.minimal {
			t = rubricTemplateMin
		}
		attrs = append(attrs, t)
	}

	if at.withExplanation {
		t := explanationTemplate
		if at.minimal {
			t = explanationTemplateMin
		}
		attrs = append(attrs, t)
	}

	joinedComments := strings.TrimSpace(strings.Join(comments, "\n"))
	joinedAttrs := strings.TrimSpace(strings.Join(attrs, "\n\n"))

	blocks := fmt.Sprintf("\n%s\n", joinedAttrs)

	if !at.minimal {
		blocks = fmt.Sprintf("%s\n\n%s", joinedComments, joinedAttrs)
		blocks = strings.TrimSpace(blocks)
	}

	content := at.idTemplate.Render()
	content = strings.ReplaceAll(content, "<optional-attributes>", blocks)
	return content
}

func (at attributeTemplate) GetName() string {
	return at.name
}

func (at attributeTemplate) GetUnrenderedContent() string {
	return at.content
}

/******************************************************************************
 * Challenge template optional blocks
 *
 * These come in three pieces:
 *   - The "full" one with the place holder/description
 *   - The "minimum" one with just the markup
 *   - The "silent" one that shows up as a comment when the option isn't
 *     specified
 */

const blockHeader = `<!-- other optional sections -->`

const rubricTemplate = `##### !rubric

[Put your rubric here specifying how to allocate points for the challenge]

##### !end-rubric`

const rubricTemplateMin = `##### !rubric

##### !end-rubric`

const rubricTemplateSilent = `<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->`

const hintTemplate = `##### !hint

[Put a single hint, here. Add more hint blocks as needed.]

##### !end-hint`

const hintTemplateMin = `##### !hint

##### !end-hint`

const hintTemplateSilent = `<!-- !hint - !end-hint (markdown, hidden, students click to view) -->`

const explanationTemplate = `##### !explanation-correct:

[Put the explanation for the CORRECT response, here.]

##### !end-explanation

##### !explanation-incorrect:

[Put the explanation for the INCORRECT response, here.]

##### !end-explanation`

const explanationTemplateMin = `##### !explanation-correct:

##### !end-explanation

##### !explanation-incorrect:

##### !end-explanation`

const explanationTemplateSilent = `<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->`

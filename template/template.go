package template

/* Template is a struct that stores a text template. Templates have several
 * fields which any string can be inserted into. It is incredibly useful for
 * making multiple pages have identical layout and styling. Fields do not have
 * to be unique, and when Eval is called all identical fields will be populated
 * with the same input value.
 */
type Template struct {
        items []*templateItem
}

type templateItem struct {
        replaceMe bool
        value     string
}

/* New creates a new template from a string of valid template syntax. Any curly
 * bracket wrapped text in a template will be treated as a field, and will be
 * able to be replaced. An example template string:
 *
 * <head><title>{title}</title></head><body><p>{content}</p></body>
 *
 * The above template string has two fields: title and content. Eval can be
 * called to return a string where these two fields are replaced with custom
 * input.
 */
func New (templateString string) (template *Template) {
        item := &templateItem{}
        template = &Template {
                []*templateItem { item },
        }

        skip := false
        for _, ch := range(templateString) {
                if skip { item.value += string(ch); skip = false }

                if item.replaceMe {
                        if ch == '}' {
                                item = &templateItem { replaceMe: false }
                                template.items = append(template.items, item)
                        } else {
                                item.value += string(ch)
                        }
                        continue
                }
        
                switch ch {
                case '\\':
                        skip = true
                        break
                case '{':
                        item = &templateItem { replaceMe: true }
                        template.items = append(template.items, item)
                        break
                default:
                        item.value += string(ch)
                }
        }

        return
}

/* Eval uses the template to format the inputs map and outputs the result as a
 * string.
 */
func (template *Template) Eval (inputs map[string] string) (output string) {
        for _, item := range(template.items) {
                if item.replaceMe {
                        value, defined := inputs[item.value]
                        if defined { output += value }
                } else {
                        output += item.value
                }
        }

        return
}

/* ByteEval wraps Eval, but returns a []byte instead of a string which can be
 * more useful for writing a response body.
 */
func (template *Template) ByteEval (inputs map[string] string) (output []byte) {
        return []byte(template.Eval(inputs))
}

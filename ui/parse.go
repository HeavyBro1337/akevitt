package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/rivo/tview"
	"golang.org/x/net/html"
)

type DOM struct {
	p               tview.Primitive
	stringCallbacks map[string]func(string)
	intCallbacks    map[string]func(int)
	boolCallbacks   map[string]func(bool)
	callbacks       map[string]func()
}

func (d *DOM) Primitive() tview.Primitive {
	return d.p
}

func ParseHTML(r io.Reader, app *tview.Application) (*DOM, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	if node.Type == html.DocumentNode {
		node = node.FirstChild
	}

	rootFlex := tview.NewFlex().SetFullScreen(true).SetDirection(tview.FlexRow)
	d := &DOM{
		p:               rootFlex,
		intCallbacks:    map[string]func(int){},
		boolCallbacks:   make(map[string]func(bool)),
		stringCallbacks: make(map[string]func(string)),
		callbacks:       make(map[string]func()),
	}
	for child := node; child != nil; child = child.NextSibling {
		if primitive := d.parseNode(child, app); primitive != nil {
			rootFlex.AddItem(primitive, 0, 1, true)
		}
	}

	return d, nil
}

func (d *DOM) parseNode(n *html.Node, app *tview.Application) tview.Primitive {
	p := func() tview.Primitive {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "form":
				return d.parseForm(n, app)
			case "div":
				return d.parseDiv(n, app)
			case "span":
				return d.parseSpan(n)
			case "p":
				return d.parseParagraph(n)
			case "ul":
				return d.parseUnorderedList(n, app)
			case "ol":
				return d.parseOrderedList(n, app)
			case "li":
				return d.parseListItem(n)
			case "table":
				return d.parseTable(n, app)
			case "td", "th":
				return d.parseTableCell(n)
			case "input":
				return d.parseInput(n)
			case "button":
				return d.parseButton(n)
			case "textarea":
				return d.parseTextArea(n)
			case "grid":
				return d.parseGrid(n, app)
			case "empty":
				return nil
			default:
				return d.parseChildren(n, app)
			}
		} else if n.Type == html.TextNode {
			return tview.NewTextView().SetText(n.Data)
		}
		return nil
	}()

	return p
}

func (d *DOM) parseDiv(n *html.Node, app *tview.Application) tview.Primitive {
	flex := makeFlex(n)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		child := d.parseNode(c, app)
		flex.AddItem(child, IntValue(fixedSizeKey, n, 0), IntValue(proportionKey, c, 1), HasKey(focusedKey, c))
	}

	return flex
}

func (d *DOM) parseGrid(n *html.Node, app *tview.Application) tview.Primitive {
	grid := tview.NewGrid().
		SetRows(IntArrValue(rowsKey, n, 0)...).
		SetColumns(IntArrValue(colsKey, n, 0)...)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		child := d.parseNode(c, app)

		if child != nil {
			grid.AddItem(child,
				IntValue(rowKey, n, 0),
				IntValue(colKey, n, 0),
				IntValue(rowSpanKey, n, 1),
				IntValue(colSpanKey, n, 1),
				IntValue(minHeightKey, n, 0), IntValue(minWidthKey, n, 0),
				HasKey(focusedKey, n),
			)
		}
	}
	return grid
}

func (d *DOM) parseForm(n *html.Node, app *tview.Application) tview.Primitive {
	form := tview.NewForm()
	form.SetTitle(Value(titleKey, n))

	if HasKey(focusedKey, n) {
		app.SetFocus(form)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data != "input" {
			continue
		}

		switch strings.ToLower(ValueFallback("type", c, "text")) {
		case "password":
			{
				form.AddPasswordField(Value(labelKey, c), "", IntValue(widthKey, c, 0), []rune(ValueFallback(maskKey, c, "*"))[0], func(text string) {
					if f, ok := d.stringCallbacks[Value(keyKey, c)]; ok {
						f(text)
					}
				})
			}
		case "text":
			{
				form.AddInputField(Value(labelKey, c), "", IntValue(widthKey, c, 0), nil, func(text string) {
					if f, ok := d.stringCallbacks[Value(keyKey, c)]; ok {
						f(text)
					}
				})
			}
		case "checkbox":
			{
				form.AddCheckbox(Value(labelKey, c), HasKey(checkedKey, c), func(checked bool) {
					if f, ok := d.boolCallbacks[Value(keyKey, c)]; ok {
						f(checked)
					}
				})
			}
		case "button":
			{
				form.AddButton(Value(labelKey, c), func() {
					fmt.Printf("d.callbacks: %v\n", d.callbacks)
					if f, ok := d.callbacks[Value(keyKey, c)]; ok {
						f()
					}
				})
			}
		}
	}

	return form
}

func (d *DOM) parseSpan(n *html.Node) tview.Primitive {
	textView := tview.NewTextView()
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			textView.SetText(textView.GetText(true) + c.Data)
		}
	}

	return textView
}

func (d *DOM) parseParagraph(n *html.Node) tview.Primitive {
	textView := tview.NewTextView()
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			textView.SetText(textView.GetText(true) + c.Data)
		}
	}

	return textView
}

func (d *DOM) parseUnorderedList(n *html.Node, app *tview.Application) tview.Primitive {
	list := tview.NewList()
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "li" {
			item := d.parseNode(c, app)
			if textView, ok := item.(*tview.TextView); ok {
				list.AddItem(textView.GetText(true), "", 0, nil)
			}
		}
	}
	return list
}

func (d *DOM) parseOrderedList(n *html.Node, app *tview.Application) tview.Primitive {
	list := tview.NewList()
	index := 1
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "li" {
			item := d.parseNode(c, app)
			if textView, ok := item.(*tview.TextView); ok {
				list.AddItem(textView.GetText(true), "", 0, nil)
				index++
			}
		}
	}
	return list
}

func (d *DOM) parseListItem(n *html.Node) tview.Primitive {
	textView := tview.NewTextView()
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			textView.SetText(textView.GetText(true) + c.Data)
		}
	}
	return textView
}

func (d *DOM) parseTable(n *html.Node, app *tview.Application) tview.Primitive {
	table := tview.NewTable()
	row := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "tr" {
			d.parseTableRowInto(table, row, c, app)
			row++
		}
	}
	return table
}

func (d *DOM) parseTableRowInto(table *tview.Table, row int, n *html.Node, app *tview.Application) {
	col := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "td" || c.Data == "th" {
			cell := d.parseNode(c, app)
			if textView, ok := cell.(*tview.TextView); ok {
				table.SetCell(row, col, tview.NewTableCell(textView.GetText(true)))
				col++
			}
		}
	}
}

func (d *DOM) parseTableCell(n *html.Node) tview.Primitive {
	textView := tview.NewTextView()
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			textView.SetText(textView.GetText(true) + c.Data)
		}
	}
	return textView
}

func (d *DOM) parseInput(n *html.Node) tview.Primitive {
	return tview.NewInputField().SetLabel(Value(labelKey, n)).SetChangedFunc(func(text string) {
		if f, ok := d.stringCallbacks[Value(keyKey, n)]; ok {
			f(text)
		}
	})
}

func (d *DOM) parseButton(n *html.Node) tview.Primitive {
	button := tview.NewButton("").SetSelectedFunc(func() {
		fmt.Printf("d.callbacks: %#v\n", d.callbacks)
		fmt.Printf("Value(labelKey, n): %v\n", Value(labelKey, n))
		if f, ok := d.callbacks[Value(keyKey, n)]; ok {
			f()
		}
	})
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			button.SetLabel(button.GetLabel() + c.Data)
		}
	}
	return button
}

func (d *DOM) parseTextArea(n *html.Node) tview.Primitive {
	return tview.NewInputField().SetFieldWidth(0).SetChangedFunc(func(text string) {
		if f, ok := d.stringCallbacks[Value(keyKey, n)]; ok {
			f(text)
		}
	})
}

func (d *DOM) parseChildren(n *html.Node, app *tview.Application) tview.Primitive {
	flex := makeFlex(n).SetFullScreen(true)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		child := d.parseNode(c, app)
		if child != nil {
			flex.AddItem(child, 0, 1, true)
		}
	}
	return flex
}

package ui

import (
	"strconv"
	"strings"

	"github.com/IvanKorchmit/akevitt"
	"github.com/rivo/tview"
	"golang.org/x/net/html"
)

func HasKey(key string, n *html.Node) bool {
	for _, v := range n.Attr {
		if strings.EqualFold(v.Key, key) {
			return true
		}
	}

	return false
}

func Value(key string, n *html.Node) string {
	for _, v := range n.Attr {
		if strings.EqualFold(v.Key, key) {
			return v.Val
		}
	}

	return ""
}

func makeFlex(n *html.Node) *tview.Flex {
	direction := tview.FlexRow
	switch strings.ToLower(Value(directionKey, n)) {
	case directionHorizontal:
		{
			direction = tview.FlexColumn
		}
	case directionVertical:
		{
			direction = tview.FlexRow
		}
	}

	fullscreen := HasKey(fullscreenKey, n)

	flex := tview.NewFlex().SetDirection(direction).SetFullScreen(fullscreen)

	if title := Value(titleKey, n); title != "" {
		flex.SetTitle(title)
	}

	flex.SetBorder(HasKey(borderKey, n))

	return flex
}

func IntValue(key string, n *html.Node, d int) int {
	v := Value(key, n)

	if v == "" {
		return d
	}

	integer, _ := strconv.Atoi(v)

	return integer
}

func Array(key string, n *html.Node) []string {
	v := Value(key, n)

	if v == "" {
		return nil
	}

	return strings.Split(v, " ")
}

func IntArrValue(key string, n *html.Node, d ...int) []int {
	arr := Array(key, n)

	if arr == nil {
		return d
	}

	return akevitt.MapSlice(arr, func(v string) int {
		i, _ := strconv.Atoi(v)

		return i
	})
}

func ValueFallback(key string, n *html.Node, d string) string {
	if v := Value(key, n); v != "" {
		return v
	}

	return d
}

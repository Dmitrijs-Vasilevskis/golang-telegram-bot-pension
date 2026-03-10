package helpers

import "regexp"

var threeOrMoreCloseBrackets = regexp.MustCompile(`\){3,}`)

func HasThreeOrMoreCloseBrackets(text string) bool {
    return threeOrMoreCloseBrackets.MatchString(text)
}
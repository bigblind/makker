package wordsplash

import (
	"math/rand"
	"strings"
)

var consonants = map[string]float64{
	"c": 0.044944183266288636,
	"b": 0.024103781967398506,
	"d": 0.06870870288696103,
	"g": 0.032553029935863266,
	"f": 0.03599411945265676,
	"h": 0.09845070194995074,
	"k": 0.012471930079645876,
	"j": 0.002471768526147434,
	"m": 0.03886977172490671,
	"l": 0.06502528312250604,
	"n": 0.10903245609783681,
	"q": 0.0015347582351895833,
	"p": 0.03116366984927059,
	"s": 0.10221489846362623,
	"r": 0.09672207951663195,
	"t": 0.14630284818817751,
	"w": 0.03812662563207807,
	"v": 0.015799932147530657,
	"y": 0.03189066059225513,
	"x": 0.002423302476615131,
	"z": 0.0011954958884634647,
}

var vowels = map[string]float64{
	"a": 0.21435695538057742,
	"i": 0.18283464566929133,
	"e": 0.3333858267716535,
	"u": 0.07238845144356955,
	"o": 0.1970341207349081,
}

func addLetter(current, action string) string {
	if action == "v" {
		current = addLetterFromSet(current, vowels)
	}

	if action == "c" {
		current = addLetterFromSet(current, consonants)
	}

	if action == "f" {
		for len(current) < 9 {
			if rand.Float64() < 0.6 {
				action = "c"
			} else {
				action = "v"
			}
			current = addLetter(current, action)
		}
	}

	n := len(current)
	nv := strings.Count(current, "a") + strings.Count(current, "e") + strings.Count(current, "i") + strings.Count(current, "o") + strings.Count(current, "u")
	nc := n - nv

	if nv < 2 {
		if 9-n == 2-nv {
			for i := 0; i < 2-nv; i++ {
				current = addLetterFromSet(current, vowels)
				return current
			}
		}

		if 9-n == 2-nc {
			for i := 0; i < 2-nc; i++ {
				current = addLetterFromSet(current, consonants)
				return current
			}
		}
	}

	return current
}

func addLetterFromSet(s string, set map[string]float64) string {
	var t float64 = 0
	target := rand.Float64()
	for l, prob := range set {
		t += prob
		if t >= target {
			return s + l
		}
	}

	return "e"
}

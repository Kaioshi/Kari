package words

import (
	"Kari/lib"
	"Kari/lib/logger"
	"bytes"
	"io/ioutil"
	"strings"
)

func getFileContents(filename string) [][]byte {
	wordList, err := ioutil.ReadFile("data/words/" + filename + ".txt")
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return bytes.Split(wordList, []byte("\n"))
}

func RandVerb() string {
	verb := string(*lib.RandSelectBytes(getFileContents("verbs")))
	if verb == "" {
		logger.Debug("Had to re-fetch Verb")
		verb = RandVerb()
	}
	return verb[0:strings.Index(verb, " ")]
}

func RandVerbs() string {
	verb := string(*lib.RandSelectBytes(getFileContents("verbs")))
	if verb == "" {
		logger.Debug("Had to re-fetch Verb[s]")
		verb = RandVerbs()
	}
	return strings.Split(verb, " ")[1]
}

func RandVerbed() string {
	verb := string(*lib.RandSelectBytes(getFileContents("verbs")))
	if verb == "" {
		logger.Debug("Had to re-fetch Verb[ed]")
		verb = RandVerbed()
	}
	return strings.Split(verb, " ")[2]
}

func RandVerbing() string {
	verb := string(*lib.RandSelectBytes(getFileContents("verbs")))
	if verb == "" {
		logger.Debug("Had to re-fetch Verb[ing]")
		verb = RandVerbing()
	}
	return strings.Split(verb, " ")[3]
}

func RandNoun() string {
	noun := string(*lib.RandSelectBytes(getFileContents("nouns")))
	if noun == "" {
		logger.Debug("Had to re-fetch Noun")
		noun = RandNoun()
	}
	return noun
}

func RandAdjective() string {
	adj := string(*lib.RandSelectBytes(getFileContents("adjectives")))
	if adj == "" {
		logger.Debug("Had to re-fetch Adjective")
		adj = RandAdjective()
	}
	return adj
}

func RandAdverb() string {
	adv := string(*lib.RandSelectBytes(getFileContents("adverbs")))
	if adv == "" {
		logger.Debug("Had to re-fetch Adverb")
		adv = RandAdverb()
	}
	return adv
}

func RandPronoun() string {
	pro := string(*lib.RandSelectBytes(getFileContents("pronouns")))
	if pro == "" {
		logger.Debug("Had to re-fetch Pronoun")
		pro = RandPronoun()
	}
	return pro
}

func RandPreposition() string {
	prep := string(*lib.RandSelectBytes(getFileContents("prepositions")))
	if prep == "" {
		logger.Debug("Had to re-fetch Preposition")
		prep = RandPreposition()
	}
	return prep
}

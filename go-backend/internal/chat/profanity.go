package chat

import (
	"strings"
	"unicode"
)

// Yomon so'zlar ro'yxati — o'zbek, rus, ingliz tillarida
var badWords = []string{
	// O'zbek
	"sik", "sikish", "sikin", "sikil", "sikaman", "siqish",
	"qo'taq", "qotaq", "kotaq", "ko'taq",
	"jo'ra", "jura",
	"amingni", "aming", "amini", "amiga",
	"ko'tak", "kotak",
	"axmoq", "ahmoq",
	"eshak", "eshakvoy",
	"harom", "haromzoda",
	"badzot", "bad zot",
	"iflos",
	"sharmanda",
	"yoqimsiz",
	"jinni", "jinnivor",
	"dangasa",
	"tentak",
	"nodon",
	"qashqir",
	"murdashuy",
	"cho'chqa",
	"it", "itvachcha",
	"kaltak",
	"maymun",

	// Rus (матерные)
	"блять", "блядь", "бля", "блят",
	"сука", "сучка", "сучара",
	"хуй", "хуя", "хуе", "хуёв", "нахуй", "нихуя", "похуй",
	"пизда", "пизд", "пиздец", "пиздат",
	"ебать", "ебан", "ебал", "заебал", "уебок", "уёб", "ёбан", "еба",
	"мудак", "мудила",
	"залупа",
	"дрочить", "дрочи",
	"шлюха", "шалава",
	"пидор", "пидар", "пидорас",
	"гандон",
	"даун",
	"дебил",
	"идиот",
	"тварь",
	"мразь",
	"урод",

	// Ingliz
	"fuck", "fucker", "fucking", "fck", "fcking", "fuk",
	"shit", "shitty", "bullshit",
	"bitch", "bitches",
	"ass", "asshole", "arsehole",
	"dick", "dickhead",
	"pussy",
	"cunt",
	"bastard",
	"whore",
	"slut",
	"nigger", "nigga",
	"faggot", "fag",
	"retard", "retarded",
	"damn", "dammit",
	"cock", "cocksucker",
	"motherfucker", "mf",
	"wtf", "stfu",
	"idiot", "moron",
}

// normalizeText — harflarni normallashtirish (lookalike harflar)
func normalizeText(s string) string {
	s = strings.ToLower(s)

	// Lookalike harflar almashtirish
	replacer := strings.NewReplacer(
		"0", "o",
		"1", "i",
		"3", "e",
		"4", "a",
		"5", "s",
		"@", "a",
		"$", "s",
		"!", "i",
	)
	s = replacer.Replace(s)

	// Bo'shliq va maxsus belgilarni olib tashlash
	var result []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '\'' || r == 0x2019 {
			result = append(result, r)
		} else {
			result = append(result, ' ')
		}
	}
	return string(result)
}

// ContainsProfanity — xabarda yomon so'z bormi tekshirish
func ContainsProfanity(message string) bool {
	normalized := normalizeText(message)
	words := strings.Fields(normalized)

	for _, word := range words {
		for _, bad := range badWords {
			if word == bad {
				return true
			}
			// So'z ichida ham tekshirish (masalan "blya" ichida "bly")
			if len(bad) >= 3 && strings.Contains(word, bad) {
				return true
			}
		}
	}

	// Butun matnda ham tekshirish (bo'shliqsiz yozilgan holat)
	noSpaces := strings.ReplaceAll(normalized, " ", "")
	for _, bad := range badWords {
		if len(bad) >= 4 && strings.Contains(noSpaces, bad) {
			return true
		}
	}

	return false
}

// CensorMessage — yomon so'zlarni *** bilan almashtirish
func CensorMessage(message string) string {
	normalized := normalizeText(message)
	result := message

	for _, bad := range badWords {
		if strings.Contains(normalized, bad) {
			// Original textdan topib almashtirish
			lower := strings.ToLower(result)
			idx := strings.Index(lower, bad)
			for idx != -1 {
				censored := strings.Repeat("*", len(bad))
				result = result[:idx] + censored + result[idx+len(bad):]
				lower = strings.ToLower(result)
				idx = strings.Index(lower, bad)
			}
		}
	}

	return result
}

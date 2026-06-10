package content

import (
	"fmt"
	"math/rand"
	"time"
)

type NumberExercise struct {
	Question string
	Answer   string
	Help     string
}

func GetNumberExercises() []NumberExercise {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var exercises []NumberExercise

	// Basic numbers (0-20)
	basics := map[int]string{
		0: "null", 1: "eins", 2: "zwei", 3: "drei", 4: "vier", 5: "fünf",
		6: "sechs", 7: "sieben", 8: "acht", 9: "neun", 10: "zehn",
		11: "elf", 12: "zwölf", 13: "dreizehn", 14: "vierzehn", 15: "fünfzehn",
		16: "sechzehn", 17: "siebzehn", 18: "achtzehn", 19: "neunzehn", 20: "zwanzig",
	}

	for i := 0; i <= 20; i++ {
		exercises = append(exercises, NumberExercise{
			Question: fmt.Sprintf("%d", i),
			Answer:   basics[i],
			Help:     "Basic number",
		})
	}

	// Tens
	tens := map[int]string{
		30: "dreißig", 40: "vierzig", 50: "fünfzig", 60: "sechzig",
		70: "siebzig", 80: "achtzig", 90: "neunzig", 100: "hundert",
	}
	for v, s := range tens {
		exercises = append(exercises, NumberExercise{
			Question: fmt.Sprintf("%d", v),
			Answer:   s,
			Help:     "Tens",
		})
	}

	// Random numbers between 21 and 99
	for i := 0; i < 20; i++ {
		n := r.Intn(79) + 21
		if n%10 == 0 {
			continue // Already covered by tens
		}
		exercises = append(exercises, NumberExercise{
			Question: fmt.Sprintf("%d", n),
			Answer:   formatGermanNumber(n),
			Help:     "Compound number (one-and-twenty style)",
		})
	}

	// Hundreds
	for i := 0; i < 10; i++ {
		n := (r.Intn(9)+1)*100 + r.Intn(100)
		exercises = append(exercises, NumberExercise{
			Question: fmt.Sprintf("%d", n),
			Answer:   formatGermanNumber(n),
			Help:     "Hundreds",
		})
	}

	// Thousands
	for i := 0; i < 10; i++ {
		n := (r.Intn(9)+1)*1000 + r.Intn(1000)
		exercises = append(exercises, NumberExercise{
			Question: fmt.Sprintf("%d", n),
			Answer:   formatGermanNumber(n),
			Help:     "Thousands",
		})
	}

	// Time exercises (increased to 30 exercises)
	for i := 0; i < 30; i++ {
		hour := r.Intn(12) + 1
		minute := []int{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55}[r.Intn(12)]
		exercises = append(exercises, NumberExercise{
			Question: fmt.Sprintf("%02d:%02d", hour, minute),
			Answer:   formatGermanTime(hour, minute),
			Help:     "Time expression",
		})
	}

	// Ordinal numbers (1-20)
	ordinalsEng := map[int]string{
		1: "first", 2: "second", 3: "third", 4: "fourth", 5: "fifth",
		6: "sixth", 7: "seventh", 8: "eighth", 9: "ninth", 10: "tenth",
		11: "eleventh", 12: "twelfth", 13: "thirteenth", 14: "fourteenth",
		15: "fifteenth", 16: "sixteenth", 17: "seventeenth", 18: "eighteenth",
		19: "nineteenth", 20: "twentieth",
	}
	for i := 1; i <= 20; i++ {
		exercises = append(exercises, NumberExercise{
			Question: ordinalsEng[i],
			Answer:   formatGermanOrdinal(i),
			Help:     "Ordinal number",
		})
	}

	// Year exercises (1900-2025)
	for i := 0; i < 15; i++ {
		year := r.Intn(126) + 1900
		exercises = append(exercises, NumberExercise{
			Question: fmt.Sprintf("Year %d", year),
			Answer:   formatGermanYear(year),
			Help:     "Year pronunciation",
		})
	}

	// Shuffle
	r.Shuffle(len(exercises), func(i, j int) {
		exercises[i], exercises[j] = exercises[j], exercises[i]
	})

	return exercises
}

func formatGermanYear(y int) string {
	if y < 1100 || y >= 2000 {
		return formatGermanNumber(y)
	}
	hundreds := y / 100
	rem := y % 100
	hStr := formatGermanNumber(hundreds) + "hundert"
	if rem == 0 {
		return hStr
	}
	return hStr + formatGermanNumber(rem)
}

func formatGermanNumber(n int) string {
	if n <= 20 {
		basics := map[int]string{
			0: "null", 1: "eins", 2: "zwei", 3: "drei", 4: "vier", 5: "fünf",
			6: "sechs", 7: "sieben", 8: "acht", 9: "neun", 10: "zehn",
			11: "elf", 12: "zwölf", 13: "dreizehn", 14: "vierzehn", 15: "fünfzehn",
			16: "sechzehn", 17: "siebzehn", 18: "achtzehn", 19: "neunzehn", 20: "zwanzig",
		}
		return basics[n]
	}

	if n < 100 {
		ones := n % 10
		tens := n / 10
		tensStr := map[int]string{
			2: "zwanzig", 3: "dreißig", 4: "vierzig", 5: "fünfzig",
			6: "sechzig", 7: "siebzig", 8: "achtzig", 9: "neunzig",
		}[tens]

		if ones == 0 {
			return tensStr
		}
		onesStr := map[int]string{
			1: "ein", 2: "zwei", 3: "drei", 4: "vier", 5: "fünf",
			6: "sechs", 7: "sieben", 8: "acht", 9: "neun",
		}[ones]

		return onesStr + "und" + tensStr
	}

	if n < 1000 {
		h := n / 100
		rem := n % 100
		hStr := "hundert"
		if h > 1 {
			hStr = formatGermanNumber(h) + "hundert"
		}
		if rem == 0 {
			return hStr
		}
		return hStr + formatGermanNumber(rem)
	}

	if n < 10000 {
		t := n / 1000
		rem := n % 1000
		tStr := "tausend"
		if t > 1 {
			tStr = formatGermanNumber(t) + "tausend"
		}
		if rem == 0 {
			return tStr
		}
		return tStr + formatGermanNumber(rem)
	}

	return fmt.Sprintf("%d", n) // Fallback for very large numbers
}

func formatGermanOrdinal(n int) string {
	if n == 1 {
		return "erste"
	}
	if n == 2 {
		return "zweite"
	}
	if n == 3 {
		return "dritte"
	}
	if n == 7 {
		return "siebte"
	}
	if n == 8 {
		return "achte"
	}
	if n < 20 {
		return formatGermanNumber(n) + "te"
	}
	return formatGermanNumber(n) + "ste"
}

func formatGermanTime(h, m int) string {
	if m == 0 {
		return fmt.Sprintf("%d Uhr", h)
	}
	if m == 30 {
		nextHour := (h % 12) + 1
		return fmt.Sprintf("halb %d", nextHour)
	}
	if m == 15 {
		return fmt.Sprintf("viertel nach %d", h)
	}
	if m == 45 {
		nextHour := (h % 12) + 1
		return fmt.Sprintf("viertel vor %d", nextHour)
	}
	if m < 30 {
		return fmt.Sprintf("%d nach %d", m, h)
	}
	// m > 30
	nextHour := (h % 12) + 1
	return fmt.Sprintf("%d vor %d", 60-m, nextHour)
}

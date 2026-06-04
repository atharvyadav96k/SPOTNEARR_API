package token

import (
	"strconv"
	"strings"
)

// UnitType groups units into semantic categories.
type UnitType string

const (
	UnitTypeWeight      UnitType = "weight"
	UnitTypeVolume      UnitType = "volume"
	UnitTypeLength      UnitType = "length"
	UnitTypeArea        UnitType = "area"
	UnitTypeStorage     UnitType = "storage"
	UnitTypePower       UnitType = "power"
	UnitTypeSpeed       UnitType = "speed"
	UnitTypeCount       UnitType = "count"
	UnitTypeClothSize   UnitType = "clothing_size"
	UnitTypeNumericSize UnitType = "numeric_size" // shoe / waist / bed size — context-free
	UnitTypeTextile     UnitType = "textile"      // tc (thread count), gsm
)

// UnitMatch is one extracted unit measurement.
type UnitMatch struct {
	Quantity float64
	RawUnit  string
	Type     UnitType
}

type unitDef struct {
	canonical string
	kind      UnitType
}

var unitDefs = map[string]unitDef{
	// ── Weight ────────────────────────────────────────────────────────────────
	"g":         {"g", UnitTypeWeight},
	"gm":        {"g", UnitTypeWeight},
	"gram":      {"g", UnitTypeWeight},
	"grams":     {"g", UnitTypeWeight},
	"grm":       {"g", UnitTypeWeight},
	"kg":        {"kg", UnitTypeWeight},
	"kgs":       {"kg", UnitTypeWeight},
	"kilo":      {"kg", UnitTypeWeight},
	"kilogram":  {"kg", UnitTypeWeight},
	"kilograms": {"kg", UnitTypeWeight},
	"lb":        {"lb", UnitTypeWeight},
	"lbs":       {"lb", UnitTypeWeight},
	"pound":     {"lb", UnitTypeWeight},
	"pounds":    {"lb", UnitTypeWeight},
	"oz":        {"oz", UnitTypeWeight},
	"ounce":     {"oz", UnitTypeWeight},
	"ounces":    {"oz", UnitTypeWeight},
	"mg":        {"mg", UnitTypeWeight},
	"milligram": {"mg", UnitTypeWeight},

	// ── Volume ────────────────────────────────────────────────────────────────
	"ml":         {"ml", UnitTypeVolume},
	"milliliter": {"ml", UnitTypeVolume},
	"millilitre": {"ml", UnitTypeVolume},
	"l":          {"l", UnitTypeVolume},
	"liter":      {"l", UnitTypeVolume},
	"litre":      {"l", UnitTypeVolume},
	"liters":     {"l", UnitTypeVolume},
	"litres":     {"l", UnitTypeVolume},
	"lt":         {"l", UnitTypeVolume},
	"ltr":        {"l", UnitTypeVolume},
	"floz":       {"fl_oz", UnitTypeVolume},
	"gallon":     {"gallon", UnitTypeVolume},
	"gallons":    {"gallon", UnitTypeVolume},

	// ── Length ────────────────────────────────────────────────────────────────
	"mm":         {"mm", UnitTypeLength},
	"millimeter": {"mm", UnitTypeLength},
	"millimetre": {"mm", UnitTypeLength},
	"cm":         {"cm", UnitTypeLength},
	"centimeter": {"cm", UnitTypeLength},
	"centimetre": {"cm", UnitTypeLength},
	"m":          {"m", UnitTypeLength},
	"meter":      {"m", UnitTypeLength},
	"metre":      {"m", UnitTypeLength},
	"inch":       {"inch", UnitTypeLength},
	"inches":     {"inch", UnitTypeLength},
	"in":         {"inch", UnitTypeLength},
	"ft":         {"ft", UnitTypeLength},
	"feet":       {"ft", UnitTypeLength},
	"foot":       {"ft", UnitTypeLength},

	// ── Area ──────────────────────────────────────────────────────────────────
	"sqft":   {"sqft", UnitTypeArea},
	"sqm":    {"sqm", UnitTypeArea},
	"sqyard": {"sqyard", UnitTypeArea},
	"sqyrd":  {"sqyard", UnitTypeArea},

	// ── Storage ───────────────────────────────────────────────────────────────
	"kb":       {"kb", UnitTypeStorage},
	"mb":       {"mb", UnitTypeStorage},
	"gb":       {"gb", UnitTypeStorage},
	"tb":       {"tb", UnitTypeStorage},
	"pb":       {"pb", UnitTypeStorage},
	"megabyte": {"mb", UnitTypeStorage},
	"gigabyte": {"gb", UnitTypeStorage},
	"terabyte": {"tb", UnitTypeStorage},

	// ── Power / Battery ───────────────────────────────────────────────────────
	"w":        {"w", UnitTypePower},
	"watt":     {"w", UnitTypePower},
	"watts":    {"w", UnitTypePower},
	"kw":       {"kw", UnitTypePower},
	"kilowatt": {"kw", UnitTypePower},
	"mah":      {"mah", UnitTypePower},
	"ah":       {"ah", UnitTypePower},
	"v":        {"v", UnitTypePower},
	"volt":     {"v", UnitTypePower},
	"volts":    {"v", UnitTypePower},
	"va":       {"va", UnitTypePower},

	// ── Speed / Performance ───────────────────────────────────────────────────
	"mbps": {"mbps", UnitTypeSpeed},
	"gbps": {"gbps", UnitTypeSpeed},
	"kbps": {"kbps", UnitTypeSpeed},
	"ghz":  {"ghz", UnitTypeSpeed},
	"mhz":  {"mhz", UnitTypeSpeed},
	"hz":   {"hz", UnitTypeSpeed},
	"rpm":  {"rpm", UnitTypeSpeed},

	// ── Count / Pack ─────────────────────────────────────────────────────────
	"pack":     {"pack", UnitTypeCount},
	"pcs":      {"pcs", UnitTypeCount},
	"pc":       {"pcs", UnitTypeCount},
	"piece":    {"pcs", UnitTypeCount},
	"pieces":   {"pcs", UnitTypeCount},
	"pair":     {"pair", UnitTypeCount},
	"pairs":    {"pair", UnitTypeCount},
	"set":      {"set", UnitTypeCount},
	"sets":     {"set", UnitTypeCount},
	"roll":     {"roll", UnitTypeCount},
	"rolls":    {"roll", UnitTypeCount},
	"sheet":    {"sheet", UnitTypeCount},
	"sheets":   {"sheet", UnitTypeCount},
	"box":      {"box", UnitTypeCount},
	"boxes":    {"box", UnitTypeCount},
	"bottle":   {"bottle", UnitTypeCount},
	"bottles":  {"bottle", UnitTypeCount},
	"tube":     {"tube", UnitTypeCount},
	"tubes":    {"tube", UnitTypeCount},
	"strip":    {"strip", UnitTypeCount},
	"strips":   {"strip", UnitTypeCount},
	"capsule":  {"capsule", UnitTypeCount},
	"capsules": {"capsule", UnitTypeCount},
	"jar":      {"jar", UnitTypeCount},
	"jars":     {"jar", UnitTypeCount},
	"tablet":   {"tablet", UnitTypeCount},
	"tablets":  {"tablet", UnitTypeCount},

	// ── Clothing sizes ────────────────────────────────────────────────────────
	"xs":   {"xs", UnitTypeClothSize},
	"xxs":  {"xxs", UnitTypeClothSize},
	"xl":   {"xl", UnitTypeClothSize},
	"xxl":  {"xxl", UnitTypeClothSize},
	"xxxl": {"xxxl", UnitTypeClothSize},
	"2xl":  {"2xl", UnitTypeClothSize},
	"3xl":  {"3xl", UnitTypeClothSize},
	"4xl":  {"4xl", UnitTypeClothSize},

	// ── Textile ───────────────────────────────────────────────────────────────
	"tc":     {"tc", UnitTypeTextile},  // thread count (300tc bedsheets)
	"gsm":    {"gsm", UnitTypeTextile}, // grams per square metre (fabric weight)
	"denier": {"denier", UnitTypeTextile},
}

// ambiguousSingleLetters are single-letter tokens that are ALSO valid unit
// abbreviations (m=meter, l=liter, v=volt, w=watt, g=gram).
// They should only be treated as units when preceded by a number.
// Without a number they are noise (e.g. "m2" chip split → "m" orphan).
var ambiguousSingleLetters = map[string]struct{}{
	"m": {}, "l": {}, "v": {}, "w": {}, "g": {},
}

// clothingSizeLetters: s/m/l need the "size" keyword before them.
var clothingSizeLetters = map[string]struct{}{
	"s": {}, "m": {}, "l": {},
}

// ExtractUnitFilter scans tokens, extracts quantity+unit pairs, and returns
// a cleaned token slice with unit tokens removed.
func ExtractUnitFilter(tokens []string) ([]UnitMatch, []string) {
	var matches []UnitMatch
	skip := make([]bool, len(tokens))

	for i, t := range tokens {
		if skip[i] {
			continue
		}

		// ── "size N"  → numeric size (shoe / waist) ───────────────────────────
		if t == "size" && i+1 < len(tokens) {
			next := tokens[i+1]
			if isNumericToken(next) {
				n, _ := strconv.ParseFloat(next, 64)
				matches = append(matches, UnitMatch{Quantity: n, RawUnit: "size", Type: UnitTypeNumericSize})
				skip[i], skip[i+1] = true, true
				continue
			}
			// "size xl/xxl/…"
			if def, ok := unitDefs[next]; ok && def.kind == UnitTypeClothSize {
				matches = append(matches, UnitMatch{RawUnit: def.canonical, Type: UnitTypeClothSize})
				skip[i], skip[i+1] = true, true
				continue
			}
			// "size s/m/l"
			if _, ok := clothingSizeLetters[next]; ok {
				matches = append(matches, UnitMatch{RawUnit: next, Type: UnitTypeClothSize})
				skip[i], skip[i+1] = true, true
				continue
			}
			// "size" with no recognised following token — keep "size" as a token
			continue
		}

		// ── Unambiguous clothing sizes (xl, xxl, 2xl…) without "size" keyword ─
		if def, ok := unitDefs[t]; ok && def.kind == UnitTypeClothSize {
			matches = append(matches, UnitMatch{RawUnit: def.canonical, Type: UnitTypeClothSize})
			skip[i] = true
			continue
		}

		// ── Numeric token ──────────────────────────────────────────────────────
		if isNumericToken(t) {
			if i+1 < len(tokens) && !skip[i+1] {
				next := tokens[i+1]
				if def, ok := resolveUnit(next); ok {
					// Don't consume ambiguous single letters unless they ARE a unit here
					// (already confirmed by resolveUnit returning ok)
					n, _ := strconv.ParseFloat(t, 64)
					matches = append(matches, UnitMatch{Quantity: n, RawUnit: def.canonical, Type: def.kind})
					skip[i], skip[i+1] = true, true
					continue
				}
			}
			// bare number — leave in tokens
			continue
		}

		// ── Stray unit token (unit with no preceding number) ──────────────────
		// Ambiguous single letters (m, l, g, v, w) without a preceding number
		// are chip/model leftovers — keep them as tokens, don't drop silently.
		if _, isAmbig := ambiguousSingleLetters[t]; isAmbig {
			// keep in tokens — could be part of a model name like "m2", "v60"
			continue
		}

		// Non-ambiguous stray unit (e.g. a lone "kg" or "ml") — drop as noise
		if def, ok := resolveUnit(t); ok && def.kind != UnitTypeClothSize {
			skip[i] = true
			continue
		}
	}

	kept := make([]string, 0, len(tokens))
	for i, t := range tokens {
		if !skip[i] {
			kept = append(kept, t)
		}
	}

	return matches, kept
}

// NormalizeUnit returns the canonical unit string for the given alias,
// or the input lowercased if the alias is not recognised.
func NormalizeUnit(unit string) string {
	lower := strings.ToLower(strings.TrimSpace(unit))
	if def, ok := unitDefs[lower]; ok {
		return def.canonical
	}
	return lower
}

func resolveUnit(t string) (unitDef, bool) {
	if def, ok := unitDefs[t]; ok {
		return def, true
	}
	trimmed := strings.TrimRight(t, ".")
	if def, ok := unitDefs[trimmed]; ok {
		return def, true
	}
	return unitDef{}, false
}

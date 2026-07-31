package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"deutsch-tui/internal/core"
)

const systemPromptExplain = `You are a patient and helpful German tutor. 
A student is reviewing a flashcard and needs a brief explanation of the grammar, usage, or cultural context.
Keep your explanation concise (2-4 sentences). 
Focus on the most important rule or pattern that helps them remember or use this correctly.
If there are common pitfalls or "false friends", mention them.
Always output plain text without any markdown formatting.`

const systemPromptExplainDict = `You are a patient and helpful German tutor.
A student looked up a dictionary entry and wants a brief pedagogical explanation.
Keep your explanation concise (2-4 sentences).
Cover the most useful grammar note (gender/article, strong/weak verb, case government), a natural collocation or usage tip, and one short example if helpful.
Mention common pitfalls or false friends when relevant.
Always output plain text without any markdown formatting.`

func ExplainCard(ctx context.Context, provider Provider, card core.Card) (string, error) {
	if provider == nil {
		return "", errors.New("AI provider is disabled; enable one in Settings to get explanations")
	}
	chat, ok := provider.(ChatProvider)
	if !ok {
		return "Offline/template provider cannot provide detailed explanations.", nil
	}
	user := userPromptForExplanation(card)
	raw, err := chat.SendChat(ctx, systemPromptExplain, user)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(raw), nil
}

func ExplainDictionaryEntry(ctx context.Context, provider Provider, entry core.DictionaryEntry) (string, error) {
	if provider == nil {
		return "", errors.New("AI provider is disabled; enable one in Settings to get explanations")
	}
	chat, ok := provider.(ChatProvider)
	if !ok {
		return "Offline/template provider cannot provide detailed explanations.", nil
	}
	user := userPromptForDictionaryExplanation(entry)
	raw, err := chat.SendChat(ctx, systemPromptExplainDict, user)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(raw), nil
}

func userPromptForExplanation(card core.Card) string {
	var b strings.Builder
	fmt.Fprintf(&b, "German: %s\n", card.Prompt)
	fmt.Fprintf(&b, "English: %s\n", card.Answer)
	if strings.TrimSpace(card.Extra) != "" {
		fmt.Fprintf(&b, "Context/Extra: %s\n", card.Extra)
	}
	if len(card.Tags) > 0 {
		fmt.Fprintf(&b, "Tags: %s\n", strings.Join(card.Tags, ", "))
	}
	b.WriteString("\nPlease explain how to use this correctly or the grammar behind it.")
	return b.String()
}

func userPromptForDictionaryExplanation(entry core.DictionaryEntry) string {
	var b strings.Builder
	fmt.Fprintf(&b, "German headword: %s\n", entry.Word)
	fmt.Fprintf(&b, "English: %s\n", entry.Translation)
	if entry.WordClass != "" {
		fmt.Fprintf(&b, "Word class: %s\n", entry.WordClass)
	}
	if entry.Gender != "" {
		fmt.Fprintf(&b, "Gender: %s\n", entry.Gender)
	}
	if entry.Forms != "" {
		fmt.Fprintf(&b, "Forms: %s\n", entry.Forms)
	}
	if len(entry.Tags) > 0 {
		fmt.Fprintf(&b, "Tags: %s\n", strings.Join(entry.Tags, ", "))
	}
	if len(entry.Examples) > 0 {
		b.WriteString("Examples:\n")
		for _, ex := range entry.Examples {
			fmt.Fprintf(&b, "- %s\n", ex)
		}
	}
	b.WriteString("\nPlease explain how to use this word correctly.")
	return b.String()
}

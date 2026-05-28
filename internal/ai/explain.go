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

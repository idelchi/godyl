// Package aimatch selects release assets with a language model when Godyl's
// deterministic matcher returns equally ranked candidates.
package aimatch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	aiconfig "github.com/idelchi/godyl/internal/config/ai"
	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/internal/match"

	"charm.land/fantasy"
	"charm.land/fantasy/providers/openai"
	"charm.land/fantasy/providers/openaicompat"
)

const (
	providerOllama = "ollama"
	providerOpenAI = "openai"

	defaultOllamaURL   = "http://localhost:11434/v1"
	defaultOllamaModel = "gemma3:4b"
	defaultOpenAIModel = "gpt-5-nano"
)

// Selector implements match.AssetSelector through Fantasy.
type Selector struct {
	config aiconfig.Config
}

// New creates a selector whose provider is initialized only for an ambiguous
// deterministic match.
func New(config aiconfig.Config) *Selector {
	return &Selector{config: config}
}

// Select asks the configured model for exactly one of the supplied candidates.
func (s *Selector) Select(ctx context.Context, request match.SelectionRequest) (string, error) {
	model, modelID, err := s.languageModel(ctx)
	if err != nil {
		return "", err
	}

	prompt, err := selectionPrompt(request)
	if err != nil {
		return "", err
	}

	const maxOutputTokens int64 = 128

	debug.Debug("deterministic asset matching is ambiguous; asking %s model %q", s.config.Provider, modelID)

	response, err := model.Generate(ctx, fantasy.Call{
		Prompt: fantasy.Prompt{
			fantasy.NewSystemMessage(
				"Select the one release asset that best matches the requested target. " +
					"Reply with only its exact candidate string and nothing else.",
			),
			fantasy.NewUserMessage(prompt),
		},
		MaxOutputTokens: new(maxOutputTokens),
		ProviderOptions: s.providerOptions(),
	})
	if err != nil {
		return "", fmt.Errorf("selecting asset with %s model %q: %w", s.config.Provider, modelID, err)
	}

	selected, err := exactCandidate(response.Content.Text(), request.Candidates)
	if err != nil {
		return "", err
	}

	debug.Debug("AI asset selector chose %q", selected)

	return selected, nil
}

func (s *Selector) languageModel(ctx context.Context) (fantasy.LanguageModel, string, error) {
	provider, modelID, err := s.provider()
	if err != nil {
		return nil, "", err
	}

	model, err := provider.LanguageModel(ctx, modelID)
	if err != nil {
		return nil, "", fmt.Errorf("creating %s model %q: %w", s.config.Provider, modelID, err)
	}

	return model, modelID, nil
}

func (s *Selector) providerOptions() fantasy.ProviderOptions {
	if strings.EqualFold(strings.TrimSpace(s.config.Provider), providerOpenAI) {
		return openai.NewProviderOptions(&openai.ProviderOptions{
			ReasoningEffort: new(openai.ReasoningEffortMinimal),
		})
	}

	return fantasy.ProviderOptions{
		providerOllama: &openaicompat.ProviderOptions{
			ReasoningEffort: new(openai.ReasoningEffortNone),
		},
	}
}

func (s *Selector) provider() (fantasy.Provider, string, error) {
	providerName := strings.ToLower(strings.TrimSpace(s.config.Provider))

	switch providerName {
	case providerOllama:
		model := valueOr(s.config.Model, defaultOllamaModel)
		url := valueOr(s.config.URL, defaultOllamaURL)
		apiKey := valueOr(s.config.APIKey, providerOllama)

		provider, err := openaicompat.New(
			openaicompat.WithName(providerOllama),
			openaicompat.WithBaseURL(url),
			openaicompat.WithAPIKey(apiKey),
			openaicompat.WithObjectMode(fantasy.ObjectModeText),
		)

		return provider, model, err

	case providerOpenAI:
		apiKey := valueOr(s.config.APIKey, os.Getenv("OPENAI_API_KEY"))
		if apiKey == "" {
			return nil, "", errors.New("OpenAI API key is empty; set GODYL_AI_API_KEY or OPENAI_API_KEY")
		}

		options := []openai.Option{
			openai.WithAPIKey(apiKey),
			openai.WithUseResponsesAPI(),
		}

		if s.config.URL != "" {
			options = append(options, openai.WithBaseURL(s.config.URL))
		}

		provider, err := openai.New(options...)

		return provider, valueOr(s.config.Model, defaultOpenAIModel), err

	default:
		return nil, "", fmt.Errorf("unsupported AI provider %q (supported: ollama, openai)", s.config.Provider)
	}
}

func selectionPrompt(request match.SelectionRequest) (string, error) {
	candidates, err := json.Marshal(request.Candidates)
	if err != nil {
		return "", fmt.Errorf("encoding asset candidates: %w", err)
	}

	return fmt.Sprintf(
		"Target: %s\nOperating system: %s\nArchitecture: %s\nDistribution: %s\nC library: %s\nCandidates (JSON): %s",
		request.Target,
		request.Platform.OS.String(),
		request.Platform.Architecture.String(),
		request.Platform.Distribution.String(),
		request.Platform.Library.String(),
		candidates,
	), nil
}

func exactCandidate(response string, candidates []string) (string, error) {
	answer := strings.TrimSpace(response)

	answer = strings.TrimSpace(strings.Trim(answer, "`\"'"))

	for _, candidate := range candidates {
		if answer == candidate {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("model returned %q, which is not an exact candidate", answer)
}

func valueOr(value, fallback string) string {
	if value != "" {
		return value
	}

	return fallback
}

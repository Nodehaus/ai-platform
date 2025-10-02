package training_datasets

import (
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/google/uuid"

	"ai-platform/cmd/web"
)

const defaultPrompt = `You are an AI assistant specializing in creating educational question/answer pairs from legal texts. Generate diverse questions and answers based on the provided EUR-Lex text.

Create questions of three complexity levels:
- Simple: Basic facts, definitions, key terms
- Medium: Relationships, implications, comparisons
- Complex: Analysis, synthesis, evaluation, critical thinking

Ensure variety in question types and complexity levels.`

func TrainingDatasetStep3Handler(w http.ResponseWriter, r *http.Request) {
	token := web.GetTokenFromCookie(r)
	if token == "" {
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	// Extract project ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 || pathParts[3] == "" {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	projectIDStr := pathParts[3]
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID format", http.StatusBadRequest)
		return
	}

	// Get parameters from query
	corpus := r.URL.Query().Get("corpus")
	language := r.URL.Query().Get("language")

	if language == "" {
		http.Redirect(w, r, "/web/projects/"+projectIDStr+"/training-datasets/step1", http.StatusSeeOther)
		return
	}

	// Fetch project details to get the project name
	projectName, err := fetchProjectName(r, token, projectID)
	if err != nil {
		web.ClearTokenCookie(w)
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	templ.Handler(TrainingDatasetStep3(projectIDStr, projectName, corpus, language, defaultPrompt)).ServeHTTP(w, r)
}
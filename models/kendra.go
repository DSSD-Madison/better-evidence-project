package models

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/kendra"
	"github.com/aws/aws-sdk-go-v2/service/kendra/types"

	"github.com/DSSD-Madison/gmu/internal"
)

var opts = kendra.Options{
	Credentials: prov,
	Region:      region,
}

var client = kendra.New(opts)

type KendraResult struct {
	Title   string
	Excerpt string
	Link    string
}

type KendraResults struct {
	QueryOutput kendra.QueryOutput
	Results     []KendraResult
	Query       string
	Count       int
}

type KendraSuggestions struct {
	SuggestionOutput kendra.GetQuerySuggestionsOutput
	Suggestions      []string
	Suggestion       types.SuggestionTextWithHighlights
}

func queryOutputToResults(out kendra.QueryOutput) KendraResults {
	results := KendraResults{
		Results: []KendraResult{},
	}

	for _, item := range out.ResultItems {
		res := KendraResult{
			Title:   internal.TrimExtension(*item.DocumentTitle.Text),
			Excerpt: *item.DocumentExcerpt.Text,
			Link:    *item.DocumentURI,
		}
		results.Results = append(results.Results, res)
		results.Count = int(*out.TotalNumberOfResults)
	}

	return results
}

func MakeQuery(query string) KendraResults {
	kendraQuery := kendra.QueryInput{
		IndexId:   &indexId,
		QueryText: &query,
	}
	out, err := client.Query(context.TODO(), &kendraQuery)

	// TODO: this needs to be fixed to a proper error
	if err != nil {
		log.Printf("Kendra Query Failed %+v", err)
	}

	results := queryOutputToResults(*out)
	results.Query = query
	return results
}

func querySuggestionsOutputToSuggestions(out kendra.GetQuerySuggestionsOutput) KendraSuggestions {
	suggestions := KendraSuggestions{
		SuggestionOutput: out,
		Suggestions:      make([]string, 0),
	}

	for _, item := range out.Suggestions {
		suggestions.Suggestions = append(suggestions.Suggestions, *item.Value.Text.Text)
	}

	return suggestions
}

func GetSuggestions(query string) KendraSuggestions {
	kendraQuery := kendra.GetQuerySuggestionsInput{
		IndexId:   &indexId,
		QueryText: &query,
	}
	out, err := client.GetQuerySuggestions(context.TODO(), &kendraQuery)
	if err != nil {
		log.Printf("Kendra Suggestions Query Failed %+v", err)
	}

	suggestions := querySuggestionsOutputToSuggestions(*out)
	return suggestions
}

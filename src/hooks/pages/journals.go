package pages

import "baby-blog/types"

func Journals(pageData map[string]interface{}, db *types.Models) map[string]interface{} {
	// Get all journals from the database
	journals, err := db.Journal.GetAll()

	if err != nil {
		pageData["Failure"] = "Failed to load journals"
		return pageData
	}
	// Convert database journals to the format expected by the template
	journalsData := make([]map[string]string, 0, len(journals))
	for _, journal := range journals {
		journalsData = append(journalsData, map[string]string{
			"title":   journal.Title,
			"content": journal.Content,
			"date":    journal.CreatedAt.Format("2006-01-02"),
		})
	}
	// Return the data from the database
	pageData["journals"] = journalsData
	pageData["journalCount"] = len(journalsData)
	return pageData

}

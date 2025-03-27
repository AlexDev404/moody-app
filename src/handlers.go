package main

import (
	"baby-blog/database/models"
	"baby-blog/forms"
	"baby-blog/forms/validator"
	"net/http"
)

func (app *Application) JournalHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.Logger.Error(FormHandlerErrorMessage, "error", err)
		http.Error(w, FormHandlerBadRequestMessage, http.StatusBadRequest)
		return
	}

	validator := validator.NewValidator()
	formData, formErrors := forms.JournalForm(w, r, validator)

	for key, value := range formErrors {
		formData[key] = value
	}

	if formErrors == nil {
		journalData := &models.Journal{
			Title:   formData["title"].(string),
			Content: formData["content"].(string),
		}
		err := app.models.Journal.Insert(journalData)
		if err != nil {
			formData["Failure"] = "✗ Failed to submit journal entry. Please try again later."
		} else {
			formData["Message"] = "✓ Your journal entry has been submitted. Thank you!"
		}
	}

	app.Render(w, r, app.templates, formData)
}

func (app *Application) TodoHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.Logger.Error(FormHandlerErrorMessage, "error", err)
		http.Error(w, FormHandlerBadRequestMessage, http.StatusBadRequest)
		return
	}

	validator := validator.NewValidator()
	formData, formErrors := forms.TodoForm(w, r, validator)

	for key, value := range formErrors {
		formData[key] = value
	}

	if formErrors == nil {
		todoData := &models.Todo{
			Task:      formData["task"].(string),
			Completed: false,
		}
		err := app.models.Todo.Insert(todoData)
		if err != nil {
			formData["Failure"] = "✗ Failed to submit todo item. Please try again later."
		} else {
			formData["Message"] = "✓ Your todo item has been submitted. Thank you!"
		}
	}

	app.Render(w, r, app.templates, formData)
}

func (app *Application) FeedbackHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.Logger.Error(FormHandlerErrorMessage, "error", err)
		http.Error(w, FormHandlerBadRequestMessage, http.StatusBadRequest)
		return
	}

	validator := validator.NewValidator()
	formData, formErrors := forms.FeedbackForm(w, r, validator)

	for key, value := range formErrors {
		formData[key] = value
	}

	if formErrors == nil {
		feedbackData := &models.Feedback{
			Fullname: formData["fullname"].(string),
			Email:    formData["email"].(string),
			Subject:  formData["subject"].(string),
			Message:  formData["message"].(string),
		}
		err := app.models.Feedback.Insert(feedbackData)
		if err != nil {
			formData["Failure"] = "✗ Failed to submit feedback. Please try again later."
		} else {
			formData["Message"] = "✓ Your feedback has been submitted. Thank you!"
		}
	}

	app.Render(w, r, app.templates, formData)
}

package dto

import (
	"encoding/json"
	"testing"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateHistoryRequest_ValidationAndDefaults(t *testing.T) {
	validate := validator.New()

	t.Run("Defaults are applied when fields are missing", func(t *testing.T) {
		jsonInput := `{"description": "test story"}`
		var req CreateHistoryRequest
		err := json.Unmarshal([]byte(jsonInput), &req)
		require.NoError(t, err)

		err = defaults.Set(&req)
		require.NoError(t, err)

		assert.Equal(t, "test story", req.Description)
		require.NotNil(t, req.ChapterSize)
		assert.Equal(t, 3, *req.ChapterSize)
		require.NotNil(t, req.ImageSize)
		assert.Equal(t, 1, *req.ImageSize)
		require.NotNil(t, req.Save)
		assert.True(t, *req.Save)

		err = validate.Struct(req)
		assert.NoError(t, err)
	})

	t.Run("Validation fails on invalid values", func(t *testing.T) {
		invalidChapter := 11
		req := CreateHistoryRequest{
			Description: "test",
			ChapterSize: &invalidChapter,
		}
		// defaults might fill others, but let's check validation on this one
		err := validate.Struct(req)
		assert.Error(t, err)
	})
    
    t.Run("Validation fails on missing description", func(t *testing.T) {
		req := CreateHistoryRequest{}
        err := defaults.Set(&req)
        require.NoError(t, err)
        
		err = validate.Struct(req)
		assert.Error(t, err)
	})
}

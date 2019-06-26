package delivery

import (
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/larien/family-tree/repository"
	"github.com/larien/family-tree/controller"
	"github.com/larien/family-tree/entity"
	"github.com/stretchr/testify/assert"
)


func TestPersonEndpoints(t *testing.T) {
	r, err := repository.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	c := controller.New(r)

	router := New(c)

	t.Run("should GET all People", func(t *testing.T) {
		r.Person.Clear()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/person", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should have created resource", func(t *testing.T) {
		r.Person.Clear()
		w := httptest.NewRecorder()

		payload := fmt.Sprintf(`[
			{
				"name": "Anakin"
			},
			{
				"name": "Luke",
				"parents": ["Anakin"]
			}
		]`)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/person", strings.NewReader(payload))
		router.ServeHTTP(w, req)

		var people *entity.Person
		json.NewDecoder(w.Body).Decode(&people)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("shouldn't create resource because of invalid payload", func(t *testing.T) {
		r.Person.Clear()
		w := httptest.NewRecorder()

		payload := fmt.Sprintf(`{
				"invalid": "parse"
			}`)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/person", strings.NewReader(payload))
		router.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should GET a Person", func(t *testing.T) {
		r.Person.Clear()
		w := httptest.NewRecorder()
		payload := fmt.Sprintf(`[
			{
				"name": "Leia",
				"parents": ["Anakin", "Padme"],
				"children": ["Ben"]
			}
		]`)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/person", strings.NewReader(payload))
		router.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)

		w = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/api/v1/person/Leia", nil)
		router.ServeHTTP(w, req)
		var people entity.Person
		assert.Nil(t, err)
		json.NewDecoder(w.Body).Decode(&people)
		assert.Equal(t, "Leia", people.Name)
		parents := []string{"Anakin", "Padme"}
		assert.Equal(t, parents, people.Parents)
		children := []string{"Ben"}
		assert.Equal(t, children, people.Children)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	r.DB.Session.Close()
	r.DB.Driver.Close()
}
package routes

import (
	"encoding/json"
	"html/template"
	"net/http"
	"spoonacular-api/db"
	"spoonacular-api/session"
	"testing"
)

func SetupUserRoutes()(err error) {
	err := nil
	//handle and serve static HTML pages(accessable to user)
	err := http.HandleFunc("/", HomePageHandler)
	err := http.HandleFunc("/profile", ProfilePageHandler)

	//handle and serve JSON data(accessed programmatically within the html pages)
	err := http.HandleFunc("/api/profile", ProfileHandler)
	err := http.HandleFunc("/api/add-favorite", AddFavoriteHandler)
	return err
}

// HomePageHandler serves the main HTML page when users visit the root URL.
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// serves the static profile page
func ProfilePageHandler(w http.ResponseWriter, r *http.Request)(err error) {
	session, _ := session.Store.Get(r, "session-name")
	userID, ok := session.Values["userID"].(int)

	//if user is not logged in, redirect to login page
	if !ok || userID == 0 {han
		tmpl, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, map[string]string{"ErrorMessage": "You must be logged in to view the profile page."})
		return
	}

	tmpl, err := template.ParseFiles("templates/profile.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// ProfileHandler fetches the user's profile data and serves it as JSON back to the frontend.
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := session.Store.Get(r, "session-name")
	userID, ok := session.Values["userID"].(int)
	if !ok || userID == 0 {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	userProfile, err := db.GetUserProfile(userID)
	if err != nil {
		http.Error(w, "Failed to load user profile", http.StatusInternalServerError)
		return
	}
	//send json response back to the frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

func TestProfilePageHandler(*t testing.T){
	var err := nil
	req,err := http.NewRequest("GET","/profile",nil)
	if err!= nil{
		t.Fatal(err)
	}
	resRec := httptest.NewRecorder()
	err := http.HandleFunc("/profile", ProfilePageHandler)
}

//Reference for writing testHandler
/*
func TestHealthCheckHandler(t *testing.T) {
    // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
    // pass 'nil' as the third parameter.
    req, err := http.NewRequest("GET", "/health-check", nil)
    if err != nil {
       t.Fatal(err)
    }
    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HealthCheckHandler)
    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method
    // directly and pass in our Request and ResponseRecorder.
    handler.ServeHTTP(rr, req)
    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    // Check the response body is what we expect.
    expected := `{"alive": true}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
 */

// AddFavoriteHandler handles adding a recipe to the user's favorites list.
func AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	session, _ := session.Store.Get(r, "session-name")
	userID, ok := session.Values["userID"].(int)
	if !ok || userID == 0 {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	var req struct {
		RecipeID int    `json:"recipe_id"`
		Title    string `json:"title"`
		Image    string `json:"image"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//add recipe to favorites table in the database(linked to user ID)
	err := db.AddRecipeToFavorites(userID, req.RecipeID, req.Title, req.Image)
	if err != nil {
		http.Error(w, "Failed to add recipe to favorites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe added to favorites successfully"})
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	enableCORS := strings.ToLower(os.Getenv("ENABLE_CORS")) == "true"
	handler := loginHandler
	if enableCORS {
		handler = corsMiddleware(handler)
	}

	http.HandleFunc("/login", handler)
	log.Println("LDAP login server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS for all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type Response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	var req Login
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ldapUrl := os.Getenv("LDAP_URL")
	baseDN := os.Getenv("LDAP_BASE_DN")
	bindDN := os.Getenv("LDAP_BIND_DN")
	bindPassword := os.Getenv("LDAP_BIND_PASSWORD")
	filterTemplate := os.Getenv("LDAP_SEARCH_FILTER")
	if filterTemplate == "" {
		filterTemplate = "(uid=%s)"
	}

	conn, err := ldap.DialURL(ldapUrl)
	if err != nil {
		http.Error(w, "LDAP connection failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	if err := conn.Bind(bindDN, bindPassword); err != nil {
		http.Error(w, "LDAP bind failed", http.StatusInternalServerError)
		return
	}

	filter := fmt.Sprintf(filterTemplate, ldap.EscapeFilter(req.Username))

	searchReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter,
		[]string{"dn"},
		nil,
	)

	result, err := conn.Search(searchReq)
	if err != nil || len(result.Entries) != 1 {
		http.Error(w, "User not found or multiple entries", http.StatusUnauthorized)
		return
	}

	userDN := result.Entries[0].DN

	if err := conn.Bind(userDN, req.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(Response{Success: true, Message: "Login successful"})
}

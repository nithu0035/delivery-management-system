package auth

import (
    "encoding/json"
    "net/http"
    "strconv"

    "delivery/internal/users"

    "github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, svc *Service, userSvc *users.Service, orderSvc interface{}) {
    r.HandleFunc("/auth/register", func(w http.ResponseWriter, req *http.Request) {
        var in struct{ Email, Password, Role string }
        _ = json.NewDecoder(req.Body).Decode(&in)
        if in.Role == "" { in.Role = "customer" }
        u, err := userSvc.Create(req.Context(), in.Email, in.Password, in.Role)
        if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        token, _ := svc.GenerateToken(req.Context(), u)
        json.NewEncoder(w).Encode(map[string]string{"token": token})
    }).Methods("POST")

    r.HandleFunc("/auth/login", func(w http.ResponseWriter, req *http.Request) {
        var in struct{ Email, Password string }
        _ = json.NewDecoder(req.Body).Decode(&in)
        u, err := userSvc.Authenticate(req.Context(), in.Email, in.Password)
        if err != nil { http.Error(w, "unauthorized", http.StatusUnauthorized); return }
        token, _ := svc.GenerateToken(req.Context(), u)
        json.NewEncoder(w).Encode(map[string]string{"token": token})
    }).Methods("POST")
}

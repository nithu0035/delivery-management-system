package orders

import (
    "encoding/json"
    "net/http"
    "strconv"

    "delivery/internal/auth"
    "delivery/internal/users"

    "github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, authSvc *auth.Service, svc *Service, userSvc *users.Service) {
    r.HandleFunc("/orders", func(w http.ResponseWriter, req *http.Request) {
        // create order (requires auth)
        token := req.Header.Get("Authorization")
        claims, err := authSvc.ParseToken(req.Context(), token)
        if err != nil { http.Error(w, "unauth", http.StatusUnauthorized); return }
        uidFloat := claims["sub"].(float64)
        uid := uint(uidFloat)
        var in struct{ Items interface{} }
        _ = json.NewDecoder(req.Body).Decode(&in)
        o, err := svc.Create(req.Context(), uid, in.Items)
        if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        json.NewEncoder(w).Encode(o)
    }).Methods("POST")

    r.HandleFunc("/orders", func(w http.ResponseWriter, req *http.Request) {
        // list for user (auth)
        token := req.Header.Get("Authorization")
        claims, err := authSvc.ParseToken(req.Context(), token)
        if err != nil { http.Error(w, "unauth", http.StatusUnauthorized); return }
        uidFloat := claims["sub"].(float64)
        uid := uint(uidFloat)
        role := claims["role"].(string)
        if role == "admin" {
            // list all
            var os []Order
            _ = svc.db.Find(&os).Error
            json.NewEncoder(w).Encode(os)
            return
        }
        os, err := svc.ListForUser(req.Context(), uid)
        if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
        json.NewEncoder(w).Encode(os)
    }).Methods("GET")

    r.HandleFunc("/orders/{id:[0-9]+}/cancel", func(w http.ResponseWriter, req *http.Request) {
        token := req.Header.Get("Authorization")
        claims, err := authSvc.ParseToken(req.Context(), token)
        if err != nil { http.Error(w, "unauth", http.StatusUnauthorized); return }
        uidFloat := claims["sub"].(float64)
        uid := uint(uidFloat)
        role := claims["role"].(string)
        vars := mux.Vars(req)
        id, _ := strconv.Atoi(vars["id"])
        if err := svc.Cancel(req.Context(), uint(id), uid, role=="admin"); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest); return
        }
        w.WriteHeader(http.StatusNoContent)
    }).Methods("POST")
}

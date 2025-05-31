package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
)

type CustomBucket struct {
	From    float64 `json:"from"`
	To      float64 `json:"to"`
	Percent float64 `json:"percent"`
}

type CustomTaxRequest struct {
	AnnualPay float64        `json:"annualPay"`
	Deduction float64        `json:"deduction"`
	Buckets   []CustomBucket `json:"buckets"`
}

type CustomTaxResponse struct {
	Tax    float64 `json:"tax"`
	Inhand float64 `json:"inhand"`
}

func RenderCustomForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "custom.html")))
	tmpl.Execute(w, nil)
}

func CalculateCustomTax(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req CustomTaxRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	salary := req.AnnualPay - req.Deduction
	tax := 0.0
	var buckets [][]float64
	var percents []float64

	for _, b := range req.Buckets {
		if b.To == 0 || b.To <= b.From {
			continue // Skip invalid bucket
		}
		if b.From == 0 {
			buckets = append(buckets, []float64{b.From, b.To})
		} else {
			buckets = append(buckets, []float64{b.From - 1, b.To})
		}

		percents = append(percents, b.Percent)
	}
	tax = GetTaxAmount(buckets, percents, salary)

	resp := CustomTaxResponse{
		Tax:    tax,
		Inhand: req.AnnualPay - tax,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

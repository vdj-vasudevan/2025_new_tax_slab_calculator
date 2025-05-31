package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type request struct {
	AnnualPay float64 `json:"annual_pay"`
}

type Result struct {
	AnnualPay float64
	Tax       float64
	Inhand    float64
	HasResult bool
}

func RenderForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))
	tmpl.Execute(w, Result{})
}

func CalculateTax(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	var req request
	// var annualPay float64
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	// if n, err := strconv.ParseFloat(req.AnnualPay, 64); err == nil {
	// 	annualPay = n
	// } else {
	// 	http.Error(w, "Invalid annual pay amount", http.StatusBadRequest)
	// 	return
	// }
	annualPay := req.AnnualPay

	bucket := [][]float64{{0.0, 400000.0}, {400000.0, 800000.0}, {800000.0, 1200000.0}, {1200000.0, 1600000.0}, {1600000.0, 2000000.0}, {2000000.0, 2400000.0}, {2400000.0, 0.0}}
	percent := []float64{0.0, 5.0, 10.0, 15.0, 20.0, 25.0, 30.0}

	sdSalary := annualPay - 75000
	var tax float64
	if sdSalary > 1200000 {
		tax = GetTaxAmount(bucket, percent, sdSalary)

		fmt.Printf("In Hand salary %f\n", (annualPay - tax))
		fmt.Printf("Tax : %f", tax)
	} else {

		fmt.Printf("In Hand salary %f\n", (annualPay - tax))
		fmt.Printf("Tax : %f", tax)
	}
	inhand := annualPay - tax
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{
		"tax":    tax,
		"inhand": inhand,
	})
}

func GetTaxAmount(bucket [][]float64, percent []float64, salary float64) float64 {
	tax := 0.0

	for i := 0; i <= len(bucket)-2; i++ {

		if (bucket[i][1] - bucket[i][0]) < salary {
			tax = tax + ((bucket[i][1] - bucket[i][0]) * percent[i] / 100)
			salary = salary - (bucket[i][1] - bucket[i][0])
		} else {
			tax = tax + (salary * percent[i] / 100)
			salary = 0
		}
	}
	tax = tax + (salary * percent[len(percent)-1] / 100)
	return tax
}

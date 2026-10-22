package domain

type Comision struct {
	Nivel      int     `json:"nivel"`
	Porcentaje float64 `json:"porcentaje"`
	Monto      float64 `json:"monto"`
}

// RepartirComisiones calcula las comisiones sin acceso a red, base de datos ni reloj.
func RepartirComisiones(total float64, tasas []float64) []Comision {
	if total <= 0 || len(tasas) == 0 {
		return []Comision{}
	}
	out := make([]Comision, 0, len(tasas))
	for index, tasa := range tasas {
		if tasa <= 0 {
			continue
		}
		out = append(out, Comision{Nivel: index + 1, Porcentaje: tasa, Monto: total * tasa / 100})
	}
	return out
}

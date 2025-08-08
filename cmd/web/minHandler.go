package web

import (
	"linn221/Requester/views"
	"net/http"
)

func HandleMin(v *views.MyTemplates, h func(v *views.MyTemplates, w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := h(v, w, r)
		if err != nil {
			w.Header().Add("HX-Reswap", "outerHTML")
			w.Header().Add("HX-Retarget", "#flash")
			v.ShowErrorBox(w, err.Error())
			return
		}
	}
}

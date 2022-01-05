package functions

import (
	"fmt"
	"math"
	"net/http"
	"quiz-3/utils"
	"runtime"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// segitiga-sama-sisi
// persegi
// persegi-panjang
// lingkaran
// jajar-genjang
func GetFunction(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	runtime.GOMAXPROCS(4)

	q := r.URL.Query()

	path := strings.Split(r.URL.Path, "/")

	var input = map[string]interface{}{}
	for i, val := range q {
		if i == "hitung" {
			input["hitung"] = val[0]
		} else {
			if path[len(path)-1] != "lingkaran" {
				convVal, _ := strconv.Atoi(val[0])
				input[i] = convVal
			} else {
				convVal, _ := strconv.ParseFloat(val[0], 64)
				input[i] = convVal
			}
		}
	}

	message := goChanMsg(path, input)

	utils.ResponseJSON(rw, message, http.StatusOK)
}

func goChanMsg(path []string, input map[string]interface{}) string {
	hitung := input["hitung"].(string)
	// channel
	var ch = make(chan int)
	var ch2 = make(chan float64)

	jenisBangun := path[len(path)-1]

	switch jenisBangun {
	case "segitiga-sama-sisi":
		go segitigaSamaSisi(input["alas"].(int), input["tinggi"].(int), hitung, ch)
	case "persegi":
		go persegi(input["sisi"].(int), hitung, ch)
	case "persegi-panjang":
		go persegiPanjang(input["panjang"].(int), input["lebar"].(int), hitung, ch)
	case "lingkaran":
		go lingkaran(input["jariJari"].(float64), hitung, ch2)
	case "jajar-genjang":
		go jajarGenjang(input["sisi"].(int), input["alas"].(int), input["tinggi"].(int), hitung, ch)
	}

	select {
	case nonlingkaran := <-ch:
		if hitung == "luas" {
			return fmt.Sprintf("Luas %s = %v", strings.Join(strings.Split(jenisBangun, "-"), " "), nonlingkaran)
		} else {
			return fmt.Sprintf("Keliling %s = %v", strings.Join(strings.Split(jenisBangun, "-"), " "), nonlingkaran)
		}
	case lingkaran := <-ch2:
		if hitung == "luas" {
			return fmt.Sprintf("Luas lingkaran = %.2f", lingkaran)
		} else {
			return fmt.Sprintf("Keliling lingkaran = %.2f", lingkaran)
		}
	}

}

func segitigaSamaSisi(alas, tinggi int, hitung string, ch chan int) {
	if hitung == "luas" {
		ch <- (alas * tinggi) / 2
	}
	ch <- alas + alas + tinggi
}
func persegi(sisi int, hitung string, ch chan int) {
	if hitung == "luas" {
		ch <- sisi * sisi
	}
	ch <- 4 * sisi
}
func persegiPanjang(panjang, lebar int, hitung string, ch chan int) {
	if hitung == "luas" {
		ch <- panjang * lebar
	}
	ch <- (2 * panjang) + (2 * lebar)
}
func lingkaran(jariJari float64, hitung string, ch chan float64) {
	if hitung == "luas" {
		ch <- math.Pi * math.Pow(jariJari, 2)
	}
	ch <- 2 * math.Pi * jariJari
}
func jajarGenjang(sisi, alas, tinggi int, hitung string, ch chan int) {
	if hitung == "luas" {
		ch <- alas * tinggi
	}
	ch <- (2 * alas) + (2 * sisi)
}

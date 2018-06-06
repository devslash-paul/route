package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/miekg/dns"

	"github.com/paulthom12345/route/data"
)

var brandRepository *data.BrandRepository

func main() {
	http.HandleFunc("/api/brand", apiBrand)
	http.HandleFunc("/", redirectHandler)

	_, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		log.Fatal("error making client from default file", err)
	}

	dns.HandleFunc("to.", func(w dns.ResponseWriter, r *dns.Msg) {
		log.Println("Serving request to test")
		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = true
		m.RecursionAvailable = true
		rr, err := dns.NewRR("to.	5	IN	CNAME	localhost")
		if err != nil {
			log.Fatal("unable to create RR ", err)
		}
		m.Answer = []dns.RR{rr}
		w.WriteMsg(m)

	})

	dns.HandleFunc(".", func(w dns.ResponseWriter, m *dns.Msg) {
		c := new(dns.Client)
		log.Println("Handling from", m.Question)
		r, _, err := c.Exchange(m, "1.1.1.1:53")
		if err != nil {
			log.Fatal(err)
		}
		w.WriteMsg(r)
	})

	go func() {
		srv := &dns.Server{Addr: ":53", Net: "udp"}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Failed to set up UDP listener", err)
		}
	}()

	db, err := Migrate()
	if err != nil {
		log.Fatal(err)
	}
	brandRepository = data.CreateBrandRepository(db)
	http.ListenAndServe(":8080", nil)
}

type RestError struct {
	Error string `json:"error"`
}

func serializeResponse(err error, result interface{}) []byte {
	if err != nil {
		body, _ := json.MarshalIndent(&RestError{
			Error: err.Error(),
		}, "", "    ")
		return body
	}
	body, _ := json.MarshalIndent(result, "", "    ")
	return body

}

func apiBrand(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		err := brandRepository.CreateBrand(req.PostFormValue("brand"))
		resp.Write(serializeResponse(err, "OK"))
	case "GET":
		res, err := brandRepository.GetBrands()
		resp.Write(serializeResponse(err, res))
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func redirectHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Location", "http://www.google.com")
	resp.WriteHeader(302)
}

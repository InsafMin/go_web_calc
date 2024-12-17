package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/InsafMin/go_web_calc/pkg/calculator"
	"io"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

//func (a *Application) Run() error {
//	for {
//		log.Println("input expression")
//		reader := bufio.NewReader(os.Stdin)
//		text, err := reader.ReadString('\n')
//		if err != nil {
//			log.Println("failed to read expression from console")
//		}
//
//		text = strings.TrimSpace(text)
//		if text == "exit" {
//			log.Println("application was successfully closed")
//			return nil
//		}
//
//		result, err := calculator.Calc(text)
//		if err != nil {
//			log.Println(text, "calculation failed wit error:", err)
//		} else {
//			log.Println(text, "=", result)
//		}
//	}
//}

type Request struct {
	Expression string `json:"expression"`
}

func isCalcError(err error) bool {
	return errors.Is(err, calculator.ErrInvalidExpression) ||
		errors.Is(err, calculator.ErrExtraOperator) ||
		errors.Is(err, calculator.ErrUnacceptableSymbol) ||
		errors.Is(err, calculator.ErrOperatorNotSupported) ||
		errors.Is(err, calculator.ErrExtraOpenBracket) ||
		errors.Is(err, calculator.ErrExtraCloseBracket) ||
		errors.Is(err, calculator.ErrDivisionByZero)
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "error: Unacceptable method", http.StatusUnprocessableEntity)
		log.Println("Error: Unacceptable method")
		return
	}

	request := new(Request)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Failed to close body:", err)
		}
	}(r.Body)

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "error with json", http.StatusInternalServerError)
		log.Println("Error with json:", err)
		return
	}

	result, err := calculator.Calc(request.Expression)
	if err != nil {
		if isCalcError(err) {
			http.Error(w, "error: "+err.Error(), http.StatusUnprocessableEntity)
			log.Println("Calculating error:", err)
		} else {
			http.Error(w, "unknown error occurred", http.StatusInternalServerError)
			log.Println("Unknown calculating error:", err)
		}
		return
	} else {
		_, err := fmt.Fprintf(w, "result: %f", result)
		if err != nil {
			http.Error(w, "unknown error occurred", http.StatusInternalServerError)
			log.Println("Failed to write response:", err)
		}
		log.Printf("Expression: %s --- result: %f\n", request.Expression, result)
	}
}

func routeHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/calculate" {
			http.Error(w, "404 page not found", http.StatusNotFound)
			log.Println("404 page not found")
			return
		}
		next.ServeHTTP(w, r)
	}
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				http.Error(w, "unknown error occurred", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *Application) RunServer() error {
	log.Println("starting server on port:", a.config.Addr)
	mux := http.NewServeMux()
	calc := http.HandlerFunc(CalcHandler)

	mux.Handle("/api/v1/calculate", routeHandler(calc))

	handler := PanicMiddleware(mux)
	err := http.ListenAndServe(":"+a.config.Addr, handler)
	if err != nil {
		return err
	}
	return nil
}

package application

import (
	"bytes"
	"github.com/InsafMin/go_web_calc/pkg/calculator"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFoundHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/invalid-path", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := routeHandler(CalcHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", status)
	}
}

func TestInvalidMethodHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/calculate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := routeHandler(CalcHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", status)
	}
}

func TestInvalidJSON(t *testing.T) {
	body := []byte(`invalid-json`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := routeHandler(CalcHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", status)
	}

	expected := "error with json\n"
	if rr.Body.String() != expected {
		t.Errorf("expected response %s, got %s", expected, rr.Body.String())
	}
}

func TestSuccessfulCalculation(t *testing.T) {
	body := []byte(`{"expression": "(2+2)*2"}`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := routeHandler(CalcHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status 200, got %d", status)
	}

	expected := "result: 8.000000"
	if rr.Body.String() != expected {
		t.Errorf("expected response %s, got %s", expected, rr.Body.String())
	}
}

func TestCalcHandler_DivisionByZero(t *testing.T) {
	body := []byte(`{"expression": "1/0"}`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := routeHandler(CalcHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", status)
	}

	expected := "error: " + calculator.ErrDivisionByZero.Error() + "\n"
	if rr.Body.String() != expected {
		t.Errorf("expected response %s, got %s", expected, rr.Body.String())
	}
}

func TestCalcHandler_ExtraOperator(t *testing.T) {
	body := []byte(`{"expression": "2++2"}`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := routeHandler(CalcHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", status)
	}

	expected := "error: " + calculator.ErrExtraOperator.Error() + "\n"
	if rr.Body.String() != expected {
		t.Errorf("expected response %s, got %s", expected, rr.Body.String())
	}
}

func TestCalcHandler_ExtraOpenBracket(t *testing.T) {
	body := []byte(`{"expression": "(2+3"}`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := routeHandler(CalcHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", status)
	}

	expected := "error: " + calculator.ErrExtraOpenBracket.Error() + "\n"
	if rr.Body.String() != expected {
		t.Errorf("expected response %s, got %s", expected, rr.Body.String())
	}
}

func TestCalcHandler_ExtraCloseBracket(t *testing.T) {
	body := []byte(`{"expression": "2+3)"}`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := routeHandler(CalcHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", status)
	}

	expected := "error: " + calculator.ErrExtraCloseBracket.Error() + "\n"
	if rr.Body.String() != expected {
		t.Errorf("expected response %s, got %s", expected, rr.Body.String())
	}
}

func TestCalcHandler_UnacceptableSymbol(t *testing.T) {
	body := []byte(`{"expression": "2@3"}`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := routeHandler(CalcHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", status)
	}

	expected := "error: " + calculator.ErrUnacceptableSymbol.Error() + "\n"
	if rr.Body.String() != expected {
		t.Errorf("expected response %s, got %s", expected, rr.Body.String())
	}
}

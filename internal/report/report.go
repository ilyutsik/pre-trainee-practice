package report

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
)

type Record struct {
	SKU        string
	Qty        int
	PriceCents int // unit price
}

type Exporter interface {
	Name() string
	Export(records []Record) ([]byte, error)
}

type JSONExporter struct{}

type CSVExporter struct {
	Comma rune // if 0, use ','
}

func (JSONExporter) Name() string {
	return "json"
}
func (JSONExporter) Export(records []Record) ([]byte, error) {
	err := checkRecords(records)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return []byte("[]"), nil
	}

	type jsonRecord struct {
		SKU        string `json:"sku"`
		Qty        int    `json:"qty"`
		PriceCents int    `json:"price_cents"`
		LineCents  int    `json:"line_cents"`
	}
	jsons := make([]jsonRecord, len(records))
	for i, record := range records {
		jsons[i] = jsonRecord{
			SKU:        record.SKU,
			Qty:        record.Qty,
			PriceCents: record.PriceCents,
			LineCents:  record.PriceCents * record.Qty}
	}

	out, err := json.Marshal(jsons)
	if err != nil {
		return nil, fmt.Errorf("error marshal json: %w", err)
	}
	return out, nil
}

func (CSVExporter) Name() string {
	return "csv"
}
func (c CSVExporter) Export(records []Record) ([]byte, error) {
	err := checkRecords(records)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	w := csv.NewWriter(&buffer)
	if c.Comma != 0 {
		w.Comma = c.Comma
	}

	line := []string{"sku", "qty", "price_cents", "line_cents"}
	err = w.Write(line)
	if err != nil {
		return nil, fmt.Errorf("csv writer error: %w", err)
	}

	for _, record := range records {
		line = []string{record.SKU, strconv.Itoa(record.Qty), strconv.Itoa(record.PriceCents), strconv.Itoa(record.PriceCents * record.Qty)}
		err = w.Write(line)
		if err != nil {
			return nil, fmt.Errorf("csv writer error: %w", err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("csv flush error: %w", err)
	}

	csvBytes := buffer.Bytes()
	return csvBytes, nil
}

func CommaOf(e Exporter) (rune, bool) {
	if exporter, ok := e.(CSVExporter); ok {
		if exporter.Comma == 0 {
			return ',', true
		}
		return exporter.Comma, true
	}
	return 0, false
}

// RenderAll runs every exporter. Keys are exporter names.
func RenderAll(records []Record, exporters []Exporter) (map[string][]byte, error) {

	out := make(map[string][]byte, len(exporters))
	for _, exporter := range exporters {
		name := exporter.Name()
		if _, exist := out[name]; exist {
			return nil, fmt.Errorf("report: duplicate exporter name %q", name)
		}
		export, err := exporter.Export(records)
		if err != nil {
			return nil, err
		}
		out[name] = export
	}
	return out, nil
}

func checkRecords(records []Record) error {
	for _, record := range records {
		if record.SKU == "" {
			return fmt.Errorf("validation error: missing squ")
		}
		if record.Qty < 0 {
			return fmt.Errorf("validation error: negative qty for SKU: %s", record.SKU)
		}
		if record.PriceCents < 0 {
			return fmt.Errorf("validation error: negative price cents for SKU: %s", record.SKU)
		}
	}
	return nil
}

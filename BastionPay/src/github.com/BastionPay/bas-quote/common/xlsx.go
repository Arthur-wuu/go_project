package common

import (
	"github.com/kataras/iris/core/errors"
	"github.com/tealeg/xlsx"
	"io"
	"reflect"
)

type Xlsx struct {
	file    *xlsx.File
	sheet   *xlsx.Sheet
	data    []interface{}
	rowNames []string
	filters map[string]func(interface{}) interface{}
}

func NewXlsx(data []interface{},rowNames []string, filters map[string]func(interface{}) interface{}) (*Xlsx, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")

	if err != nil {
		return nil, err
	}

	return &Xlsx{
		data:    data,
		file:    file,
		sheet:   sheet,
		filters: filters,
		rowNames: rowNames,
	}, nil
}

func (e *Xlsx) Generate() error {
	if len(e.data) == 0 {
		return errors.New("EMPTY_DATA")
	}

	e.header()
	e.body()

	return nil

}


func (e *Xlsx) File(path string) error {
	if len(e.data) == 0 {
		return errors.New("EMPTY_DATA")
	}

	if err := e.file.Save(path); err != nil {
		return err
	}

	return nil

}

func (e *Xlsx) IoWriter(w io.Writer) error {
	if len(e.data) == 0 {
		return errors.New("EMPTY_DATA")
	}

	if err := e.file.Write(w); err != nil {
		return err
	}

	return nil

}

func (e *Xlsx) header() {
	var (
		to    reflect.Type
		tag   string
		elem  reflect.Value
		count int
	)
	if len(e.rowNames) != 0 {
		row := e.sheet.AddRow()
		for i:=0; i < len(e.rowNames); i++ {
			row.AddCell().Value = e.rowNames[i]
		}
	}

	elem = reflect.ValueOf(e.data[0])

	switch elem.Kind() {
	case reflect.Struct:
		to = elem.Type()
		count = elem.NumField()
		row := e.sheet.AddRow()

		for i := 0; i < count; i++ {
			fieldName := ""
			tag = string(to.Field(i).Tag.Get("xlsx"))
			if tag == "" {
				fieldName = to.Field(i).Name
			} else if tag == "-" {
				continue
			} else {
				fieldName = tag
			}

			row.AddCell().Value = fieldName
		}
	case reflect.Ptr:
		to = elem.Elem().Type()
		count = elem.Elem().NumField()
		row := e.sheet.AddRow()

		for i := 0; i < count; i++ {
			fieldName := ""
			tag = string(to.Field(i).Tag.Get("xlsx"))
			if tag == "" {
				fieldName = to.Field(i).Name
			} else if tag == "-" {
				continue
			} else {
				fieldName = tag
			}

			row.AddCell().Value = fieldName
		}
	case reflect.Array, reflect.Slice:
		row := e.sheet.AddRow()
		for i := 0; i < elem.Len(); i++ {
			row.AddCell().SetValue(elem.Index(i))
		}
	}
}

func (e *Xlsx) body() {
	var (
		elem  reflect.Value
		count int
	)

	elem = reflect.ValueOf(e.data[0])

	switch elem.Kind() {
	case reflect.Struct:
		for _, v := range e.data {
			elem = reflect.ValueOf(v)

			count = elem.NumField()
			row := e.sheet.AddRow()
			for i := 0; i < count; i++ {
				tag := string(elem.Type().Field(i).Tag.Get("xlsx"))
				if tag == "-" {
					continue
				}

				if f, ok := e.filters[elem.Type().Field(i).Name]; ok {
					row.AddCell().SetValue(f(elem.Field(i)))
				} else {
					row.AddCell().SetValue(elem.Field(i))
				}
			}
		}
	case reflect.Ptr:
		for _, v := range e.data {
			elem = reflect.ValueOf(v).Elem()

			count = elem.NumField()
			row := e.sheet.AddRow()

			for i := 0; i < count; i++ {
				tag := string(elem.Type().Field(i).Tag.Get("xlsx"))
				if tag == "-" {
					continue
				}

				if f, ok := e.filters[elem.Type().Field(i).Name]; ok {
					row.AddCell().SetValue(f(elem.Field(i).Interface()))
				} else {
					row.AddCell().SetValue(elem.Field(i))
				}
			}
		}
	case reflect.Array, reflect.Slice:
		for _, v := range e.data {
			elem = reflect.ValueOf(v)
			row := e.sheet.AddRow()
			for i := 0; i < elem.Len(); i++ {
				row.AddCell().SetValue(elem.Index(i))
			}
		}
	}
}

package importer

import (
	"esTool/pkg/utils"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type FiledType string
const (
	FiledTypeTimePrefix = FiledType("@time")
	FiledTypeInt    = FiledType("@int")
	FiledTypeFloat  = FiledType("@float")
	FiledTypeString  = FiledType("@string")

)


type Field struct {
	Name string `validate:"required"`
	Type FiledType
	// for time
	format string
}
func (field *Field) isTime()bool{
	return strings.HasPrefix(string(field.Type), string(FiledTypeTimePrefix))
}
func (field *Field) isInt()bool{
	return field.Type == FiledTypeInt
}
func (field *Field) isFloat()bool{
	return field.Type == FiledTypeFloat
}
func (field *Field) isString()bool{
	return field.Type == FiledTypeString
}

func (field *Field)getFormat()string{
	if field.isTime(){
		if field.format == ""{
			tmps := strings.Split(string(field.Type), "@")
			if len(tmps) != 3{
				return ""
			}
			field.format =  tmps[2]
		}
		return  field.format
	}
	return ""
}

func (field *Field) Value(str string) (interface{}, error)  {
	if field.isString(){
		return str, nil
	}else if field.isInt(){
		return utils.S(str).Int64()
	}else if field.isFloat(){
		return utils.S(str).Float64()
	}else if field.isTime(){
		format := field.getFormat()
		if format == ""{
			return nil, fmt.Errorf("time format error type :%s", field.Type)
		}
		res, err :=   time.Parse(format, str)
		if err != nil{
			err = errors.Wrapf(err,"str:%s", str)
		}
		return res, err
	}
	return str, nil
}
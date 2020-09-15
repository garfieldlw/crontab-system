package valuate

import (
	"github.com/Knetic/govaluate"
	"github.com/shopspring/decimal"
)

func Valuate(exp string, param map[string]interface{}) (interface{}, error) {
	expression, err := govaluate.NewEvaluableExpression(exp)
	if err != nil {
		return nil, err
	}

	return expression.Evaluate(param)
}

func ValuateToInt64(exp string, param map[string]interface{}) (int64, error) {
	val, errValue := Valuate(exp, param)
	if errValue != nil {
		return 0, errValue
	}

	if val == nil {
		return 0, nil
	}

	return val.(int64), nil
}

func ValuateToFloat64(exp string, param map[string]interface{}) (float64, error) {
	val, errValue := Valuate(exp, param)
	if errValue != nil {
		return 0, errValue
	}

	if val == nil {
		return 0, nil
	}

	return val.(float64), nil
}

func ValuateToDecimal(exp string, param map[string]interface{}) (decimal.Decimal, error) {
	val, errValue := ValuateToFloat64(exp, param)
	if errValue != nil {
		return decimal.Zero, errValue
	}

	return decimal.NewFromFloat(val), nil
}

package gostats

type CorrelationType string

const CorrelationPearson CorrelationType = "Pearson"
const CorrelationSpearman CorrelationType = "SpearMan"
const CorrelationCustom CorrelationType = "Custom"
const CorrelationKendall CorrelationType = "Kendal"

func (ct CorrelationType) String() string {
	return "Correlation" + string(ct)
}

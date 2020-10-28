package crmgo

// LONG LIVE CONVENIENCE!

type Q map[string]interface{}

// GT is greater than
func (q Q) GT(key string, val interface{}) Q {
	q[key] = GT(val)
	return q
}

// GTE is greater than or equal
func (q Q) GTE(key string, val interface{}) Q {
	q[key] = GTE(val)
	return q
}

// LT is less than
func (q Q) LT(key string, val interface{}) Q {
	q[key] = LT(val)
	return q
}

// LTE is less than or equal
func (q Q) LTE(key string, val interface{}) Q {
	q[key] = LTE(val)
	return q
}

// GT is greater than
func GT(val interface{}) Q {
	return Q{"$gt": val}
}

// GTE is greater than or equal
func GTE(val interface{}) Q {
	return Q{"$gte": val}
}

// LT is less than
func LT(val interface{}) Q {
	return Q{"$lt": val}
}

// LTE is less than or equal
func LTE(val interface{}) Q {
	return Q{"$lte": val}
}

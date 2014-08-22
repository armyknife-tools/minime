package schema

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/mapstructure"
)

// getSource represents the level we want to get for a value (internally).
// Any source less than or equal to the level will be loaded (whichever
// has a value first).
type getSource byte

const (
	getSourceState getSource = iota
	getSourceConfig
	getSourceDiff
	getSourceSet
)

// getResult is the internal structure that is generated when a Get
// is called that contains some extra data that might be used.
type getResult struct {
	Value  interface{}
	Exists bool
	Schema *Schema
}

var getResultEmpty getResult

// ResourceData is used to query and set the attributes of a resource.
type ResourceData struct {
	schema  map[string]*Schema
	config  *terraform.ResourceConfig
	state   *terraform.ResourceState
	diff    *terraform.ResourceDiff
	diffing bool

	setMap   map[string]string
	newState *terraform.ResourceState
	once     sync.Once
}

// Get returns the data for the given key, or nil if the key doesn't exist
// in the schema.
//
// If the key does exist in the schema but doesn't exist in the configuration,
// then the default value for that type will be returned. For strings, this is
// "", for numbers it is 0, etc.
//
// If you also want to test if something is set at all, use GetOk.
func (d *ResourceData) Get(key string) interface{} {
	v, _ := d.GetOk(key)
	return v
}

// GetChange returns the old and new value for a given key.
//
// If there is no change, then old and new will simply be the same.
func (d *ResourceData) GetChange(key string) (interface{}, interface{}) {
	o, n := d.getChange(key, getSourceConfig, getSourceDiff)
	return o.Value, n.Value
}

// GetOk returns the data for the given key and whether or not the key
// existed or not in the configuration. The second boolean result will also
// be false if a key is given that isn't in the schema at all.
func (d *ResourceData) GetOk(key string) (interface{}, bool) {
	var parts []string
	if key != "" {
		parts = strings.Split(key, ".")
	}

	r := d.getObject("", parts, d.schema, getSourceSet)
	return r.Value, r.Exists
}

// HasChange returns whether or not the given key has been changed.
func (d *ResourceData) HasChange(key string) bool {
	o, n := d.GetChange(key)
	return !reflect.DeepEqual(o, n)
}

// Set sets the value for the given key.
//
// If the key is invalid or the value is not a correct type, an error
// will be returned.
func (d *ResourceData) Set(key string, value interface{}) error {
	if d.setMap == nil {
		d.setMap = make(map[string]string)
	}

	parts := strings.Split(key, ".")
	return d.setObject("", parts, d.schema, value)
}

// Id returns the ID of the resource.
func (d *ResourceData) Id() string {
	var result string

	if d.state != nil {
		result = d.state.ID
	}

	if d.newState != nil {
		result = d.newState.ID
	}

	return result
}

// ConnInfo returns the connection info for this resource.
func (d *ResourceData) ConnInfo() map[string]string {
	if d.newState != nil {
		return d.newState.ConnInfo
	}

	if d.state != nil {
		return d.state.ConnInfo
	}

	return nil
}

// Dependencies returns the dependencies in this state.
func (d *ResourceData) Dependencies() []terraform.ResourceDependency {
	if d.newState != nil {
		return d.newState.Dependencies
	}

	if d.state != nil {
		return d.state.Dependencies
	}

	return nil
}

// SetId sets the ID of the resource. If the value is blank, then the
// resource is destroyed.
func (d *ResourceData) SetId(v string) {
	d.once.Do(d.init)
	d.newState.ID = v
}

// SetConnInfo sets the connection info for a resource.
func (d *ResourceData) SetConnInfo(v map[string]string) {
	d.once.Do(d.init)
	d.newState.ConnInfo = v
}

// SetDependencies sets the dependencies of a resource.
func (d *ResourceData) SetDependencies(ds []terraform.ResourceDependency) {
	d.once.Do(d.init)
	d.newState.Dependencies = ds
}

// State returns the new ResourceState after the diff and any Set
// calls.
func (d *ResourceData) State() *terraform.ResourceState {
	var result terraform.ResourceState
	result.ID = d.Id()

	// If we have no ID, then this resource doesn't exist and we just
	// return nil.
	if result.ID == "" {
		return nil
	}

	result.Attributes = d.stateObject("", d.schema)
	result.ConnInfo = d.ConnInfo()
	result.Dependencies = d.Dependencies()

	if v := d.Id(); v != "" {
		result.Attributes["id"] = d.Id()
	}

	return &result
}

func (d *ResourceData) init() {
	var copyState terraform.ResourceState
	if d.state != nil {
		copyState = *d.state
	}

	d.newState = &copyState
}

func (d *ResourceData) diffChange(k string) (interface{}, interface{}, bool) {
	// Get the change between the state and the config.
	o, n := d.getChange(k, getSourceState, getSourceConfig)
	if !o.Exists {
		o.Value = nil
	}
	if !n.Exists {
		n.Value = nil
	}

	if n.Exists && n.Schema.StateFunc != nil {
		n.Value = n.Schema.StateFunc(n.Value)
	}

	// Return the old, new, and whether there is a change
	return o.Value, n.Value, !reflect.DeepEqual(o.Value, n.Value)
}

func (d *ResourceData) getChange(
	key string,
	oldLevel getSource,
	newLevel getSource) (getResult, getResult) {
	var parts, parts2 []string
	if key != "" {
		parts = strings.Split(key, ".")
		parts2 = strings.Split(key, ".")
	}

	o := d.getObject("", parts, d.schema, oldLevel)
	n := d.getObject("", parts2, d.schema, newLevel)
	return o, n
}

func (d *ResourceData) get(
	k string,
	parts []string,
	schema *Schema,
	source getSource) getResult {
	switch schema.Type {
	case TypeList:
		return d.getList(k, parts, schema, source)
	case TypeMap:
		return d.getMap(k, parts, schema, source)
	case TypeSet:
		return d.getSet(k, parts, schema, source)
	case TypeBool:
		fallthrough
	case TypeInt:
		fallthrough
	case TypeString:
		return d.getPrimitive(k, parts, schema, source)
	default:
		panic(fmt.Sprintf("%s: unknown type %s", k, schema.Type))
	}
}

func (d *ResourceData) getSet(
	k string,
	parts []string,
	schema *Schema,
	source getSource) getResult {
	s := &Set{F: schema.Set}
	result := getResult{Schema: schema, Value: s}
	raw := d.getList(k, nil, schema, source)
	if !raw.Exists {
		if len(parts) > 0 {
			return d.getList(k, parts, schema, source)
		}

		return result
	}

	list := raw.Value.([]interface{})
	if len(list) == 0 {
		if len(parts) > 0 {
			return d.getList(k, parts, schema, source)
		}

		return result
	}

	// This is a reverse map of hash code => index in config used to
	// resolve direct set item lookup for turning into state. Confused?
	// Read on...
	//
	// To create the state (the state* functions), a Get call is done
	// with a full key such as "ports.0". The index of a set ("0") doesn't
	// make a lot of sense, but we need to deterministically list out
	// elements of the set like this. Luckily, same sets have a deterministic
	// List() output, so we can use that to look things up.
	//
	// This mapping makes it so that we can look up the hash code of an
	// object back to its index in the REAL config.
	var indexMap map[int]int
	if len(parts) > 0 {
		indexMap = make(map[int]int)
	}

	// Build the set from all the items using the given hash code
	for i, v := range list {
		code := s.add(v)
		if indexMap != nil {
			indexMap[code] = i
		}
	}

	// If we're trying to get a specific element, then rewrite the
	// index to be just that, then jump direct to getList.
	if len(parts) > 0 {
		index := parts[0]
		indexInt, err := strconv.ParseInt(index, 0, 0)
		if err != nil {
			return getResultEmpty
		}

		codes := s.listCode()
		if int(indexInt) >= len(codes) {
			return getResultEmpty
		}
		code := codes[indexInt]
		realIndex := indexMap[code]

		parts[0] = strconv.FormatInt(int64(realIndex), 10)
		return d.getList(k, parts, schema, source)
	}

	result.Exists = true
	return result
}

func (d *ResourceData) getMap(
	k string,
	parts []string,
	schema *Schema,
	source getSource) getResult {
	elemSchema := &Schema{Type: TypeString}

	result := make(map[string]interface{})
	resultSet := false
	prefix := k + "."

	if d.state != nil && source >= getSourceState {
		for k, _ := range d.state.Attributes {
			if !strings.HasPrefix(k, prefix) {
				continue
			}

			single := k[len(prefix):]
			result[single] = d.getPrimitive(k, nil, elemSchema, source).Value
			resultSet = true
		}
	}

	if d.config != nil && source == getSourceConfig {
		// For config, we always set the result to exactly what was requested
		if m, ok := d.config.Get(k); ok {
			result = m.(map[string]interface{})
			resultSet = true
		} else {
			result = nil
		}
	}

	if d.diff != nil && source >= getSourceDiff {
		for k, v := range d.diff.Attributes {
			if !strings.HasPrefix(k, prefix) {
				continue
			}
			resultSet = true

			single := k[len(prefix):]

			if v.NewRemoved {
				delete(result, single)
			} else {
				result[single] = d.getPrimitive(k, nil, elemSchema, source).Value
			}
		}
	}

	if d.setMap != nil && source >= getSourceSet {
		cleared := false
		for k, _ := range d.setMap {
			if !strings.HasPrefix(k, prefix) {
				continue
			}
			resultSet = true

			if !cleared {
				// We clear the results if they are in the set map
				result = make(map[string]interface{})
				cleared = true
			}

			single := k[len(prefix):]
			result[single] = d.getPrimitive(k, nil, elemSchema, source).Value
		}
	}

	// If we're requesting a specific element, return that
	var resultValue interface{} = result
	if len(parts) > 0 {
		resultValue = result[parts[0]]
	}

	return getResult{
		Value:  resultValue,
		Exists: resultSet,
		Schema: schema,
	}
}

func (d *ResourceData) getObject(
	k string,
	parts []string,
	schema map[string]*Schema,
	source getSource) getResult {
	if len(parts) > 0 {
		// We're requesting a specific key in an object
		key := parts[0]
		parts = parts[1:]
		s, ok := schema[key]
		if !ok {
			return getResultEmpty
		}

		if k != "" {
			// If we're not at the root, then we need to append
			// the key to get the full key path.
			key = fmt.Sprintf("%s.%s", k, key)
		}

		return d.get(key, parts, s, source)
	}

	// Get the entire object
	result := make(map[string]interface{})
	for field, _ := range schema {
		result[field] = d.getObject(k, []string{field}, schema, source).Value
	}

	return getResult{
		Value:  result,
		Exists: true,
		Schema: &Schema{
			Elem: schema,
		},
	}
}

func (d *ResourceData) getList(
	k string,
	parts []string,
	schema *Schema,
	source getSource) getResult {
	if len(parts) > 0 {
		// We still have parts left over meaning we're accessing an
		// element of this list.
		idx := parts[0]
		parts = parts[1:]

		// Special case if we're accessing the count of the list
		if idx == "#" {
			schema := &Schema{Type: TypeInt}
			return d.get(k+".#", parts, schema, source)
		}

		key := fmt.Sprintf("%s.%s", k, idx)
		switch t := schema.Elem.(type) {
		case *Resource:
			return d.getObject(key, parts, t.Schema, source)
		case *Schema:
			return d.get(key, parts, t, source)
		}
	}

	// Get the entire list.
	count := d.getList(k, []string{"#"}, schema, source)
	result := make([]interface{}, count.Value.(int))
	for i, _ := range result {
		is := strconv.FormatInt(int64(i), 10)
		result[i] = d.getList(k, []string{is}, schema, source).Value
	}

	return getResult{
		Value:  result,
		Exists: count.Exists,
		Schema: schema,
	}
}

func (d *ResourceData) getPrimitive(
	k string,
	parts []string,
	schema *Schema,
	source getSource) getResult {
	var result string
	var resultSet bool
	if d.state != nil && source >= getSourceState {
		result, resultSet = d.state.Attributes[k]
	}

	if d.config != nil && source == getSourceConfig {
		// For config, we always return the exact value
		if v, ok := d.config.Get(k); ok {
			if err := mapstructure.WeakDecode(v, &result); err != nil {
				panic(err)
			}

			resultSet = true
		} else {
			result = ""
			resultSet = false
		}

	}

	if d.diff != nil && source >= getSourceDiff {
		attrD, ok := d.diff.Attributes[k]
		if ok && !attrD.NewComputed {
			result = attrD.New
			resultSet = true
		}
	}

	if d.setMap != nil && source >= getSourceSet {
		if v, ok := d.setMap[k]; ok {
			result = v
			resultSet = true
		}
	}

	if !resultSet {
		result = ""
	}

	var resultValue interface{}
	switch schema.Type {
	case TypeBool:
		if result == "" {
			resultValue = false
			break
		}

		v, err := strconv.ParseBool(result)
		if err != nil {
			panic(err)
		}

		resultValue = v
	case TypeString:
		// Use the value as-is. We just put this case here to be explicit.
		resultValue = result
	case TypeInt:
		if result == "" {
			resultValue = 0
			break
		}

		v, err := strconv.ParseInt(result, 0, 0)
		if err != nil {
			panic(err)
		}

		resultValue = int(v)
	default:
		panic(fmt.Sprintf("Unknown type: %s", schema.Type))
	}

	return getResult{
		Value:  resultValue,
		Exists: resultSet,
		Schema: schema,
	}
}

func (d *ResourceData) set(
	k string,
	parts []string,
	schema *Schema,
	value interface{}) error {
	switch schema.Type {
	case TypeList:
		return d.setList(k, parts, schema, value)
	case TypeMap:
		return d.setMapValue(k, parts, schema, value)
	case TypeSet:
		return d.setSet(k, parts, schema, value)
	case TypeBool:
		fallthrough
	case TypeInt:
		fallthrough
	case TypeString:
		return d.setPrimitive(k, schema, value)
	default:
		panic(fmt.Sprintf("%s: unknown type %s", k, schema.Type))
	}
}

func (d *ResourceData) setList(
	k string,
	parts []string,
	schema *Schema,
	value interface{}) error {
	if len(parts) > 0 {
		// We're setting a specific element
		idx := parts[0]
		parts = parts[1:]

		// Special case if we're accessing the count of the list
		if idx == "#" {
			return fmt.Errorf("%s: can't set count of list", k)
		}

		key := fmt.Sprintf("%s.%s", k, idx)
		switch t := schema.Elem.(type) {
		case *Resource:
			return d.setObject(key, parts, t.Schema, value)
		case *Schema:
			return d.set(key, parts, t, value)
		}
	}

	var vs []interface{}
	if err := mapstructure.Decode(value, &vs); err != nil {
		return fmt.Errorf("%s: %s", k, err)
	}

	// Set the entire list.
	var err error
	for i, elem := range vs {
		is := strconv.FormatInt(int64(i), 10)
		err = d.setList(k, []string{is}, schema, elem)
		if err != nil {
			break
		}
	}
	if err != nil {
		for i, _ := range vs {
			is := strconv.FormatInt(int64(i), 10)
			d.setList(k, []string{is}, schema, nil)
		}

		return err
	}

	d.setMap[k+".#"] = strconv.FormatInt(int64(len(vs)), 10)
	return nil
}

func (d *ResourceData) setMapValue(
	k string,
	parts []string,
	schema *Schema,
	value interface{}) error {
	elemSchema := &Schema{Type: TypeString}
	if len(parts) > 0 {
		return fmt.Errorf("%s: full map must be set, no a single element", k)
	}

	// Delete any prior map set
	/*
		v := d.getMap(k, nil, schema, getSourceSet)
		for subKey, _ := range v.(map[string]interface{}) {
			delete(d.setMap, fmt.Sprintf("%s.%s", k, subKey))
		}
	*/

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Map {
		return fmt.Errorf("%s: must be a map", k)
	}
	if v.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("%s: keys must strings", k)
	}
	vs := make(map[string]interface{})
	for _, mk := range v.MapKeys() {
		mv := v.MapIndex(mk)
		vs[mk.String()] = mv.Interface()
	}

	for subKey, v := range vs {
		err := d.set(fmt.Sprintf("%s.%s", k, subKey), nil, elemSchema, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *ResourceData) setObject(
	k string,
	parts []string,
	schema map[string]*Schema,
	value interface{}) error {
	if len(parts) > 0 {
		// We're setting a specific key in an object
		key := parts[0]
		parts = parts[1:]

		s, ok := schema[key]
		if !ok {
			return fmt.Errorf("%s (internal): unknown key to set: %s", k, key)
		}

		if k != "" {
			// If we're not at the root, then we need to append
			// the key to get the full key path.
			key = fmt.Sprintf("%s.%s", k, key)
		}

		return d.set(key, parts, s, value)
	}

	// Set the entire object. First decode into a proper structure
	var v map[string]interface{}
	if err := mapstructure.Decode(value, &v); err != nil {
		return fmt.Errorf("%s: %s", k, err)
	}

	// Set each element in turn
	var err error
	for k1, v1 := range v {
		err = d.setObject(k, []string{k1}, schema, v1)
		if err != nil {
			break
		}
	}
	if err != nil {
		for k1, _ := range v {
			d.setObject(k, []string{k1}, schema, nil)
		}
	}

	return err
}

func (d *ResourceData) setPrimitive(
	k string,
	schema *Schema,
	v interface{}) error {
	if v == nil {
		delete(d.setMap, k)
		return nil
	}

	var set string
	switch schema.Type {
	case TypeBool:
		var b bool
		if err := mapstructure.Decode(v, &b); err != nil {
			return fmt.Errorf("%s: %s", k, err)
		}

		set = strconv.FormatBool(b)
	case TypeString:
		if err := mapstructure.Decode(v, &set); err != nil {
			return fmt.Errorf("%s: %s", k, err)
		}
	case TypeInt:
		var n int
		if err := mapstructure.Decode(v, &n); err != nil {
			return fmt.Errorf("%s: %s", k, err)
		}

		set = strconv.FormatInt(int64(n), 10)
	default:
		return fmt.Errorf("Unknown type: %s", schema.Type)
	}

	d.setMap[k] = set
	return nil
}

func (d *ResourceData) setSet(
	k string,
	parts []string,
	schema *Schema,
	value interface{}) error {
	if len(parts) > 0 {
		return fmt.Errorf("%s: can only set the full set, not elements", k)
	}

	if s, ok := value.(*Set); ok {
		value = s.List()
	}

	return d.setList(k, nil, schema, value)
}

func (d *ResourceData) stateList(
	prefix string,
	schema *Schema) map[string]string {
	countRaw := d.get(prefix, []string{"#"}, schema, getSourceSet)
	if !countRaw.Exists {
		return nil
	}
	count := countRaw.Value.(int)

	result := make(map[string]string)
	if count > 0 {
		result[prefix+".#"] = strconv.FormatInt(int64(count), 10)
	}
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%s.%d", prefix, i)

		var m map[string]string
		switch t := schema.Elem.(type) {
		case *Resource:
			m = d.stateObject(key, t.Schema)
		case *Schema:
			m = d.stateSingle(key, t)
		}

		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

func (d *ResourceData) stateMap(
	prefix string,
	schema *Schema) map[string]string {
	v := d.getMap(prefix, nil, schema, getSourceSet)
	if !v.Exists {
		return nil
	}

	elemSchema := &Schema{Type: TypeString}
	result := make(map[string]string)
	for mk, _ := range v.Value.(map[string]interface{}) {
		mp := fmt.Sprintf("%s.%s", prefix, mk)
		for k, v := range d.stateSingle(mp, elemSchema) {
			result[k] = v
		}
	}

	return result
}

func (d *ResourceData) stateObject(
	prefix string,
	schema map[string]*Schema) map[string]string {
	result := make(map[string]string)
	for k, v := range schema {
		key := k
		if prefix != "" {
			key = prefix + "." + key
		}

		for k1, v1 := range d.stateSingle(key, v) {
			result[k1] = v1
		}
	}

	return result
}

func (d *ResourceData) statePrimitive(
	prefix string,
	schema *Schema) map[string]string {
	v := d.Get(prefix)
	if v == nil {
		return nil
	}

	var vs string
	switch schema.Type {
	case TypeBool:
		vs = strconv.FormatBool(v.(bool))
	case TypeString:
		vs = v.(string)
	case TypeInt:
		vs = strconv.FormatInt(int64(v.(int)), 10)
	default:
		panic(fmt.Sprintf("Unknown type: %s", schema.Type))
	}

	return map[string]string{
		prefix: vs,
	}
}

func (d *ResourceData) stateSet(
	prefix string,
	schema *Schema) map[string]string {
	raw := d.get(prefix, nil, schema, getSourceSet)
	if !raw.Exists {
		return nil
	}

	set := raw.Value.(*Set)
	list := set.List()
	result := make(map[string]string)
	result[prefix+".#"] = strconv.FormatInt(int64(len(list)), 10)
	for i := 0; i < len(list); i++ {
		key := fmt.Sprintf("%s.%d", prefix, i)

		var m map[string]string
		switch t := schema.Elem.(type) {
		case *Resource:
			m = d.stateObject(key, t.Schema)
		case *Schema:
			m = d.stateSingle(key, t)
		}

		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

func (d *ResourceData) stateSingle(
	prefix string,
	schema *Schema) map[string]string {
	switch schema.Type {
	case TypeList:
		return d.stateList(prefix, schema)
	case TypeMap:
		return d.stateMap(prefix, schema)
	case TypeSet:
		return d.stateSet(prefix, schema)
	case TypeBool:
		fallthrough
	case TypeInt:
		fallthrough
	case TypeString:
		return d.statePrimitive(prefix, schema)
	default:
		panic(fmt.Sprintf("%s: unknown type %s", prefix, schema.Type))
	}
}

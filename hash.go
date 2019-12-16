package hash

import (
	"fmt"
	"reflect"
	"unsafe"
)

//nolint:gocyclo,funlen
func Interface(i interface{}, seed uintptr) uintptr {
	switch k := i.(type) {
	case Hashable:
		return k.Hash()
	case bool:
		return Bool(k, seed)
	case int:
		return Int(k, seed)
	case int8:
		return Int8(k, seed)
	case int16:
		return Int16(k, seed)
	case int32:
		return Int32(k, seed)
	case int64:
		return Int64(k, seed)
	case uint:
		return Uint(k, seed)
	case uint8:
		return Uint8(k, seed)
	case uint16:
		return Uint16(k, seed)
	case uint32:
		return Uint32(k, seed)
	case uint64:
		return Uint64(k, seed)
	case uintptr:
		return Uintptr(k, seed)
	case float32:
		return Float32(k, seed)
	case float64:
		return Float64(k, seed)
	case complex64:
		return Complex64(k, seed)
	case complex128:
		return Complex128(k, seed)
	case string:
		return String(k, seed)
	case []interface{}:
		return sliceInterfaceHash(k, seed)
	default:
		return Value(reflect.ValueOf(k), seed)
	}
}

func Bool(b bool, seed uintptr) uintptr {
	return algarray[algMEM8](noescape(unsafe.Pointer(&b)), seed)
}

func Int(x int, seed uintptr) uintptr {
	return algarray[algINT](noescape(unsafe.Pointer(&x)), seed)
}

func Int8(x int8, seed uintptr) uintptr {
	return algarray[algMEM8](noescape(unsafe.Pointer(&x)), seed)
}

func Int16(x int16, seed uintptr) uintptr {
	return algarray[algMEM16](noescape(unsafe.Pointer(&x)), seed)
}

func Int32(x int32, seed uintptr) uintptr {
	return algarray[algMEM32](noescape(unsafe.Pointer(&x)), seed)
}

func Int64(x int64, seed uintptr) uintptr {
	return algarray[algMEM64](noescape(unsafe.Pointer(&x)), seed)
}

func Uint(x uint, seed uintptr) uintptr {
	return algarray[algUINT](noescape(unsafe.Pointer(&x)), seed)
}

func Uint8(x uint8, seed uintptr) uintptr {
	return algarray[algMEM8](noescape(unsafe.Pointer(&x)), seed)
}

func Uint16(x uint16, seed uintptr) uintptr {
	return algarray[algMEM16](noescape(unsafe.Pointer(&x)), seed)
}

func Uint32(x uint32, seed uintptr) uintptr {
	return algarray[algMEM32](noescape(unsafe.Pointer(&x)), seed)
}

func Uint64(x uint64, seed uintptr) uintptr {
	return algarray[algMEM64](noescape(unsafe.Pointer(&x)), seed)
}

func Uintptr(x, seed uintptr) uintptr {
	return algarray[algPTR](noescape(unsafe.Pointer(&x)), seed)
}

func Float32(f float32, seed uintptr) uintptr {
	return algarray[algFLOAT32](noescape(unsafe.Pointer(&f)), seed)
}

func Float64(f float64, seed uintptr) uintptr {
	return algarray[algFLOAT64](noescape(unsafe.Pointer(&f)), seed)
}

func Complex64(c complex64, seed uintptr) uintptr {
	return algarray[algCPLX64](noescape(unsafe.Pointer(&c)), seed)
}

func Complex128(c complex128, seed uintptr) uintptr {
	return algarray[algCPLX128](noescape(unsafe.Pointer(&c)), seed)
}

func String(s string, seed uintptr) uintptr {
	return algarray[algSTRING](noescape(unsafe.Pointer(&s)), seed)
}

func UnsafePointer(p unsafe.Pointer, seed uintptr) uintptr {
	return Uintptr(uintptr(p), seed)
}

//nolint:funlen
func Value(v reflect.Value, seed uintptr) uintptr {
	switch v.Kind() {
	case reflect.Bool:
		return Bool(v.Bool(), seed)
	case reflect.Int:
		return Int(int(v.Int()), seed)
	case reflect.Int8:
		return Int8(int8(v.Int()), seed)
	case reflect.Int16:
		return Int16(int16(v.Int()), seed)
	case reflect.Int32:
		return Int32(int32(v.Int()), seed)
	case reflect.Int64:
		return Int64(v.Int(), seed)
	case reflect.Uint:
		return Uint(uint(v.Uint()), seed)
	case reflect.Uint8:
		return Uint8(uint8(v.Uint()), seed)
	case reflect.Uint16:
		return Uint16(uint16(v.Uint()), seed)
	case reflect.Uint32:
		return Uint32(uint32(v.Uint()), seed)
	case reflect.Uint64:
		return Uint64(v.Uint(), seed)
	case reflect.Uintptr:
		return Uintptr(uintptr(v.Uint()), seed)
	case reflect.UnsafePointer:
		return UnsafePointer(unsafe.Pointer(v.Pointer()), seed)
	case reflect.Float32:
		return Float32(float32(v.Float()), seed)
	case reflect.Float64:
		return Float64(v.Float(), seed)
	case reflect.Complex64:
		return Complex64(complex64(v.Complex()), seed)
	case reflect.Complex128:
		return Complex128(v.Complex(), seed)
	case reflect.String:
		return String(v.String(), seed)
	case reflect.Struct:
		return structHash(v, seed)
	case reflect.Array:
		return arrayHash(v, seed)
	case reflect.Ptr:
		return Uintptr(v.Pointer(), seed)
	}
	panic(fmt.Sprintf("value %v has unhashable type %v", v, v.Type()))
}

func structHash(v reflect.Value, seed uintptr) uintptr {
	t := v.Type()
	h := seed
	for i := v.NumField(); i > 0; {
		i--
		h = String(t.Field(i).Name, h)
		h = Interface(v.Field(i).Interface(), h)
	}
	return h
}

func arrayHash(v reflect.Value, seed uintptr) uintptr {
	h := seed
	for i := v.Len(); i > 0; {
		i--
		h = Value(v.Index(i), h)
	}
	return h
}

func sliceInterfaceHash(slice []interface{}, seed uintptr) uintptr {
	h := seed
	for _, elem := range slice {
		h = Interface(elem, h)
	}
	return h
}

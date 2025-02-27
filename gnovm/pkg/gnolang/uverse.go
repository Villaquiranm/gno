package gnolang

import (
	"fmt"
	"strings"

	bm "github.com/gnolang/gno/gnovm/pkg/benchops"
)

// ----------------------------------------
// non-primitive builtin types

var gErrorType = &DeclaredType{
	PkgPath: uversePkgPath,
	Name:    "error",
	Base: &InterfaceType{
		PkgPath: uversePkgPath,
		Methods: []FieldType{
			{
				Name: "Error",
				Type: &FuncType{
					Params: nil,
					Results: []FieldType{
						{
							// Name: "",
							Type: StringType,
						},
					},
				},
			},
		},
	},
	sealed: true,
}

var gStringerType = &DeclaredType{
	PkgPath: uversePkgPath,
	Name:    "stringer",
	Base: &InterfaceType{
		PkgPath: uversePkgPath,
		Methods: []FieldType{
			{
				Name: "String",
				Type: &FuncType{
					Params: nil,
					Results: []FieldType{
						{
							// Name: "",
							Type: StringType,
						},
					},
				},
			},
		},
	},
	sealed: true,
}

// ----------------------------------------
// Uverse package

var (
	uverseNode  *PackageNode
	uverseValue *PackageValue
	uverseInit  = uverseUninitialized
)

const (
	uverseUninitialized = iota
	uverseInitializing
	uverseInitialized
)

func init() {
	// Skip Uverse init during benchmarking to load stdlibs in the benchmark main function.
	if !(bm.OpsEnabled || bm.StorageEnabled) {
		// Call Uverse() so we initialize the Uverse node ahead of any calls to the package.
		Uverse()
	}
}

const uversePkgPath = ".uverse"

// UverseNode returns the uverse PackageValue.
// If called while initializing the UverseNode itself, it will return an empty
// PackageValue.
func Uverse() *PackageValue {
	switch uverseInit {
	case uverseUninitialized:
		uverseInit = uverseInitializing
		makeUverseNode()
		uverseInit = uverseInitialized
	case uverseInitializing:
		return &PackageValue{}
	}

	return uverseValue
}

// UverseNode returns the uverse PackageNode.
// If called while initializing the UverseNode itself, it will return an empty
// PackageNode.
func UverseNode() *PackageNode {
	switch uverseInit {
	case uverseUninitialized:
		uverseInit = uverseInitializing
		makeUverseNode()
		uverseInit = uverseInitialized
	case uverseInitializing:
		return &PackageNode{}
	}

	return uverseNode
}

func makeUverseNode() {
	// NOTE: uverse node is hidden, thus the leading dot in pkgPath=".uverse".
	uverseNode = NewPackageNode("uverse", uversePkgPath, nil)

	// temporary convenience functions.
	def := func(n Name, tv TypedValue) {
		uverseNode.Define(n, tv)
	}
	defNative := uverseNode.DefineNative

	// Primitive types
	undefined := TypedValue{}
	def("._", undefined)   // special, path is zero.
	def("iota", undefined) // special
	def("nil", undefined)
	def("bool", asValue(BoolType))
	def("byte", asValue(Uint8Type))
	def("float32", asValue(Float32Type))
	def("float64", asValue(Float64Type))
	def("int", asValue(IntType))
	def("int8", asValue(Int8Type))
	def("int16", asValue(Int16Type))
	def("int32", asValue(Int32Type))
	def("int64", asValue(Int64Type))
	def("rune", asValue(Int32Type))
	def("string", asValue(StringType))
	def("uint", asValue(UintType))
	def("uint8", asValue(Uint8Type))
	def("uint16", asValue(Uint16Type))
	def("uint32", asValue(Uint32Type))
	def("uint64", asValue(Uint64Type))
	def("any", asValue(&InterfaceType{PkgPath: uversePkgPath}))
	// NOTE on 'typeval': We can't call the type of a TypeValue a
	// "type", even though we want to, because it conflicts with
	// the pre-existing syntax for type-switching, `switch
	// x.(type) {case SomeType:...}`, for if x.(type) were not a
	// type-switch but a type-assertion, and the resulting value
	// could be any type, such as an IntType; whereas as the .X of
	// a SwitchStmt, the type of an IntType value is not IntType
	// but always a TypeType (all types are of type TypeType).
	//
	// The ideal solution is to keep the syntax consistent for
	// type-assertions, but for backwards compatibility, the
	// keyword that represents the TypeType type is not "type" but
	// "typeval".  The value of a "typeval" value is represented
	// by a TypeValue.
	def("typeval", asValue(gTypeType))
	def("error", asValue(gErrorType))
	def("any", asValue(&InterfaceType{}))

	// Values
	def("true", untypedBool(true))
	def("false", untypedBool(false))

	// Functions
	defNative("append",
		Flds( // params
			"x", GenT("X", nil), // args[0]
			"args", MaybeNativeT(Vrd(GenT("X.Elem()", nil))), // args[1]
		),
		Flds( // results
			"res", GenT("X", nil), // res
		),
		func(m *Machine) {
			arg0, arg1 := m.LastBlock().GetParams2()
			// As a special case, if arg1 is a string type, first convert it into
			// a data slice type.
			if arg1.TV.T != nil && arg1.TV.T.Kind() == StringKind {
				arg1String := arg1.TV.GetString()
				// NOTE: this hack works because
				// arg1 PointerValue is not a pointer,
				// so the modification here is only local.
				newArrayValue := m.Alloc.NewDataArray(len(arg1String))
				copy(newArrayValue.Data, []byte(arg1String))
				arg1.TV = &TypedValue{
					T: m.Alloc.NewType(&SliceType{ // TODO: reuse
						Elt: Uint8Type,
						Vrd: true,
					}),
					V: m.Alloc.NewSlice(newArrayValue, 0, len(arg1String), len(arg1String)), // TODO: pool?
				}
			}
			arg0Type := arg0.TV.T
			arg1Type := arg1.TV.T
			switch arg0Value := arg0.TV.V.(type) {
			// ----------------------------------------------------------------
			// append(nil, ???)
			case nil:
				switch arg1Value := arg1.TV.V.(type) {
				// ------------------------------------------------------------
				// append(nil, nil)
				case nil: // no change
					m.PushValue(TypedValue{
						T: arg0Type,
						V: nil,
					})
					return

				// ------------------------------------------------------------
				// append(nil, *SliceValue)
				case *SliceValue:
					arg1Length := arg1Value.Length
					arg1Offset := arg1Value.Offset
					arg1Base := arg1Value.GetBase(m.Store)
					arg1EndIndex := arg1Offset + arg1Length

					if arg1Length == 0 { // no change
						m.PushValue(TypedValue{
							T: arg0Type,
							V: nil,
						})
						return
					} else if arg0Type.Elem().Kind() == Uint8Kind {
						// append(nil, *SliceValue) new data bytes ---
						arrayValue := m.Alloc.NewDataArray(arg1Length)
						if arg1Base.Data == nil {
							copyListToData(
								arrayValue.Data[:arg1Length],
								arg1Base.List[arg1Offset:arg1EndIndex])
						} else {
							copy(
								arrayValue.Data[:arg1Length],
								arg1Base.Data[arg1Offset:arg1EndIndex])
						}
						m.PushValue(TypedValue{
							T: arg0Type,
							V: m.Alloc.NewSlice(arrayValue, 0, arg1Length, arg1Length),
						})
						return
					} else {
						// append(nil, *SliceValue) new list ---------
						arrayValue := m.Alloc.NewListArray(arg1Length)
						if arg1Length > 0 {
							for i := 0; i < arg1Length; i++ {
								arrayValue.List[i] = arg1Base.List[arg1Offset+i].unrefCopy(m.Alloc, m.Store)
							}
						}
						m.PushValue(TypedValue{
							T: arg0Type,
							V: m.Alloc.NewSlice(arrayValue, 0, arg1Length, arg1Length),
						})
						return
					}
				default:
					panic("should not happen")
				}

			// ----------------------------------------------------------------
			// append(*SliceValue, ???)
			case *SliceValue:
				arg0Length := arg0Value.Length
				arg0Offset := arg0Value.Offset
				arg0Capacity := arg0Value.Maxcap
				arg0Base := arg0Value.GetBase(m.Store)
				switch arg1Value := arg1.TV.V.(type) {
				// ------------------------------------------------------------
				// append(*SliceValue, nil)
				case nil: // no change
					m.PushValue(TypedValue{
						T: arg0Type,
						V: arg0Value,
					})
					return

				// ------------------------------------------------------------
				// append(*SliceValue, *SliceValue)
				case *SliceValue:
					arg1Length := arg1Value.Length
					arg1Offset := arg1Value.Offset
					arg1Base := arg1Value.GetBase(m.Store)
					if arg0Length+arg1Length <= arg0Capacity {
						// append(*SliceValue, *SliceValue) w/i capacity -----
						if 0 < arg1Length { // implies 0 < xvc
							if arg0Base.Data == nil {
								// append(*SliceValue.List, *SliceValue) ---------
								list := arg0Base.List
								if arg1Base.Data == nil {
									for i := 0; i < arg1Length; i++ {
										oldElem := list[arg0Offset+arg0Length+i]
										// unrefCopy will resolve references and copy their values
										// to copy by value rather than by reference.
										newElem := arg1Base.List[arg1Offset+i].unrefCopy(m.Alloc, m.Store)
										list[arg0Offset+arg0Length+i] = newElem

										m.Realm.DidUpdate(
											arg0Base,
											oldElem.GetFirstObject(m.Store),
											newElem.GetFirstObject(m.Store),
										)
									}
								} else {
									copyDataToList(
										list[arg0Offset+arg0Length:arg0Offset+arg0Length+arg1Length],
										arg1Base.Data[arg1Offset:arg1Offset+arg1Length],
										arg0Type.Elem())
									m.Realm.DidUpdate(arg1Base, nil, nil)
								}
							} else {
								// append(*SliceValue.Data, *SliceValue) ---------
								data := arg0Base.Data
								if arg1Base.Data == nil {
									copyListToData(
										data[arg0Offset+arg0Length:arg0Offset+arg0Length+arg1Length],
										arg1Base.List[arg1Offset:arg1Offset+arg1Length])
									m.Realm.DidUpdate(arg0Base, nil, nil)
								} else {
									copy(
										data[arg0Offset+arg0Length:arg0Offset+arg0Length+arg1Length],
										arg1Base.Data[arg1Offset:arg1Offset+arg1Length])
								}
							}
							m.PushValue(TypedValue{
								T: arg0Type,
								V: m.Alloc.NewSlice(arg0Base, arg0Offset, arg0Length+arg1Length, arg0Capacity),
							})
							return
						} else { // no change
							m.PushValue(TypedValue{
								T: arg0Type,
								V: arg0Value,
							})
							return
						}
					} else if arg0Type.Elem().Kind() == Uint8Kind {
						// append(*SliceValue, *SliceValue) new data bytes ---
						newLength := arg0Length + arg1Length
						arrayValue := m.Alloc.NewDataArray(newLength)
						if 0 < arg0Length {
							if arg0Base.Data == nil {
								copyListToData(
									arrayValue.Data[:arg0Length],
									arg0Base.List[arg0Offset:arg0Offset+arg0Length])
							} else {
								copy(
									arrayValue.Data[:arg0Length],
									arg0Base.Data[arg0Offset:arg0Offset+arg0Length])
							}
						}
						if 0 < arg1Length {
							if arg1Base.Data == nil {
								copyListToData(
									arrayValue.Data[arg0Length:newLength],
									arg1Base.List[arg1Offset:arg1Offset+arg1Length])
							} else {
								copy(
									arrayValue.Data[arg0Length:newLength],
									arg1Base.Data[arg1Offset:arg1Offset+arg1Length])
							}
						}
						m.PushValue(TypedValue{
							T: arg0Type,
							V: m.Alloc.NewSlice(arrayValue, 0, newLength, newLength),
						})
						return
					} else {
						// append(*SliceValue, *SliceValue) new list ---------
						arrayLen := arg0Length + arg1Length
						arrayValue := m.Alloc.NewListArray(arrayLen)
						if arg0Length > 0 {
							if arg0Base.Data == nil {
								for i := 0; i < arg0Length; i++ {
									arrayValue.List[i] = arg0Base.List[arg0Offset+i].unrefCopy(m.Alloc, m.Store)
								}
							} else {
								panic("should not happen")
							}
						}

						if arg1Length > 0 {
							if arg1Base.Data == nil {
								for i := 0; i < arg1Length; i++ {
									arrayValue.List[arg0Length+i] = arg1Base.List[arg1Offset+i].unrefCopy(m.Alloc, m.Store)
								}
							} else {
								copyDataToList(
									arrayValue.List[arg0Length:arg0Length+arg1Length],
									arg1Base.Data[arg1Offset:arg1Offset+arg1Length],
									arg1Type.Elem(),
								)
							}
						}
						m.PushValue(TypedValue{
							T: arg0Type,
							V: m.Alloc.NewSlice(arrayValue, 0, arrayLen, arrayLen),
						})
						return
					}
				// ------------------------------------------------------------
				default:
					panic("should not happen")
				}
			// ----------------------------------------------------------------
			// append(?!!, ???)
			default:
				panic("should not happen")
			}
		},
	)
	defNative("cap",
		Flds( // params
			"x", AnyT(),
		),
		Flds( // results
			"", "int",
		),
		func(m *Machine) {
			arg0 := m.LastBlock().GetParams1()
			res0 := TypedValue{
				T: IntType,
				V: nil,
			}
			res0.SetInt(int64(arg0.TV.GetCapacity()))
			m.PushValue(res0)
			return
		},
	)
	defNative("copy",
		Flds( // params
			"dst", GenT("X", nil),
			"src", GenT("Y", nil),
		),
		Flds( // results
			"", "int",
		),
		func(m *Machine) {
			arg0, arg1 := m.LastBlock().GetParams2()
			dst, src := arg0, arg1
			switch bdt := baseOf(dst.TV.T).(type) {
			case *SliceType:
				switch bst := baseOf(src.TV.T).(type) {
				case PrimitiveType:
					if debug {
						debug.Println("copy(<%s>,<%s>)", bdt.String(), bst.String())
					}
					if bst.Kind() != StringKind {
						panic("should not happen")
					}
					if bdt.Elt != Uint8Type {
						panic("should not happen")
					}
					// NOTE: this implementation is almost identical to the next one.
					// note that in some cases optimization
					// is possible if dstv.Data != nil.
					dstl := dst.TV.GetLength()
					srcl := src.TV.GetLength()
					minl := dstl
					if srcl < dstl {
						minl = srcl
					}
					if minl == 0 {
						// return 0.
						m.PushValue(defaultTypedValue(m.Alloc, IntType))
						return
					}
					dstv := dst.TV.V.(*SliceValue)
					// TODO: consider an optimization if dstv.Data != nil.
					for i := 0; i < minl; i++ {
						dstev := dstv.GetPointerAtIndexInt2(m.Store, i, bdt.Elt)
						srcev := src.TV.GetPointerAtIndexInt(m.Store, i)
						dstev.Assign2(m.Alloc, m.Store, m.Realm, srcev.Deref(), false)
					}
					res0 := TypedValue{
						T: IntType,
						V: nil,
					}
					res0.SetInt(int64(minl))
					m.PushValue(res0)
					return
				case *SliceType:
					dstl := dst.TV.GetLength()
					srcl := src.TV.GetLength()
					minl := dstl
					if srcl < dstl {
						minl = srcl
					}
					if minl == 0 {
						// return 0.
						m.PushValue(defaultTypedValue(m.Alloc, IntType))
						return
					}
					dstv := dst.TV.V.(*SliceValue)
					srcv := src.TV.V.(*SliceValue)
					for i := 0; i < minl; i++ {
						dstev := dstv.GetPointerAtIndexInt2(m.Store, i, bdt.Elt)
						srcev := srcv.GetPointerAtIndexInt2(m.Store, i, bst.Elt)
						dstev.Assign2(m.Alloc, m.Store, m.Realm, srcev.Deref(), false)
					}
					res0 := TypedValue{
						T: IntType,
						V: nil,
					}
					res0.SetInt(int64(minl))
					m.PushValue(res0)
					return
				default:
					panic("should not happen")
				}
			default:
				panic("should not happen")
			}
		},
	)
	defNative("delete",
		Flds( // params
			"m", MapT(GenT("K", nil), GenT("V", nil)), // map type
			"k", GenT("K", nil), // map key
		),
		nil, // results
		func(m *Machine) {
			arg0, arg1 := m.LastBlock().GetParams2()
			itv := arg1.Deref()
			switch baseOf(arg0.TV.T).(type) {
			case *MapType:
				mv := arg0.TV.V.(*MapValue)
				val, ok := mv.GetValueForKey(m.Store, &itv)
				if !ok {
					return
				}

				// delete
				mv.DeleteForKey(m.Store, &itv)

				if m.Realm != nil {
					// mark key as deleted
					keyObj := itv.GetFirstObject(m.Store)
					m.Realm.DidUpdate(mv, keyObj, nil)

					// mark value as deleted
					valObj := val.GetFirstObject(m.Store)
					m.Realm.DidUpdate(mv, valObj, nil)
				}

				return
			default:
				panic(fmt.Sprintf(
					"unexpected map type %s",
					arg0.TV.T.String()))
			}
		},
	)
	defNative("len",
		Flds( // params
			"x", AnyT(),
		),
		Flds( // results
			"", "int",
		),
		func(m *Machine) {
			arg0 := m.LastBlock().GetParams1()
			res0 := TypedValue{
				T: IntType,
				V: nil,
			}
			res0.SetInt(int64(arg0.TV.GetLength()))
			m.PushValue(res0)
			return
		},
	)
	defNative("make",
		Flds( // params
			"t", GenT("T.(type)", nil),
			"z", Vrd(AnyT()),
		),
		Flds( // results
			"", GenT("T", nil),
		),
		func(m *Machine) {
			arg0, arg1 := m.LastBlock().GetParams2()
			vargs := arg1
			vargsl := vargs.TV.GetLength()
			tt := arg0.TV.GetType()
			switch bt := baseOf(tt).(type) {
			case *SliceType:
				et := bt.Elem()
				if vargsl == 1 {
					lv := vargs.TV.GetPointerAtIndexInt(m.Store, 0).Deref()
					li := int(lv.ConvertGetInt())
					if et.Kind() == Uint8Kind {
						arrayValue := m.Alloc.NewDataArray(li)
						m.PushValue(TypedValue{
							T: tt,
							V: m.Alloc.NewSlice(arrayValue, 0, li, li),
						})
						return
					} else {
						arrayValue := m.Alloc.NewListArray(li)
						if et.Kind() == InterfaceKind {
							// leave as is
						} else {
							// init zero elements with concrete type.
							for i := 0; i < li; i++ {
								arrayValue.List[i] = defaultTypedValue(m.Alloc, et)
							}
						}
						m.PushValue(TypedValue{
							T: tt,
							V: m.Alloc.NewSlice(arrayValue, 0, li, li),
						})
						return
					}
				} else if vargsl == 2 {
					lv := vargs.TV.GetPointerAtIndexInt(m.Store, 0).Deref()
					li := int(lv.ConvertGetInt())
					cv := vargs.TV.GetPointerAtIndexInt(m.Store, 1).Deref()
					ci := int(cv.ConvertGetInt())

					if ci < li {
						panic(&Exception{Value: typedString(`makeslice: cap out of range`)})
					}

					if et.Kind() == Uint8Kind {
						arrayValue := m.Alloc.NewDataArray(ci)
						m.PushValue(TypedValue{
							T: tt,
							V: m.Alloc.NewSlice(arrayValue, 0, li, ci),
						})
						return
					} else {
						arrayValue := m.Alloc.NewListArray(ci)
						if et := bt.Elem(); et.Kind() == InterfaceKind {
							// leave as is
						} else {
							// Initialize all elements within capacity with default
							// type values. These need to be initialized because future
							// slice operations could get messy otherwise. Simple capacity
							// expansions like `a = a[:cap(a)]` would make it trivial to
							// initialize zero values at the time of the slice operation.
							// But sequences of operations like:
							// 		a := make([]int, 1, 10)
							// 		a = a[7:cap(a)]
							// 		a = a[3:5]
							//
							// require a bit more work to handle correctly, requiring that
							// all new TypedValue slice elements be checked to ensure they have
							// a value for every slice operation, which is not desirable.
							for i := 0; i < ci; i++ {
								arrayValue.List[i] = defaultTypedValue(m.Alloc, et)
							}
						}
						m.PushValue(TypedValue{
							T: tt,
							V: m.Alloc.NewSlice(arrayValue, 0, li, ci),
						})
						return
					}
				} else {
					panic("make() of slice type takes 2 or 3 arguments")
				}
			case *MapType:
				// NOTE: the type is not used.
				if vargsl == 0 {
					m.PushValue(TypedValue{
						T: tt,
						V: m.Alloc.NewMap(0),
					})
					return
				} else if vargsl == 1 {
					lv := vargs.TV.GetPointerAtIndexInt(m.Store, 0).Deref()
					li := int(lv.ConvertGetInt())
					m.PushValue(TypedValue{
						T: tt,
						V: m.Alloc.NewMap(li),
					})
					return
				} else {
					panic("make() of map type takes 1 or 2 arguments")
				}
			case *ChanType:
				if vargsl == 0 {
					panic("not yet implemented")
				} else if vargsl == 1 {
					panic("not yet implemented")
				} else {
					panic("make() of chan type takes 1 or 2 arguments")
				}
			default:
				panic(fmt.Sprintf(
					"cannot make type %s kind %v",
					tt.String(), tt.Kind()))
			}
		},
	)
	defNative("new",
		Flds( // params
			"t", GenT("T.(type)", nil),
		),
		Flds( // results
			"", GenT("*T", nil),
		),
		func(m *Machine) {
			arg0 := m.LastBlock().GetParams1()
			tt := arg0.TV.GetType()
			vv := defaultValue(m.Alloc, tt)
			m.Alloc.AllocatePointer()
			hi := m.Alloc.NewHeapItem(TypedValue{
				T: tt,
				V: vv,
			})
			m.PushValue(TypedValue{
				T: m.Alloc.NewType(&PointerType{
					Elt: tt,
				}),
				V: PointerValue{
					TV:    &hi.Value,
					Base:  hi,
					Index: 0,
				},
			})
			return
		},
	)
	// NOTE: panic is its own statement type, and is not defined as a function.
	defNative("print",
		Flds( // params
			"xs", Vrd(AnyT()), // args[0]
		),
		nil, // results
		func(m *Machine) {
			arg0 := m.LastBlock().GetParams1()
			xv := arg0
			xvl := xv.TV.GetLength()
			ss := make([]string, xvl)
			for i := 0; i < xvl; i++ {
				ev := xv.TV.GetPointerAtIndexInt(m.Store, i).Deref()
				ss[i] = ev.Sprint(m)
			}
			rs := strings.Join(ss, " ")
			if debug {
				print(rs)
			}
			m.Output.Write([]byte(rs))
		},
	)
	defNative("println",
		Flds( // param
			"xs", Vrd(AnyT()), // args[0]
		),
		nil, // results
		func(m *Machine) {
			arg0 := m.LastBlock().GetParams1()
			xv := arg0
			xvl := xv.TV.GetLength()
			ss := make([]string, xvl)
			for i := 0; i < xvl; i++ {
				ev := xv.TV.GetPointerAtIndexInt(m.Store, i).Deref()
				ss[i] = ev.Sprint(m)
			}
			rs := strings.Join(ss, " ") + "\n"
			if debug {
				println("DEBUG/stdout: " + rs)
			}
			m.Output.Write([]byte(rs))
		},
	)
	defNative("println",
		Flds( // param
			"xs", Vrd(AnyT()), // args[0]
		),
		nil, // results
		func(m *Machine) {
			arg0 := m.LastBlock().GetParams1()
			xv := arg0
			xvl := xv.TV.GetLength()
			ss := make([]string, xvl)
			for i := 0; i < xvl; i++ {
				ev := xv.TV.GetPointerAtIndexInt(m.Store, i).Deref()
				ss[i] = ev.Sprint(m)
			}
			rs := strings.Join(ss, " ") + "\n"
			if debug {
				println("DEBUG/stdout: " + rs)
			}
			m.Output.Write([]byte(rs))
		},
	)
	defNative("recover",
		nil, // params
		Flds( // results
			"exception", AnyT(),
		),
		func(m *Machine) {
			exception := m.Recover()
			if exception == nil {
				m.PushValue(TypedValue{})
			} else {
				m.PushValue(exception.Value)
			}
		},
	)
	uverseValue = uverseNode.NewPackage()
}

func copyDataToList(dst []TypedValue, data []byte, et Type) {
	for i := 0; i < len(data); i++ {
		dst[i] = TypedValue{T: et}
		dst[i].SetUint8(data[i])
	}
}

func copyListToData(dst []byte, tvs []TypedValue) {
	for i := 0; i < len(tvs); i++ {
		dst[i] = tvs[i].GetUint8()
	}
}

func copyListToRunes(dst []rune, tvs []TypedValue) {
	for i := 0; i < len(tvs); i++ {
		dst[i] = tvs[i].GetInt32()
	}
}

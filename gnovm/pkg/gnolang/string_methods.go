// Code generated by "stringer -type=Kind,Op,TransCtrl,TransField,VPType,Word -output string_methods.go ."; DO NOT EDIT.

package gnolang

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[InvalidKind-0]
	_ = x[BoolKind-1]
	_ = x[StringKind-2]
	_ = x[IntKind-3]
	_ = x[Int8Kind-4]
	_ = x[Int16Kind-5]
	_ = x[Int32Kind-6]
	_ = x[Int64Kind-7]
	_ = x[UintKind-8]
	_ = x[Uint8Kind-9]
	_ = x[Uint16Kind-10]
	_ = x[Uint32Kind-11]
	_ = x[Uint64Kind-12]
	_ = x[Float32Kind-13]
	_ = x[Float64Kind-14]
	_ = x[BigintKind-15]
	_ = x[BigdecKind-16]
	_ = x[ArrayKind-17]
	_ = x[SliceKind-18]
	_ = x[PointerKind-19]
	_ = x[StructKind-20]
	_ = x[PackageKind-21]
	_ = x[InterfaceKind-22]
	_ = x[ChanKind-23]
	_ = x[FuncKind-24]
	_ = x[MapKind-25]
	_ = x[TypeKind-26]
	_ = x[BlockKind-27]
	_ = x[HeapItemKind-28]
	_ = x[TupleKind-29]
	_ = x[RefTypeKind-30]
}

const _Kind_name = "InvalidKindBoolKindStringKindIntKindInt8KindInt16KindInt32KindInt64KindUintKindUint8KindUint16KindUint32KindUint64KindFloat32KindFloat64KindBigintKindBigdecKindArrayKindSliceKindPointerKindStructKindPackageKindInterfaceKindChanKindFuncKindMapKindTypeKindBlockKindHeapItemKindTupleKindRefTypeKind"

var _Kind_index = [...]uint16{0, 11, 19, 29, 36, 44, 53, 62, 71, 79, 88, 98, 108, 118, 129, 140, 150, 160, 169, 178, 189, 199, 210, 223, 231, 239, 246, 254, 263, 275, 284, 295}

func (i Kind) String() string {
	if i >= Kind(len(_Kind_index)-1) {
		return "Kind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpInvalid-0]
	_ = x[OpHalt-1]
	_ = x[OpNoop-2]
	_ = x[OpExec-3]
	_ = x[OpPrecall-4]
	_ = x[OpCall-5]
	_ = x[OpCallNativeBody-6]
	_ = x[OpReturn-7]
	_ = x[OpReturnFromBlock-8]
	_ = x[OpReturnToBlock-9]
	_ = x[OpDefer-10]
	_ = x[OpCallDeferNativeBody-11]
	_ = x[OpGo-12]
	_ = x[OpSelect-13]
	_ = x[OpSwitchClause-14]
	_ = x[OpSwitchClauseCase-15]
	_ = x[OpTypeSwitch-16]
	_ = x[OpIfCond-17]
	_ = x[OpPopValue-18]
	_ = x[OpPopResults-19]
	_ = x[OpPopBlock-20]
	_ = x[OpPopFrameAndReset-21]
	_ = x[OpPanic1-22]
	_ = x[OpPanic2-23]
	_ = x[OpUpos-32]
	_ = x[OpUneg-33]
	_ = x[OpUnot-34]
	_ = x[OpUxor-35]
	_ = x[OpUrecv-37]
	_ = x[OpLor-38]
	_ = x[OpLand-39]
	_ = x[OpEql-40]
	_ = x[OpNeq-41]
	_ = x[OpLss-42]
	_ = x[OpLeq-43]
	_ = x[OpGtr-44]
	_ = x[OpGeq-45]
	_ = x[OpAdd-46]
	_ = x[OpSub-47]
	_ = x[OpBor-48]
	_ = x[OpXor-49]
	_ = x[OpMul-50]
	_ = x[OpQuo-51]
	_ = x[OpRem-52]
	_ = x[OpShl-53]
	_ = x[OpShr-54]
	_ = x[OpBand-55]
	_ = x[OpBandn-56]
	_ = x[OpEval-64]
	_ = x[OpBinary1-65]
	_ = x[OpIndex1-66]
	_ = x[OpIndex2-67]
	_ = x[OpSelector-68]
	_ = x[OpSlice-69]
	_ = x[OpStar-70]
	_ = x[OpRef-71]
	_ = x[OpTypeAssert1-72]
	_ = x[OpTypeAssert2-73]
	_ = x[OpStaticTypeOf-74]
	_ = x[OpCompositeLit-75]
	_ = x[OpArrayLit-76]
	_ = x[OpSliceLit-77]
	_ = x[OpSliceLit2-78]
	_ = x[OpMapLit-79]
	_ = x[OpStructLit-80]
	_ = x[OpFuncLit-81]
	_ = x[OpConvert-82]
	_ = x[OpFieldType-112]
	_ = x[OpArrayType-113]
	_ = x[OpSliceType-114]
	_ = x[OpPointerType-115]
	_ = x[OpInterfaceType-116]
	_ = x[OpChanType-117]
	_ = x[OpFuncType-118]
	_ = x[OpMapType-119]
	_ = x[OpStructType-120]
	_ = x[OpAssign-128]
	_ = x[OpAddAssign-129]
	_ = x[OpSubAssign-130]
	_ = x[OpMulAssign-131]
	_ = x[OpQuoAssign-132]
	_ = x[OpRemAssign-133]
	_ = x[OpBandAssign-134]
	_ = x[OpBandnAssign-135]
	_ = x[OpBorAssign-136]
	_ = x[OpXorAssign-137]
	_ = x[OpShlAssign-138]
	_ = x[OpShrAssign-139]
	_ = x[OpDefine-140]
	_ = x[OpInc-141]
	_ = x[OpDec-142]
	_ = x[OpValueDecl-144]
	_ = x[OpTypeDecl-145]
	_ = x[OpSticky-208]
	_ = x[OpBody-209]
	_ = x[OpForLoop-210]
	_ = x[OpRangeIter-211]
	_ = x[OpRangeIterString-212]
	_ = x[OpRangeIterMap-213]
	_ = x[OpRangeIterArrayPtr-214]
	_ = x[OpReturnCallDefers-215]
	_ = x[OpVoid-255]
}

const (
	_Op_name_0 = "OpInvalidOpHaltOpNoopOpExecOpPrecallOpCallOpCallNativeBodyOpReturnOpReturnFromBlockOpReturnToBlockOpDeferOpCallDeferNativeBodyOpGoOpSelectOpSwitchClauseOpSwitchClauseCaseOpTypeSwitchOpIfCondOpPopValueOpPopResultsOpPopBlockOpPopFrameAndResetOpPanic1OpPanic2"
	_Op_name_1 = "OpUposOpUnegOpUnotOpUxor"
	_Op_name_2 = "OpUrecvOpLorOpLandOpEqlOpNeqOpLssOpLeqOpGtrOpGeqOpAddOpSubOpBorOpXorOpMulOpQuoOpRemOpShlOpShrOpBandOpBandn"
	_Op_name_3 = "OpEvalOpBinary1OpIndex1OpIndex2OpSelectorOpSliceOpStarOpRefOpTypeAssert1OpTypeAssert2OpStaticTypeOfOpCompositeLitOpArrayLitOpSliceLitOpSliceLit2OpMapLitOpStructLitOpFuncLitOpConvert"
	_Op_name_4 = "OpFieldTypeOpArrayTypeOpSliceTypeOpPointerTypeOpInterfaceTypeOpChanTypeOpFuncTypeOpMapTypeOpStructType"
	_Op_name_5 = "OpAssignOpAddAssignOpSubAssignOpMulAssignOpQuoAssignOpRemAssignOpBandAssignOpBandnAssignOpBorAssignOpXorAssignOpShlAssignOpShrAssignOpDefineOpIncOpDec"
	_Op_name_6 = "OpValueDeclOpTypeDecl"
	_Op_name_7 = "OpStickyOpBodyOpForLoopOpRangeIterOpRangeIterStringOpRangeIterMapOpRangeIterArrayPtrOpReturnCallDefers"
	_Op_name_8 = "OpVoid"
)

var (
	_Op_index_0 = [...]uint16{0, 9, 15, 21, 27, 36, 42, 58, 66, 83, 98, 105, 126, 130, 138, 152, 170, 182, 190, 200, 212, 222, 240, 248, 256}
	_Op_index_1 = [...]uint8{0, 6, 12, 18, 24}
	_Op_index_2 = [...]uint8{0, 7, 12, 18, 23, 28, 33, 38, 43, 48, 53, 58, 63, 68, 73, 78, 83, 88, 93, 99, 106}
	_Op_index_3 = [...]uint8{0, 6, 15, 23, 31, 41, 48, 54, 59, 72, 85, 99, 113, 123, 133, 144, 152, 163, 172, 181}
	_Op_index_4 = [...]uint8{0, 11, 22, 33, 46, 61, 71, 81, 90, 102}
	_Op_index_5 = [...]uint8{0, 8, 19, 30, 41, 52, 63, 75, 88, 99, 110, 121, 132, 140, 145, 150}
	_Op_index_6 = [...]uint8{0, 11, 21}
	_Op_index_7 = [...]uint8{0, 8, 14, 23, 34, 51, 65, 84, 102}
)

func (i Op) String() string {
	switch {
	case i <= 23:
		return _Op_name_0[_Op_index_0[i]:_Op_index_0[i+1]]
	case 32 <= i && i <= 35:
		i -= 32
		return _Op_name_1[_Op_index_1[i]:_Op_index_1[i+1]]
	case 37 <= i && i <= 56:
		i -= 37
		return _Op_name_2[_Op_index_2[i]:_Op_index_2[i+1]]
	case 64 <= i && i <= 82:
		i -= 64
		return _Op_name_3[_Op_index_3[i]:_Op_index_3[i+1]]
	case 112 <= i && i <= 120:
		i -= 112
		return _Op_name_4[_Op_index_4[i]:_Op_index_4[i+1]]
	case 128 <= i && i <= 142:
		i -= 128
		return _Op_name_5[_Op_index_5[i]:_Op_index_5[i+1]]
	case 144 <= i && i <= 145:
		i -= 144
		return _Op_name_6[_Op_index_6[i]:_Op_index_6[i+1]]
	case 208 <= i && i <= 215:
		i -= 208
		return _Op_name_7[_Op_index_7[i]:_Op_index_7[i+1]]
	case i == 255:
		return _Op_name_8
	default:
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TRANS_CONTINUE-0]
	_ = x[TRANS_SKIP-1]
	_ = x[TRANS_EXIT-2]
}

const _TransCtrl_name = "TRANS_CONTINUETRANS_SKIPTRANS_EXIT"

var _TransCtrl_index = [...]uint8{0, 14, 24, 34}

func (i TransCtrl) String() string {
	if i >= TransCtrl(len(_TransCtrl_index)-1) {
		return "TransCtrl(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TransCtrl_name[_TransCtrl_index[i]:_TransCtrl_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TRANS_ROOT-0]
	_ = x[TRANS_BINARY_LEFT-1]
	_ = x[TRANS_BINARY_RIGHT-2]
	_ = x[TRANS_CALL_FUNC-3]
	_ = x[TRANS_CALL_ARG-4]
	_ = x[TRANS_INDEX_X-5]
	_ = x[TRANS_INDEX_INDEX-6]
	_ = x[TRANS_SELECTOR_X-7]
	_ = x[TRANS_SLICE_X-8]
	_ = x[TRANS_SLICE_LOW-9]
	_ = x[TRANS_SLICE_HIGH-10]
	_ = x[TRANS_SLICE_MAX-11]
	_ = x[TRANS_STAR_X-12]
	_ = x[TRANS_REF_X-13]
	_ = x[TRANS_TYPEASSERT_X-14]
	_ = x[TRANS_TYPEASSERT_TYPE-15]
	_ = x[TRANS_UNARY_X-16]
	_ = x[TRANS_COMPOSITE_TYPE-17]
	_ = x[TRANS_COMPOSITE_KEY-18]
	_ = x[TRANS_COMPOSITE_VALUE-19]
	_ = x[TRANS_FUNCLIT_TYPE-20]
	_ = x[TRANS_FUNCLIT_HEAP_CAPTURE-21]
	_ = x[TRANS_FUNCLIT_BODY-22]
	_ = x[TRANS_FIELDTYPE_TYPE-23]
	_ = x[TRANS_FIELDTYPE_TAG-24]
	_ = x[TRANS_ARRAYTYPE_LEN-25]
	_ = x[TRANS_ARRAYTYPE_ELT-26]
	_ = x[TRANS_SLICETYPE_ELT-27]
	_ = x[TRANS_INTERFACETYPE_METHOD-28]
	_ = x[TRANS_CHANTYPE_VALUE-29]
	_ = x[TRANS_FUNCTYPE_PARAM-30]
	_ = x[TRANS_FUNCTYPE_RESULT-31]
	_ = x[TRANS_MAPTYPE_KEY-32]
	_ = x[TRANS_MAPTYPE_VALUE-33]
	_ = x[TRANS_STRUCTTYPE_FIELD-34]
	_ = x[TRANS_MAYBENATIVETYPE_TYPE-35]
	_ = x[TRANS_ASSIGN_LHS-36]
	_ = x[TRANS_ASSIGN_RHS-37]
	_ = x[TRANS_BLOCK_BODY-38]
	_ = x[TRANS_DECL_BODY-39]
	_ = x[TRANS_DEFER_CALL-40]
	_ = x[TRANS_EXPR_X-41]
	_ = x[TRANS_FOR_INIT-42]
	_ = x[TRANS_FOR_COND-43]
	_ = x[TRANS_FOR_POST-44]
	_ = x[TRANS_FOR_BODY-45]
	_ = x[TRANS_GO_CALL-46]
	_ = x[TRANS_IF_INIT-47]
	_ = x[TRANS_IF_COND-48]
	_ = x[TRANS_IF_BODY-49]
	_ = x[TRANS_IF_ELSE-50]
	_ = x[TRANS_IF_CASE_BODY-51]
	_ = x[TRANS_INCDEC_X-52]
	_ = x[TRANS_RANGE_X-53]
	_ = x[TRANS_RANGE_KEY-54]
	_ = x[TRANS_RANGE_VALUE-55]
	_ = x[TRANS_RANGE_BODY-56]
	_ = x[TRANS_RETURN_RESULT-57]
	_ = x[TRANS_PANIC_EXCEPTION-58]
	_ = x[TRANS_SELECT_CASE-59]
	_ = x[TRANS_SELECTCASE_COMM-60]
	_ = x[TRANS_SELECTCASE_BODY-61]
	_ = x[TRANS_SEND_CHAN-62]
	_ = x[TRANS_SEND_VALUE-63]
	_ = x[TRANS_SWITCH_INIT-64]
	_ = x[TRANS_SWITCH_X-65]
	_ = x[TRANS_SWITCH_CASE-66]
	_ = x[TRANS_SWITCHCASE_CASE-67]
	_ = x[TRANS_SWITCHCASE_BODY-68]
	_ = x[TRANS_FUNC_RECV-69]
	_ = x[TRANS_FUNC_TYPE-70]
	_ = x[TRANS_FUNC_BODY-71]
	_ = x[TRANS_IMPORT_PATH-72]
	_ = x[TRANS_CONST_TYPE-73]
	_ = x[TRANS_CONST_VALUE-74]
	_ = x[TRANS_VAR_NAME-75]
	_ = x[TRANS_VAR_TYPE-76]
	_ = x[TRANS_VAR_VALUE-77]
	_ = x[TRANS_TYPE_TYPE-78]
	_ = x[TRANS_FILE_BODY-79]
}

const _TransField_name = "TRANS_ROOTTRANS_BINARY_LEFTTRANS_BINARY_RIGHTTRANS_CALL_FUNCTRANS_CALL_ARGTRANS_INDEX_XTRANS_INDEX_INDEXTRANS_SELECTOR_XTRANS_SLICE_XTRANS_SLICE_LOWTRANS_SLICE_HIGHTRANS_SLICE_MAXTRANS_STAR_XTRANS_REF_XTRANS_TYPEASSERT_XTRANS_TYPEASSERT_TYPETRANS_UNARY_XTRANS_COMPOSITE_TYPETRANS_COMPOSITE_KEYTRANS_COMPOSITE_VALUETRANS_FUNCLIT_TYPETRANS_FUNCLIT_HEAP_CAPTURETRANS_FUNCLIT_BODYTRANS_FIELDTYPE_TYPETRANS_FIELDTYPE_TAGTRANS_ARRAYTYPE_LENTRANS_ARRAYTYPE_ELTTRANS_SLICETYPE_ELTTRANS_INTERFACETYPE_METHODTRANS_CHANTYPE_VALUETRANS_FUNCTYPE_PARAMTRANS_FUNCTYPE_RESULTTRANS_MAPTYPE_KEYTRANS_MAPTYPE_VALUETRANS_STRUCTTYPE_FIELDTRANS_MAYBENATIVETYPE_TYPETRANS_ASSIGN_LHSTRANS_ASSIGN_RHSTRANS_BLOCK_BODYTRANS_DECL_BODYTRANS_DEFER_CALLTRANS_EXPR_XTRANS_FOR_INITTRANS_FOR_CONDTRANS_FOR_POSTTRANS_FOR_BODYTRANS_GO_CALLTRANS_IF_INITTRANS_IF_CONDTRANS_IF_BODYTRANS_IF_ELSETRANS_IF_CASE_BODYTRANS_INCDEC_XTRANS_RANGE_XTRANS_RANGE_KEYTRANS_RANGE_VALUETRANS_RANGE_BODYTRANS_RETURN_RESULTTRANS_PANIC_EXCEPTIONTRANS_SELECT_CASETRANS_SELECTCASE_COMMTRANS_SELECTCASE_BODYTRANS_SEND_CHANTRANS_SEND_VALUETRANS_SWITCH_INITTRANS_SWITCH_XTRANS_SWITCH_CASETRANS_SWITCHCASE_CASETRANS_SWITCHCASE_BODYTRANS_FUNC_RECVTRANS_FUNC_TYPETRANS_FUNC_BODYTRANS_IMPORT_PATHTRANS_CONST_TYPETRANS_CONST_VALUETRANS_VAR_NAMETRANS_VAR_TYPETRANS_VAR_VALUETRANS_TYPE_TYPETRANS_FILE_BODY"

var _TransField_index = [...]uint16{0, 10, 27, 45, 60, 74, 87, 104, 120, 133, 148, 164, 179, 191, 202, 220, 241, 254, 274, 293, 314, 332, 358, 376, 396, 415, 434, 453, 472, 498, 518, 538, 559, 576, 595, 617, 643, 659, 675, 691, 706, 722, 734, 748, 762, 776, 790, 803, 816, 829, 842, 855, 873, 887, 900, 915, 932, 948, 967, 988, 1005, 1026, 1047, 1062, 1078, 1095, 1109, 1126, 1147, 1168, 1183, 1198, 1213, 1230, 1246, 1263, 1277, 1291, 1306, 1321, 1336}

func (i TransField) String() string {
	if i >= TransField(len(_TransField_index)-1) {
		return "TransField(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TransField_name[_TransField_index[i]:_TransField_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[VPUverse-0]
	_ = x[VPBlock-1]
	_ = x[VPField-2]
	_ = x[VPValMethod-3]
	_ = x[VPPtrMethod-4]
	_ = x[VPInterface-5]
	_ = x[VPSubrefField-6]
	_ = x[VPDerefField-18]
	_ = x[VPDerefValMethod-19]
	_ = x[VPDerefPtrMethod-20]
	_ = x[VPDerefInterface-21]
	_ = x[VPNative-32]
}

const (
	_VPType_name_0 = "VPUverseVPBlockVPFieldVPValMethodVPPtrMethodVPInterfaceVPSubrefField"
	_VPType_name_1 = "VPDerefFieldVPDerefValMethodVPDerefPtrMethodVPDerefInterface"
	_VPType_name_2 = "VPNative"
)

var (
	_VPType_index_0 = [...]uint8{0, 8, 15, 22, 33, 44, 55, 68}
	_VPType_index_1 = [...]uint8{0, 12, 28, 44, 60}
)

func (i VPType) String() string {
	switch {
	case i <= 6:
		return _VPType_name_0[_VPType_index_0[i]:_VPType_index_0[i+1]]
	case 18 <= i && i <= 21:
		i -= 18
		return _VPType_name_1[_VPType_index_1[i]:_VPType_index_1[i+1]]
	case i == 32:
		return _VPType_name_2
	default:
		return "VPType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-0]
	_ = x[NAME-1]
	_ = x[INT-2]
	_ = x[FLOAT-3]
	_ = x[IMAG-4]
	_ = x[CHAR-5]
	_ = x[STRING-6]
	_ = x[ADD-7]
	_ = x[SUB-8]
	_ = x[MUL-9]
	_ = x[QUO-10]
	_ = x[REM-11]
	_ = x[BAND-12]
	_ = x[BOR-13]
	_ = x[XOR-14]
	_ = x[SHL-15]
	_ = x[SHR-16]
	_ = x[BAND_NOT-17]
	_ = x[ADD_ASSIGN-18]
	_ = x[SUB_ASSIGN-19]
	_ = x[MUL_ASSIGN-20]
	_ = x[QUO_ASSIGN-21]
	_ = x[REM_ASSIGN-22]
	_ = x[BAND_ASSIGN-23]
	_ = x[BOR_ASSIGN-24]
	_ = x[XOR_ASSIGN-25]
	_ = x[SHL_ASSIGN-26]
	_ = x[SHR_ASSIGN-27]
	_ = x[BAND_NOT_ASSIGN-28]
	_ = x[LAND-29]
	_ = x[LOR-30]
	_ = x[ARROW-31]
	_ = x[INC-32]
	_ = x[DEC-33]
	_ = x[EQL-34]
	_ = x[LSS-35]
	_ = x[GTR-36]
	_ = x[ASSIGN-37]
	_ = x[NOT-38]
	_ = x[NEQ-39]
	_ = x[LEQ-40]
	_ = x[GEQ-41]
	_ = x[DEFINE-42]
	_ = x[BREAK-43]
	_ = x[CASE-44]
	_ = x[CHAN-45]
	_ = x[CONST-46]
	_ = x[CONTINUE-47]
	_ = x[DEFAULT-48]
	_ = x[DEFER-49]
	_ = x[ELSE-50]
	_ = x[FALLTHROUGH-51]
	_ = x[FOR-52]
	_ = x[FUNC-53]
	_ = x[GO-54]
	_ = x[GOTO-55]
	_ = x[IF-56]
	_ = x[IMPORT-57]
	_ = x[INTERFACE-58]
	_ = x[MAP-59]
	_ = x[PACKAGE-60]
	_ = x[RANGE-61]
	_ = x[RETURN-62]
	_ = x[SELECT-63]
	_ = x[STRUCT-64]
	_ = x[SWITCH-65]
	_ = x[TYPE-66]
	_ = x[VAR-67]
}

const _Word_name = "ILLEGALNAMEINTFLOATIMAGCHARSTRINGADDSUBMULQUOREMBANDBORXORSHLSHRBAND_NOTADD_ASSIGNSUB_ASSIGNMUL_ASSIGNQUO_ASSIGNREM_ASSIGNBAND_ASSIGNBOR_ASSIGNXOR_ASSIGNSHL_ASSIGNSHR_ASSIGNBAND_NOT_ASSIGNLANDLORARROWINCDECEQLLSSGTRASSIGNNOTNEQLEQGEQDEFINEBREAKCASECHANCONSTCONTINUEDEFAULTDEFERELSEFALLTHROUGHFORFUNCGOGOTOIFIMPORTINTERFACEMAPPACKAGERANGERETURNSELECTSTRUCTSWITCHTYPEVAR"

var _Word_index = [...]uint16{0, 7, 11, 14, 19, 23, 27, 33, 36, 39, 42, 45, 48, 52, 55, 58, 61, 64, 72, 82, 92, 102, 112, 122, 133, 143, 153, 163, 173, 188, 192, 195, 200, 203, 206, 209, 212, 215, 221, 224, 227, 230, 233, 239, 244, 248, 252, 257, 265, 272, 277, 281, 292, 295, 299, 301, 305, 307, 313, 322, 325, 332, 337, 343, 349, 355, 361, 365, 368}

func (i Word) String() string {
	if i < 0 || i >= Word(len(_Word_index)-1) {
		return "Word(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Word_name[_Word_index[i]:_Word_index[i+1]]
}

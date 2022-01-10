package deepcopy

/*
Feature support list:

	DeepCopy	Kind			Remarks

		o		Bool
		o		Int
		o		Int8
		o		Int16
		o		Int32
		o		Int64
		o		Uint
		o		Uint8
		o		Uint16
		o		Uint32
		o		Uint64
				Uintptr			Bitwise copy
		o		Float32
		o		Float64
		o		Complex64
		o		Complex128
		o		Array
				Chan			Share (bitwise copy) or create new
				Func			Bitwise copy or ignore (leaf it nil)
				Interface		Share (bitwise copy) or [Unsafe] DeepCopy
		o		Map
		o		Ptr
		o		Slice
		o		String
				Struct			DeepCopy exported fields or [Unsafe] all fields
				UnsafePointer	Bitwise copy

* If Remarks looks like "xxx or yyy", xxx is the default behaviour. Use WithOption()
	to specify it.
*/

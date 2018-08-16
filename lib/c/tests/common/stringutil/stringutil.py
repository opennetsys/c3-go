import unittest
from ctypes import *

# note: this file must first be built. see the make file
stringutil = CDLL('../../../common/stringutil/stringutil.so')

class Test():
    def __init__(self, inp, exp):
        self.inp=inp
        self.exp=exp

class CollectionOfTests():
    def __init__(self, StringTestObjects):
        self.collection = list(StringTestObjects) #the cast here is to indicate that some checks would be nice

class BytesResponse(Structure):
    _fields_ = [
        ("r0", c_void_p),
        ("r1", c_int),
    ]

stringutil.CompactJSON.restype = BytesResponse

CompactJSONTests = CollectionOfTests([
    Test("\x00\x01{\"foo\": \"bar\",\"hello\" :\"world\",\"a\": {\"x\":\"b\"}}\x00", "{\"foo\":\"bar\",\"hello\":\"world\",\"a\":{\"x\":\"b\"}}"),
    Test("\x00\x01[\"foo\", \"bar\",\"hello\" ,\"world\",[\"sub\"]]\x00", "[\"foo\",\"bar\",\"hello\",\"world\",[\"sub\"]]"),
    Test("{}", "{}"),
    Test("[]", "[]"),
    Test("", "{}"),
])

class TestStringUtil(unittest.TestCase):
    def test_compactJSON(self):
        for t in CompactJSONTests.collection:
            b = bytearray()
            b.extend(map(ord, t.inp))
            arr = (c_byte * len(b))(*b)

            res = stringutil.CompactJSON(arr, c_int(len(b)))
            ArrayType = c_ubyte*(c_int(res.r1).value)
            pa = cast(c_void_p(res.r0), POINTER(ArrayType))

            b1 = bytearray(t.exp, "utf8")
            self.assertEqual(c_int(res.r1).value, len(list(b1)))
            self.assertEqual(pa.contents[:], list(b1))

if __name__ == '__main__':
    unittest.main()

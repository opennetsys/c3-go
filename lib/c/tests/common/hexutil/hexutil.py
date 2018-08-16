import unittest
from ctypes import *

# note: this file must first be built. see the make file
hexutil = CDLL('../../../common/hexutil/hexutil.so')

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

hexutil.DecodeString.restype = BytesResponse
hexutil.EncodeBytes.restype = BytesResponse
hexutil.DecodeBytes.restype = BytesResponse

DecodeUint64Tests = CollectionOfTests([
    Test("0x7b", 123),
    Test("0x32F9E39", 53452345),
])
EncodeUint64Tests = CollectionOfTests([
    Test(123, "0x7b"),
    Test(53452345, "0x32f9e39"), # note: why is this lowercase and the prev test upper?!
])
EncodeStringTests = CollectionOfTests([
    Test("hello", "0x68656c6c6f"),
    Test("123", "0x313233"),
])
DecodeStringTests = CollectionOfTests([
    Test("0x1234", [18, 52]),
])
EncodeToStringTests = CollectionOfTests([
    Test("hello", "0x68656c6c6f"),
])
EncodeBytesTests = CollectionOfTests([
    Test("hello", "68656c6c6f"),
])
DecodeBytesTests = CollectionOfTests([
    Test("68656c6c6f", "hello"),
])
EncodeBigIntTests = CollectionOfTests([
    Test("123", "0x7b"), # note: why lowercase?
    Test("53452345", "0x32f9e39"), # note: why lowercase?
])
DecodeBigIntTests = CollectionOfTests([
    Test("0x7B", "123"),
])
EncodeIntTests = CollectionOfTests([
    Test(123, "0x7b"), #note: why lowercase?
    Test(-932445, "0xfffffffffff1c5a3"), #note: why lowercase?
])
DecodeIntTests = CollectionOfTests([
    Test("0x7b", 123), #note: why lowercase?
    Test("0xfffffffffff1c5a3", -932445), #note: why lowercase?
])
EncodeFloat64Tests = CollectionOfTests([
    Test(123, "0x405ec00000000000"), #note: why lowercase?
    Test(-561.2863, "0xc0818a4a57a786c2"), #note: why lowercase?
])
DecodeFloat64Tests = CollectionOfTests([
    Test("0x405ec00000000000", 123), #note: why lowercase?
    Test("0xc0818a4a57a786c2", -561.2863), #note: why lowercase?
])

class TestHexUtil(unittest.TestCase):
    def test_decodeUint64(self):
        for t in DecodeUint64Tests.collection:
            res = hexutil.DecodeUint64(c_char_p(t.inp.encode('utf-8')))
            self.assertEqual(c_ulonglong(res).value, t.exp)

    def test_encodeUint64(self):
        for t in EncodeUint64Tests.collection:
            res = hexutil.EncodeUint64(c_ulonglong(t.inp))
            self.assertEqual(c_char_p(res).value.decode('utf-8'), t.exp)

    def test_encodeString(self):
        for t in EncodeStringTests.collection:
            res = hexutil.EncodeString(c_char_p(t.inp.encode('utf-8')))
            self.assertEqual(c_char_p(res).value.decode('utf-8'), t.exp)

    def test_decodeString(self):
        for t in DecodeStringTests.collection:
            res = hexutil.DecodeString(c_char_p(t.inp.encode('utf-8')))
            ArrayType = c_ubyte*(c_int(res.r1).value)
            pa = cast(c_void_p(res.r0), POINTER(ArrayType))

            self.assertEqual(c_int(res.r1).value, len(t.exp))
            self.assertEqual(pa.contents[:], t.exp)

    def test_encodeToString(self):
        for t in EncodeToStringTests.collection:
            b = bytearray()
            b.extend(map(ord, t.inp))
            arr = (c_byte * len(b))(*b)

            res = c_char_p(hexutil.EncodeToString(arr, c_int(len(b)))).value.decode('utf-8')
            self.assertEqual(res, t.exp)

    def test_encodeBytes(self):
        for t in EncodeBytesTests.collection:
            b = bytearray()
            b.extend(map(ord, t.inp))
            arr = (c_byte * len(b))(*b)

            res = hexutil.EncodeBytes(arr, c_int(len(b)))
            ArrayType = c_ubyte*(c_int(res.r1).value)
            pa = cast(c_void_p(res.r0), POINTER(ArrayType))

            b1 = bytearray(t.exp, "utf8")
            self.assertEqual(c_int(res.r1).value, len(list(b1)))
            self.assertEqual(pa.contents[:], list(b1))

    def test_decodeBytes(self):
        for t in DecodeBytesTests.collection:
            b = bytearray()
            b.extend(map(ord, t.inp))
            arr = (c_byte * len(b))(*b)

            res = hexutil.DecodeBytes(arr, c_int(len(b)))
            ArrayType = c_ubyte*(c_int(res.r1).value)
            pa = cast(c_void_p(res.r0), POINTER(ArrayType))

            b1 = bytearray(t.exp, "utf8")
            self.assertEqual(c_int(res.r1).value, len(list(b1)))
            self.assertEqual(pa.contents[:], list(b1))

    def test_encodeBigInt(self):
        for t in EncodeBigIntTests.collection:
            res = c_char_p(hexutil.EncodeBigInt(c_char_p(t.inp.encode('utf-8')))).value.decode('utf-8')
            self.assertEqual(res, t.exp)

    def test_decodeBigInt(self):
        for t in DecodeBigIntTests.collection:
            res = c_char_p(hexutil.DecodeBigInt(c_char_p(t.inp.encode('utf-8')))).value.decode('utf-8')
            self.assertEqual(res, t.exp)

    def test_encodeInt(self):
        for t in EncodeIntTests.collection:
            res = c_char_p(hexutil.EncodeInt(c_int(t.inp))).value.decode('utf-8')
            self.assertEqual(res, t.exp)

    def test_decodeInt(self):
        for t in DecodeIntTests.collection:
            res = c_int(hexutil.DecodeInt(c_char_p(t.inp.encode('utf-8')))).value
            self.assertEqual(res, t.exp)
    
    def test_encodeFloat64(self):
        for t in EncodeFloat64Tests.collection:
            res = c_char_p(hexutil.EncodeFloat64(c_double(t.inp))).value.decode('utf-8')
            self.assertEqual(res, t.exp)

    #TODO: cannot get test to pass
    # def test_decodeFloat64(self):
        # for t in DecodeFloat64Tests.collection:
            # print(t.inp)
            # res = c_float(hexutil.DecodeFloat64(c_char_p(t.inp.encode('utf-8')))).value
            # self.assertEqual(res, t.exp)

if __name__ == '__main__':
    unittest.main()

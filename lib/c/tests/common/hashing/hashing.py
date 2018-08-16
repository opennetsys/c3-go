import unittest
from ctypes import *

# note: this file must first be built. see the make file
hashing = CDLL('../../../common/hashing/hashing.so')
hexutil = CDLL('../../../common/hexutil/hexutil.so')

class Test():
    def __init__(self, inp, exp):
        self.inp=inp
        self.exp=exp

class CollectionOfTests():
    def __init__(self, StringTestObjects):
        self.collection = list(StringTestObjects) #the cast here is to indicate that some checks would be nice

HashTests = CollectionOfTests([
    Test("hello world", "0x0ac561fac838104e3f2e4ad107b4bee3e938bf15f2b15f009ccccd61a913f017"),
])

class BytesResponse(Structure):
    _fields_ = [
        ("r0", c_void_p),
        ("r1", c_int),
    ]

hashing.Hash.restype = BytesResponse
hexutil.EncodeBytes.restype = BytesResponse

class TestHashing(unittest.TestCase):
    def test_hash(self):
        for t in HashTests.collection:
            b = bytearray()
            b.extend(map(ord, t.inp))
            arr = (c_byte * len(b))(*b)

            res = hashing.Hash(arr, c_int(len(arr)))
            ArrayType = c_ubyte*(c_int(res.r1).value)
            pa = cast(c_void_p(res.r0), POINTER(ArrayType))

            res1 = hexutil.EncodeBytes(c_void_p(res.r0), len(pa.contents[:]))
            ArrayType = c_ubyte*(c_int(res1.r1).value)
            pa1 = cast(c_void_p(res1.r0), POINTER(ArrayType))

            arr1 = bytes(t.exp, 'utf-8')
            self.assertEqual(c_int(res1.r1).value, len(arr1)-2)
            self.assertEqual("0x" + "".join(map(chr, pa1.contents[:])), arr1.decode("latin-1"))

    def test_hashToHexString(self):
        for t in HashTests.collection:
            b = bytearray()
            b.extend(map(ord, t.inp))
            arr = (c_byte * len(b))(*b)

            res = hashing.HashToHexString(arr, c_int(len(arr)))
            s = c_char_p(res).value

            self.assertEqual(s.decode('utf-8'), t.exp)

    def test_isEqual(self):
        b = bytearray()
        b.extend(map(ord, "foo"))
        arr = (c_byte * len(b))(*b)

        eq = hashing.IsEqual(c_char_p("0x1234".encode('utf-8')), arr, c_int(len(arr)))

        self.assertEqual(0, c_int(eq).value)

        b = bytearray()
        b.extend(map(ord, "hello world"))
        arr = (c_byte * len(b))(*b)

        eq = hashing.IsEqual(c_char_p("0x0ac561fac838104e3f2e4ad107b4bee3e938bf15f2b15f009ccccd61a913f017".encode('utf-8')), arr, c_int(len(arr)))

        self.assertEqual(1, c_int(eq).value)

if __name__ == '__main__':
    unittest.main()

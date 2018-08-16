# Copyright (c) 2013-2014 Conformal Systems LLC.
# Use of this source code is governed by an ISC
# license that can be found in the LICENSE file.

import unittest
from ctypes import *

# note: this file must first be built. see the make file
base58 = CDLL('../../../common/base58/base58.so')

class StringTest():
    def __init__(self, inp, exp):
        self.inp=inp
        self.exp=exp

class CollectionOfStringTests():
    def __init__(self, StringTestObjects):
        self.collection = list(StringTestObjects) #the cast here is to indicate that some checks would be nice

StringTests = CollectionOfStringTests([
    StringTest("", ""),
    StringTest(" ", "Z"),
    StringTest("-", "n"),
    StringTest("0", "q"),
    StringTest("1", "r"),
    StringTest("-1", "4SU"),
    StringTest("11", "4k8"),
    StringTest("abc", "ZiCa"),
    StringTest("1234598760", "3mJr7AoUXx2Wqd"),
    StringTest("abcdefghijklmnopqrstuvwxyz", "3yxU3u1igY8WkgtjK92fbJQCd4BZiiT1v25f"),
    StringTest("00000000000000000000000000000000000000000000000000000000000000", "3sN2THZeE9Eh9eYrwkvZqNstbHGvrxSAM7gXUXvyFQP8XvQLUqNCS27icwUeDT7ckHm4FUHM2mTVh1vbLmk7y"),
])

HexTests = CollectionOfStringTests([
    StringTest("2g", "61"),
    StringTest("a3gV", "626262"),
    StringTest("aPEr", "636363"),
    StringTest("2cFupjhnEsSn59qHXstmK2ffpLv2", "73696d706c792061206c6f6e6720737472696e67"),
    StringTest("1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L", "00eb15231dfceb60925886b67d065299925915aeb172c06647"),
    StringTest("ABnLTmg", "516b6fcd0f"),
    StringTest("3SEo3LWLoPntC", "bf4f89001e670274dd"),
    StringTest("3EFU7m", "572e4794"),
    StringTest("EJDM8drfXA6uyA", "ecac89cad93923c02321"),
    StringTest("Rt5zm", "10c8511e"),
    StringTest("1111111111", "00000000000000000000"),
])

InvalidStringTests = CollectionOfStringTests([
    StringTest("0", ""),
    StringTest("O", ""),
    StringTest("I", ""),
    StringTest("l", ""),
    StringTest("3mJr0", ""),
    StringTest("O3yxU", ""),
    StringTest("3sNI", ""),
    StringTest("4kl8", ""),
    StringTest("0OIl", ""),
    StringTest("!@#$%^&*()-_=+~`", ""),
])

class DecodeResponse(Structure):
    _fields_ = [
        ("r0", c_void_p),
        ("r1", c_int),
    ]

base58.Decode.restype = DecodeResponse

class TestBase58(unittest.TestCase):
    def test_encode(self):
        for t in StringTests.collection:
            b = bytearray()
            b.extend(map(ord, t.inp))
            arr = (c_byte * len(b))(*b)

            actual = c_char_p(base58.Encode(arr, c_int(len(b)))).value.decode('utf-8')
            self.assertEqual(actual, t.exp)

    def test_decode(self):
        for t in HexTests.collection:
            res = base58.Decode(c_char_p(t.inp.encode('utf-8')))
            ArrayType = c_ubyte*(c_int(res.r1).value)
            pa = cast(c_void_p(res.r0), POINTER(ArrayType))

            arr = bytes.fromhex(t.exp)
            self.assertEqual(c_int(res.r1).value, len(arr))
            self.assertEqual("".join(map(chr, pa.contents[:])), arr.decode("latin-1"))

    def test_invalid_strings(self):
        for t in InvalidStringTests.collection:
            res = base58.Decode(c_char_p(t.inp.encode('utf-8')))
            ArrayType = c_ubyte*(c_int(res.r1).value)
            pa = cast(c_void_p(res.r0), POINTER(ArrayType))

            arr = bytes.fromhex(t.exp)
            self.assertEqual(c_int(res.r1).value, len(arr))
            self.assertEqual("".join(map(chr, pa.contents[:])), arr.decode("latin-1"))

if __name__ == '__main__':
    unittest.main()

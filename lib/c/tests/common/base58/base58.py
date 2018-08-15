import unittest
from ctypes import *

base58 = CDLL('../../../common/base58/base58.h')

class StringTest():
    def __init__(self, inp, exp):
        self.inp=inp
        self.exp=exp

class CollectionOfStringTests():
    def __init__(self, StringTestObjects):
        self.tests = list(StringTestObjects) #the cast here is to indicate that some checks would be nice

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

class TestBase58(unittest.TestCase):
    def test_encode(self):
        for t in StringTests.tests:
            b = bytearray()
            b.extend(map(ord, t.inp))

            self.assertEqual(base58.Encode(b, len(b)), t.exp)

if __name__ == '__main__':
    unittest.main()

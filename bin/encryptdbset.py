# Copyright (c) 2019 - 2023, t.

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
#
# Original by triangle
# Modified by arina
#
# https://discord.com/channels/922182394323292170/1102114858440327269/1137939760149704746

import binascii
import os
import struct
import sys
import hashlib
import zlib

import hwdecrypt

# keyspec_311 = [0x06856c49, 0x3aa19541, 0x5f13a7c1]
keyspec_312 = [0x49e66da3, 0x59e1e89a, 0x24ebb207]

class FileReference:
    def __init__(self, version, name, sha):
        self.version = version
        self.name = name
        self.sha = sha
        self.encrypted_sha = None
        self.size = None

    def getkeys(self, init_keys):
        keys = [int(self.sha[:8], 16), int(self.sha[8:16], 16), int(self.sha[16:24], 16)]
        return [a ^ b for a, b in zip(init_keys, keys)]

    def __repr__(self):
        return "<FileReference {0} {1} {2} {3}>".format(
            self.name, self.sha, self.encrypted_sha, self.size
        )

class Manifest:
    def __init__(self, fobj):
        sha1hash = self.eatbytes(fobj, 20)
        self.version = self.prefixstring(fobj)
        self.lang = self.prefixstring(fobj)

        self.files = []
        nfentries = self.ubyte(fobj)
        for i in range(nfentries):
            na = self.prefixstring(fobj)
            ha = self.prefixstring(fobj)
            self.files.append(FileReference(self.version, na, ha))

        for i in range(nfentries):
            sha1 = binascii.hexlify(self.eatbytes(fobj, 20)).decode("ascii")
            size = self.uint(fobj)
            self.files[i].encrypted_sha = sha1
            self.files[i].size = size

    def eatbytes(self, stream, n):
        return stream.read(n)

    def prefixstring(self, stream):
        ln = self.ubyte(stream)
        return stream.read(ln).decode("ascii")

    def ubyte(self, stream):
        return struct.unpack("<B", stream.read(1))[0]

    def uint(self, stream):
        b1, b2, b3, b4 = struct.unpack("<BBBB", stream.read(4))
        f = b4 << 8
        f = (f + b3) << 8
        f = (f + b2) << 7
        f = b1 + f
        return f

    def s_ubyte(self, b):
        return struct.pack("B", b)

    def s_uint(self, i):
        # wtf???
        b1 = i & 0xff
        if not (b1 & 0x80):
            b1 = b1 | 0x80
        rest = (i - b1)
        b2 = (rest >> 7) & 0xff
        rest = (rest >> 7) - b2
        b3 = (rest >> 8) & 0xff
        b4 = (rest >> 16) & 0xff
        return bytes([b1, b2, b3, b4])

    def s_pstring(self, s):
        b = s.encode("utf8")
        return b"".join([struct.pack("B", len(b)), b])

    def serialize(self, stream):
        buf = []

        buf.append(self.s_pstring(self.version))
        buf.append(self.s_pstring(self.lang))
        buf.append(self.s_ubyte(len(self.files)))

        for f in self.files:
            buf.append(self.s_pstring(f.name))
            buf.append(self.s_pstring(f.sha))
        
        for f in self.files:
            buf.append(binascii.unhexlify(f.encrypted_sha))
            buf.append(self.s_uint(f.size))
        
        blob = b"".join(buf)
        digest = hashlib.sha1(blob).digest()
        stream.write(digest)
        stream.write(blob)

    def __repr__(self):
        return "<Manifest {0} {1} {2} files>".format(
            self.version, self.lang, len(self.files)
        )

def hash_file(in_dir, from_file: FileReference):
    h = hashlib.sha1()
    with open(os.path.join(in_dir, from_file.name), "rb") as f:
        h.update(f.read())
    return h.hexdigest()

def rekey(in_dir, from_file: FileReference, out_dir: str, to_keyspec):
    # 1. read clear db file
    with open(os.path.join(in_dir, from_file.name), "rb") as f:
        clr_buf = f.read()
    assert clr_buf[:16] == b"SQLite format 3\x00", "Missing SQLite file signature. Is it already encrypted?"
    # 2. zlib compress with no header (requires python 3.11+)
    crypt_buf = bytearray(zlib.compress(clr_buf, wbits=-zlib.MAX_WBITS))
    # 3. encrypt compressed db for as client
    ks1 = from_file.getkeys(to_keyspec)
    keys = hwdecrypt.Keyset(ks1[0], ks1[1], ks1[2])
    hwdecrypt.decrypt(keys, crypt_buf)

    with open(os.path.join(out_dir, from_file.name), "wb") as f:
        f.write(crypt_buf)

def main():
    out_dir = sys.argv[1]
    update_original = sys.argv[2]
    masterdata_refs = sys.argv[3:]
    print(out_dir, update_original, masterdata_refs)
    new_version = hashlib.sha1()
    for masterdata_ref in masterdata_refs:
        basedir = os.path.dirname(masterdata_ref)
        with open(masterdata_ref, "rb") as f:
            mf = Manifest(f)

        for file in mf.files:
            new_version.update(bytes(hash_file(basedir, file), 'ascii'))

    new_version = new_version.hexdigest()
    print(f'Version: {new_version}')
    out_dir = os.path.join(out_dir, new_version)
    os.makedirs(out_dir, exist_ok=True)

    for masterdata_ref in masterdata_refs:
        basedir = os.path.dirname(masterdata_ref)
        with open(masterdata_ref, "rb") as f:
            mf = Manifest(f)
        mf.version = new_version
        for file in mf.files:
            file.sha = hash_file(basedir, file)
            print(file)
            rekey(basedir, file, out_dir, keyspec_312)
            file.size = os.path.getsize(os.path.join(out_dir, file.name))
            file.encrypted_sha = hash_file(out_dir, file)

        with open(os.path.join(out_dir, os.path.basename(masterdata_ref)), "wb") as f:
            mf.serialize(f)
        if update_original == 'update':
            # update input file
            file_name = os.path.basename(masterdata_ref)
            # 2 backups
            backup_path = os.path.join(basedir, f'{file_name}.backup')
            backup_backup_path = os.path.join(basedir, f'{file_name}.backup.backup')
            if os.path.exists(backup_path):
                if os.path.exists(backup_backup_path):
                    os.remove(backup_backup_path)
                os.rename(backup_path, backup_backup_path)
            os.rename(masterdata_ref, backup_path)
            with open(masterdata_ref, "wb") as f:
                mf.serialize(f)


if __name__ == "__main__":
    main()
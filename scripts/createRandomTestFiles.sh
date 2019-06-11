#!/bin/bash
# file: createRandomTestFiles.sh

dd if=/dev/urandom of=random1k.bin bs=1K count=1 iflag=fullblock
dd if=/dev/urandom of=random3k.bin bs=1K count=3 iflag=fullblock
dd if=/dev/urandom of=random5k.bin bs=1K count=5 iflag=fullblock

dd if=/dev/urandom of=random1M.bin bs=1M count=1 iflag=fullblock
dd if=/dev/urandom of=random3M.bin bs=1M count=3 iflag=fullblock
dd if=/dev/urandom of=random5M.bin bs=1M count=5 iflag=fullblock
dd if=/dev/urandom of=random20M.bin bs=1M count=20 iflag=fullblock
dd if=/dev/urandom of=random21M.bin bs=1M count=21 iflag=fullblock

dd if=/dev/urandom of=random1G.bin bs=64M count=16 iflag=fullblock

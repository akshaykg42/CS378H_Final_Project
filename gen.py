import random
import string
from random import randint

blocksizes = {1, 10, 100, 1000, 2000, 5000, 10000}
inputsizes = {10000000}

for inputsize in inputsizes:
	for blocksize in blocksizes:
		s = ""
		s += str(blocksize)
		s += "\n"
		for i in range(0, inputsize):
			s += str(randint(0, 1))
		s += "\n"
		f = open("ctr_input_" + str(inputsize) + "_" + str(blocksize) + ".txt", "x")
		f.write(s)


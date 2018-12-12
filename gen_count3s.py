import random
import string
from random import randint

blocksizes = {1, 10, 100, 1000, 2000, 5000, 10000}
inputsizes = {10000, 100000, 1000000, 10000000}

for inputsize in inputsizes:
	for blocksize in blocksizes:
		s = str(blocksize) + "\n"
		for i in range(inputsize):
			s += str(randint(1, 11)) + "\n"
		f = open("count3s_input_" + str(inputsize) + "_" + str(blocksize) + ".txt", "x")
		f.write(s)
		f.close()
import demucs.separate as demucs
import sys

runArgs = sys.argv[1:]

demucs.main(runArgs)

'''
May be worth adding more commands to this file
    i.e. calling to check if cuda is available

For now, it's just a wrapper for demucs.separate
'''
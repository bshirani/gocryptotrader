import glob
import os


def last_file_in_dir(dirname):
    return max(glob.glob(dirname), key=os.path.getctime)


def last_dir_in_dir(dirname):
    return max(glob.glob(dirname), key=os.path.getctime)

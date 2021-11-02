import glob
import os


def last_file_in_dir(dirname):
    list_of_files = glob.glob(dirname)
    return max(list_of_files, key=os.path.getctime)

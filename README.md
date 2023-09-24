# TimelineTools

TimelineTools allows for a user to easily backup files 
from all folders specified in inputFolders.csv, and 
copy them over to a folder structure containing 
subfolders per date. It consists of two programs, 
one to do the backup, and one to provide a list of
numbers of folders per date and filetype.

## files2timeline

This program reads the inputFolders.csv and the 
extensions.csv files, and then proceeds to iterate
over all files in the folders named in the first 
file. 

If the file has an extension matching any of 
the ones in the second file, and the file was 
modified within the last 30 days, then it is copied 
to the target dated folder. The filename is extended
with a suffix containing part of the hash of the 
contents, to be able to store several versions of 
the file in the same dated folder.

## timelinecounter

This program iterates over the specified timeline 
folder and sums up the numbers of files per dated 
subfolder and extension. A CSV file named 
file_count.csv is written to timeline 

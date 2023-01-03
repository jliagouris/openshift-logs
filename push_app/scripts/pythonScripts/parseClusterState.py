import sys

#print(f"Arguments of the script : {sys.argv[1:]=}")
len = len(sys.argv)
for i in range(len):
    if sys.argv[i] == 'State:':
        print(sys.argv[i + 1])
#open text file in read mode
import sys

filename = sys.path[0] + "/../createAdmin-"
filename += sys.argv[1]
filename += "-" + str(sys.argv[2])
filename += ".txt"
text_file = open(filename, "r")
 
#read whole file to a string
admin_login_cmd = text_file.read()
 
#close file
text_file.close()
 
#print(admin_login_cmd, '\n')

splitRes = admin_login_cmd.splitlines()

print(splitRes[4].strip())


#open text file in read mode
text_file = open("scripts/createAdmin.txt", "r")
 
#read whole file to a string
admin_login_cmd = text_file.read()
 
#close file
text_file.close()
 
#print(admin_login_cmd, '\n')

splitRes = admin_login_cmd.splitlines()

print(splitRes[4].strip())


with open('b64shellcode.txt', 'r') as f:
    line = f.readline()

chunk_size = 60
out = ""

for i in range(len(line)):
    out += line[i]
    if (i % chunk_size == 0 or i == len(line) -1) and i != 0:
        if i == len(line) - 1:
            print("\"" + out + "\"")
            break
        print("\"" + out + "\" +")
        out = ""

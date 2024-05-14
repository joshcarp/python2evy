for i in range(5):
    print("for", i)

count = 0
while count < 5:
    print("while", count)
    count += 1

for i in range(1, 4):
    for j in range(1, 4):
        if i != j:
            print(f"({i}, {j})")

# Open and read a result file
tab = []
with open('.\\result\\result.csv', 'r') as file:

    for line in file:
        tab.append(line.strip().split(','))

# print(tab[1])

instances = {}

for i in range(1, len(tab)):
    instance_name = tab[i][0]
    if instance_name not in instances:
        instances[instance_name] = tab[i][5]
    else:
        instances[instance_name] = max(instances[instance_name], tab[i][5])

print(instances)
# save in file
with open('times.csv', 'w') as file:
    file.write('instance,best_result\n')
    for instance_name, best_result in instances.items():
        file.write(f'{instance_name},{best_result}\n')
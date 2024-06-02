import matplotlib.pyplot as plt

fig, ax = plt.subplots()
points = []
with open("./results.txt", "r") as f:
    data = f.read().split("\n")[:-1]
    for i in data:
        x, y = i.split(' ')
        points.append([float(x), int(y)])



x_axis = [point[0] - points[0][0] for point in points]
y_axis = [point[1] for point in points]
ax.plot(x_axis, y_axis)
ax.set(xlabel='Seconds', ylabel='Number of jobs', title='NOJ/Seconds')
plt.show()

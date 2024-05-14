def greet(name):
    print("Hello,", name)
greet("Alice")

def concat(a: str, b: str) -> str:
    return a + b

a = "foo"
b = "bar"
print(concat(a, b))

def calculate_area(length, width):
    area = length * width
    return area

result = calculate_area(5, 8)
print("Area of the rectangle:", result)
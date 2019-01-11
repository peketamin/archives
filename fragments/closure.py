# coding: utf-8

def accumulator(y):
    def inner(x):
        return x + y
    return inner

if __name__ == "__main__":
    a = accumulator(10)
    for i in range(5):
        print a(i)

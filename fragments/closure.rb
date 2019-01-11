def accumulator(x)
    lambda {|y|  return x + y }
end


i = 0
loop{
    a = accumulator(10)
    puts a.call(i)
    i += 1
    if i > 5
        break
    end
}


misc = {
    "fact4": ("""
        letrec fact(n) = if (isz n) then 1 else * (n, (fact -(n, 1)))
        in
        (fact 4)
    """, 24)
}

lazy = {
    "infinite": ("""
        letrec infinite-loop(x) = ' ( infinite-loop(x,1) )
            in let f = proc(z) 11
                in (f (infinite-loop 0))
    """, 11)
}

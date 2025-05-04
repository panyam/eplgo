package chapter3

/*
letreclang = {
    "double": ("""
        letrec double(x) = if isz(x) then 0 else - ((double -(x,1)), -2)
        in (double 6)
    """, 12),
    "oddeven": ("""
            letrec
                even(x) = if isz(x) then 1 else (odd -(x,1))
                odd(x) = if isz(x) then 0 else (even -(x,1))
            in (odd 13)
        """, 1),
    "currying": ("""
            letrec f(x,y) = if (isz y)
                            then x
                            else (f +(x,y))
            in
            (f 1 2 3 4 5 0)
        """, 15)
}
*/


letlang = {
    "num": (""" 3 """, 3),
    "var": (""" x """, 5),
    "isz_true": ("isz 0", True),
    "isz_false": ("isz 1", False),
    "diff": ("-(-(x,3),-(v,i))", 3),
    "if": (" if isz -(x,11) then -(y,2) else -(y,4) ", 18),
    "let": ("let x = 5 in -(x,3)", 2),
    "letnested": ("let z = 5 in let x = 3 in let y = -(x, 1) in let x = 4 in -(z, -(x,y))", 3),
    "let3": ("""
        let x = 7 in
            let y = 2 in
                let y = let x = -(x,1) in -(x,y)
                in -(-(x,8), y)
        """, -5),
    "letmultiargs": ("""
        let x = 7 y = 2 in
            let y = let x = -(x,1) in -(x,y)
            in -(-(x, 8),y)
        """, -5),
}

proclang = {
    "proc1": (""" let f = proc (x) -(x,11) in (f (f 77)) """, 55),
    "proc2": ("""
        let x = 200 in
            let f = proc(z) -(z,x) in
                let x = 100 in
                    let g = proc(z) -(z,x) in
                        -((f 1), (g 1))
        """, -100),
    "proc_multiargs": ("""
        let f = proc(x,y) -(x,y) in
            -((f 1 10), (f 10 5))
    """, -14),
    "proc_currying": (""" let f = proc(x,y) -(x,y) in ((f 5) 3) """, 2),
    "proc_currying2": ("""
        let f = proc(x,y)
                if (isz y)
                then x
                else proc(a,b) (if isz b then +(a,x,y) else +(a,b,x,y))
        in
        (f 1 2 2 0)
    """, 5),
}

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


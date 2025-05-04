
exprefs = {
    "oddeven": ("""
        let x = newref(0) in
            letrec
                even(dummy)
                    = if isz(deref(x))
                      then 1
                      else begin
                        setref(x, -(deref(x), 1));
                        (odd 888)
                      end
                odd(dummy)
                    = if isz(deref(x))
                      then 0
                      else begin
                        setref(x, -(deref(x), 1));
                        (even 888)
                      end
            in begin setref(x, 13) ; (odd 888) end
            """, 1),
    "counter": ("""
        let g = let counter = newref(0)
                in proc(dummy)
                    begin
                        setref(counter, -(deref(counter), -1)) ;
                        deref(counter)
                    end
        in let a = (g 11)
            in let b = (g 11)
                in -(a,b)
    """, -1)
}

imprefs = {
    "oddeven": ("""
        let x = 0 in
            letrec
                even(dummy)
                    = if isz(x)
                      then 1
                      else begin
                        set x = -(x,1) ;
                        (odd 888)
                      end
                odd(dummy)
                    = if isz(x)
                      then 0
                      else begin
                        set x = -(x,1) ;
                        (even 888)
                      end
            in begin set x = 13 ; (odd 888) end
        """, 1),
    "counter": ("""
            let g = let counter = 0
                    in proc(dummy)
                        begin
                            set counter = -(counter, -1) ;
                            counter
                        end
            in let a = (g 11)
                in let b = (g 11)
                    in -(a,b)
        """, -1),

    "recproc": ("""
            let times4 = 0 in
                begin
                    set times4 = proc(x)
                                    if isz(x)
                                    then 0
                                else -((times4 -(x,1)), -4) ;
                    (times4 3)
                end
        """, 12),
    "callbyref": ("""
            let a = 3
            in let b = 4
                in let swap = proc(x,y)
                                let temp = deref(x)
                                in begin
                                    setref(x, deref(y));
                                    setref(y,temp)
                                end
                    in begin ((swap ref a) ref b) ; -(a,b) end
        """, 1),
}

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

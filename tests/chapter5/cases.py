
exceptions = {
        "listindex": ("""
        proc(str) 
            letrec inner (lst) =
                if (isnull lst)
                then raise ("ListIndexFailed")
                else 
                    ( if (stringequals (car lst) str)
                    then 0 
                    else - ((inner (cdr lst)), -1) )
            in 5
    """, 5)
        
}

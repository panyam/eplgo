from epl.chapter7.typed import Type

basic_checks = {
    "non_bool_test":        ( "if 3 then 88 else 99", False ),
    "proc_val_rator":       ( "proc (x) (x 3)", True ),
    "non_proc_val_rator1":  ( "proc (x) (3 x)", False ),
    "non_proc_val_rator2":  ( "let x = 4 in (x 3)", False ),
    "non_proc_val_rator3":  ( "(proc (x) (x 3) 4)", False ),
    "non_int_diff_arg":     ( "let x = isz(0) in -(3,x)", False ),
    "non_int_diff_arg2":    ( "proc(x) -(3,x) isz(0)", False ),
    "non_proc_val_rator4":  ( "let f = 3 in proc (x) (f x)", False ),
    "non_proc_val_rator5":  ( "( proc (f) proc (x) (f x) 3 )", False ),
    # "safe_non_terminating": ( "letrec f(x) = (f -(x, -1)) in f 1", False ),
    "3":        ( "3", Type.as_leaf('int') ),
    "diff":     ( "-(10,20)", Type.as_leaf('int') ),
}

checked = {
    "proc":     ( "proc (x : int) -> int -(x, 11)",
                   Type.as_func([Type.as_leaf('int')],
                                 Type.as_leaf('int')) ),
    "proc2":    ( "proc (x : int) -(x, 11)",
                   Type.as_func([Type.as_leaf('int')], None) ),
}

inferred = {
    "case1":  (
        """
        letrec
            even (x : int) -> ? = if isz(x) then 1 else (odd -(x, 1))
            odd (x : ?) -> bool = if isz(x) then 0 else (even -(x, 1))
            in (odd 13)
        """, True )
}

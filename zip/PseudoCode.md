# 1. Recursive versions of insertion and deletion.
```
insert(x, root):
    if root = null then {x.left ← x.right ← null; x.rank ← RandomRank; return x}
    if x.key < root.key then
        if insert(x, root.left) = x then
            if x.rank < root.rank then root.left ← x
            else {root.left ← x.right; x.right ← root; return x}
    else
        if insert(x, root.right) = x then
            if x.rank ≤ root.rank then root.right ← x
            else {root.right ← x.left; x.left ← root; return x}
return root

zip(x, y):
    if x = null then return y
    if y = null then return x
    if x.rank < y.rank then {y.left ← zip(x, y.left); return y}
    else {x.right ← zip(x.right , y); return x}

delete(x, root):
    if x.key = root.key then return zip(root.left , root.right)
    if x.key < root.key then
        if x.key = root.left.key then
            root.left ← zip(root.left.left , root.left.right)
        else delete(x, root.left)
    else
        if x.key = root.right.key then
            root.right ← zip(root.right.left , root.right.right)
        else delete(x, root.right)
return root
```

# 2. Iterative Insertion
```
insert(x)
    rank ← x.rank ← RandomRank
    key ← x.key
    cur ← root
    while cur 6= null and (rank < cur .rank or (rank = cur .rank and key > cur .key)) do
        prev ← cur
        cur ← if key < cur .key then cur .left else cur .right
    if cur = root then root ← x
    else if key < prev.key then prev.left ← x
    else prev.right ← x
    if cur = null then {x.left ← x.right ← null; return}
    if key < cur .key then x.right ← cur else x.left ← cur
    prev ← x
    while cur 6= null do
        fix ← prev
        if cur .key < key then
            repeat {prev ← cur ; cur ← cur .right}
            until cur = null or cur .key > key
        else
            repeat {prev ← cur ; cur ← cur .left}
            until cur = null or cur .key < key
        if fix .key > key or (fix = x and prev.key > key) then
            fix .left ← cur
        else
            fix .right ← cur
```


# 3. Iterative Deletion
```
delete(x)
    key ← x.key
    cur ← root
    while key 6= cur .key do
        prev ← cur
        cur ← if key < cur .key then cur .left else cur .right
    left ← cur .left; right ← cur .right
    if left = null then cur ← right
    else if right = null then cur ← left
    else if left.rank ≥ right.rank then cur ← left
    else cur ← right
    if root = x then root ← cur
    else if key < prev.key then prev.left ← cur
    else prev.right ← cur
    while left 6= null and right 6= null do
        if left.rank ≥ right.rank then
            repeat {prev ← left; left ← left.right}
            until left = null or left.rank < right.rank
            prev.right ← right
        else
            repeat {prev ← right; right ← right.left}
            until right = null or left.rank ≥ right.rank
            prev.left ← left
```

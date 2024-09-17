## Reference
"Zip-zip Trees: Making Zip Trees More Balanced,
Biased, Compact, or Persistent"  <br>
\<Ofek Gila1,  Michael T. Goodrich1,
and Robert E. Tarjan\>

# 1. Insertion
```
function Insert(x)
    rank ← x.rank ← RandomRank
    key ← x.key
    cur ← root
    while cur ̸= null and (rank < cur.rank or (rank = cur.rank and key > cur.key)) do 
        prev ← cur
        cur ← if key < cur.key then cur.lef t else cur.right
    if cur = root then root ← x
    else if key < prev.key then prev.lef t ← x
    else prev.right ← x
    
    if cur = null then { x.lef t ← x.right ← null;  return }
    if key < cur.key then x.right ← cur else x.lef t ← cur
    prev ← x
    
    while cur ̸= null do
        f ix ← prev
        
        if cur.key < key then
            repeat { prev ← cur;  cur ← cur.right }
            until cur = null or cur.key > key
        else
            repeat { prev ← cur;  cur ← cur.lef t }
            until cur = null or cur.key < key
            
        if f ix.key > key or (f ix = x and prev.key > key) then
            f ix.lef t ← cur
        else
            f ix.right ← cur
```

# 2. Deletion
```
function Delete(x)
key ← x.key
cur ← root
while key ̸= cur.key do
prev ← cur
cur ← if key < cur.key then cur.lef t else cur.right lef t ← cur.lef t;  right ← cur.right
if lef t = null then cur ← right
else if right = null then cur ← lef t
else if lef t.rank ≥ right.rank then cur ← lef t else cur ← right
if root = x then root ← cur
else if key < prev.key then prev.lef t ← cur
else prev.right ← cur
while lef t ̸= null and right ̸= null do
if lef t.rank ≥ right.rank then
repeat { prev ← lef t; lef t ← lef t.right } until lef t = null or lef t.rank < right.rank prev.right ← right
else
repeat { prev ← right; right ← right.lef t } until right = null or lef t.rank ≥ right.rank prev.lef t ← lef 
```

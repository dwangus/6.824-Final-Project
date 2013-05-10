package document

/*
# Operations are lists of components.
# Each component either inserts or deletes at a specified position in the document.
#
# Components are either:
#  {i:'str', p:100}: Insert 'str' at position 100 in the document
#  {d:'str', p:100}: Delete 'str' at position 100 in the document
#
# Components in an operation are executed sequentially, so the position of components
# assumes previous components have already executed.
#
# Eg: This op:
#   [{Insert:'abc', Position:0}]
# is equivalent to this op:
#   [{Insert:'a', Position:0}, {Insert:'b', Position:1}, {Insert:'c', Position:2}]
*/

type Component struct {
  Insert string
  Delete string
  Position int
}

type TextOp []Component

func transformPosition(pos int, comp Component) (newpos int) {
  if(comp.Insert != ""){
    if comp.Position <= pos {
      newpos = pos + len(comp.Insert)
    } else {
      newpos = pos
    }
  } else{
    if(pos <= comp.Position){
      newpos = pos
    } else if pos <= comp.Position + len(comp.Delete){
      newpos = comp.Position
    } else {
      newpos = pos - len(comp.Delete)
    }
  }
  return
}

func (op TextOp) Append(comp Component) {
  op = append(op, comp)
}

func (comp1 Component) transform(dest *TextOp, comp2 Component) {
  pos1 := comp1.Position
  if comp1.Insert != "" { //Insert
    comp1.Position = transformPosition(pos1, comp2)
  } else { //Delete
    if comp2.Insert != "" { // Delete vs Insert
      deleted := comp1.Delete
      if pos1 < comp2.Position {
        (*dest).Append(Component{Position: comp1.Position, Delete:   deleted[:comp2.Position-pos1]})
        deleted = deleted[comp2.Position-pos1:]
      }
      if deleted != "" {
        (*dest).Append(Component{Position: comp1.Position+len(comp2.Insert), Delete: deleted})
      }
    }
  }
  return
}

func (op1 TextOp) transform(op2 TextOp) TextOp {
  for _, comp2 := range op2 {
    for _, comp1 := range op1 {
      comp1.transform(&op1, comp2)
    }
  }
  return op1
}
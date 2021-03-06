package main

import (

    "fmt"
    "math/rand"
    "time"

)

const (
    Unknown = 0
    RED = 1
    BLACK = 2 
)


const (
    
    Neutron = 0
    LEFT_CHILD = 1
    RIGHT_CHILD = 2 

)

const (
    
    FAIL = 0 
    SUCCESS = 1
    EXISTED = 2

)

type CMPFun func(x,y interface{}) int        // 函数指针



type RB_NODE struct {
    parent *RB_NODE    
    left   *RB_NODE
    right  *RB_NODE
    col     int    
    value   interface{}                // 空接口：赋值任意类型变量
}


type RB struct{
    root *RB_NODE
    cmp  CMPFun
    cnt  int
    delerr int
    adderr int
}

func RBTreeInit(cmp CMPFun) *RB{   //参数是一个函数指针

    var hdl *RB = new(RB)
    
    if hdl == nil || cmp == nil{
        return nil
    }
    hdl.root = nil
    hdl.cmp = cmp
    hdl.cnt = 0
    
    return hdl
}

func RBTreeDestroy(hdl *RB) int{

    return 1
}

func rotate_left(this *RB,nd *RB_NODE){

    var p *RB_NODE
    
    if nil == nd || nil == nd.right{
    
        fmt.Printf("OOP! current node or right child is nil\n")
        return 
    }
    
    p = nd.right
    nd.right = p.left
    
    if p.left != nil {
        p.left.parent = nd
    }
    
    p.left = nd 
    
    
    
    p.parent = nd.parent
    nd.parent = p


    if p.parent == nil {
        this.root = p
    }else {
        
        if p.parent.left == nd {
            p.parent.left = p
        } else {
            p.parent.right = p
        }
    }

}


func rotate_right(this *RB,nd *RB_NODE){

    var p *RB_NODE
    
    if nil == nd || nil == nd.left{
    
        fmt.Printf("OOP! current node or left child is nil\n")
        
        return 
    }
    
    p = nd.left
    nd.left = p.right
    
    if p.right != nil {
        p.right.parent = nd
    }
    p.right = nd 
    
    p.parent = nd.parent
    nd.parent = p


    if p.parent == nil {
        this.root = p
    }else {
        
        if p.parent.left == nd {
            p.parent.left = p
        } else {
            p.parent.right = p
        }
    }

}


func insert_fixup(this *RB,nd *RB_NODE) {
    
    var uncle *RB_NODE
    
    for nd.parent != nil && nd.parent.col == RED { // red-red conflict
    
        if nd.parent == nd.parent.parent.left{
            uncle = nd.parent.parent.right
            if uncle != nil && uncle.col == RED {
            
                 nd.parent.col = BLACK
                 uncle.col = BLACK
                 nd.parent.parent.col = RED
                 nd = nd.parent.parent
            
            }else {
            
                if(nd == nd.parent.right){
                    nd = nd.parent
                    rotate_left(this, nd)
                }
                
                nd.parent.col = BLACK
                nd.parent.parent.col = RED
                rotate_right(this,nd.parent.parent)
            
            }
        
        }else {
            uncle = nd.parent.parent.left
            if uncle != nil && uncle.col == RED {
                nd.parent.col = BLACK
                uncle.col = BLACK
                nd.parent.parent.col = RED
                nd = nd.parent.parent
            }else {
                
                if(nd == nd.parent.left){
                   nd= nd.parent
                   rotate_right(this,nd)
                }
                nd.parent.col = BLACK
                nd.parent.parent.col = RED
                rotate_left(this,nd.parent.parent)
            
            }
        
        }

    }

    this.root.col = BLACK
    
}


func insert(this *RB,nd *RB_NODE) int {

    var p *RB_NODE = this.root
    
    if this.cmp == nil{
        return 0
    }

    for nil != p {

        if this.cmp(nd.value,p.value) > 0 {
            
            if p.right == nil {
                p.right = nd
                nd.parent = p
                break
            } else {
                p = p.right
            }
            
        } else if this.cmp(nd.value,p.value) < 0 {
            
            if p.left == nil {
                p.left = nd
                nd.parent = p
                break
                
            }else {
                p = p.left
            }
        } else {
            this.adderr++
            return 2;
        }
    }
    
    return 1;

}


func (this *RB) insert_nd(val interface{}) int{  //指针接收器

    var nd *RB_NODE = new(RB_NODE)  //用户不用关心释放问题 且内存已清零
    
    if nil == nd {
        return 0
    }
    
    nd.value = val
    
    if this.cnt == 0 {
    
        this.root = nd      
        
    } else {
        var rt int = insert(this,nd) 
        if rt == 2 { 
            return 2        //节点已存在
        }
    }    

    
    this.cnt++
        
    //step: fixup 
    
    nd.col = RED
    
    insert_fixup(this,nd)
    
    return 1
}




func get(this *RB,v interface{}) (int, *RB_NODE){

    var p *RB_NODE = this.root

    for nil != p {

        if this.cmp(v,p.value) > 0 {    
        
                p = p.right
            
        } else if this.cmp(v,p.value) < 0 {
        
                p = p.left
        } else {
            
            return 1,p
        }
    }

    return 0,nil

}




func (this *RB) get_nd( v interface{}) (int,interface{}){

    var rt int
    var p *RB_NODE
  
    
    rt,p = get(this,v) 

    return rt,p.value

}


func successor (nd *RB_NODE) *RB_NODE{
    
    for nd = nd.right;nd.left != nil;{
        nd = nd.left
    }    
    return nd
}


func rm_fixup(this *RB,nd *RB_NODE){

    var bro *RB_NODE

    if this == nil {
    
        fmt.Printf("OOP! nil node doesnot to rm_fixup\n")
        return 
    }
    
    if nd == nil {
        fmt.Printf("INFO! nil node doesnot to rm_fixup\n")
        return 
    }
    
    for (nd != this.root && nd.col == BLACK ) {
        
        
        if (nd.parent == nil){
            fmt.Printf("value: %v;\n",nd.value)
            return
        }
        
        if nd.parent.left == nd {  
            bro = nd.parent.right
            
            if bro.col == RED {  //case 1：x的兄弟是红色的。 
                bro.col = BLACK
                nd.parent.col = RED
                rotate_left(this,nd.parent)
                bro = nd.parent.right;
                
            } else {    
            
                if (bro.left == nil || bro.left.col == BLACK) && (bro.right == nil || bro.right.col == BLACK){ //case 2：x的兄弟bro是黑色的
                   
                    bro.col = RED
                    nd = nd.parent
                }else {
                    if bro.right == nil || bro.right.col == BLACK{   //case 3：x的兄弟bro是黑色的，bro的右孩子是黑色（w的左孩子是红色）。 
                        bro.left.col = BLACK
                        bro.col = RED
                        rotate_right(this,bro)
                        bro = nd.parent.right
                    }
                    bro.col = nd.parent.col         // case 4 : x的兄弟bro是黑色的,bro的右孩子是红色
                    nd.parent.col = BLACK
                    bro.right.col = BLACK
                    rotate_left(this,nd.parent)
                    nd = this.root                    
                    
                }
            }
            
        } else {
        
            bro = nd.parent.left
            
            if bro.col == RED {  //case 1：x的兄弟是红色的。 
                bro.col = BLACK
                nd.parent.col = RED
                rotate_right(this,nd.parent)
                bro = nd.parent.left
                
            } else {    
            
                if (bro.left == nil || bro.left.col == BLACK) && (bro.right == nil || bro.right.col == BLACK) { //case 2：x的兄弟是黑色的，
                    
                    bro.col = RED
                    nd = nd.parent
                    
                }else {
                
                    if bro.left == nil || bro.left.col == BLACK{   //case 3：x的兄弟w是黑色的，w的右孩子是黑色（w的左孩子是红色）。 
                        bro.right.col = BLACK
                        bro.col = RED
                        rotate_left(this,bro)
                        bro = nd.parent.left
                    }
                    bro.col = nd.parent.col         // case 4 :x的兄弟w是黑色的，w的右孩子是红色（w的左孩子是任意色）
                    nd.parent.col = BLACK
                    bro.left.col = BLACK
                    rotate_right(this,nd.parent)
                    nd = this.root                    
                    
                }
            }
        
        }
    
    }
    
    nd.col = BLACK
}


//remove the node from the tree, and if the node is not in the tree,return good also
func (this *RB) remove_nd(v interface{}) int{

    var p *RB_NODE
    var nd *RB_NODE
    var chld *RB_NODE
    var nil_black_nd *RB_NODE
    
    _,p = get(this,v)  // p 需要删除的节点
    
    
    if p == nil {
        
        fmt.Printf("OOP! this node is not in the tree\n")
        this.delerr++;
        return 1 
    }
    
    nd = p  
    
    if p.left != nil && p.right != nil{ // 如果p有两个子节点，则后继节点为删除节点
        nd = successor(p)   //not-leaf must have a successor
        p.value = nd.value   // 颜色不变
    }

    
    //nd 是需要的释放的节点；nd 只可能是有一个子树 或 没有；nd 是黑色节点，则需要调节
    if nd.left == nil && nd.right == nil {//叶子节点

        if nd.parent == nil {//删除只有root节点的树    
            this.root = nil
        }else {
        
            chld = new(RB_NODE)
            chld.col = BLACK
            
            if nd == nd.parent.left {
                nd.parent.left = chld
                
            } else {
                nd.parent.right = chld
            }
            
            chld.parent = nd.parent
            nil_black_nd = chld
        }

    }else { // 一个子节点
        
        chld = nd.right
        
        if chld == nil {
            chld = nd.left
        }
        
        if nd.parent == nil {
            this.root = chld
        }else if nd == nd.parent.left {
            nd.parent.left = chld
        } else {
            nd.parent.right = chld
        }
        
        chld.parent = nd.parent

    }

    if nd.col == BLACK {
        rm_fixup(this,chld)
    }
    
    if nil_black_nd != nil {
        if nil_black_nd == nil_black_nd.parent.left {
            nil_black_nd.parent.left = nil
        } else {
            nil_black_nd.parent.right = nil
        }    
    }

    // nd waiting for GC。
    this.cnt--
    return 1
}

func walkthrough( nd *RB_NODE) int {

     if nd == nil{
        return 1
     }
     if nd.left != nil {
        walkthrough(nd.left)
     }
     fmt.Printf("value: %v\n",nd.value)
     
     if nd.right != nil {
        walkthrough(nd.right)
     }
    
    return 1
}

// go through all tree nodes
func (this *RB) walk_nd() int{
    
    walkthrough(this.root)
    fmt.Printf("+++total %v/ adderr %v / delerr %v+++\n",this.cnt,this.adderr,this.delerr)
    return 1
}


func check (nd *RB_NODE, d int, c int) bool {


    if nd == nil { 
        return true
    }
    
    if(nd.col == RED) && (nd.parent == nil || nd.parent.col == RED){
        return false
    }
    
    if nd.col== BLACK {
        c++
    }
    
    if nd.left == nil && nd.right == nil {
        
        if d == c {
            return true
        } else {
            return false
        }
    }
    
    return check(nd.left,d,c) && check(nd.right,d,c)
}

// check the tree is red-black or not 
func (this *RB) check_rbt() bool {

    var d int
    var nd *RB_NODE
    
    if (this == nil) || this.root == nil || (this.root != nil && this.root.col == RED) {
    
        return false
    }
    
    nd = this.root
    
    for nd != nil {
        
        if(nd.col == BLACK){
            d++
        }
        nd = nd.left
    }
    fmt.Printf("The black dep is  %v\n",d+1)
    return check(this.root,d,0)

}

/*==============RTtree user code as below ==================*/

func RBTCMP (x,y interface{}) int {

    var a,b, rt int
    
    a,OKa:= x.(int)           // 类型断言(类型转换)
    b,OKb:= y.(int)
    
    if (OKa !=  true || OKb !=  true)  {
        fmt.Printf("OOP! variable type is not matched\n")
    }
    
    if (a > b) {
        rt = 1
    }else if a < b {
        rt = -1
    } else {
        rt = 0
    }
    
    return rt;
}

// 数组添加
func RBTreeTest1(){

    var rbt *RB = RBTreeInit(RBTCMP)
    var rt int 
    
    if nil == rbt {
        fmt.Printf("OOP!RB tree init fail\n")
    }
    
    var arr1 = []int {563,2,1000,23,45,789,43,234,69,416,30,356,12,7,9,431}
    
    for _,key:=range arr1{

        rt = rbt.insert_nd(key)
        
        if (rt == EXISTED) {
            fmt.Printf("OOP!RB tree node existed (%v)\n",key)
        }    
    }

    rbt.walk_nd()

    
    
    var arr2 = []int {563,2,1000,23,45,789,43,234,69,416,30,356,12,7,9,431}
    
    for _,key:=range arr2{
        rt = rbt.remove_nd(key)
        
    }
    rbt.walk_nd()
    
    if rbt.check_rbt(){
        fmt.Printf("This is a RED BLACK tree\n")
    } else {
        fmt.Printf("OOP! It is not a RED BLACK tree\n")
    }
}

func RBTreeTest2(){

    var rbt *RB = RBTreeInit(RBTCMP)
    var i int
    var key int 
    var rt int 
    
    rand.Seed(time.Now().Unix())
    
    for i= 0; i< 1000;i++ {
    
        key = rand.Int()
        rt = rbt.insert_nd(key)
        
        if (rt == EXISTED) {
            fmt.Printf("OOP!RB tree node existed (%v)\n",key)
        }    
    
    }

    rbt.walk_nd()
    
    if rbt.check_rbt(){
        fmt.Printf("This is a RED BLACK tree\n")
    }
    
}

//随机添加
func RBTreeTest3(){

    var rbt *RB = RBTreeInit(RBTCMP)
    var i int
    var key int 
    var rt int 
    
    
    rand.Seed(time.Now().Unix())
    
    for i= 0; i< 500;i++ {
    
        key = rand.Intn(500)
        rt = rbt.insert_nd(key)
        
        if (rt == EXISTED) {
            fmt.Printf("OOP!RB tree node existed (%v)\n",key)
        }    
    }

    for i = 0; i<2000;i++ {
    
        key = rand.Intn(500)
        rt = rbt.remove_nd(key)
        
    }
    
    rbt.walk_nd()
    
    if rbt.check_rbt(){
        fmt.Printf("This is a RED BLACK tree\n")
    }
    
}

func main() {

    RBTreeTest1()
    RBTreeTest2()
    RBTreeTest3()
    

}
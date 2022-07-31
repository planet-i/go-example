package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mikespook/gorbac"
)

func init() {
	// 设置logger的输出选项
	// LstdFlags 标准logger的初始值；Logger类型表示一个活动状态的记录日志的对象
	// Lshortfile 文件无路径名+行号
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func LoadJson(filename string, v interface{}) error {
	// os.Open 打开一个文件用于读取
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	// 创建一个从filename读取并解码json对象的f,从输入流读取下一个json编码值并保存在v指向的值里
	return json.NewDecoder(f).Decode(v)
}

func SaveJson(filename string, v interface{}) error {
	// os.OpenFile使用指定的选项、指定的模式打开指定名称的文件
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// 将v的json编码写入f，并会写入一个换行符
	return json.NewEncoder(f).Encode(v)
}
func main() {
	// 获取一个RBAC实例
	rbac := gorbac.New()
	// 获取一个新角色
	rA := gorbac.NewStdRole("role-a")
	rB := gorbac.NewStdRole("role-b")
	rC := gorbac.NewStdRole("role-c")
	rD := gorbac.NewStdRole("role-d")
	rE := gorbac.NewStdRole("role-e")
	// 获取一些新权限
	pA := gorbac.NewStdPermission("permission-a")
	pB := gorbac.NewStdPermission("permission-b")
	pC := gorbac.NewStdPermission("permission-c")
	pD := gorbac.NewStdPermission("permission-d")
	pE := gorbac.NewStdPermission("permission-e")
	// 向角色添加权限
	rA.Assign(pA)
	rB.Assign(pB)
	rC.Assign(pC)
	rD.Assign(pD)
	rE.Assign(pE)
	// 将角色添加到实例
	rbac.Add(rA)
	rbac.Add(rB)
	rbac.Add(rC)
	rbac.Add(rD)
	rbac.Add(rE)
	// 设置角色继承，后者为父
	rbac.SetParent("role-a", "role-b")
	rbac.SetParents("role-b", []string{"role-c", "role-d"})
	rbac.SetParent("role-e", "role-d")
	// 判断角色是否具有某种权限
	if rbac.IsGranted("role-a", pA, nil) &&
		rbac.IsGranted("role-a", pB, nil) &&
		rbac.IsGranted("role-a", pC, nil) &&
		rbac.IsGranted("role-a", pD, nil) {
		fmt.Println("The role-a has been granted permis-a, b, c and d.")
	}
	rbac.SetParent("role-c", "role-a")
	// 判断实例中是否存在角色继承关系成环情况
	if err := gorbac.InherCircle(rbac); err != nil {
		fmt.Println("A circle inheratance occurred.")
	}
	roles := []string{"role-a", "role-b", "role-c", "role-d", "role-e"}
	// 所有角色都设置了此权限则返回true
	flag1 := gorbac.AllGranted(rbac, roles, pA, nil)
	// 只要大于等于一个角色有此权限则返回true
	flag2 := gorbac.AnyGranted(rbac, roles, pA, nil)
	fmt.Println(flag1)
	fmt.Println(flag2)
	// 从JSON文件中加载goRBAC实例以及如何将实例保存回去
	var jsonRoles map[string][]string
	var jsonInher map[string][]string
	if err := LoadJson("roles.json", &jsonRoles); err != nil {
		//Print err 然后调用 os.Exit(1)
		log.Fatal(err)
	}
	if err := LoadJson("inher.json", &jsonInher); err != nil {
		log.Fatal(err)
	}
	// 初始化实例
	rbac2 := gorbac.New()
	// 可以制作“角色”并将其添加到goRBAC实例中。
	permissions := make(gorbac.Permissions)
	for rid, pids := range jsonRoles {
		role := gorbac.NewStdRole(rid)
		for _, pid := range pids {
			_, ok := permissions[pid]
			if !ok {
				permissions[pid] = gorbac.NewStdPermission(pid)
			}
			role.Assign(permissions[pid])
		}
		rbac2.Add(role)
	}
	// 分配角色之间的继承关系
	for rid, parents := range jsonInher {
		if err := rbac2.SetParents(rid, parents); err != nil {
			log.Fatal(err)
		}
	}
	// 增加一个nobody角色并为它指定read-taxt权限
	nobody := gorbac.NewStdRole("nobody")
	permissions["read-text"] = gorbac.NewStdPermission("read-text")
	nobody.Assign(permissions["read-text"])
	rbac.Add(nobody)
	// 存储对json的更改
	jsonOutputRoles := make(map[string][]string)
	jsonOutputInher := make(map[string][]string)
	// 这里的闭包是一个WalkHandler类型，是用户定义的用于处理角色的函数、是Readonly进程，不实现则不访问goRBAC实例
	SaveJsonHandler := func(r gorbac.Role, parents []string) error {
		// WARNING: Don't use gorbac.RBAC instance in the handler,otherwise it causes deadlock.
		permissions := make([]string, 0)
		for _, p := range r.(*gorbac.StdRole).Permissions() {
			permissions = append(permissions, p.ID())
		}
		jsonOutputRoles[r.ID()] = permissions
		jsonOutputInher[r.ID()] = parents
		return nil
	}
	// gorbac.Walk 将rbac2中的每个角色传递给SaveJsonHandler
	if err := gorbac.Walk(rbac2, SaveJsonHandler); err != nil {
		log.Fatalln(err)
	}
	if err := SaveJson("new-roles.json", &jsonOutputRoles); err != nil {
		log.Fatal(err)
	}
	if err := SaveJson("new-inher.json", &jsonOutputInher); err != nil {
		log.Fatal(err)
	}

}

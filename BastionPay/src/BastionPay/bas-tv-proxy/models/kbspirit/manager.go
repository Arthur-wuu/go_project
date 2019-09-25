package kbspirit

import "fmt"

// 搜索器管理器
type Manager struct {
	searcherOn  map[string]Searcher // 打开的searcher
	searcherOff map[string]Searcher // 关闭的searcher
}

// NewManager 新建管理器
func NewManager() *Manager {
	return &Manager{
		searcherOn:  make(map[string]Searcher),
		searcherOff: make(map[string]Searcher),
	}
}

// Add 添加并激活搜索器
func (self *Manager) Add(name string, searcher Searcher) error {
	if nil == searcher {
		return fmt.Errorf("Nil searcher")
	}
	if _, ok := self.searcherOn[name]; ok {
		return fmt.Errorf("Searcher named \"%v\" exist,State->On", name)
	} else if _, ok = self.searcherOn[name]; ok {
		return fmt.Errorf("Searcher named \"%v\" exist,State->Off", name)
	}
	self.searcherOn[name] = searcher
	return nil
}

// Get 获取名称为name的搜索器，返回这个搜索器及其状态，不存在返回nil
func (self *Manager) Get(name string) (Searcher, bool) {
	if searcher, ok := self.searcherOn[name]; ok {
		return searcher, true
	} else if searcher, ok := self.searcherOff[name]; ok {
		return searcher, false
	}
	return nil, false
}

// Load 重新加载搜索器
func (self *Manager) Load(name string) Searcher {
	searcher := self.searcherOff[name]
	if nil != searcher {
		self.searcherOn[name] = searcher
		delete(self.searcherOff, name)
		return searcher
	} else {
		return nil
	}
}

// UnLoad 卸载搜索器
func (self *Manager) UnLoad(name string) Searcher {
	searcher := self.searcherOn[name]
	if nil != searcher {
		self.searcherOff[name] = searcher
		delete(self.searcherOn, name)
		return searcher
	} else {
		return nil
	}
}

// ForEach 遍历搜索器,f为回调函数
// f 函数的参数name代表搜索器名称，searcher当前遍历到的搜索器，on标识当前的搜索起是否处于开的状态
// f 函数的返回值表示是否继续遍历
func (self *Manager) ForEach(f func(name string, searcher Searcher, on bool) bool) {
	// 处于开状态的Searcher
	for name, searcher := range self.searcherOn {
		if !f(name, searcher, true) {
			return
		}
	}

	// 处于关状态的Searcher
	for name, searcher := range self.searcherOff {
		if !f(name, searcher, false) {
			return
		}
	}
}

package ahocorasick

import (
	"testing"
)

var (
	zhSensitiveWords = []string{"暗杀", "罢工", "罢课", "暴乱", "暴政", "出售假钞", "出售枪支", "出售手枪",
		"弹劾", "党禁", "党内分裂", "颠覆中国政权", "颠覆中华人民共和国政", "东北独立",
		"东突", "东土耳其斯坦", "独裁", "独裁政治", "反党", "反共", "反华", "反社会",
		"反政府", "仿真假钞", "废统", "国研新闻邮件", "国贼", "海外护法", "和平修炼",
		"红色恐怖", "回民暴动", "集体上访", "疆独", "警察殴打", "绝食抗暴", "开天目",
		"两岸三地论坛", "两个中国", "卖国", "蒙独", "蒙古独立", "民族矛盾", "全新假钞",
		"上海帮", "省委大门集合", "省政府大门集合", "示威", "事实独立", "台独", "台湾18DY电影",
		"台湾独立", "台湾狗", "台湾建国运动组织", "台湾青年独立联盟", "台湾政论区", "台湾自由联盟",
		"天安门录影带", "天安门母亲", "天安门事件", "天安门屠杀", "天安门一代", "天府广场集会",
		"新观察论坛", "新疆独立", "新生网", "新唐人", "新闻封锁", "新语丝", "找政府评理",
		"政府软弱", "政府无能", "支那", "中共政权", "中共小丑", "中国社会的艾滋病", "中国威胁论",
		"中国真实内容", "中国之春", "中国猪", "中國當局", "中华联邦政府", "中华民国", "中华人民实话实说",
		"中华人民正邪", "中华养生益智功", "中华真实报道", "专政机器", "转法轮", "自由亚洲", "自由运动",
		"宗教压迫", "阻止中华人民共和国统", "保*钓*抗*日", "迫害法轮功", "疆 独", "法 轮 功", "真&善&忍",
		"疆.独", "买枪", "教徒", "暴行", "枪模", "涉日", "反法游行", "反法示威", "抵制奥运", "抵制游行",
		"抵制示威", "抵制邪恶", "拥护台独", "敌对分子", "殉道圣人", "法轮大法", "游行示威", "爆炸装置",
		"狱中举报", "示威游行", "管制刀具", "网络封锁", "东突厥斯坦", "中共当权者", "伊斯兰运动",
		"反共游击队", "中共特权阶层", "中共统治集团", "网络活动颠覆", "雪灾西藏地震", "东突厥斯坦伊斯兰",
		"阿扁推翻", "藏獨", "独立台湾", "对日强硬", "法车仑", "法伦功", "法轮佛", "分裂中华人民共和国",
		"分裂中国", "颠覆中华人民共和国", "颠覆国家", "香港独立", "港独"}

	enWords = []string{"his", "hers", "he", "she", "her", "jordan", "kobe", "lebron"}
)

func TestSearch(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		t.Parallel()
		ac := NewMatcher()

		ac.BuildWithPatterns(zhSensitiveWords)
		m1 := ac.Search("你这个反社会分子，我要没收你的管制刀具！")
		if m1[0] != "反社会" || m1[1] != "管制刀具" {
			t.Fatalf("expected `反社会`, `管制刀具`, but got: %s, %s", m1[0], m1[1])
		}

		for _, w := range enWords {
			ac.AddPattern(w)
		}
		ac.Build()

		m2 := ac.Search("shers got jordan")
		if m2[0] != "she" || m2[1] != "he" || m2[2] != "her" || m2[3] != "hers" || m2[4] != "jordan" {
			t.Fatalf("expected `she`, `he`, `her`, `hers`, `jordan`, but got: %v", m2)
		}
	})

	t.Run("corner case", func(t *testing.T) {
		t.Parallel()
		ac := NewMatcher()
		ac.BuildWithPatterns(nil)
		m := ac.Search("foo bar baz")
		if m != nil {
			t.Fatalf("you should got nothing, but got: %v", m)
		}
		if ac.Match("aho corasick") {
			t.Fatalf("should not matched")
		}

		ac.AddPattern("foo")
		ac.AddPattern("foo")
		ac.BuildWithPatterns([]string{" "})
		m = ac.Search("foo bar baz")
		if len(m) != 3 || m[0] != "foo" || m[1] != " " || m[2] != " " {
			t.Fatalf("expected got `foo`, but got: %v", m)
		}

		s1 := "ahishershe"
		ac1 := NewMatcher()
		ac1.BuildWithPatterns([]string{"his", "hers", "he", "she"})
		m1 := ac1.Search(s1)
		if m1[0] != "his" || m1[1] != "she" || m1[2] != "he" || m1[3] != "hers" || m1[4] != "she" || m1[5] != "he" {
			t.Fatalf("expected `his`, `she`, `he`, `hers`, `she`, `he`, but got: %v", m)
		}

		si := ac1.SearchIndexed(s1)
		if len(m1) != len(si) {
			t.Fatalf("Search and SearchIndexed is not consistent for data: %s", s1)
		}

		chars := []rune(s1)
		for i := 0; i < len(m1); i++ {
			start := si[i].Start
			end := si[i].Start + si[i].Len
			x := string(chars[start:end])
			if x != m1[i] {
				t.Fatalf("the %dth matched of SearchIndexed(%s) and Search(%s) are not equal", i+1, x, m1[i])
			}
		}

		ac3 := NewMatcher()
		ac3.BuildWithPatterns([]string{"1", "21", "321", "4321", "54321", "数字9", "987654321"})
		s3 := "数字987654321"
		chars3 := []rune(s3)
		exp3 := []string{"数字9", "987654321", "54321", "4321", "321", "21", "1"}

		m3 := ac3.Search(s3)
		si3 := ac3.SearchIndexed(s3)

		if len(m3) != len(si3) {
			t.Fatalf("Search and SearchIndexed is not consistent for data: %s", s3)
		}

		for i := 0; i < len(m3); i++ {
			start := si3[i].Start
			end := si3[i].Start + si3[i].Len
			x := string(chars3[start:end])
			if x != m3[i] || x != exp3[i] {
				t.Fatalf("%dth matched of SearchIndexed(%s) and Search(%s) are not equal, and expected: %s", i+1, x, m3[i], exp3[i])
			}
		}
	})
}

func TestSearchIndexed(t *testing.T) {
	ac := NewMatcher()

	ac.BuildWithPatterns(zhSensitiveWords)
	s := "你这个反社会分子，我要没收你的管制刀具！"
	m := ac.SearchIndexed(s)
	if len(m) != 2 {
		t.Fatalf("expected two word matched, but got: %d", len(m))
	}

	chars := []rune(s)
	matched := make([]string, 0, len(m))
	for _, hit := range m {
		s := string(chars[hit.Start:(hit.Start + hit.Len)])
		matched = append(matched, s)
	}

	if matched[0] != "反社会" || matched[1] != "管制刀具" {
		t.Fatalf("expected `反社会`, `管制刀具`, but got: %s, %s", matched[0], matched[1])
	}
}

func TestMatch(t *testing.T) {
	ac := NewMatcher()
	ac.BuildWithPatterns(zhSensitiveWords)
	cases := []struct {
		q string
		m bool
	}{
		{"独裁", true},
		{"罢工", true},
		{"共", false},
		{"喜提", false},
		{"shit", false},
		{"港独分子", true},
	}

	for _, c := range cases {
		if ac.Match(c.q) != c.m {
			t.Fatalf("expected matched result: %t for %s", c.m, c.q)
		}
	}

	m := NewMatcher()
	m.BuildWithPatterns([]string{"俄罗斯", "war", "Ukraine", "😭", "こんにちは", "¿puedes", "침략"})
	if m.Match("2022年2月24日开始，俄白联盟以“非军事化、去纳粹化”为由，大规模入侵乌克兰") {
		t.Fatal("should not matched")
	}

	ac1 := NewMatcher()
	ac1.BuildWithPatterns([]string{"foo", "bar", "baz"})
	if !ac1.Match("xxxxxxx你好 yyyyybar") {
		t.Fatalf("expected matched result: %s for %s", "xxxxxxx你好 yyyyybar", "bar")
	}

	ac2 := NewMatcher()
	ac2.BuildWithPatterns([]string{"国人", "中国人", "新中国"})
	exp := []string{"新中国", "中国人", "国人"}
	ret := ac2.Search("我是新中国人")
	if ret[0] != exp[0] || ret[1] != exp[1] || ret[2] != exp[2] {
		t.Fatalf("expected matched result: %v for %s, but got %+v", exp, "我是新中国人", ret)
	}
}

func TestNotBuild(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic with msg: `you should use `Build() or BuildWithPatterns()` before searching`")
		}
	}()

	ac := NewMatcher()
	ac.AddPattern("foo")
	ac.AddPattern("bar")
	// ac.Build()
	// ac.BuildWithPatterns(nil)
	_ = ac.Search("foo bar baz")
}

func BenchmarkAC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ac := NewMatcher()
		ac.BuildWithPatterns(zhSensitiveWords)
		_ = ac.Search("你这个反社会分子，我要没收你的管制刀具！")
	}
}

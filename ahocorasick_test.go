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

// TestSearchSuffixPatterns 验证 AC 算法能否正确匹配同一个位置的所有重叠 pattern。
// 验证修复：直接子节点匹配时，需遍历 n 的 fail 链收集所有以当前位置结尾的 pattern。
func TestSearchSuffixPatterns(t *testing.T) {
	t.Run("suffix pattern on fail chain", func(t *testing.T) {
		// "ab" 匹配后，其 fail 链上的 "b"（ab 后缀）也应匹配。
		// 即使下一个字符 'c' 是 ab 的直接子节点（abc 模式），
		// 也应通过 fail 链找到 "b"。
		ac := NewMatcher()
		ac.BuildWithPatterns([]string{"ab", "b", "abc"})
		ret := ac.Search("abc")
		if len(ret) != 3 {
			t.Fatalf("expected 3 matches, but got %d: %v", len(ret), ret)
		}
		// 验证所有 pattern 都找到了（顺序不重要）
		patternSet := map[string]bool{"ab": true, "b": true, "abc": true}
		for _, s := range ret {
			if !patternSet[s] {
				t.Fatalf("unexpected match: %s", s)
			}
		}

		si := ac.SearchIndexed("abc")
		if len(si) != 3 {
			t.Fatalf("SearchIndexed: expected 3 hits, but got %d: %v", len(si), si)
		}
	})

	t.Run("multiple suffix patterns lost in chain", func(t *testing.T) {
		// "abc" 匹配后，fail 链上有 "bc" 和 "c"。
		// 下一字符 'd' 是 abc 的直接子节点（abcd 模式），
		// 导致 "bc" 和 "c" 都丢失。
		ac := NewMatcher()
		ac.BuildWithPatterns([]string{"abc", "bc", "c", "abcd"})
		ret := ac.Search("abcd")
		if len(ret) != 4 {
			t.Fatalf("expected 4 matches [abc bc c abcd], but got %d: %v", len(ret), ret)
		}
		// 注意：顺序取决于实现，但数量必须正确
		// 验证每个结果的内容都是有效 pattern
		patternSet := map[string]bool{"abc": true, "bc": true, "c": true, "abcd": true}
		for _, s := range ret {
			if !patternSet[s] {
				t.Fatalf("unexpected match: %s", s)
			}
		}

		si := ac.SearchIndexed("abcd")
		if len(si) != 4 {
			t.Fatalf("SearchIndexed: expected 4 hits, but got %d: %v", len(si), si)
		}
	})

	t.Run("overlapping short patterns at every position", func(t *testing.T) {
		// "aaaa" 中每个位置都应该匹配 "a"，
		// 同时还有 "aa"/"aaa"/"aaaa" 等重叠 pattern。
		ac := NewMatcher()
		ac.BuildWithPatterns([]string{"a", "aa", "aaa", "aaaa"})
		ret := ac.Search("aaaa")
		// 预期: a@0, aa@0, aaa@0, aaaa@0, a@1, aa@1, aaa@1, a@2, aa@2, a@3
		// 总共 10 个匹配
		if len(ret) != 10 {
			t.Fatalf("expected 10 matches for 'aaaa', but got %d: %v", len(ret), ret)
		}

		// SearchIndexed 验证
		si := ac.SearchIndexed("aaaa")
		if len(si) != 10 {
			t.Fatalf("SearchIndexed: expected 10 hits, but got %d: %v", len(si), si)
		}

		// 按位置验证每个 start 位置出现的匹配数
		posCount := make(map[int]int)
		for _, h := range si {
			posCount[h.Start]++
		}
		// 位置0: a, aa, aaa, aaaa → 4
		// 位置1: a, aa, aaa          → 3
		// 位置2: a, aa               → 2
		// 位置3: a                   → 1
		expectedPosCount := map[int]int{0: 4, 1: 3, 2: 2, 3: 1}
		for pos, cnt := range expectedPosCount {
			if posCount[pos] != cnt {
				t.Fatalf("position %d: expected %d matches, got %d", pos, cnt, posCount[pos])
			}
		}
	})

	t.Run("suffix pattern found only via fail chain (not at end)", func(t *testing.T) {
		// "ab" 的 fail 指向 "b"，但 "b" 一直要到下一个不匹配的字符（或尾随检查）才被找到。
		// 即使行尾检查能补救，但顺序可能不对，且 SearchIndexed 应给出正确的 start 位置。
		ac := NewMatcher()
		ac.BuildWithPatterns([]string{"ab", "b"})

		// 不带后续匹配：尾随检查能补救
		ret := ac.Search("ab")
		// 预期: ab@0, b@1
		if len(ret) != 2 {
			t.Fatalf("expected 2 matches [ab b], got %d: %v", len(ret), ret)
		}

		// 带后续不匹配字符：fail 链会在下一字符被遍历
		ret2 := ac.Search("abx")
		if len(ret2) != 2 {
			t.Fatalf("expected 2 matches [ab b], got %d: %v", len(ret2), ret2)
		}
	})

	t.Run("deep fail chain in unicode text", func(t *testing.T) {
		// 中文场景：重叠 pattern
		ac := NewMatcher()
		ac.BuildWithPatterns([]string{"中国人", "国人", "人", "中国"})

		ret := ac.Search("中国人")
		if len(ret) != 4 {
			t.Fatalf("expected 4 matches [中国人 国人 人 中国] or similar, got %d: %v", len(ret), ret)
		}

		si := ac.SearchIndexed("中国人")
		if len(si) != 4 {
			t.Fatalf("SearchIndexed: expected 4 hits, got %d: %v", len(si), si)
		}

		chars := []rune("中国人")
		for _, h := range si {
			if h.Start < 0 || h.Start+h.Len > len(chars) {
				t.Fatalf("invalid hit range: Start=%d Len=%d", h.Start, h.Len)
			}
		}
	})

	t.Run("Match also misses suffix patterns", func(t *testing.T) {
		// Match 方法也有同样的 bug：
		// 当 "abc" 匹配时，应返回 true，不管是否遍历了 fail 链
		ac := NewMatcher()
		ac.BuildWithPatterns([]string{"ab", "b"})
		if !ac.Match("ab") {
			t.Fatal("Match('ab') should be true")
		}

		// 更复杂的场景：
		ac2 := NewMatcher()
		ac2.BuildWithPatterns([]string{"abc", "bc", "c", "abcd"})
		if !ac2.Match("abcd") {
			t.Fatal("Match('abcd') should be true")
		}
	})
}

func BenchmarkAC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ac := NewMatcher()
		ac.BuildWithPatterns(zhSensitiveWords)
		_ = ac.Search("你这个反社会分子，我要没收你的管制刀具！")
	}
}

package chrgg

func calcFee(lastV, currV int64, unitFen int64) (
	valuePer100 int64, /* 表差 1/100 计量单位*/
	feeFen int64, /* 1 计量单位多少分 */
) {
	// 注意这里的表差是 1/100 计量单位
	// 为了规避误差，这里先做乘法再做除法
	// 最后除以100是把 1/100 计量单位，转成 1 计量单位
	// 不是分转元! 所以最后得出的结果是 1 个计量单位多少分
	valuePer100 = currV - lastV
	feeFen = (valuePer100 * unitFen) / 100
	return
}

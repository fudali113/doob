package doob

func (this *urlInfo) addUrlPara(v urlMacthPara) {
	this.urlParas = append(this.urlParas, v)
}

func (this *urlInfo) len() int {
	return len(this.urlParas)
}

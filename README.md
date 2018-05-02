# nlp-corpus-tools
NLP语料库处理工具

该语料库处理工具以csv格式为基础。

## Todo list

- [ ] field: 字段操作
  - [ ] rename: 字段重命名
  - [ ] delete: 删除字段
  - [x] keep: 只保留若干字段
  - [ ] 
- [ ] text: 文本操作相关
  - [x] replace: 字符串替换
  - [x] filter-space: 过滤掉文本中的空格，制表符等
  - [x] filter-newline-char: 过滤掉文本中的换行符
  - [x] iconv: 编码转换
  - [x] dbc2sbc: 全角转半角
  - [x] trim: 替换字符串前后的空格字符
- [ ] filter: 过滤，例如满足某些条件的文本应该去掉
  - [x] length: 按字符串的长度进行过滤
  - [x] null: 过滤掉空字符串的记录
  - [x] in: 如果字段值等于某些值，则过滤
- [ ] corpus: 语料库
  - [x] batch: 把一个大文件分拆成小的批量文件，可以方便进行标注
  - [x] merge: 例如将一个若干小的批量文件合并成大文件，例如训练集
  - [x] split: 将文件按比例拆分成训练集和测试集
  - [x] stat: 统计csv文件的记录数, 字段数，计算文件编码等
  - [ ] 
- [ ] check: 标注结果是否正确，例如标注预测的结果是否正确
  - [ ] succ: 判断预测结果是否正确
  - [ ] 

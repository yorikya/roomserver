package devices

import (
	"log"
	"strconv"
	"strings"
	"sort"

	"github.com/smira/go-statsd"
)


var (
	airCoolCodes = []IRCode{
		IRCode{
			"CLEAN",
			[]int{9036, 4452, 672, 1612, 616, 1668, 616, 524, 624, 496, 648, 492, 592, 524, 644, 1636, 652, 1660, 588, 1696, 620, 1664, 620, 1664, 648, 492, 676, 440, 648, 496, 616, 496, 648, 496, 620, 1664, 620, 520, 648, 1636, 652, 464, 676, 1636, 620, 492, 620, 524, 648, 492, 596, 520, 620, 1664, 648, 1636, 620, 1664, 648, 496, 616, 496, 620, 524, 564, 580, 612, 500, 648, 496, 616, 500, 648, 492, 588, 528, 648, 1636, 672, 468, 648, 1664, 620, 496, 648, 492, 592, 524, 648, 468, 620, 520, 644, 472, 592, 548, 620, 524, 644, 472, 676, 464, 564, 552, 648, 496, 616, 496, 620, 524, 620, 496, 644, 496, 648, 496, 620, 492, 648, 496, 620, 496, 648, 492, 676, 440, 644, 500, 564, 576, 620, 496, 672, 468, 592, 524, 676, 440, 648, 492, 588, 528, 648, 496, 620, 520, 644, 472, 620, 524, 644, 468, 620, 524, 564, 552, 644, 1640, 644, 496, 620, 520, 648, 468, 652, 492, 560, 556, 648, 492, 620, 496, 648, 496, 644, 468, 620, 524, 648, 1636, 672, 468, 616, 1668, 644, 500, 536, 576, 644, 472, 620, 524, 648, 492, 620, 496, 620, 1688, 620, 496, 620, 524, 620, 1664, 620, 1664, 644, 472, 644, 1664, 652},
		},
		IRCode{
			"COND",
			[]int{9120, 4368, 648, 1636, 664, 1648, 676, 436, 704, 440, 648, 468, 644, 472, 676, 1632, 640, 1644, 648, 1636, 648, 1664, 648, 1636, 616, 496, 648, 496, 620, 496, 644, 496, 556, 1728, 648, 1636, 676, 464, 672, 1612, 676, 468, 664, 1620, 620, 496, 672, 468, 676, 464, 648, 468, 648, 496, 648, 468, 672, 468, 624, 1660, 648, 468, 676, 468, 616, 524, 652, 464, 644, 496, 648, 468, 648, 492, 652, 464, 644, 1640, 704, 440, 620, 1664, 672, 468, 620, 496, 676, 468, 648, 464, 676, 468, 676, 440, 648, 492, 676, 468, 648, 464, 676, 468, 632, 484, 648, 492, 624, 492, 672, 444, 676, 1632, 620, 524, 624, 492, 672, 468, 620, 496, 700, 444, 620, 492, 704, 412, 648, 496, 704, 436, 644, 472, 704, 436, 644, 472, 708, 436, 612, 500, 704, 440, 704, 412, 676, 464, 676, 468, 616, 500, 700, 440, 676, 440, 676, 464, 652, 1632, 676, 440, 676, 464, 652, 492, 644, 472, 676, 468, 616, 496, 704, 440, 648, 468, 676, 464, 672, 468, 708, 1576, 732, 384, 676, 1608, 704, 436, 652, 464, 648, 496, 644, 468, 708, 436, 676, 464, 620, 496, 680, 1604, 704, 440, 552, 1728, 704, 1580, 708, 1580, 704, 436, 700},
		},
		IRCode{
			"OFF",
			[]int{9088, 4452, 644, 1640, 672, 1612, 620, 520, 620, 496, 620, 520, 652, 464, 732, 1552, 680, 1628, 648, 1636, 648, 1636, 708, 1576, 620, 524, 672, 444, 672, 468, 612, 1672, 676, 464, 700, 1584, 648, 1636, 648, 1636, 676, 464, 612, 1672, 648, 468, 648, 496, 648, 492, 672, 444, 676, 1636, 620, 492, 644, 1640, 620, 1664, 672, 472, 648, 464, 680, 464, 676, 464, 704, 412, 648, 492, 652, 464, 652, 492, 620, 1664, 700, 416, 644, 496, 648, 496, 700, 412, 648, 496, 672, 444, 676, 464, 652, 464, 644, 496, 648, 496, 620, 496, 728, 412, 676, 440, 644, 496, 652, 464, 700, 1584, 704, 436, 640, 504, 704, 412, 732, 408, 648, 468, 676, 440, 644, 496, 648, 468, 648, 496, 672, 468, 644, 472, 704, 436, 640, 476, 676, 464, 644, 472, 620, 524, 648, 464, 648, 496, 676, 464, 652, 464, 680, 464, 620, 496, 672, 468, 648, 468, 704, 436, 676, 468, 676, 436, 680, 464, 648, 468, 644, 496, 648, 468, 672, 444, 676, 464, 676, 468, 672, 1612, 616, 524, 736, 1548, 652, 464, 648, 492, 624, 492, 676, 464, 648, 476, 696, 440, 672, 468, 648, 468, 672, 468, 652, 464, 700, 416, 676, 468, 704, 1604, 648}},
		
		IRCode{
			"HEAT_16",
			[]int{9116, 4392, 680, 1604, 704, 1580, 676, 464, 676, 440, 648, 496, 708, 404, 708, 1576, 736, 1576, 676, 1608, 708, 1576, 676, 1608, 708, 432, 588, 528, 680, 464, 640, 1640, 684, 460, 704, 1580, 680, 460, 680, 1604, 684, 432, 736, 1548, 708, 436, 648, 464, 708, 436, 732, 408, 676, 1608, 676, 440, 620, 1692, 680, 432, 732, 388, 704, 436, 676, 464, 700, 416, 672, 472, 728, 384, 704, 440, 588, 528, 672, 1608, 616, 1696, 588, 556, 660, 456, 704, 436, 648, 468, 676, 468, 616, 496, 704, 412, 676, 468, 700, 468, 648, 440, 700, 440, 704, 412, 732, 412, 672, 468, 648, 492, 592, 524, 708, 1604, 680, 408, 644, 524, 676, 440, 644, 496, 648, 468, 676, 468, 588, 528, 672, 468, 648, 492, 560, 556, 648, 496, 688, 428, 644, 496, 588, 528, 616, 524, 592, 552, 672, 440, 616, 528, 592, 524, 648, 492, 592, 524, 704, 1580, 644, 500, 588, 552, 620, 496, 668, 472, 592, 524, 620, 524, 588, 524, 648, 468, 616, 528, 644, 496, 676, 1608, 648, 496, 528, 1752, 684, 436, 616, 524, 620, 496, 616, 524, 564, 580, 648, 468, 676, 1604, 680, 1604, 648, 1664, 616, 500, 680, 1580, 728, 412, 700, 440, 652},
		},
		IRCode{
			"HEAT_17",
			[]int{8948, 4536, 564, 1720, 532, 1752, 560, 560, 612, 528, 532, 584, 560, 580, 532, 1752, 560, 1748, 564, 1720, 532, 1752, 560, 1724, 560, 1724, 560, 584, 560, 552, 564, 1720, 564, 580, 560, 1748, 536, 580, 532, 1752, 560, 584, 532, 1752, 560, 552, 564, 580, 532, 608, 536, 1748, 560, 556, 560, 584, 556, 1728, 532, 580, 564, 580, 560, 556, 560, 580, 560, 588, 556, 556, 560, 580, 560, 556, 560, 580, 536, 1748, 532, 1752, 560, 584, 556, 584, 536, 580, 532, 584, 560, 580, 532, 584, 560, 584, 556, 556, 564, 580, 560, 580, 560, 556, 560, 584, 560, 552, 564, 580, 536, 580, 560, 580, 532, 1752, 588, 556, 556, 560, 612, 528, 612, 504, 640, 500, 592, 524, 560, 584, 528, 612, 560, 556, 616, 524, 616, 500, 616, 528, 532, 580, 560, 584, 536, 580, 560, 580, 560, 584, 532, 584, 556, 556, 620, 524, 532, 584, 616, 1692, 592, 524, 644, 500, 644, 496, 620, 496, 532, 584, 560, 580, 532, 584, 560, 580, 560, 556, 592, 552, 616, 1696, 588, 524, 560, 1724, 616, 524, 592, 524, 564, 580, 588, 528, 604, 536, 560, 1724, 616, 528, 588, 1692, 564, 552, 564, 1720, 560, 1752, 532, 584, 560, 580, 536},
		},
		IRCode{
			"HEAT_18",
			[]int{9004, 4480, 620, 1664, 592, 1692, 560, 556, 648, 496, 616, 496, 644, 500, 564, 1720, 644, 1664, 540, 1744, 564, 1720, 616, 1668, 616, 528, 564, 1716, 592, 524, 564, 1720, 648, 496, 644, 1640, 564, 576, 616, 1668, 588, 528, 644, 1664, 620, 496, 592, 552, 564, 576, 540, 1744, 564, 1720, 592, 1692, 644, 472, 564, 576, 564, 552, 620, 524, 616, 524, 564, 552, 620, 524, 560, 552, 564, 580, 560, 556, 644, 1640, 644, 1664, 560, 584, 560, 552, 620, 524, 592, 524, 560, 580, 536, 580, 560, 556, 616, 524, 564, 580, 556, 560, 560, 580, 560, 556, 564, 576, 564, 552, 564, 580, 560, 552, 564, 1748, 532, 584, 560, 580, 532, 584, 564, 580, 560, 552, 564, 580, 560, 560, 556, 580, 560, 584, 560, 556, 560, 580, 560, 556, 560, 580, 536, 580, 560, 584, 560, 580, 536, 580, 560, 580, 536, 580, 560, 556, 560, 580, 532, 1752, 564, 580, 560, 580, 536, 580, 560, 584, 532, 584, 556, 556, 564, 580, 556, 560, 560, 580, 560, 584, 528, 1752, 564, 580, 536, 1748, 536, 580, 560, 580, 536, 580, 560, 580, 564, 580, 536, 1748, 560, 1724, 532, 584, 560, 1748, 536, 1748, 536, 1748, 560, 556, 560, 580, 532},
		},
		IRCode{
			"HEAT_19",
			[]int{9084, 4448, 680, 1604, 648, 1640, 640, 500, 676, 440, 676, 464, 560, 556, 680, 1604, 648, 1660, 648, 1636, 648, 1636, 644, 1668, 652, 1632, 600, 1684, 560, 552, 652, 1632, 680, 464, 648, 1664, 620, 492, 676, 1608, 652, 492, 620, 1664, 588, 528, 672, 468, 680, 460, 652, 464, 620, 1664, 620, 1692, 676, 440, 648, 468, 644, 496, 680, 436, 644, 496, 672, 468, 652, 440, 672, 468, 672, 472, 676, 464, 644, 1640, 680, 1604, 644, 500, 616, 524, 648, 468, 648, 468, 676, 464, 648, 468, 648, 448, 696, 440, 644, 524, 644, 496, 652, 464, 648, 468, 672, 468, 676, 472, 556, 556, 648, 492, 652, 1632, 648, 492, 680, 440, 648, 492, 588, 500, 648, 492, 592, 524, 700, 472, 700, 412, 620, 496, 680, 488, 684, 404, 676, 468, 648, 468, 644, 496, 648, 468, 672, 468, 704, 440, 648, 468, 680, 460, 648, 468, 700, 416, 680, 1632, 676, 436, 704, 436, 676, 468, 680, 436, 676, 464, 624, 492, 704, 440, 728, 388, 704, 408, 680, 492, 676, 1608, 648, 464, 676, 1608, 708, 436, 712, 404, 704, 436, 624, 492, 704, 436, 680, 464, 680, 1604, 680, 436, 704, 436, 732, 384, 708, 436, 648, 1636, 672, 468, 680},
		},
		IRCode{
			"HEAT_20",
			[]int{9060, 4456, 616, 1668, 616, 1664, 620, 524, 560, 556, 616, 524, 560, 556, 644, 1640, 588, 1720, 568, 1716, 564, 1720, 644, 1668, 536, 580, 560, 580, 592, 1692, 620, 1664, 564, 580, 556, 1728, 556, 556, 568, 1744, 536, 580, 560, 1724, 560, 580, 536, 580, 564, 576, 564, 580, 536, 580, 560, 1724, 560, 580, 560, 556, 560, 580, 536, 580, 564, 580, 560, 580, 536, 580, 560, 580, 536, 580, 560, 584, 536, 1748, 536, 1748, 560, 580, 560, 556, 560, 580, 560, 556, 560, 580, 564, 532, 584, 580, 536, 576, 564, 580, 564, 580, 532, 580, 564, 580, 536, 580, 560, 580, 536, 580, 556, 560, 588, 1720, 564, 580, 536, 580, 560, 580, 536, 580, 560, 584, 532, 580, 560, 556, 588, 580, 536, 580, 556, 560, 560, 580, 532, 584, 560, 584, 560, 552, 564, 580, 556, 584, 560, 556, 560, 580, 560, 556, 564, 580, 560, 556, 560, 1724, 560, 580, 560, 580, 532, 584, 564, 584, 556, 556, 560, 580, 560, 556, 560, 580, 536, 580, 560, 584, 560, 1720, 564, 580, 560, 1724, 532, 584, 560, 580, 532, 584, 560, 584, 560, 580, 560, 556, 560, 580, 560, 556, 560, 1724, 560, 584, 532, 580, 564, 1720, 588, 556, 560},
		},
		IRCode{
			"HEAT_21",
			[]int{9140, 4368, 708, 1580, 676, 1632, 680, 436, 676, 464, 680, 440, 640, 472, 676, 1636, 552, 1732, 672, 1612, 676, 1632, 680, 1604, 648, 1636, 672, 444, 672, 1612, 700, 1608, 704, 440, 648, 1636, 672, 440, 676, 1608, 676, 468, 608, 1676, 640, 476, 640, 500, 676, 464, 672, 1612, 676, 1608, 644, 500, 588, 528, 612, 528, 592, 524, 616, 524, 568, 576, 588, 528, 560, 580, 592, 524, 616, 524, 652, 464, 616, 1668, 648, 1636, 620, 524, 648, 492, 564, 552, 648, 492, 624, 492, 652, 492, 620, 492, 676, 468, 592, 548, 628, 488, 652, 492, 652, 460, 648, 468, 652, 492, 588, 524, 680, 464, 680, 1632, 624, 488, 680, 464, 624, 492, 672, 444, 648, 492, 560, 556, 676, 464, 680, 464, 532, 580, 652, 492, 532, 584, 648, 492, 560, 556, 644, 496, 564, 552, 648, 496, 676, 464, 560, 556, 616, 528, 560, 556, 640, 500, 592, 1692, 560, 556, 612, 528, 588, 552, 564, 552, 588, 556, 588, 528, 560, 580, 532, 584, 560, 580, 560, 584, 552, 1732, 560, 552, 564, 1720, 616, 528, 532, 584, 560, 580, 560, 556, 560, 580, 560, 1724, 560, 1752, 532, 1752, 532, 1752, 560, 556, 556, 584, 560, 1724, 584, 556, 536},
		},
		IRCode{
			"HEAT_22",
			[]int{9116,4396,620,1664,648,1636,668,448,672,468,648,468,648,496,644,1640,612,1696,620,1664,620,1664,580,1704,700,416,672,1636,620,1664,564,1720,648,492,612,1672,644,500,648,1636,672,444,644,1636,676,468,644,472,644,496,676,468,676,436,644,500,592,524,644,496,620,496,700,440,648,496,672,444,640,500,620,496,644,472,644,496,668,1616,644,1640,676,468,644,496,648,468,700,440,648,468,672,472,620,492,676,468,644,496,648,468,672,472,648,464,648,468,648,496,644,468,648,496,644,1664,648,468,672,472,620,496,612,500,616,528,640,476,672,468,644,500,640,472,676,468,640,476,644,496,648,468,648,492,648,468,644,500,672,468,648,468,616,524,624,492,616,528,616,1664,676,440,644,500,616,524,648,468,644,500,588,524,676,468,676,440,644,496,648,496,616,1668,620,492,648,1636,644,500,672,444,644,496,648,468,672,468,564,580,620,496,644,1636,620,524,672,1612,672,444,644,1664,648,496,592},
		},
		IRCode{
			"HEAT_23",
			[]int{9120,4392,624,1660,552,1732,640,476,708,432,696,420,704,436,556,1728,708,1604,680,1604,540,1744,644,1636,680,1608,676,1608,644,1636,704,1608,644,496,656,460,704,440,676,1608,676,436,652,1632,676,468,672,444,676,464,680,460,652,464,676,468,676,1608,608,1676,644,1640,704,436,616,528,652,464,672,468,676,440,676,464,648,468,676,1608,688,1596,680,460,676,468,648,468,704,436,648,468,704,436,676,440,676,468,672,468,704,412,676,468,648,464,676,468,648,468,644,468,708,436,676,1632,652,464,704,440,648,468,700,440,648,468,700,416,616,524,676,468,672,444,616,524,668,448,704,436,616,500,672,472,616,496,676,468,676,464,620,496,676,464,680,436,704,440,676,1608,648,468,728,412,680,464,672,440,704,440,704,412,648,492,648,468,704,436,588,556,648,1636,616,500,644,1636,680,464,672,444,644,496,648,468,676,464,704,1580,708,1584,668,468,704,436,620,1664,648,468,648,496,556,1724,672},
		},
		IRCode{
			"HEAT_24",
			[]int{9040,4448,652,1628,684,1600,684,460,560,552,684,460,560,556,676,1608,680,1628,656,1628,684,1600,680,1604,656,484,532,588,676,464,560,556,680,1628,656,488,624,492,644,1636,680,464,624,1660,560,556,680,460,560,584,652,460,656,1628,684,460,560,556,652,1632,652,1656,628,488,680,464,648,468,652,488,648,468,652,488,616,500,680,1604,656,1656,584,556,640,476,680,460,628,488,684,460,652,460,676,468,652,464,676,464,680,464,624,488,616,504,680,460,556,560,680,460,560,556,680,1628,680,464,624,492,644,472,680,460,560,556,652,488,564,552,680,464,680,460,560,556,656,484,564,552,680,464,616,500,652,488,532,612,560,552,656,488,656,460,676,464,628,488,680,1604,652,488,532,612,588,528,680,460,628,488,676,468,624,488,680,436,656,488,680,460,652,1632,680,464,624,1660,560,552,652,492,560,552,656,488,532,608,656,1628,560,556,680,1604,680,464,556,1728,652,488,652,464,652,1656,656},
		},
		IRCode{
			"HEAT_25",
			[]int{9064,4420,676,1608,680,1604,676,464,652,464,680,464,648,468,704,1576,676,1636,672,1612,676,1608,676,1608,676,1632,680,436,704,412,676,464,676,1636,652,464,676,464,680,1604,708,408,676,1608,708,432,680,436,680,464,704,1580,704,436,696,420,676,468,664,1616,704,1580,680,464,644,496,684,432,732,412,648,464,704,440,680,436,700,1584,676,1604,704,440,648,496,704,408,708,436,648,468,676,464,680,436,680,464,700,440,680,436,704,436,676,440,732,412,680,432,704,412,676,468,700,1608,680,436,728,416,648,464,680,464,676,436,648,468,652,496,728,408,704,412,708,432,672,444,708,436,620,496,676,464,648,468,704,436,704,440,704,412,680,460,652,464,704,436,652,1636,648,464,676,468,676,468,676,436,704,440,676,436,704,440,648,468,704,436,648,496,652,1632,676,440,676,1608,676,464,620,496,676,464,676,440,704,436,680,464,648,468,676,1608,732,1576,624,1660,652,464,704,436,640,1644,704},
		},
		IRCode{
			"HEAT_26",
			[]int{9096, 4392, 676, 1632, 648, 1636, 656, 460, 648, 496, 648, 468, 676, 464, 620, 1664, 700, 1584, 676, 1632, 648, 1636, 680, 1604, 620, 496, 648, 1636, 708, 436, 552, 560, 676, 1636, 644, 496, 652, 464, 676, 1608, 680, 464, 652, 1632, 644, 468, 708, 436, 580, 560, 648, 1636, 556, 1728, 700, 1584, 704, 1580, 648, 492, 652, 1632, 588, 528, 708, 436, 700, 444, 644, 468, 732, 408, 620, 496, 704, 440, 620, 1664, 556, 1724, 708, 436, 676, 464, 652, 468, 644, 468, 676, 468, 556, 556, 708, 436, 616, 500, 648, 492, 708, 436, 672, 444, 676, 464, 616, 500, 672, 468, 708, 408, 676, 468, 580, 1704, 616, 524, 700, 416, 672, 468, 676, 440, 704, 440, 676, 436, 676, 472, 552, 584, 648, 468, 704, 436, 652, 464, 652, 492, 652, 464, 700, 440, 736, 380, 704, 436, 732, 412, 620, 496, 700, 416, 704, 436, 668, 448, 676, 1636, 676, 436, 676, 468, 648, 492, 652, 464, 672, 444, 676, 464, 672, 444, 732, 412, 636, 476, 680, 464, 648, 1660, 624, 492, 656, 1628, 704, 440, 680, 436, 648, 492, 680, 436, 708, 432, 676, 468, 648, 1632, 624, 492, 652, 492, 648, 468, 644, 1640, 732, 408, 612, 1672, 648},
		},
		IRCode{
			"HEAT_27",
			[]int{9088,4424,676,1608,664,1620,704,412,704,436,704,412,728,412,612,1672,708,1604,704,1580,640,1644,704,1580,704,1576,680,1604,708,436,588,528,680,1628,704,440,676,440,704,1580,704,436,652,1632,648,468,648,492,640,504,648,1636,664,448,708,1580,648,1660,680,436,704,1580,676,464,556,588,704,408,680,464,648,468,708,432,620,496,676,1608,676,1608,648,496,704,436,672,444,676,464,620,496,708,432,648,468,680,464,608,532,648,468,680,460,680,436,708,436,676,436,648,496,648,468,708,1604,672,440,676,468,652,460,680,464,648,468,672,468,652,464,648,496,672,468,620,496,644,472,644,496,616,500,620,520,620,496,648,496,644,496,588,528,676,464,648,468,648,496,644,1640,616,496,620,524,644,500,588,524,592,552,616,496,648,472,644,496,644,500,612,528,588,1696,564,548,648,1636,616,528,560,556,616,524,560,556,648,492,620,524,564,552,588,552,592,1692,616,500,620,1688,592,524,620,1692,592},
		},
		IRCode{
			"HEAT_28",
			[]int{9060,4480,616,1668,616,1668,616,524,592,524,620,520,592,524,672,1612,644,1668,616,1668,616,1664,620,1664,648,496,588,528,676,1608,644,496,620,1692,564,548,620,524,592,1692,620,496,616,1668,588,552,620,496,632,508,620,524,616,500,644,496,648,1636,616,500,644,1640,616,524,620,524,616,500,672,468,612,504,616,524,620,496,648,1636,620,1688,648,496,592,524,616,524,592,524,672,444,620,524,640,472,648,496,644,496,616,500,648,492,616,500,648,496,616,496,648,496,620,496,620,1688,588,528,620,520,564,552,620,524,588,528,644,496,620,496,620,520,620,524,588,528,616,524,592,524,620,520,592,524,620,524,588,552,620,496,616,524,624,492,644,500,620,492,616,1668,620,524,644,500,588,524,620,524,616,500,668,472,620,496,612,500,620,524,648,492,616,1668,648,496,620,1664,592,524,616,524,592,524,644,496,564,580,620,1664,564,1696,584,552,620,1664,620,524,612,1672,616,524,560,1724,560},
		},
		IRCode{
			"HEAT_29",
			[]int{9068,4396,672,1612,676,1608,700,440,676,440,676,468,620,492,648,1636,620,1692,672,1612,672,1612,672,1608,676,468,704,412,644,496,620,1664,620,1688,652,464,648,496,620,1664,648,468,676,1608,672,468,676,440,648,492,680,464,616,496,680,1604,676,468,620,496,676,1608,644,496,704,440,640,472,676,468,640,476,648,492,648,468,672,1612,676,1632,648,496,648,468,676,464,592,524,616,500,672,468,644,472,644,496,676,468,612,504,672,468,640,476,676,464,620,496,676,468,676,436,676,1636,672,444,700,440,672,444,648,496,644,468,648,496,704,412,672,468,676,464,652,464,676,468,648,464,680,464,648,468,676,468,644,496,620,496,676,464,676,440,644,472,676,464,672,1612,620,524,648,492,620,496,648,492,676,440,700,444,648,464,676,440,676,468,676,464,672,1612,648,492,652,1632,648,468,648,492,624,492,648,496,644,496,652,1632,676,1608,640,1644,700,416,676,464,672,444,648,1660,648,1636,648},
		},
		IRCode{
			"HEAT_30",
			[]int{9036,4480,592,1692,560,1724,564,552,620,520,568,548,644,496,592,1696,616,1692,592,1692,592,1692,560,1724,616,500,644,1664,564,1720,620,496,588,1724,612,1672,560,1720,620,524,592,524,616,1668,588,552,592,524,620,524,644,1636,648,1636,676,1636,592,524,616,1668,616,1668,588,552,588,556,564,552,616,524,616,500,616,524,564,552,620,1664,592,1692,616,552,620,496,616,500,616,524,564,552,620,524,564,548,620,524,588,552,564,552,620,524,564,552,616,524,592,524,616,524,592,524,644,1668,564,548,592,552,564,552,616,528,560,552,620,524,588,524,592,552,620,520,592,524,592,552,588,528,640,500,592,524,560,556,616,552,588,524,588,528,620,524,560,556,616,524,560,1724,564,552,644,524,620,496,560,556,616,524,564,552,592,548,564,552,620,524,644,496,564,1720,564,552,648,1664,592,520,592,528,616,524,560,556,588,552,616,1696,592,520,620,524,592,1692,588,528,616,524,596,1688,616,1668,620},
		},
		IRCode{
			"HEAT_31",
			[]int{9092,4372,700,1608,648,1636,676,464,676,440,652,492,616,500,672,1608,648,1664,676,1608,672,1612,704,1580,616,1692,652,1632,592,1692,676,440,620,1688,676,468,648,468,640,1644,652,488,620,1664,676,440,704,436,676,468,620,1664,672,444,648,492,644,472,676,1608,644,500,644,468,620,520,708,436,676,440,700,440,676,440,648,492,648,1636,608,1676,680,464,644,472,704,436,644,472,676,464,616,500,704,440,644,468,676,468,648,492,648,468,648,496,672,444,672,468,676,440,644,496,640,1644,676,468,700,412,620,524,676,440,704,436,620,496,672,468,696,448,652,464,676,464,652,464,704,440,588,524,676,444,700,440,672,468,700,440,680,436,672,444,648,492,660,456,676,1636,620,496,704,436,648,496,592,524,672,440,620,524,556,560,700,440,672,444,688,456,704,1604,648,468,668,1616,676,464,648,468,672,468,624,492,704,436,648,496,648,468,668,1616,648,1616,724,436,620,1664,648,468,648,1664,676},
		},
		IRCode{
			"HEAT_32",
			[]int{9056,4452,564,1720,620,1664,620,524,564,552,616,524,592,524,616,1668,592,1716,564,1720,592,1692,620,1664,616,528,564,548,592,552,560,1724,588,1720,564,1720,568,548,620,1664,616,528,560,1720,564,552,588,556,564,576,560,1724,588,1696,644,1640,644,1640,588,552,560,1724,572,572,564,576,564,552,560,552,592,552,564,552,644,496,564,1720,616,1668,620,524,616,524,564,552,588,552,540,576,616,528,588,524,588,528,592,576,620,496,564,552,588,552,564,552,620,524,560,552,648,496,588,1724,560,552,648,468,620,520,564,552,648,496,560,556,644,496,616,528,560,556,616,524,560,556,616,524,564,552,616,524,592,524,592,552,588,552,564,552,644,500,588,524,616,524,592,1696,560,552,620,524,588,556,588,524,620,524,564,552,588,552,564,552,612,504,588,580,588,1692,564,552,592,1692,616,528,564,552,588,552,588,528,620,520,620,1664,620,1692,592,524,560,552,644,1668,620,496,612,1672,588,1720,592},
		},

		IRCode{
			"COOL_16",
			[]int{9032,4480,648,1636,560,1724,616,500,592,548,568,548,620,524,560,1724,616,1692,620,1664,564,1720,564,1720,616,528,616,500,584,532,616,1692,616,524,596,520,648,1636,620,1664,648,496,560,1724,612,500,592,552,620,524,616,1664,620,524,564,1720,616,1668,616,500,644,496,620,496,592,548,648,496,564,552,700,440,616,500,648,492,592,1692,620,496,620,524,620,520,620,496,620,524,588,528,616,524,592,524,616,524,596,548,592,524,644,496,620,496,616,524,592,524,588,1696,592,552,564,576,592,524,620,520,592,524,620,524,616,500,612,500,620,524,616,524,616,500,620,524,560,556,616,524,564,552,644,496,592,524,620,524,588,552,592,524,644,496,620,496,648,496,564,1720,588,528,616,524,644,496,592,524,620,524,644,472,616,524,564,552,616,524,568,576,592,1692,588,528,620,1664,648,492,564,552,620,524,616,500,616,524,616,524,564,1720,616,500,620,524,560,1700,608,532,588,552,588,1720,596},
		},
		IRCode{
			"COOL_17",
			[]int{9088,4424,592,1692,652,1632,568,548,592,552,644,472,672,468,620,1664,616,1696,592,1692,592,1688,564,1720,564,1720,592,552,620,496,564,1720,672,468,592,552,564,1716,564,1724,612,504,588,1720,592,524,616,500,644,524,620,1664,648,1636,648,468,616,1668,616,524,588,528,644,496,620,524,620,496,644,496,620,496,620,524,588,524,620,1664,620,524,588,556,616,496,620,524,644,468,624,520,564,552,616,524,592,524,620,524,616,524,592,524,588,552,592,524,648,496,620,1664,620,496,672,468,648,496,620,492,620,524,620,496,616,524,592,524,616,500,648,520,648,468,616,500,616,524,672,444,620,524,612,500,620,524,644,496,616,500,620,524,648,464,676,468,620,496,648,1636,672,468,648,496,616,500,644,496,620,496,620,520,624,492,676,468,592,524,672,468,620,1664,676,464,620,1664,616,500,676,468,620,496,648,492,676,468,616,500,644,496,648,468,648,1636,648,1660,648,468,644,472,676,1636,612},
		},
		IRCode{
			"COOL_18",
			[]int{9148,4396,616,1668,640,1644,672,468,680,436,668,448,704,436,672,1612,676,1636,644,1640,664,1620,644,1640,672,468,600,1684,620,496,620,1664,648,496,648,492,620,1664,648,1636,680,464,620,1664,560,552,648,496,644,496,564,1720,560,556,680,464,616,1668,612,504,648,492,560,556,676,464,676,468,556,560,648,492,560,556,676,468,560,1720,532,584,648,496,620,520,560,556,644,500,560,552,648,496,560,556,616,524,532,612,560,552,564,580,536,580,560,584,584,528,560,1724,560,584,528,612,560,556,560,580,536,584,556,584,532,584,556,584,532,584,560,580,560,584,532,584,560,556,556,584,556,560,560,584,528,584,560,584,560,580,532,584,560,584,556,556,560,584,560,1724,532,584,560,580,560,584,528,588,556,584,560,556,560,580,560,556,588,556,528,612,560,1724,532,588,580,1724,560,556,560,556,588,556,556,556,588,556,612,532,528,1752,588,1696,616,1668,616,1668,616,528,556,560,612,1696,588},
		},
		IRCode{
			"COOL_19",
			[]int{9064,4448,680,1604,680,1604,700,416,704,436,708,408,676,468,680,1604,704,1604,680,1604,676,1608,552,1732,672,1612,648,1636,680,460,680,1608,672,468,612,504,708,1576,676,1636,648,464,704,1580,680,464,680,436,704,436,732,408,652,468,704,436,652,1632,704,412,676,464,680,436,676,464,680,468,672,440,704,436,652,464,676,468,648,1636,676,436,680,464,680,460,680,436,708,436,652,464,700,440,652,464,644,472,704,464,676,440,732,384,676,464,644,472,680,460,556,1728,672,444,732,436,680,436,704,412,732,412,612,504,672,468,608,504,708,436,728,416,644,468,708,436,616,500,676,464,708,408,704,436,624,492,676,468,704,436,656,460,708,436,648,468,732,408,704,1580,652,464,704,436,680,464,652,464,676,464,676,440,704,436,652,464,676,440,760,408,624,1660,652,464,680,1604,680,464,704,408,676,468,680,436,676,464,704,1580,704,440,540,1744,640,472,708,436,664,1620,704,412,732,1576,672},
		},
		IRCode{
			"COOL_20",
			[]int{9088,4452,680,1604,648,1636,644,472,676,464,680,436,680,464,648,1636,708,1600,652,1632,708,1576,612,1672,676,468,624,492,552,1732,672,1608,704,440,676,464,652,1636,608,1672,672,444,676,1636,652,464,676,464,676,468,648,1632,684,1604,644,1640,668,444,680,464,612,504,704,436,704,440,640,476,676,464,680,436,704,440,676,436,708,1580,704,436,676,464,612,504,676,464,620,496,680,464,648,468,676,464,680,436,708,432,704,440,680,436,676,464,652,464,704,440,652,1632,648,464,680,464,708,432,680,436,680,464,648,468,676,464,676,440,648,468,704,464,648,468,644,496,652,464,676,440,676,464,672,444,704,440,676,464,668,448,676,464,652,464,704,440,648,468,728,1556,704,436,676,468,556,556,680,464,644,472,708,432,648,468,704,440,708,404,704,440,676,1608,676,464,556,1728,700,416,708,436,672,444,704,436,704,440,636,476,704,440,552,1732,672,1612,704,436,680,1604,652,464,680,1632,648},
		},
		IRCode{
			"COOL_21",
			[]int{8948,4540,556,1728,556,1724,564,580,532,584,560,580,536,580,560,1724,560,1752,532,1752,556,1728,560,1720,564,1724,612,528,532,1752,560,1724,560,584,556,584,560,1724,532,1752,556,560,584,1724,536,580,560,584,560,580,536,580,560,1724,612,1672,588,552,564,552,560,584,532,584,584,556,560,584,560,552,564,580,532,584,560,580,588,1696,564,552,588,556,584,556,560,556,560,580,588,528,560,584,532,584,556,560,584,584,584,528,560,556,588,556,532,584,584,556,532,1752,560,584,560,580,532,584,560,556,560,580,560,556,584,560,528,588,556,584,584,556,532,584,560,584,560,556,560,580,560,556,560,580,564,552,616,528,560,580,560,556,560,584,532,580,560,584,532,1752,560,556,560,580,560,584,532,584,556,584,532,584,560,580,536,580,560,584,560,580,532,1752,560,556,560,1724,560,580,560,556,560,584,560,556,560,580,560,1724,560,1752,532,584,528,588,556,1752,532,1752,560,556,560,1748,536},
		},
		IRCode{
			"COOL_22",
			[]int{9064,4424,644,1640,672,1612,672,468,624,492,648,496,648,464,676,1608,648,1664,644,1640,672,1612,704,1580,676,464,648,1636,616,1668,676,1608,648,496,644,496,680,1604,560,1724,672,444,672,1640,616,500,616,524,564,580,616,1664,648,468,616,1668,620,524,644,472,616,524,540,576,644,500,644,496,620,496,616,524,592,524,648,496,592,1692,644,468,648,496,648,492,596,520,648,496,620,496,648,492,620,496,648,492,564,580,620,496,672,444,676,464,620,496,676,464,620,1664,648,496,560,580,652,464,616,500,680,464,556,556,652,492,532,584,676,464,648,496,532,584,616,524,560,556,616,524,560,556,560,584,560,556,588,552,616,528,560,552,560,584,532,584,560,580,540,1744,536,580,560,584,556,584,532,584,560,580,536,580,560,584,532,580,560,584,560,584,532,1748,564,552,560,1724,560,584,560,556,556,584,532,584,560,584,556,584,560,1724,560,556,560,1724,556,1752,532,1752,532,584,560,1752,532},
		},
		IRCode{
			"COOL_23",
			[]int{9092,4420,640,1644,672,1612,676,468,648,468,700,416,672,468,672,1616,672,1632,556,1728,700,1584,676,1608,704,1580,676,1608,704,1608,620,1664,668,472,620,496,648,1636,676,1636,648,464,644,1644,644,496,680,436,676,464,676,468,676,1608,648,468,672,468,676,440,620,520,652,464,648,496,620,520,624,492,672,472,648,468,672,444,644,1664,648,468,676,464,676,468,648,468,672,468,648,468,672,444,676,464,644,472,680,488,680,436,672,444,676,464,620,496,680,464,616,1668,640,476,704,436,676,468,668,444,676,468,668,448,648,492,676,440,672,468,648,496,648,468,676,464,676,440,648,496,648,468,676,464,620,496,700,440,648,496,620,496,644,496,620,548,680,436,648,1636,648,468,648,496,648,492,648,468,648,492,648,468,648,496,648,464,648,496,644,496,620,1664,648,468,648,1636,704,440,644,472,672,468,676,440,652,488,676,1608,704,1608,676,1608,648,1636,556,1728,672,1612,672,468,556,1728,672},
		},
		IRCode{
			"COOL_24",
			[]int{8972, 4540, 536, 1748, 588, 1696, 532, 584, 560, 584, 532, 584, 560, 580, 560, 1724, 560, 1752, 532, 1752, 588, 1696, 532, 1752, 560, 584, 532, 580, 560, 556, 560, 584, 560, 1748, 536, 580, 560, 1724, 564, 1724, 560, 580, 532, 1752, 532, 584, 560, 584, 556, 584, 556, 560, 560, 1748, 536, 584, 556, 584, 532, 584, 532, 1752, 560, 580, 560, 584, 532, 584, 560, 580, 536, 580, 560, 584, 532, 580, 560, 560, 556, 1752, 532, 612, 532, 584, 560, 580, 532, 584, 560, 580, 536, 580, 560, 560, 556, 584, 560, 580, 532, 584, 560, 584, 532, 584, 560, 580, 532, 584, 560, 1724, 560, 584, 560, 580, 556, 560, 560, 580, 532, 584, 560, 584, 532, 584, 560, 580, 560, 556, 560, 584, 560, 580, 532, 584, 560, 584, 560, 552, 560, 584, 536, 580, 560, 580, 564, 580, 532, 584, 560, 580, 536, 580, 560, 584, 532, 584, 556, 1728, 560, 580, 564, 580, 532, 580, 564, 580, 532, 584, 560, 580, 536, 580, 532, 584, 560, 584, 560, 580, 560, 1724, 560, 584, 532, 1752, 560, 556, 560, 580, 564, 552, 560, 584, 588, 552, 536, 1748, 532, 1752, 532, 1752, 560, 584, 532, 584, 532, 584, 560, 580, 560, 584, 532},
		},
		IRCode{
			"COOL_25",
			[]int{8952, 4532, 564, 1720, 564, 1720, 564, 580, 532, 584, 560, 580, 532, 584, 564, 1720, 612, 1696, 620, 1664, 564, 1724, 560, 1748, 536, 1748, 588, 528, 564, 576, 536, 580, 564, 1748, 528, 584, 588, 556, 584, 532, 564, 576, 536, 556, 588, 576, 564, 552, 588, 556, 564, 1716, 564, 1724, 560, 1748, 536, 580, 588, 552, 540, 576, 560, 556, 564, 604, 588, 528, 532, 588, 584, 552, 532, 584, 564, 580, 528, 1752, 536, 580, 564, 604, 540, 576, 532, 584, 564, 576, 536, 580, 564, 580, 560, 556, 564, 576, 564, 576, 532, 584, 564, 580, 584, 532, 564, 576, 588, 528, 564, 1720, 616, 524, 564, 580, 560, 556, 560, 580, 588, 528, 564, 576, 540, 576, 564, 580, 536, 576, 588, 556, 564, 576, 540, 576, 564, 580, 536, 580, 560, 580, 536, 580, 560, 556, 564, 604, 560, 556, 556, 560, 560, 580, 532, 584, 564, 576, 532, 1752, 560, 556, 564, 604, 556, 560, 560, 556, 564, 576, 532, 584, 564, 576, 536, 580, 588, 556, 564, 576, 532, 1752, 560, 556, 564, 1748, 536, 576, 564, 556, 560, 580, 556, 560, 612, 556, 536, 580, 556, 1724, 564, 1720, 564, 1720, 564, 1748, 536, 1748, 536, 580, 564, 1744, 536},
		},
		IRCode{
			"COOL_26",
			[]int{9004, 4480, 620, 1668, 616, 1664, 620, 524, 588, 528, 620, 520, 564, 552, 644, 1640, 616, 1696, 616, 1664, 620, 1664, 592, 1720, 592, 524, 616, 1668, 616, 524, 592, 524, 620, 1692, 588, 524, 648, 496, 592, 524, 620, 520, 620, 496, 648, 496, 588, 524, 620, 524, 588, 1696, 620, 520, 564, 1720, 564, 1720, 616, 500, 620, 1688, 592, 524, 620, 524, 616, 524, 620, 496, 616, 528, 592, 520, 620, 524, 564, 1720, 592, 524, 616, 524, 620, 524, 588, 524, 620, 524, 564, 552, 588, 552, 564, 552, 560, 556, 588, 552, 592, 552, 560, 552, 592, 552, 612, 504, 616, 524, 560, 1724, 564, 552, 588, 580, 592, 524, 616, 500, 588, 552, 564, 552, 616, 528, 560, 552, 620, 524, 616, 524, 564, 552, 620, 520, 568, 548, 620, 524, 588, 528, 616, 524, 564, 552, 616, 524, 620, 524, 592, 524, 620, 520, 564, 552, 616, 524, 592, 1692, 564, 552, 592, 552, 644, 496, 592, 524, 616, 524, 592, 524, 616, 528, 560, 556, 560, 552, 620, 552, 560, 1720, 564, 552, 644, 1640, 620, 524, 564, 548, 620, 524, 564, 552, 588, 552, 616, 528, 592, 520, 620, 1664, 616, 1668, 616, 528, 532, 1752, 560, 1724, 588, 1720, 564},
		},
		IRCode{
			"COOL_27",
			[]int{9044,4420,684,1600,612,1672,700,416,672,472,692,420,676,468,640,1644,676,1632,652,1632,552,1732,696,1588,704,1580,732,1552,704,436,620,496,704,1608,700,1584,732,408,648,1636,732,384,704,1580,700,440,644,472,676,464,708,1604,652,464,644,472,676,464,556,1728,704,1580,704,436,672,472,704,412,704,436,680,436,676,464,652,464,704,1580,728,416,644,496,764,352,680,460,652,464,732,412,652,464,672,468,680,436,676,464,676,468,676,436,672,444,736,408,608,508,704,1580,676,464,704,436,676,440,652,492,640,476,672,468,640,476,648,492,696,420,644,496,676,468,700,416,704,436,648,468,644,500,592,524,616,524,644,500,620,492,652,492,620,496,648,492,624,492,652,1632,652,488,560,584,624,492,672,468,620,496,588,552,560,556,560,556,588,552,588,556,560,1724,560,580,588,1696,560,556,560,580,564,552,560,584,532,608,532,1752,560,556,584,1700,560,1748,568,548,560,584,556,556,588,556,588},
		},
		IRCode{
			"COOL_28",
			[]int{8944,4540,532,1752,532,1752,560,580,536,580,560,584,532,580,536,1752,556,1752,532,1752,556,1728,560,1724,560,580,588,528,560,1724,560,584,556,1752,536,1748,536,580,560,1724,560,584,560,1720,560,556,560,584,556,584,532,1752,532,1752,560,580,536,580,560,1724,560,1724,560,584,556,584,560,556,560,580,560,556,560,584,532,584,556,1724,560,584,532,608,560,556,560,584,584,528,560,584,532,584,560,580,532,584,560,584,556,584,532,584,560,580,536,580,560,584,532,1748,536,584,556,584,560,580,536,580,560,584,532,580,560,556,560,584,532,580,560,584,560,580,532,584,560,584,532,584,556,584,560,556,560,580,532,612,556,560,556,584,584,532,560,580,536,580,560,1724,560,584,532,608,560,556,560,580,588,528,560,584,532,580,560,584,532,584,560,580,560,1724,560,584,584,1700,532,584,556,584,560,556,560,580,532,612,532,1752,532,1752,556,1728,556,584,532,1752,560,556,560,580,532,612,532},
		},
		IRCode{
			"COOL_29",
			[]int{9056,4408,640,1668,564,1716,644,472,676,468,644,472,672,468,648,1636,700,1608,620,1668,616,1664,560,1724,644,1640,676,468,648,1636,672,444,672,1636,672,1612,700,444,648,1636,644,468,648,1636,648,496,700,416,672,468,672,472,644,472,668,1612,676,468,644,1640,640,1644,644,496,644,500,620,496,668,472,652,464,700,440,676,440,668,1616,672,468,620,524,648,468,672,468,620,496,672,444,704,436,672,444,648,496,676,464,648,468,672,468,644,472,648,496,672,444,644,1640,644,496,676,468,640,472,648,496,616,500,672,468,648,468,644,500,672,440,676,468,672,468,648,468,644,500,672,440,676,468,620,496,644,496,644,500,648,464,676,468,592,524,644,472,672,468,672,1612,616,524,676,468,648,468,644,496,648,468,672,444,644,496,648,468,648,496,672,468,644,1640,644,496,648,1640,648,464,648,496,620,496,644,496,620,524,648,464,676,468,648,468,668,448,644,496,640,1644,644,496,648,496,648},
		},
		IRCode{
			"COOL_30",
			[]int{9032,4480,592,1720,588,1692,596,520,604,540,592,524,560,556,588,1720,560,1724,620,1664,620,1688,592,1692,568,548,616,1668,620,1664,620,524,616,1692,592,1692,564,552,620,1664,620,524,560,1720,644,472,644,500,616,524,564,552,620,1692,644,1636,620,496,620,1668,616,1664,620,524,616,524,616,500,676,468,612,500,648,496,564,552,588,1696,616,524,620,520,564,556,588,552,560,556,644,496,588,528,648,492,592,524,672,472,620,520,620,496,616,524,624,496,644,496,620,1664,568,548,616,524,592,552,644,468,620,524,616,500,644,496,592,524,644,500,616,524,620,496,644,496,592,524,620,524,588,532,608,500,620,524,644,496,564,552,620,524,612,504,616,524,584,532,616,1668,616,524,648,496,616,496,648,496,560,556,588,552,560,556,620,524,644,468,620,524,616,1668,620,520,564,1720,620,524,592,524,612,500,620,524,620,520,616,500,620,1688,596,524,616,1664,620,524,616,1668,564,552,644,496,616},
		},
		IRCode{
			"COOL_31",
			[]int{9060,4480,588,1696,620,1664,620,520,648,468,676,468,536,580,644,1636,620,1692,644,1640,648,1636,644,1640,564,1748,620,1664,592,1692,616,496,648,1664,672,1612,652,488,596,1688,624,492,620,1664,620,524,588,524,652,492,644,496,648,468,648,496,620,1664,644,1636,648,1636,676,468,644,496,596,520,620,524,620,496,616,524,592,524,588,1696,620,520,648,496,648,468,644,496,592,524,648,496,620,492,588,528,648,492,652,492,648,468,672,468,644,472,648,496,588,524,620,1664,652,492,648,492,648,472,616,524,644,468,652,492,648,468,644,496,620,496,648,496,620,520,616,500,620,520,648,472,644,496,620,496,616,524,564,576,624,496,644,496,620,496,644,472,644,496,648,1636,588,552,620,524,592,524,644,496,620,496,648,496,592,524,644,468,648,496,620,520,648,1636,648,496,620,1660,592,524,592,552,620,496,616,524,620,524,620,496,644,496,644,1640,648,492,620,1664,620,1664,648,496,620,496,648},
		},
		IRCode{
			"COOL_32",
			[]int{9088,4452,648,1636,676,1608,620,496,672,468,680,436,708,436,648,1632,648,1664,620,1664,620,1664,620,1664,668,448,728,412,640,476,676,1608,676,1636,648,1636,648,492,648,1636,608,508,664,1620,648,492,672,444,704,464,620,496,672,1612,648,492,624,1664,644,1636,556,1728,644,472,680,488,680,436,672,444,676,468,664,448,676,468,556,1728,644,472,672,496,620,496,640,476,700,440,676,440,704,436,648,468,676,468,644,496,648,468,708,432,708,408,676,468,620,496,644,1636,648,496,696,444,648,468,676,468,644,472,648,492,620,496,676,468,672,440,648,496,676,464,624,492,620,524,620,492,616,500,620,524,700,416,676,464,704,440,640,472,648,496,668,448,648,492,556,1728,700,416,680,464,648,492,640,476,676,464,640,476,648,496,644,468,648,496,672,468,620,1664,672,444,676,1636,620,492,676,444,672,468,668,448,676,468,700,440,640,1644,704,1576,620,1668,672,1612,648,1660,620,496,676,468,552},
		},
	}
)

type IRCode struct {
	tag  string
	code []int
}

func (c IRCode) GetTag() string {
	return c.tag
}

func (c IRCode) ToCMD() string {
	var sb strings.Builder
	for i, co := range c.code {
		sb.WriteString(strconv.Itoa(co))
		if (i < len(c.code)-1) {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

type IRACAirCool struct {
	ID       string
	Name     string
	Sensor   string
	ValueStr string
	value    float64
	codes    []IRCode
}

func (_ *IRACAirCool) InRangeThreshold() bool {
	return false
}

func (s *IRACAirCool) getIRCode(tag string) (IRCode, bool) {
	if !strings.HasPrefix(tag, "COOL") && !strings.HasPrefix(tag, "HEAT") {
		tag = strings.Split(tag, "_")[0]
	} 
	
	for _, irc := range s.codes {
		if irc.tag == tag {
			return irc, true
		}
	}
	return IRCode{}, false
}

func (s *IRACAirCool) CreateCMD(cmd string) (string, string, error) {
	code, ok := s.getIRCode(cmd)
	if ok {
		return code.ToCMD(), code.GetTag(), nil
	}
	return cmd, CUSTOM, nil
}

func (s *IRACAirCool) GetID() string {
	return s.ID
}

func (s *IRACAirCool) GetSensor() string {
	return s.Sensor
}

func (s *IRACAirCool) GetName() string {
	return s.Name
}

func (s *IRACAirCool) SetValue(newValstr string) error {
	s.ValueStr = newValstr
	return nil
}

func (s *IRACAirCool) GetValueStr() string {
	return s.ValueStr
}

func (s *IRACAirCool) GetOptions(str string) []string {
	switch str {
	case "mode":
		set := make(map[string]bool) 
		for _, m := range airCoolCodes {
			set[strings.Split(m.GetTag(), "_")[0]] = true
		}
		modes := []string{}
		for k := range set {         
			modes = append(modes, k)
		}
		sort.Strings(modes)
		return modes 
	case "temp":
		set := make(map[string]bool) 
		for _, m := range airCoolCodes {
			temp := strings.Split(m.GetTag(), "_")
			if len(temp) > 1 {
				set[temp[1]] = true
			}
		}
		temps := []string{}
		for k := range set {
			temps = append(temps, k)
		}
		sort.Strings(temps)
		return temps 
	}
	return []string{}
}

func (s *IRACAirCool) SendStats(c *statsd.Client) {
	log.Println("IR AC Air Cool, SendStats need implement this function")
}

func NewIRACAirCool(id, sensor string) *IRACAirCool {
	return &IRACAirCool{
		ID:     id,
		Name:   IR_ac_aircool,
		Sensor: sensor,
		codes:  airCoolCodes,
		ValueStr: "UNSET,UNSET",
	}
}

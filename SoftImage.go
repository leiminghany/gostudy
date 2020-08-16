package main

import (
	"fmt"
	"math"
	"sync"
)

const UNITX = 2
const UNITY = 2

type ImageUnit [UNITX][UNITY]byte
type ImageTransfer struct {
	Point     [][]ImageUnit
	SoftPoint [][]byte
}

var wg sync.WaitGroup

func softImage(ppoints *[][]byte) {
	//defer wg.Done()
	var points = *ppoints
	fmt.Printf("%v", points)

	lenx := len(points)
	if lenx == 0 {
		return
	}
	leny := len(points[0])
	if leny == 0 {
		return
	}

	for i := 0; i < lenx; i++ {
		for j := 0; j < leny; j++ {
			var relatedArr []byte
			relatedArr = append(relatedArr, points[i][j])
			//fmt.Printf("i:%d, j:%d, before: [%v]",i, j, points[i][j])
			//left
			if i > 0 {
				relatedArr = append(relatedArr, points[i-1][j])
				if j > 0 {
					relatedArr = append(relatedArr, points[i-1][j-1])
				}
				if j < leny-1 {
					relatedArr = append(relatedArr, points[i-1][j+1])
				}
			}
			//right
			if i < lenx-1 {
				relatedArr = append(relatedArr, points[i+1][j])
				if j > 0 {
					relatedArr = append(relatedArr, points[i+1][j-1])
				}
				if j < leny-1 {
					relatedArr = append(relatedArr, points[i+1][j+1])
				}
			}
			//top one
			if j > 0 {
				relatedArr = append(relatedArr, points[i][j-1])
			}
			// bottom one
			if j < leny-1 {
				relatedArr = append(relatedArr, points[i][j+1])
			}
			var sum = 0
			for _, v := range relatedArr {
				sum += int(v)
			}
			avg := math.Floor(float64(sum) / float64(len(relatedArr)))
			points[i][j] = byte(avg)
			//fmt.Printf("  after: %v\n", points[i][j])
			//fmt.Println("len: ", len(relatedArr));

		}
	}

}

func NewImageTransfer(points [][]byte) *ImageTransfer {
	var transfer = new(ImageTransfer)
	var x = len(points)
	var y = len(points[0])

	groupx := int(math.Ceil(float64(x) / UNITX))
	groupy := int(math.Ceil(float64(y) / UNITX))
	//wg.Add(groupx * groupy)
	for i := 0; i < groupx; i++ {
		for j := 0; j < groupy; j++ {
			beginx := i * UNITX
			beginy := i * UNITY
			endx := (i + 1) * UNITX
			endy := (j + 1) * UNITX
			if endx > x {
				endx = x
			}
			if endy > y {
				endy = y
			}
			fmt.Printf("%v,%v,%v,%v\n", beginx, beginy, endx, endy)
			unitx := points[beginx:endx]
			unitPoint := unitx[beginy:endy]
			fmt.Println("%+v", unitPoint)
			//softImage(&unitPoint);
		}
	}

	//wg.Wait()

	// softImage(&points);
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			fmt.Printf("%v ", points[i][j])
		}
		fmt.Println()
	}

	return transfer
}
func main() {
	NewImageTransfer([][]byte{{1, 1, 1}, {1, 0, 1}, {1, 1, 1}, {2, 2, 2}})
}

/*
设计一个“imageTransfer”类型用来处理图片柔化。一张图片可以用整型矩阵来标识，像素点为矩阵中的元素，取值范围为[0, 255]. 柔化一张照片其实是创建一张新的图片，图片中的每个像素点都为原图对应像素点和其周围8个像素点的平均值的floor。
1.如果周围没有8个点，则取尽可能多的点[下面的例子有点问题]。
2.需要自定义imageTransfer结构体与对外接口，考虑并发状况一次只能处理一张图片
3.目前主流相机的像素可以达到千万级，所以请尽可能使效率最高（比如多线程）
例子：
输入:
[
    [1,1,1],
    [1,0,1],
    [1,1,1]
]
输出:
 [
   [0, 0, 0],
   [0, 0, 0],
   [0, 0, 0]
]
// 点 (0,0), (0,2), (2,0), (2,2),周围3个点加自己: floor((1 + 1 + 0 + 1)/4) = floor(0.75) = 0
// 点 (0,1), (1,0), (1,2), (2,1),周围5个点加自己: floor((1 + 1 + 1 + 1 + 1 + 0)/6) = floor(0.83333333) = 0
// 点 (1,1),周围8个点加自己: floor(8/9) = floor(0.88888889) = 0*/

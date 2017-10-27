package main

import (
    "fmt"
    "math/rand"
    "time"
    //"bytes"
)

type Gene struct {
    chrom []byte
    fit   int
}
//////// Utils /////////////////////////
func CopyPop(p []Gene) []Gene {
    tmpgenes := []Gene{}
    for _, gene := range p {
        tmpchrom := make([]byte,len(gene.chrom))
        copy(tmpchrom,gene.chrom)
        tmpgene := Gene{tmpchrom,gene.fit}
        tmpgenes = append(tmpgenes,tmpgene)
    }
    return tmpgenes
}

func CreateGenes(size int, length int) []Gene {
    rand.Seed(time.Now().UnixNano())
    genes := []Gene{}
    for i:=0;i<size;i++ {
        chrom := []byte{}
        for j:=0;j<length;j++ {
            chrom = append(chrom,byte(rand.Intn(255)))
        }
        gene := Gene{chrom,1}
        genes = append(genes,gene)
    }
    return genes
}

//////// Caluclation Fitness ///////////
func BitCount(buf byte) int{
    b:= (buf & 0x55)+((buf & 0xaa)>>1)
    b = (b & 0x33)+((b & 0xcc)>>2)
    b = (b & 0x0f)+((b & 0xf0)>>4)
    return int(b)
}

func MatchCount(buf1 byte, buf2 byte) int{
    b := buf1 ^ buf2
    return BitCount(^b)
}

func ArrayMatchCount(buf1 []byte, buf2 []byte) int {
    sum := 0
    for i:=0;i<len(buf1);i++ {
        sum = sum + MatchCount(buf1[i],buf2[i])
    }
    return sum
}

/////// Sort Algorithm /////////////////

func BubbleSort(p []Gene) {
    for i:=0;i<len(p);i++{
        for j:=len(p)-1;j>i;j--{
            if(p[j].fit < p[j-1].fit){
                tmpfit := p[j].fit
                tmpchrom := make([]byte,len(p[j].chrom))
                copy(tmpchrom,p[j].chrom)
                p[j].fit = p[j-1].fit
                copy(p[j].chrom, p[j-1].chrom)
                p[j-1].fit = tmpfit
                copy(p[j-1].chrom, tmpchrom)
            }
        }
    }
}

/////// Elete Select /////////////////

func EleteSelection(p []Gene, num int) []Gene {
    pop := CopyPop(p)
    BubbleSort(pop)
    return pop[len(p)-num:]
}

/////// Crossover    /////////////////

func OnePointCrossover(chrom1 []byte, chrom2 []byte) ([]byte, []byte){
    rand.Seed(time.Now().UnixNano())
    ran := rand.Intn(len(chrom1)*8)
    locate := int(float32(ran)/8)
    point := ran%8
    mask := []byte{0x00,0x80,0xc0,0xe0,0xf0,0xf8,0xfc,0xfe}
    next_chrom1 := make([]byte,len(chrom1))
    next_chrom2 := make([]byte,len(chrom2))
    for i:=0;i<locate;i++ {
        next_chrom1[i] = chrom2[i]
        next_chrom2[i] = chrom1[i]
    }
    next_chrom2[locate] = (chrom1[locate] & mask[point]) + (chrom2[locate] &(^mask[point]))
    next_chrom1[locate] = (chrom2[locate] & mask[point]) + (chrom1[locate] &(^mask[point]))
    for i:=locate+1;i<len(chrom1);i++ {
        next_chrom1[i] = chrom1[i]
        next_chrom2[i] = chrom2[i]
    }
    return next_chrom1, next_chrom2
}


/////// Mutation    /////////////////

func Mutation(chrom []byte) []byte{
    rand.Seed(time.Now().UnixNano())
    onebit := []byte{ 0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01 }
    ran1 := rand.Intn(len(chrom))
    ran2 := rand.Intn(8)
    chrom[ran1] ^= onebit[ran2]
    return chrom
}

////// Loop All Process ////////////
func Cycle(p []Gene,ans []byte, cycle int, elete_rate float32, cross_rate float32, mutate_rate float32) []Gene{
    rand.Seed(time.Now().UnixNano())
    p_len := len(p)
    elete_num := int(elete_rate*float32(len(p)))
    cross_num := int(cross_rate*float32(len(p)))
    for i:=0;i<cycle;i++ {
        //Calculate Fitness
        for j:=0;j<len(p);j++ {
            p[j].fit = ArrayMatchCount(ans, p[j].chrom)
        }
        elete := EleteSelection(p,elete_num)
        for j:=0;j<cross_num;j++ {
            ran1 := rand.Intn(len(elete))
            ran2 := rand.Intn(len(elete))
            next_chrom1, next_chrom2 := OnePointCrossover(elete[ran1].chrom,elete[ran2].chrom)
            p = append(p, Gene{next_chrom1,1})
            p = append(p, Gene{next_chrom2,1})
        }
        for j:=0;j<len(p);j++ {
            p[j].fit = ArrayMatchCount(ans, p[j].chrom)
        }
        for j:=0;j<len(p);j++ {
            ranmute := rand.Float32()
            if ranmute < mutate_rate {
                Mutation(p[j].chrom)
            }
        }
        p = EleteSelection(p,p_len)
    }
    return p
}

func main(){
    pop := CreateGenes(50,6)
    for i:=0;i<len(pop);i++{
        fmt.Printf("%08b\n",pop[i].chrom)
    }
    ans := []byte{}
    for i:=0;i<len(pop[0].chrom);i++{
        ans = append(ans,0xff)
    }

    pop = Cycle(pop, ans, 1000, 0.3, 0.8, 0.01)
    for i:=0;i<len(pop);i++{
        fmt.Printf("%08b,%d\n",pop[i].chrom,pop[i].fit)
    }
}

package service

//
//import (
//	"fmt"
//	"math"
//	"packCalculator/repository"
//	"sort"
//	"sync"
//)
//
//type PackService struct {
//	repo    *repository.PackRepository
//	current []int
//	//packSizes []int
//	mu sync.RWMutex
//}
//
//func NewPackService(repo *repository.PackRepository) (*PackService, error) {
//	sizes, err := repo.GetAll()
//	if err != nil {
//		return nil, err
//	}
//
//	// Ensure sizes are sorted in descending order
//	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
//
//	return &PackService{
//		repo:    repo,
//		current: sizes,
//	}, nil
//}
//
//func (s *PackService) GetCurrentSizes() []int {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	copied := make([]int, len(s.current))
//	copy(copied, s.current)
//	return copied
//}
//
//func (s *PackService) AddPackSize(size int) error {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	if err := s.repo.Add(size); err != nil {
//		return err
//	}
//
//	// Update local cache
//	fmt.Println("Before update", s.current)
//	s.current = append(s.current, size)
//	fmt.Println("After update", s.current)
//	sort.Sort(sort.Reverse(sort.IntSlice(s.current)))
//
//	return nil
//}
//
//func (s *PackService) RemovePackSize(size int) error {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	if err := s.repo.Remove(size); err != nil {
//		return err
//	}
//
//	// Update local cache
//	newSizes := make([]int, 0, len(s.current)-1)
//	for _, s := range s.current {
//		if s != size {
//			newSizes = append(newSizes, s)
//		}
//	}
//	s.current = newSizes
//
//	return nil
//}
//
//func (s *PackService) UpdatePackSizes(sizes []int) error {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	if err := s.repo.ReplaceAll(sizes); err != nil {
//		return err
//	}
//
//	// Update local cache
//	s.current = make([]int, len(sizes))
//	copy(s.current, sizes)
//	sort.Sort(sort.Reverse(sort.IntSlice(s.current)))
//
//	return nil
//}
//
////	func (s *PackService) CalculatePacks(itemsOrdered int) map[int]int {
////		if itemsOrdered <= 0 {
////			return make(map[int]int)
////		}
////
////		// Make a copy of pack sizes to work with
////		packSizes := make([]int, len(s.current))
////		copy(packSizes, s.current)
////
////		fmt.Println("In calculating ", packSizes)
////
////		result := make(map[int]int)
////		remaining := itemsOrdered
////
////		// First pass: try to use largest packs first
////		for _, size := range packSizes {
////			if remaining >= size {
////				count := remaining / size
////				result[size] = count
////				remaining -= count * size
////			}
////		}
////
////		// Second pass: optimize by checking if we can reduce total items by using a larger pack
////		// instead of multiple smaller ones (rule #2 takes precedence over rule #3)
////		if remaining > 0 {
////			for i := 0; i < len(packSizes)-1; i++ {
////				// Check if combining with next larger pack would reduce total items
////				nextSize := packSizes[i]
////				currentSize := packSizes[i+1]
////
////				totalWithCurrent := (result[currentSize] + 1) * currentSize
////				totalWithNext := ((result[nextSize] + (totalWithCurrent+nextSize-1)/nextSize) * nextSize)
////
////				if totalWithNext <= totalWithCurrent {
////					delete(result, currentSize)
////					result[nextSize] = (totalWithCurrent + nextSize - 1) / nextSize
////					remaining = 0
////					break
////				}
////			}
////
////			// If we still have remaining, add the smallest pack
////			if remaining > 0 {
////				smallestPack := packSizes[len(packSizes)-1]
////				result[smallestPack] = result[smallestPack] + 1
////			}
////		}
////
////		// Third pass: try to minimize number of packs (rule #3) without increasing total items
////		for i := 0; i < len(packSizes)-1; i++ {
////			currentSize := packSizes[i]
////			nextSize := packSizes[i+1]
////
////			// Check if we can replace multiple smaller packs with one larger pack
////			if result[nextSize] > 0 && currentSize%nextSize == 0 {
////				ratio := currentSize / nextSize
////				if result[nextSize] >= ratio {
////					result[currentSize] = result[currentSize] + 1
////					result[nextSize] = result[nextSize] - ratio
////				}
////			}
////		}
////
////		// Remove any packs with zero quantity
////		for size, count := range result {
////			if count == 0 {
////				delete(result, size)
////			}
////		}
////
////		return result
////	}
//func (s *PackService) CalculatePacks(itemsOrdered int) map[int]int {
//	if itemsOrdered <= 0 {
//		return make(map[int]int)
//	}
//
//	// Make a copy of pack sizes to work with
//	packSizes := make([]int, len(s.current))
//	copy(packSizes, s.current)
//
//	// Sort pack sizes in ascending order for the coin change algorithm
//	sort.Ints(packSizes)
//
//	// Initialize dp array to store the minimum number of packs needed for each amount
//	dp := make([]int, itemsOrdered+1)
//	// Initialize a map to track the last pack size used for each amount
//	lastPack := make([]int, itemsOrdered+1)
//
//	// Fill dp array with a large value
//	for i := range dp {
//		dp[i] = math.MaxInt
//	}
//	dp[0] = 0 // Base case: 0 packs needed to make 0 items
//
//	// Dynamic programming to calculate minimum packs
//	for _, size := range packSizes {
//		for x := size; x <= itemsOrdered; x++ {
//			if dp[x-size] != math.MaxInt && dp[x-size]+1 < dp[x] {
//				dp[x] = dp[x-size] + 1
//				lastPack[x] = size
//			}
//		}
//	}
//
//	// If it's not possible to fulfill the order exactly, select the smallest pack size that can cover the order
//	if dp[itemsOrdered] == math.MaxInt {
//		for _, size := range packSizes {
//			if size >= itemsOrdered {
//				return map[int]int{size: 1}
//			}
//		}
//		return make(map[int]int)
//	}
//
//	// Backtrack to determine the pack sizes used
//	result := make(map[int]int)
//	remaining := itemsOrdered
//	for remaining > 0 {
//		pack := lastPack[remaining]
//		result[pack]++
//		remaining -= pack
//	}
//
//	return result
//}

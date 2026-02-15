package chunker

// Chunker manages splitting followers into chunks for efficient diffing
type Chunker struct {
	chunkSize int
}

// New creates a new chunker with the specified chunk size
func New(chunkSize int) *Chunker {
	if chunkSize <= 0 {
		chunkSize = 1000 // Default chunk size
	}
	return &Chunker{chunkSize: chunkSize}
}

// Chunk splits a slice of follower hashes into chunks
func (c *Chunker) Chunk(hashes []string) [][]string {
	if len(hashes) == 0 {
		return nil
	}

	var chunks [][]string
	for i := 0; i < len(hashes); i += c.chunkSize {
		end := i + c.chunkSize
		if end > len(hashes) {
			end = len(hashes)
		}
		chunks = append(chunks, hashes[i:end])
	}
	return chunks
}

// ChunkCount returns the number of chunks needed for a given follower count
func (c *Chunker) ChunkCount(followerCount int) int {
	if followerCount <= 0 {
		return 0
	}
	count := followerCount / c.chunkSize
	if followerCount%c.chunkSize > 0 {
		count++
	}
	return count
}

// GetChunkRange returns the start and end indices for a specific chunk
func (c *Chunker) GetChunkRange(chunkID, totalCount int) (start, end int) {
	start = chunkID * c.chunkSize
	end = start + c.chunkSize
	if end > totalCount {
		end = totalCount
	}
	if start > totalCount {
		start = totalCount
	}
	return
}

// RotatingSchedule determines which chunks to scan based on rotation
// This implements the "scan 5 chunks per hour, full coverage every ~20 hours" strategy
type RotatingSchedule struct {
	chunksPerCycle int
}

// NewRotatingSchedule creates a new rotating schedule
func NewRotatingSchedule(chunksPerCycle int) *RotatingSchedule {
	if chunksPerCycle <= 0 {
		chunksPerCycle = 5
	}
	return &RotatingSchedule{chunksPerCycle: chunksPerCycle}
}

// GetChunksToScan returns which chunks should be scanned in this cycle
func (r *RotatingSchedule) GetChunksToScan(totalChunks, cycleNumber int) []int {
	if totalChunks == 0 {
		return nil
	}

	// If total chunks <= chunks per cycle, scan all
	if totalChunks <= r.chunksPerCycle {
		chunks := make([]int, totalChunks)
		for i := range chunks {
			chunks[i] = i
		}
		return chunks
	}

	// Otherwise, rotate through chunks
	startChunk := (cycleNumber * r.chunksPerCycle) % totalChunks
	chunks := make([]int, r.chunksPerCycle)
	for i := range chunks {
		chunks[i] = (startChunk + i) % totalChunks
	}
	return chunks
}

// CyclesForFullCoverage returns how many cycles needed for full coverage
func (r *RotatingSchedule) CyclesForFullCoverage(totalChunks int) int {
	if totalChunks <= r.chunksPerCycle {
		return 1
	}
	cycles := totalChunks / r.chunksPerCycle
	if totalChunks%r.chunksPerCycle > 0 {
		cycles++
	}
	return cycles
}

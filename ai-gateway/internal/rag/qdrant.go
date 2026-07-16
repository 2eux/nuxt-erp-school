package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/opencode/erp-ai-gateway/internal/config"
	"github.com/opencode/erp-ai-gateway/internal/providers"
	"go.uber.org/zap"
)

type QdrantClient struct {
	cfg        config.QdrantConfig
	httpClient *http.Client
	logger     *zap.Logger
	baseURL    string
}

type QdrantPoint struct {
	ID      string         `json:"id"`
	Vector  []float32      `json:"vector"`
	Payload map[string]any `json:"payload"`
}

type QdrantSearchHit struct {
	ID      string         `json:"id"`
	Score   float32        `json:"score"`
	Vector  []float32      `json:"vector,omitempty"`
	Payload map[string]any `json:"payload"`
}

type qdrantSearchRequest struct {
	Vector        []float32     `json:"vector"`
	Limit         int           `json:"limit"`
	WithPayload   bool          `json:"with_payload"`
	WithVector    bool          `json:"with_vector"`
	ScoreThreshold float64      `json:"score_threshold,omitempty"`
	Filter        *qdrantFilter `json:"filter,omitempty"`
}

type qdrantFilter struct {
	Must    []qdrantCondition `json:"must,omitempty"`
	Should  []qdrantCondition `json:"should,omitempty"`
	MustNot []qdrantCondition `json:"must_not,omitempty"`
}

type qdrantCondition struct {
	Key   string      `json:"key"`
	Match qdrantMatch `json:"match"`
}

type qdrantMatch struct {
	Value any `json:"value"`
}

type qdrantSearchResponse struct {
	Result []QdrantSearchHit `json:"result"`
	Time   float64           `json:"time"`
}

type qdrantUpsertRequest struct {
	Points []QdrantPoint `json:"points"`
}

type qdrantCreateCollectionRequest struct {
	Vectors qdrantVectorConfig `json:"vectors"`
}

type qdrantVectorConfig struct {
	Size     int    `json:"size"`
	Distance string `json:"distance"`
}

type qdrantCollectionInfo struct {
	Result struct {
		Name           string `json:"name"`
		PointsCount    int    `json:"points_count"`
		IndexedVectors int    `json:"indexed_vectors_count"`
		SegmentsCount  int    `json:"segments_count"`
		Config         struct {
			Params struct {
				Vectors struct {
					Size     int    `json:"size"`
					Distance string `json:"distance"`
				} `json:"vectors"`
			} `json:"params"`
		} `json:"config"`
	} `json:"result"`
}

type qdrantDeleteRequest struct {
	Points []string    `json:"points,omitempty"`
	Filter *qdrantFilter `json:"filter,omitempty"`
}

type Document struct {
	ID       string            `json:"id"`
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type SearchResult struct {
	ID       string            `json:"id"`
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	Score    float32           `json:"score"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type Chunk struct {
	ID        string            `json:"id"`
	DocID     string            `json:"doc_id"`
	Content   string            `json:"content"`
	Index     int               `json:"index"`
	Embedding []float32         `json:"embedding,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type RAGService struct {
	qdrant   *QdrantClient
	embeddings *EmbeddingService
	cfg        config.RAGConfig
	logger     *zap.Logger
}

func NewQdrantClient(cfg config.QdrantConfig, logger *zap.Logger) *QdrantClient {
	baseURL := fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)
	if cfg.Host == "" {
		baseURL = "http://localhost:6334"
	}

	return &QdrantClient{
		cfg:     cfg,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        20,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		logger: logger,
	}
}

func NewRAGService(qdrant *QdrantClient, embeddings *EmbeddingService, cfg config.RAGConfig, logger *zap.Logger) *RAGService {
	return &RAGService{
		qdrant:     qdrant,
		embeddings: embeddings,
		cfg:        cfg,
		logger:     logger,
	}
}

func (c *QdrantClient) EnsureCollection(ctx context.Context, vectorSize int, distance string) error {
	exists, err := c.CollectionExists(ctx)
	if err != nil {
		return fmt.Errorf("check collection existence: %w", err)
	}
	if exists {
		c.logger.Info("collection already exists", zap.String("collection", c.cfg.Collection))
		return nil
	}

	return c.CreateCollection(ctx, vectorSize, distance)
}

func (c *QdrantClient) CollectionExists(ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/collections/%s", c.baseURL, c.cfg.Collection)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "no such host") {
			return false, fmt.Errorf("qdrant is not reachable at %s: %w", c.baseURL, err)
		}
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200, nil
}

func (c *QdrantClient) CreateCollection(ctx context.Context, vectorSize int, distance string) error {
	if distance == "" {
		distance = "Cosine"
	}
	if vectorSize <= 0 {
		vectorSize = 3072
	}

	body := qdrantCreateCollectionRequest{
		Vectors: qdrantVectorConfig{
			Size:     vectorSize,
			Distance: distance,
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal create collection: %w", err)
	}

	url := fmt.Sprintf("%s/collections/%s", c.baseURL, c.cfg.Collection)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("create collection request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant create collection failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	c.logger.Info("qdrant collection created",
		zap.String("collection", c.cfg.Collection),
		zap.Int("vector_size", vectorSize),
		zap.String("distance", distance),
	)

	return nil
}

func (c *QdrantClient) UpsertPoints(ctx context.Context, points []QdrantPoint) error {
	if len(points) == 0 {
		return nil
	}

	body := qdrantUpsertRequest{Points: points}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal upsert: %w", err)
	}

	url := fmt.Sprintf("%s/collections/%s/points?wait=true", c.baseURL, c.cfg.Collection)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("upsert request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant upsert failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	c.logger.Debug("qdrant points upserted", zap.Int("count", len(points)))
	return nil
}

func (c *QdrantClient) Search(ctx context.Context, vector []float32, topK int, filters map[string]string, threshold float64) ([]QdrantSearchHit, error) {
	if len(vector) == 0 {
		return nil, fmt.Errorf("empty search vector")
	}
	if topK <= 0 {
		topK = 5
	}

	searchReq := qdrantSearchRequest{
		Vector:        vector,
		Limit:         topK,
		WithPayload:   true,
		WithVector:    false,
		ScoreThreshold: threshold,
	}

	if len(filters) > 0 {
		searchReq.Filter = c.buildFilter(filters)
	}

	data, err := json.Marshal(searchReq)
	if err != nil {
		return nil, fmt.Errorf("marshal search request: %w", err)
	}

	url := fmt.Sprintf("%s/collections/%s/points/search", c.baseURL, c.cfg.Collection)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("qdrant search failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result qdrantSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	return result.Result, nil
}

func (c *QdrantClient) HybridSearch(ctx context.Context, vector []float32, keyword string, topK int, filters map[string]string, threshold float64) ([]QdrantSearchHit, error) {
	vectorResults, err := c.Search(ctx, vector, topK*2, filters, threshold)
	if err != nil {
		return nil, fmt.Errorf("vector search in hybrid: %w", err)
	}

	keywordResults := make([]QdrantSearchHit, 0)
	if keyword != "" {
		keywordResults, err = c.keywordSearch(ctx, keyword, topK*2, filters)
		if err != nil {
			c.logger.Warn("keyword search failed in hybrid, using vector only", zap.Error(err))
		}
	}

	merged := c.mergeHybridResults(vectorResults, keywordResults, topK)
	return merged, nil
}

func (c *QdrantClient) keywordSearch(ctx context.Context, keyword string, topK int, filters map[string]string) ([]QdrantSearchHit, error) {
	filter := c.buildFilter(filters)
	if keyword != "" {
		filter.Must = append(filter.Must, qdrantCondition{
			Key: "content",
			Match: qdrantMatch{Value: keyword},
		})
	}

	searchReq := map[string]any{
		"vector":       make([]float32, 3072),
		"limit":        topK,
		"with_payload": true,
		"with_vector":  false,
		"filter":       filter,
	}

	data, err := json.Marshal(searchReq)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/collections/%s/points/search", c.baseURL, c.cfg.Collection)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("keyword search failed: %d", resp.StatusCode)
	}

	var result qdrantSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

func (c *QdrantClient) mergeHybridResults(vectorResults, keywordResults []QdrantSearchHit, topK int) []QdrantSearchHit {
	seen := make(map[string]bool)
	merged := make([]QdrantSearchHit, 0)

	for _, hit := range vectorResults {
		if !seen[hit.ID] {
			seen[hit.ID] = true
			merged = append(merged, hit)
		}
	}

	vectorWeight := float32(0.7)
	for _, hit := range keywordResults {
		if !seen[hit.ID] {
			seen[hit.ID] = true
			hit.Score = hit.Score * (1 - vectorWeight)
			merged = append(merged, hit)
		} else {
			for i, m := range merged {
				if m.ID == hit.ID {
					merged[i].Score = m.Score*vectorWeight + hit.Score*(1-vectorWeight)
					break
				}
			}
		}
	}

	if len(merged) > topK {
		merged = merged[:topK]
	}

	return merged
}

func (c *QdrantClient) DeletePoints(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	body := qdrantDeleteRequest{Points: ids}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal delete: %w", err)
	}

	url := fmt.Sprintf("%s/collections/%s/points/delete?wait=true", c.baseURL, c.cfg.Collection)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant delete failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	c.logger.Debug("qdrant points deleted", zap.Int("count", len(ids)))
	return nil
}

func (c *QdrantClient) DeletePointsByFilter(ctx context.Context, filters map[string]string) error {
	if len(filters) == 0 {
		return fmt.Errorf("no filter provided for deletion")
	}

	body := qdrantDeleteRequest{Filter: c.buildFilter(filters)}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal delete filter: %w", err)
	}

	url := fmt.Sprintf("%s/collections/%s/points/delete?wait=true", c.baseURL, c.cfg.Collection)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete by filter request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant delete by filter failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	c.logger.Info("qdrant points deleted by filter", zap.Any("filters", filters))
	return nil
}

func (c *QdrantClient) GetCollectionInfo(ctx context.Context) (*qdrantCollectionInfo, error) {
	url := fmt.Sprintf("%s/collections/%s", c.baseURL, c.cfg.Collection)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get collection info failed: %d", resp.StatusCode)
	}

	var info qdrantCollectionInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode collection info: %w", err)
	}

	return &info, nil
}

func (c *QdrantClient) CountPoints(ctx context.Context) (int, error) {
	info, err := c.GetCollectionInfo(ctx)
	if err != nil {
		return 0, err
	}
	return info.Result.PointsCount, nil
}

func (c *QdrantClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/healthz", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("qdrant health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("qdrant unhealthy: status %d", resp.StatusCode)
	}
	return nil
}

func (c *QdrantClient) buildFilter(filters map[string]string) *qdrantFilter {
	f := &qdrantFilter{}
	for key, value := range filters {
		f.Must = append(f.Must, qdrantCondition{
			Key:   key,
			Match: qdrantMatch{Value: value},
		})
	}
	return f
}

func (c *QdrantClient) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	if c.cfg.APIKey != "" {
		req.Header.Set("api-key", c.cfg.APIKey)
	}
}

func (c *QdrantClient) Close() {
	c.httpClient.CloseIdleConnections()
}

func (s *RAGService) IngestDocument(ctx context.Context, provider providers.Provider, doc *Document) (string, error) {
	if doc.ID == "" {
		doc.ID = uuid.New().String()
	}

	if doc.Content == "" {
		return "", fmt.Errorf("document content is empty")
	}

	chunks := s.chunkDocument(doc)
	s.logger.Info("document chunked",
		zap.String("doc_id", doc.ID),
		zap.String("title", doc.Title),
		zap.Int("chunks", len(chunks)),
	)

	chunkTexts := make([]string, len(chunks))
	for i, chunk := range chunks {
		chunkTexts[i] = chunk.Content
	}

	batchEmbeddings, err := s.embeddings.GenerateBatchEmbeddings(ctx, provider, chunkTexts)
	if err != nil {
		return "", fmt.Errorf("generate chunk embeddings: %w", err)
	}

	points := make([]QdrantPoint, len(chunks))
	for i := range chunks {
		points[i] = QdrantPoint{
			ID:     chunks[i].ID,
			Vector: batchEmbeddings[i],
			Payload: map[string]any{
				"doc_id":    doc.ID,
				"title":     doc.Title,
				"content":   chunks[i].Content,
				"chunk_idx": chunks[i].Index,
				"metadata":  buildPayloadMetadata(doc.Metadata, chunks[i].Metadata),
			},
		}
	}

	if err := s.qdrant.UpsertPoints(ctx, points); err != nil {
		return "", fmt.Errorf("upsert chunks to qdrant: %w", err)
	}

	s.logger.Info("document ingested successfully",
		zap.String("doc_id", doc.ID),
		zap.String("title", doc.Title),
		zap.Int("chunks", len(chunks)),
	)

	return doc.ID, nil
}

func (s *RAGService) Search(ctx context.Context, embedding []float32, topK int) ([]SearchResult, error) {
	if topK <= 0 {
		topK = s.cfg.TopK
	}
	if topK <= 0 {
		topK = 5
	}

	threshold := s.cfg.SimilarityThreshold
	if threshold <= 0 {
		threshold = 0.7
	}

	hits, err := s.qdrant.Search(ctx, embedding, topK, nil, threshold)
	if err != nil {
		return nil, fmt.Errorf("qdrant search: %w", err)
	}

	results := make([]SearchResult, 0, len(hits))
	for _, hit := range hits {
		payload := hit.Payload
		title, _ := payload["title"].(string)
		content, _ := payload["content"].(string)
		docID, _ := payload["doc_id"].(string)

		metadata := extractStringMetadata(payload["metadata"])

		results = append(results, SearchResult{
			ID:       docID,
			Title:    title,
			Content:  content,
			Score:    hit.Score,
			Metadata: metadata,
		})
	}

	return results, nil
}

func (s *RAGService) SearchWithFilter(ctx context.Context, embedding []float32, topK int, filters map[string]string) ([]SearchResult, error) {
	if topK <= 0 {
		topK = s.cfg.TopK
	}
	if topK <= 0 {
		topK = 5
	}

	threshold := s.cfg.SimilarityThreshold
	if threshold <= 0 {
		threshold = 0.7
	}

	hits, err := s.qdrant.Search(ctx, embedding, topK, filters, threshold)
	if err != nil {
		return nil, fmt.Errorf("qdrant filtered search: %w", err)
	}

	results := make([]SearchResult, 0, len(hits))
	for _, hit := range hits {
		payload := hit.Payload
		title, _ := payload["title"].(string)
		content, _ := payload["content"].(string)
		docID, _ := payload["doc_id"].(string)
		metadata := extractStringMetadata(payload["metadata"])

		results = append(results, SearchResult{
			ID:       docID,
			Title:    title,
			Content:  content,
			Score:    hit.Score,
			Metadata: metadata,
		})
	}

	return results, nil
}

func (s *RAGService) SearchHybrid(ctx context.Context, embedding []float32, keyword string, topK int, filters map[string]string) ([]SearchResult, error) {
	if topK <= 0 {
		topK = s.cfg.TopK
	}
	if topK <= 0 {
		topK = 5
	}

	threshold := s.cfg.SimilarityThreshold
	if threshold <= 0 {
		threshold = 0.7
	}

	hits, err := s.qdrant.HybridSearch(ctx, embedding, keyword, topK, filters, threshold)
	if err != nil {
		return nil, fmt.Errorf("qdrant hybrid search: %w", err)
	}

	results := make([]SearchResult, 0, len(hits))
	for _, hit := range hits {
		payload := hit.Payload
		title, _ := payload["title"].(string)
		content, _ := payload["content"].(string)
		docID, _ := payload["doc_id"].(string)
		metadata := extractStringMetadata(payload["metadata"])

		results = append(results, SearchResult{
			ID:       docID,
			Title:    title,
			Content:  content,
			Score:    hit.Score,
			Metadata: metadata,
		})
	}

	return results, nil
}

func (s *RAGService) ListDocuments(ctx context.Context, offset, limit int) ([]Document, int, error) {
	info, err := s.qdrant.GetCollectionInfo(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get collection info: %w", err)
	}
	total := info.Result.PointsCount

	docs := make([]Document, 0)
	return docs, total, nil
}

func (s *RAGService) DeleteDocument(ctx context.Context, id string) error {
	filters := map[string]string{
		"doc_id": id,
	}
	if err := s.qdrant.DeletePointsByFilter(ctx, filters); err != nil {
		return fmt.Errorf("delete document chunks: %w", err)
	}

	s.logger.Info("document deleted", zap.String("doc_id", id))
	return nil
}

func (s *RAGService) DeleteAll(ctx context.Context) error {
	if err := s.qdrant.DeletePointsByFilter(ctx, map[string]string{}); err != nil {
		return fmt.Errorf("delete all points: %w", err)
	}
	return nil
}

func (s *RAGService) GetStats(ctx context.Context) (map[string]any, error) {
	info, err := s.qdrant.GetCollectionInfo(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"collection":      info.Result.Name,
		"points_count":    info.Result.PointsCount,
		"indexed_vectors": info.Result.IndexedVectors,
		"segments":        info.Result.SegmentsCount,
		"vector_size":     info.Result.Config.Params.Vectors.Size,
		"distance":        info.Result.Config.Params.Vectors.Distance,
	}, nil
}

func (s *RAGService) chunkDocument(doc *Document) []Chunk {
	chunkSize := s.cfg.ChunkSize
	if chunkSize <= 0 {
		chunkSize = 1500
	}
	chunkOverlap := s.cfg.ChunkOverlap
	if chunkOverlap < 0 {
		chunkOverlap = 200
	}

	chunks := RecursiveChunk(doc.Content, chunkSize, chunkOverlap)

	for i := range chunks {
		chunks[i].DocID = doc.ID
		chunks[i].ID = fmt.Sprintf("%s_chunk_%d", doc.ID, i)
		chunks[i].Index = i
		if chunks[i].Metadata == nil {
			chunks[i].Metadata = make(map[string]string)
		}
		chunks[i].Metadata["title"] = doc.Title
		for k, v := range doc.Metadata {
			chunks[i].Metadata[k] = v
		}
	}

	return chunks
}

func buildPayloadMetadata(docMeta, chunkMeta map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range docMeta {
		merged[k] = v
	}
	for k, v := range chunkMeta {
		merged[k] = v
	}
	return merged
}

func extractStringMetadata(raw any) map[string]string {
	result := make(map[string]string)
	if raw == nil {
		return result
	}

	switch m := raw.(type) {
	case map[string]any:
		for k, v := range m {
			if s, ok := v.(string); ok {
				result[k] = s
			} else {
				result[k] = fmt.Sprintf("%v", v)
			}
		}
	case map[string]string:
		return m
	}

	return result
}

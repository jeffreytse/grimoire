---
name: design-rag-pipeline
description: Use when designing a retrieval-augmented generation (RAG) pipeline for grounding LLM outputs in external knowledge
source: Lewis et al. "Retrieval-Augmented Generation for Knowledge-Intensive NLP Tasks" (Facebook AI, NeurIPS 2020); LangChain documentation; Pinecone "RAG best practices" guide
tags: [rag, llm, embeddings, vector-database, ai, retrieval, langchain]
verified: true
---

# Design RAG Pipeline

Design a retrieval-augmented generation pipeline that retrieves relevant context and grounds LLM responses in verifiable source documents.

## Why This Is Best Practice

**Adopted by:** OpenAI (GPT with Retrieval), Microsoft (Azure AI Search + OpenAI), Anthropic (Claude with tool use for retrieval), LangChain ecosystem
**Impact:** RAG reduces LLM hallucination rates by 40-60% on knowledge-intensive tasks compared to vanilla generation (Lewis et al., 2020); enables knowledge cutoff extension and source citation without fine-tuning.

**Why best:** RAG separates parametric knowledge (what the model learned) from non-parametric knowledge (what can be retrieved). This allows updating the knowledge base without retraining, enables source attribution, and grounds outputs in verifiable documents — critical for enterprise and compliance use cases.

## Steps

1. **Design the ingestion pipeline** — Load documents → chunk → embed → store in vector DB. Chunking strategy critically affects retrieval quality; prefer semantic/paragraph chunking over fixed-size token chunking.
2. **Choose chunk size and overlap** — For dense technical text: 512 tokens, 10% overlap. For narrative text: 256 tokens. Test with your retrieval evaluation set — chunk size is a hyperparameter, not a constant.
3. **Select the embedding model** — Match embedding model to domain: `text-embedding-3-large` (OpenAI, general), `bge-large-en-v1.5` (BAAI, strong retrieval), domain-specific fine-tuned embeddings for specialized corpora. Embedding model and retrieval model must use the same model.
4. **Choose a vector database** — Pinecone (managed, production), Weaviate (open-source, hybrid search), pgvector (Postgres extension, simple stack), Chroma (local dev). Choose based on scale, hybrid search need, and operational complexity.
5. **Implement hybrid search** — Combine dense (vector) and sparse (BM25/keyword) retrieval with Reciprocal Rank Fusion (RRF). Pure vector search misses exact keyword matches; hybrid consistently outperforms either alone.
6. **Apply re-ranking** — After retrieving top-K candidates (K=20), re-rank with a cross-encoder (Cohere Rerank, bge-reranker) to return the top-N (N=5) most relevant chunks. Improves precision significantly.
7. **Construct the prompt** — Inject retrieved context with source attribution. Include instructions to cite sources and to say "I don't know" when context is insufficient. Set `max_tokens` for context window budget.

## Rules

- Always evaluate retrieval quality separately from generation quality — a good generator cannot compensate for poor retrieval.
- Store document metadata (source URL, last-updated date, section) alongside embeddings for citation and filtering.
- Implement a re-indexing strategy — stale documents produce stale answers; trigger re-indexing on document update.
- Never pass the full document to the LLM — retrieve the minimum context needed; context window bloat increases cost and reduces attention on the relevant passage.

## Examples

Pipeline: S3 bucket (PDFs) → `unstructured` parser → recursive character splitter (512T/50T overlap) → `text-embedding-3-large` → Pinecone (with BM25 hybrid) → Cohere Rerank top-5 → Claude with system prompt instructing citation.

Evaluation: RAGAS framework measuring faithfulness, answer relevancy, context precision, and context recall on a golden QA set.

## Common Mistakes

- **Fixed chunk size without evaluation** — chunk size is domain-specific; test against real queries before fixing.
- **No metadata filtering** — users querying 2024 docs get 2019 docs mixed in; filter by date, source, or category before vector search.
- **Skipping re-ranking** — returning the raw top-K from vector search includes many irrelevant chunks; re-ranking is low-cost, high-impact.

## When NOT to Use

- When the LLM's training data already covers the knowledge domain comprehensively and the knowledge does not change frequently, fine-tuning or prompt engineering is simpler and cheaper than building a retrieval pipeline.
- When the corpus is fewer than a few hundred short documents, a full vector database + embedding pipeline is over-engineered; load all documents directly into context or use keyword search instead.
- When latency requirements are under 200ms end-to-end, the retrieval + re-ranking round-trip typically adds 300-800ms and RAG should not be applied without explicit latency budgeting.

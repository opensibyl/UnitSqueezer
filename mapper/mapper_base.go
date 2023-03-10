package mapper

import (
	"context"

	"github.com/dominikbraun/graph"
	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/squ/indexer"
	"github.com/opensibyl/squ/log"
	"github.com/opensibyl/squ/object"
)

type BaseMapper struct {
	indexerRef indexer.Indexer
}

func (baseMapper *BaseMapper) SetIndexer(i indexer.Indexer) {
	baseMapper.indexerRef = i
}

func (baseMapper *BaseMapper) GetRelatedCaseSignatures(_ context.Context, targetSignature string) (map[string]interface{}, error) {
	g := baseMapper.indexerRef.GetSibylCache().ReverseCallGraph

	// all the vertexes related to this signature
	vertexes := baseMapper.indexerRef.GetVertexesWithSignature(targetSignature)

	caseSet := baseMapper.indexerRef.GetCaseSet()
	matchedCases := make(map[string]interface{}, 0)
	for _, eachV := range vertexes {
		err := graph.BFS(g.Graph, eachV, func(k string) bool {
			functionWithPath, err := g.Graph.Vertex(k)
			if err != nil {
				log.Log.Warnf("access %v failed: %v", k, err)
				return true
			}
			s := functionWithPath.GetSignature()

			if _, ok := caseSet[s]; ok {
				// reach
				matchedCases[s] = nil
			}
			return false
		})
		if err != nil {
			return nil, err
		}
	}
	if len(matchedCases) == 0 {
		log.Log.Warnf("no cases related to %v", targetSignature)
	}
	return matchedCases, nil
}

func (baseMapper *BaseMapper) Diff2Cases(ctx context.Context, diffMap object.DiffFuncMap) ([]*openapi.ObjectFunctionWithSignature, error) {
	caseSetToRun := make(map[string]interface{})
	for fileName, eachFunctionList := range diffMap {
		log.Log.Infof("handle modified file: %s, functions: %d", fileName, len(eachFunctionList))
		for _, eachFunc := range eachFunctionList {
			cases, err := baseMapper.GetRelatedCaseSignatures(ctx, eachFunc.GetSignature())
			if err != nil {
				return nil, err
			}
			// no cases can reach
			if len(cases) == 0 {
				eachFunc.Reachable = false
			} else {
				eachFunc.Reachable = true
			}

			for k := range cases {
				// merge
				caseSetToRun[k] = nil
				// record
				eachFunc.ReachBy = append(eachFunc.ReachBy, k)
			}
		}
	}
	casesToRun := make([]*openapi.ObjectFunctionWithSignature, 0)
	for eachCase := range caseSetToRun {
		functionWithSignature, err := baseMapper.indexerRef.GetFuncWithSignature(ctx, eachCase)
		if err != nil {
			return nil, err
		}
		casesToRun = append(casesToRun, functionWithSignature)
	}
	return casesToRun, nil
}

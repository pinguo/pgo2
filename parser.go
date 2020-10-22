package pgo2

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/pinguo/pgo2/util"
)

var controllerActionInfo = make(map[string]map[string][]*ActionInfo)

type ActionInfo struct {
	ControllerName string                      // 控制器名
	Name           string                      // action名
	PkgPath        string                      // 包路径
	PkgName        string                      // 包名
	Desc           string                      // action描述
	ParamsDesc     map[string]*ActionInfoParam // 参数描述
}

type ActionInfoParam struct {
	Name     string
	DftValue string
	NameType string
	Usage    string
}

var gParse *Parser

func NewParser() *Parser {
	if gParse != nil {
		return gParse
	}
	return &Parser{}
}

type Parser struct {
}

func (p *Parser) Dir(path string, filterFunc ...func(os.FileInfo) bool) map[string]*ast.Package {
	fileToken := token.NewFileSet()
	var filter func(os.FileInfo) bool
	if len(filterFunc) > 0 {
		filter = filterFunc[0]
	}

	pkgs, err := parser.ParseDir(fileToken, path, filter, parser.ParseComments)
	if err != nil {
		panic("parserDir:" + path + ",err:" + err.Error())
	}

	return pkgs
}

func (p *Parser) GetActionInfo(pkgPath, controllerName, actionName string) *ActionInfo {
	pkgRealPath := p.pkgRealPath(App().BasePath(), pkgPath)
	p.InitActionInfo(pkgRealPath, pkgPath)

	pkgInfo := controllerActionInfo[pkgPath]
	if pkgInfo == nil {
		return nil
	}

	cInfo, has := controllerActionInfo[pkgPath][controllerName]
	if !has {
		return nil
	}

	// Prepare params
	var prepareParams *ActionInfo
	// Action params
	var aInfo *ActionInfo
	for _, tmpInfo := range cInfo {
		switch tmpInfo.Name {
		case PrepareMethod:
			prepareParams = tmpInfo
		case actionName:
			aInfo = tmpInfo
		}

		if prepareParams != nil && aInfo != nil {
			break
		}
	}

	if prepareParams == nil || prepareParams.ParamsDesc == nil {
		return aInfo
	}

	// merge
	for k, v := range prepareParams.ParamsDesc {
		aInfo.ParamsDesc[k] = v
	}

	return aInfo
}

// pkgRealPath get real path
func (p *Parser) pkgRealPath(appPath, pkgPath string) string {

	d := appPath + pkgPath[strings.Index(pkgPath, "/"):]
	if exist, err := util.Exist(d); exist == false || err != nil {
		return ""
	}

	return d
}

func (p *Parser) InitActionInfo(pkgRealPath, pkgPath string) {
	if _, has := controllerActionInfo[pkgPath]; has {
		return
	}

	controllerActionInfo[pkgPath] = p.ActionInfo(pkgRealPath, pkgPath)

	return
}

func (p *Parser) ActionInfo(pkgRealPath, pkgPath string) map[string][]*ActionInfo {
	filter := func(f os.FileInfo) bool {
		if f.Name() == "init.go" {
			return false
		}

		return true
	}

	ret := make(map[string][]*ActionInfo)
	pkgS := p.Dir(pkgRealPath, filter)
	for pkgName, pkg := range pkgS {
		//fmt.Println("pkgName", pkgName)
		for _, pFile := range pkg.Files {
			for _, decl := range pFile.Decls {
				switch specDecl := decl.(type) {
				case *ast.FuncDecl:
					if specDecl.Recv == nil {
						break
					}

					exp, ok := specDecl.Recv.List[0].Type.(*ast.StarExpr) // Check that the type is correct first beforing throwing to parser
					if !ok {
						break
					}

					controllerName := fmt.Sprint(exp.X)
					methodName := specDecl.Name.String()

					desc := p.parserActionDesc(specDecl.Doc)
					params := p.parserCommentParams(specDecl.Doc)
					if len(params) == 0 {
						params = p.parserFlagParams(specDecl.Body)
					}

					if _, has := ret[controllerName]; !has {
						ret[controllerName] = make([]*ActionInfo, 0, 1)
					}
					ret[controllerName] = append(ret[controllerName], &ActionInfo{
						ControllerName: controllerName,
						Name:           methodName,
						PkgPath:        pkgPath,
						PkgName:        pkgName,
						Desc:           desc,
						ParamsDesc:     params,
					})
				}
			}

		}

	}

	return ret
}

// parserActionDesc
func (p *Parser) parserActionDesc(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}

	keyWord := "@ActionDesc"
	for _, v := range doc.List {
		if pos := strings.Index(v.Text, keyWord); pos >= 0 {
			return strings.Trim(v.Text[pos+len(keyWord):], " ")
		}
	}

	return ""
}

// parserCommentParams
func (p *Parser) parserCommentParams(doc *ast.CommentGroup) map[string]*ActionInfoParam {
	ret := make(map[string]*ActionInfoParam)
	if doc == nil {
		return ret
	}

	keyWord := "@Params"
	for _, v := range doc.List {
		if pos := strings.Index(v.Text, keyWord); pos >= 0 {
			msg := strings.Trim(v.Text[pos+len(keyWord):], " ")
			tmpArr := strings.Split(msg, " ")
			name := strings.TrimLeft(tmpArr[0], "-")
			name = strings.TrimLeft(name, "-")
			name = strings.TrimLeft(name, " ")
			if name == "" {
				continue
			}

			usage := strings.TrimLeft(msg, tmpArr[0])
			ret[name] = &ActionInfoParam{
				Name:  name,
				Usage: usage,
			}

		}
	}

	return ret
}

// parserFlagParams parse flag Params
func (p *Parser) parserFlagParams(body *ast.BlockStmt) map[string]*ActionInfoParam {
	if body == nil {
		return nil
	}
	ret := make(map[string]*ActionInfoParam)
	for _, stmt := range body.List {
		switch expr := stmt.(type) {
		case *ast.AssignStmt:
			for _, expr := range expr.Rhs {
				if actionInfoParam := p.ParseCallExpr(expr, "flag", 3); actionInfoParam != nil {
					ret[actionInfoParam.Name] = actionInfoParam
				}
			}
		case *ast.ExprStmt:
			if actionInfoParam := p.ParseCallExpr(expr.X, "flag", 3); actionInfoParam != nil {
				ret[actionInfoParam.Name] = actionInfoParam
			}
		}
	}

	return ret
}

func (p *Parser) ParseCallExpr(expr ast.Expr, pkgName string, minLenArg int) *ActionInfoParam {
	callExpr, callOk := expr.(*ast.CallExpr)
	if !callOk {
		return nil
	}

	SelectorExpr, selectorOk := callExpr.Fun.(*ast.SelectorExpr)
	if !selectorOk {
		return nil
	}

	pkgIdent, pkgOk := SelectorExpr.X.(*ast.Ident)
	if !pkgOk {
		return nil
	}

	if pkgName != "" && pkgName != pkgIdent.Name {
		return nil
	}

	lenArg := len(callExpr.Args)
	if minLenArg > 0 && lenArg < minLenArg {
		return nil
	}

	funcName := SelectorExpr.Sel.Name
	funcFilterVar := "Var"
	namePos, dftValuePos, usagePos := 3, 2, 1
	nameType := strings.Replace(funcName, funcFilterVar, "", 1)

	tmpArgs := make([]string, lenArg)
	tmpType := make([]string, lenArg)

	parerNameType := func(argVIdent *ast.Ident) (tmpTypeValue string) {
		if argVIdent.Obj != nil {
			if argVValueSpec, ok := argVIdent.Obj.Decl.(*ast.ValueSpec); ok {
				if argvIdent2, ok := argVValueSpec.Type.(*ast.Ident); ok {
					tmpTypeValue = argvIdent2.Name
					if argvIdent2.Obj != nil {
						if argvTypeSpec, ok := argvIdent2.Obj.Decl.(*ast.TypeSpec); ok {
							if argVIdent3, ok := argvTypeSpec.Type.(*ast.Ident); ok {
								tmpTypeValue = argVIdent3.Name
							}

						}
					}
				}
			}

		}
		tmpTypeValue = strings.Title(tmpTypeValue)
		return
	}

	for tmpArgKey, arg := range callExpr.Args {
		tmpValue := ""
		tmpTypeValue := ""
		switch argV := arg.(type) {
		case *ast.BasicLit:
			tmpValue = argV.Value
			tmpTypeValue = argV.Kind.String()
		case *ast.Ident:
			tmpValue = argV.Name
			tmpTypeValue = parerNameType(argV)
		case *ast.UnaryExpr:
			argVIdent, ok := argV.X.(*ast.Ident)
			if ok {
				tmpTypeValue = parerNameType(argVIdent)
			}

		}

		tmpArgs[tmpArgKey] = strings.Trim(tmpValue, "\"")
		tmpType[tmpArgKey] = tmpTypeValue

	}

	if funcName == funcFilterVar {
		nameType = tmpType[0]
		tmpArgs[lenArg-namePos] = tmpArgs[lenArg-dftValuePos]
		tmpArgs[lenArg-dftValuePos] = "'unknown'"

	}

	return &ActionInfoParam{
		Name:     tmpArgs[lenArg-namePos],
		NameType: nameType,
		DftValue: tmpArgs[lenArg-dftValuePos],
		Usage:    tmpArgs[lenArg-usagePos],
	}
}

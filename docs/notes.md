# GoLisp 開発ノート

## 開発方針

- **実装は自分でタイピング**: 学習目的のため、コードは自分で書く
- **Claudeの役割**: 設計レビュー、相談、デバッグ支援、アーキテクチャ提案
- **ドキュメント**: SPECIFICATION.mdとdesign/M1-architecture.mdは参照・設計用

## 進捗

- [x] 仕様書作成 (SPECIFICATION.md)
- [x] M1アーキテクチャ設計 (design/M1-architecture.md)
- [x] プロジェクト初期化 (go mod init)
- [x] internal/types/types.go 実装完了
  - Number, Symbol, Nil, String, Boolean, Cons
  - 各型のString()メソッド実装
- [X] internal/reader/lexer.go 実装
- [X] internal/reader/lexer.goに対応したテストを書く
- [x] internal/reader/parser.go 実装完了
  - parseExpr(): NUMBER, STRING, SYMBOL, LPAREN, QUOTE, エラー処理
  - parseList(): リストのパース、空リスト対応
  - Boolean特殊処理: 't' → Boolean, 'nil' → Nil
  - Quote展開: 'x → (quote x)
- [x] internal/reader/parser_test.go 実装完了
  - 全テストパス
  - Number, String, Symbol, Boolean, List, NestedList, Quote, Errors, MultilineString
- [x] internal/repl/repl.go 実装完了
  - Read-Eval-Print-Loop (M1ではEvalは恒等関数)
  - エラー時も継続実行
  - Ctrl+Dで終了
- [x] cmd/gospl/main.go 実装完了


詳細は `design/M1-architecture.md` を参照。

## メモ

- M1完成: S式のパース・表示までできるようになった
- 次はEvaluatorとEnvironmentの実装

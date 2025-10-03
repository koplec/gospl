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
- [ ] internal/reader/lexer.go 実装
- [ ] internal/reader/parser.go 実装
- [ ] internal/repl/repl.go 実装
- [ ] main.go 実装

## 次回やること

1. `internal/reader/lexer.go` の実装
   - TokenType, Token, Position 型の定義
   - Lexer構造体の実装
   - NextToken()メソッドの実装
2. `internal/reader/parser.go` の実装
3. `internal/repl/repl.go` の実装
4. `main.go` の実装

詳細は `design/M1-architecture.md` を参照。

## メモ

- (ここに気づいたことや疑問点を記録)

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
- [ ] internal/reader/parser.go 実装
- [ ] internal/reader/parser.goに対応したテストを書く
- [ ] internal/repl/repl.go 実装
- [ ] internal/reader/repl.goに対応したテストを書く
- [ ] main.go 実装
- [ ] main.goに対応したテストを書く

詳細は `design/M1-architecture.md` を参照。

## メモ

- (ここに気づいたことや疑問点を記録)

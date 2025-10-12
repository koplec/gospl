# GoLisp 開発ノート

## 開発方針

- **実装は自分でタイピング**: 学習目的のため、コードは自分で書く
- **Claudeの役割**: 設計レビュー、相談、デバッグ支援、アーキテクチャ提案
- **ドキュメント**: SPECIFICATION.mdとdesign/M1-architecture.mdは参照・設計用

## 進捗

### M1: 基本的なS式とREPL ✅ 完了 (2025-10-12)

- [x] 仕様書作成 (SPECIFICATION.md)
- [x] M1アーキテクチャ設計 (design/M1-architecture.md)
- [x] プロジェクト初期化 (go mod init)
- [x] internal/types/types.go 実装完了
  - Number, Symbol, Nil, String, Boolean, Cons
  - 各型のString()メソッド実装
- [x] internal/reader/lexer.go 実装完了
- [x] internal/reader/lexer_test.go 実装完了
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

**動作確認**: REPLが正常に起動し、S式の入力・出力が動作することを確認

### M2: 基本的な評価器と算術 ✅ 完了 (2025-10-13)

- [x] M2アーキテクチャ設計 (design/M2-architecture.md)
- [x] internal/eval/env.go 実装完了
  - Environment構造体
  - Set/Getメソッド
  - NewGlobalEnvironment
- [x] internal/eval/builtins.go 実装完了
  - BuiltinFn型エイリアス
  - BuiltinFunc構造体
  - 算術関数: +, -, *, /
  - Common Lisp準拠（引数チェック、ゼロ除算）
- [x] internal/eval/builtins_test.go 実装完了
  - 正常系テスト
  - エラーケーステスト（型エラー、ゼロ除算）
- [x] internal/eval/eval.go 実装完了
  - Eval関数（自己評価、シンボル参照、関数適用）
  - evalList、evalArgs、apply
- [x] internal/eval/eval_test.go 実装完了
  - 統合テスト（数値、算術演算、ネスト）
  - エラーケーステスト
- [x] internal/repl/repl.go 更新
  - 評価器の統合

**動作確認**: `(+ 1 2)` → `3` など、算術演算が正常に動作

**学び**:
- evalArgsでの無限ループバグ（current = cons.Cdrを忘れていた）
- テスト駆動開発の重要性を実感

### M3: 関数とクロージャ (次のステップ)

詳細は `design/M1-architecture.md`, `design/M2-architecture.md` を参照。

## メモ

- M1完成: S式のパース・表示までできるようになった
- M2完成: 評価器と算術演算が動くようになった
- 次はlambda、defun、クロージャの実装
